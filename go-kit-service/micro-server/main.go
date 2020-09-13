package main

import (
	"flag"
	"fmt"
	"log"
	"golang.org/x/time/rate"
	kitlog "github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"micro-server/Services"
	"micro-server/util"

	"github.com/gorilla/mux"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	name := flag.String("name", "","服务名称")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名")
	}
	if *port == 0 {
		log.Fatal("请指定端口号")
	}
	util.SetServicePortAndName(*name, *port)

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "micro-srv", "1.0")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	//accessservice
	accessService := Services.AccessService{}
	accessEndpoint := Services.AccessEndpoint(accessService)
	accessHandler := httptransport.NewServer(accessEndpoint, Services.DecodeAccessRequest, Services.EncodeUserResponse)

	//userservice
	userService := Services.UserService{}
	limit := rate.NewLimiter(1, 3)
	UserendPoint := Services.RateLimit(limit)(Services.SrvLogger(logger)(Services.UserAuth()(Services.GenUserEndpoint(userService))))
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(Services.AppErrorEncoder),
	}
	serverHandler := httptransport.NewServer(UserendPoint, Services.DecodeUserRequest, Services.EncodeUserResponse, options...)

	r := mux.NewRouter()

	r.Methods("POST").Path("/access/token").Handler(accessHandler)

	r.Methods("GET", "DELETE").Path(`/user/{userid:\d+}`).Handler(serverHandler)

	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	errC := make(chan error)

	go func() {
		util.RegisterService()
		err := http.ListenAndServe(":"+strconv.Itoa(*port), r)
		if err != nil {
			errC <- err
		}
	}()

	go func() {
		signC := make(chan os.Signal)
		signal.Notify(signC, syscall.SIGINT, syscall.SIGTERM)
		errC <- fmt.Errorf("%s", <-signC)
	}()

	getErr := <-errC
	util.UnRegisterService()
	fmt.Println(getErr)
}