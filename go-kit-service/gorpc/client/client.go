package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/client/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"time"
)

func main()  {
	cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{ cert },
		ServerName: "localhost",
		RootCAs: certPool,
	})
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := prod.NewOrderServiceClient(conn)
	t := timestamp.Timestamp{ Seconds: time.Now().Unix()}
	response, err := client.NewOrder(context.Background(), &prod.OrderMain{
		OrderId: 101,
		OrderNo: "nd908323",
		UserId: 123,
		OrderMoney: 22.5,
		OrderTime: &t,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
