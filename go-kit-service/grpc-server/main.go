package main

import (
	"github.com/TyrellJing/Hermes/go-kit-service/grpc-server/pb"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"net"
)

func main() {

	BookServer := new(BookServiceServer)

	bookListHandler := kitgrpc.NewServer(
		BookListEndpoint(),
		dencodeServiceRequest,
		encodeServiceResponse,
	)
	BookServer.BookListHandler = bookListHandler

	bookInfoHandler := kitgrpc.NewServer(
		BookInfoEndPoint(),
		dencodeServiceRequest,
		encodeServiceResponse,
	)
	BookServer.BookInfoHandler = bookInfoHandler

	ln, _ := net.Listen("tcp", ":8081")
	gs := grpc.NewServer()
	pb.RegisterBookServiceServer(gs, BookServer)
	gs.Serve(ln)
}
