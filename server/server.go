package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"github.com/TyrellJing/Hermes/router"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"syscall"
	"time"
)

const (
	GRACEFUL_ENVIRON_KEY 	= "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD 	= 3
	REQUEST_TIMEOUT 		= 300 * time.Second
)

type Server struct {
	httpServer  	*http.Server
	listener 		net.Listener

	isGraceful 		bool
	signalChan 		chan os.Signal
	shutdownChan 	chan bool
}

func NewServer(handler router.Router, port string) *Server {
	server := new(Server)

	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true
	}

	server.httpServer = &http.Server{
		Handler: http.TimeoutHandler(&handler, REQUEST_TIMEOUT, "time out"),
		Addr:	 ":" + port,
	}
	server.isGraceful = isGraceful
	server.signalChan = make(chan os.Signal)

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
		syscall.SIGUSR1,
		syscall.SIGHUP,
		syscall.SIGUSR2,
	)

	for {
		sig = <-srv.signalChan
		switch sig {
		case syscall.SIGHUP:
			srv.logf("received SIGHUP, graceful restarting HTTP server")

			if pid, err := srv.startNew(); err != nil {
				srv.logf("start new process failed: " + err.Error() + ", continue serving")
			} else {
				ppid := strconv.FormatUint(uint64(pid), 10)
				srv.logf("start new process successed, the new pid is " + ppid)

				srv.shutDown()
			}
		case syscall.SIGUSR2:
			buf := make([]byte, 1638400)
			buf = buf[:runtime.Stack(buf, true)]
			fmt.Println(string(buf))
			srv.logf("goroutine stack output to stack.log")

			memory, err := os.OpenFile("/dev/shm/memory.log", os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				srv.logf(err.Error())
			}
			defer memory.Close()
			pprof.WriteHeapProfile(memory)

			cpu, err := os.OpenFile("/dev/shm/cpu.log", os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				srv.logf(err.Error())
			}
			pprof.StartCPUProfile(cpu)

			time.Sleep(30 * time.Second)
			pprof.StopCPUProfile()
			cpu.Close()
		default:
		}
	}
}

func (srv *Server) shutDown() {
	if err := srv.httpServer.Shutdown(context.Background()); err != nil {
		srv.logf("HTTP server shutdown error: " + err.Error())
	} else {
		srv.logf("HTTP server shutdown success.")
		srv.shutdownChan <- true
	}
}

func (srv *Server) startNew() (uintptr, error) {
	file, err := srv.listener.(*net.TCPListener).File()
	if err != nil {
		return 0, err
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
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), file.Fd()},
	}

	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return uintptr(fork), nil
}

func (srv *Server) logf(format string) {
	pid := strconv.Itoa(os.Getpid())
	format = "[pid:" + pid + "] " + format

	fmt.Println(format)
}






