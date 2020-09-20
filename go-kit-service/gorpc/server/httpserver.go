package main

import (
	"context"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
	"github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/util"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	gwmux := runtime.NewServeMux()
	opt := []grpc.DialOption {
		grpc.WithTransportCredentials(util.GetClientCert()),
	}
	grpcEndPoint := "localhost:8081"
	err := prod.RegisterProdServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndPoint, opt)
	if err != nil {
		log.Fatal(err)
	}
	err = prod.RegisterOrderServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndPoint, opt)
	if err != nil {
		log.Fatal(err)
	}
	httpSrv := http.Server{
		Addr: ":8080",
		Handler: gwmux,
	}
	httpSrv.ListenAndServe()
}

