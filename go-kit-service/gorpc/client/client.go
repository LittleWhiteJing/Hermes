package main

import (
	"context"
	"fmt"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/client/proto"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := prod.NewProdServiceClient(conn)
	response, err := client.GetProdStock(context.Background(), &prod.ProdRequest{
		ProdId: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.ProdStock)
}
