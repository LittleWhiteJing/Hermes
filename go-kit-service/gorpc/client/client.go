package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
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

	client := prod.NewProdServiceClient(conn)
	response, err := client.GetProdStock(context.Background(), &prod.ProdRequest{
		ProdId: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.ProdStock)
}
