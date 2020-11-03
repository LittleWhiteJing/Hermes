package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/TyrellJing/Hermes/radix-tree"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	GRACEFUL_ENVIRON_KEY 	= "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD 	= 3
)

type Server struct {
	httpServer  	*http.Server
	listener 		net.Listener

	isGraceful 		bool
	signalChan 		chan os.Signal
	shutdownChan 	chan bool

	TLSConfig		*tls.Config

}

func NewServer(handler radix_tree.Router, addr string, readTimeout, writeTimeout time.Duration) *Server {
	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true
	}

	server := new(Server)
	server.httpServer = &http.Server{
		Addr:	 		addr,
		Handler: 		&handler,
		ReadTimeout:	readTimeout,
		WriteTimeout:	writeTimeout,
	}
	server.isGraceful = isGraceful
	server.signalChan = make(chan os.Signal)
	server.shutdownChan = make(chan bool)

	return server
}

func (srv *Server) ListenAndServe() (err error) {
	addr := srv.httpServer.Addr

	var ln net.Listener
	ln, err = srv.getNetListener(addr)

	if err != nil {
		return
	}

	srv.listener = ln
	srv.logf("server started!")

	return srv.Serve()
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) (err error) {
	addr := srv.httpServer.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}

	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return
	}
	srv.listener = tls.NewListener(ln, config)
	return srv.Serve()
}


func (srv *Server) getNetListener(addr string) (ln net.Listener, err error) {
	if srv.isGraceful {
		file := os.NewFile(GRACEFUL_LISTENER_FD, "")
		ln, err = net.FileListener(file)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return
		}
	}
	return
}

func (srv *Server) Serve() (err error) {
	go srv.handleSignals()

	err = srv.httpServer.Serve(srv.listener)

	srv.logf("waiting for connections closed.")
	<-srv.shutdownChan
	srv.logf("all connections closed bye!")

	return
}


func (srv *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)

	for {
		sig = <-srv.signalChan
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			srv.logf("received SIGTERM, graceful shutting down HTTP server.")
			srv.shutdownHTTPServer()
		case syscall.SIGUSR2:
			srv.logf("received SIGUSR2, graceful restarting HTTP server.")

			if pid, err := srv.startNewProcess(); err != nil {
				srv.logf("start new process failed: %v, continue serving.", err)
			} else {
				srv.logf("start new process successed, the new pid is %d", pid)
				srv.shutdownHTTPServer()
			}
		default:
		}
	}
}

func (srv *Server) shutdownHTTPServer() {
	if err := srv.httpServer.Shutdown(context.Background()); err != nil {
		srv.logf("HTTP server shutdown error: " + err.Error())
	} else {
		srv.logf("HTTP server shutdown success.")
		srv.shutdownChan <- true
	}
}

func (srv *Server) startNewProcess() (fd uintptr, err error) {
	listenerFd, err := srv.getTCPListenerFd()
	if err != nil {
		return
	}

	envs := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			envs = append(envs, value)
		}
	}
	envs = append(envs, GRACEFUL_ENVIRON_STRING)

	execSpec := &syscall.ProcAttr{
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return uintptr(fork), nil
}

func (srv *Server) getTCPListenerFd() (fd uintptr, err error) {
	file, err := srv.listener.(*net.TCPListener).File()
	if err != nil {
		return
	}
	return file.Fd(), nil
}

func (srv *Server) logf(format string, args ...interface{}) {
	pid := strconv.Itoa(os.Getpid())
	format = "[pid:" + pid + "] " + format

	log.Printf(format, args...)
}






