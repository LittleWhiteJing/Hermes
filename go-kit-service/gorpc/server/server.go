package main

import (
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
	"github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/service"
	"github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/util"
	"google.golang.org/grpc"
	"net"
)


func main()  {
	rpcSrv := grpc.NewServer(grpc.Creds(util.GetServerCert()))
	prod.RegisterProdServiceServer(rpcSrv, new(service.ProdService))

	listener, _ := net.Listen("tcp", ":8081")
	rpcSrv.Serve(listener)
}



