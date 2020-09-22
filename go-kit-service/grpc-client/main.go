package main

import (
	"context"
	"fmt"
	"github.com/TyrellJing/Hermes/go-kit-service/grpc-client/pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	serviceAddr := "127.0.0.1:8081"
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	BookClient := pb.NewBookServiceClient(conn)
	BookListRequest := &pb.BookListRequest {
		PageNum: 1,
		PerPage: 20,
	}
	BookList, err := BookClient.GetBookList(context.Background(), BookListRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(BookList)

	BookDetailRequest := &pb.BookDetailRequest{
		BookId: 1,
	}
	BookDetail, err := BookClient.GetBookInfo(context.Background(), BookDetailRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(BookDetail)
}
