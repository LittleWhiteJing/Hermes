package main

import (
	"google.golang.org/grpc"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
	"net"
)


func main()  {
	rpcSrv := grpc.NewServer()
	prod.RegisterProdServiceServer(rpcSrv, new(ProdService))

	listener, _ := net.Listen("tcp", ":8081")
	rpcSrv.Serve(listener)
}


