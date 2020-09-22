package main

import (
	"flag"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"log"
	"micro-server/internal/endpoint"
	"micro-server/internal/service"
	"micro-server/internal/transport"
	"micro-server/middleware"
	"micro-server/registry"
	"micro-server/registry/consul"
	"micro-server/util"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	name := flag.String("name", "","Service Name")
	port := flag.Int("port", 0, "Service Port")
	flag.Parse()
	if *name == "" {
		log.Fatal("Please Set Service Name")
	}
	if *port == 0 {
		log.Fatal("Please Set Service Port")
	}
	logger := util.GetLogger()

	//accessservice
	accessService := service.AccessService{}
	accessEndpoint := endpoint.AccessEndpoint(accessService)
	accessHandler := httptransport.NewServer(accessEndpoint, transport.DecodeAccessRequest, transport.EncodeUserResponse)

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "raysonxin",
		Subsystem: "arithmetic_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "raysonxin",
		Subsystem: "arithemetic_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	//userservice
	userService := service.UserService{}
	limit := rate.NewLimiter(1, 3)
	UserendPoint := middleware.SrvMonitor(requestCount, requestLatency)(middleware.RateLimit(limit)(middleware.SrvLogger(logger)(middleware.JwtAuth()(endpoint.GenUserEndpoint(userService)))))
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(endpoint.AppErrorEncoder),
	}
	serverHandler := httptransport.NewServer(UserendPoint, transport.DecodeUserRequest, transport.EncodeUserResponse, options...)

	r := mux.NewRouter()

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("POST").Path("/access/token").Handler(accessHandler)

	r.Methods("GET", "DELETE").Path(`/user/{userid:\d+}`).Handler(serverHandler)

	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	errC := make(chan error)

	serviceID := string(time.Now().UnixNano())
	consulOpt, err := consul.NewConsulRegister(consul.SetConsulClient("127.0.0.1:8500"))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		var err error
		srv := registry.Service{
			ServiceID: serviceID,
			ServiceName: *name,
			ServicePort: *port,
			ServiceAddr: "192.168.1.104",
			ServiceTags: []string{"primary"},
		}
		srvCheck := consul.ServiceCheck{
			CheckAddr: fmt.Sprintf("http://%s:%d/health", srv.ServiceAddr, srv.ServicePort),
			CheckIntval: "5s",
		}
		err = consulOpt.RegisterService(srv, srvCheck)
		err = http.ListenAndServe(":"+strconv.Itoa(*port), r)
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
	consulOpt.UnRegisterService(serviceID)
	fmt.Println(getErr)
}