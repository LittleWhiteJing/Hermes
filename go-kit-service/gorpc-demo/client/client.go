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
	"io"
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

	getServerClientStreamResponse(conn)
}

func getGrpcTimestamp(conn *grpc.ClientConn) {
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

func getServerStreamResponse (conn *grpc.ClientConn) {
	client := prod.NewUserServiceClient(conn)
	var i int32
	req := prod.UserScoreRequest{}
	req.Users = make([]*prod.UserInfo, 0)

	for i = 1; i < 6; i++ {
		req.Users = append(req.Users, &prod.UserInfo{UserId: i})
	}
	stream, err := client.GetUserScoreByServerStream(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}
	for  {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}
}

func getClientStreamResponse (conn *grpc.ClientConn) {
	client := prod.NewUserServiceClient(conn)
	var i int32
	stream, err := client.GetUserScoreByClientStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for j := 1; j < 4; j++ {
		req := prod.UserScoreRequest{}
		req.Users = make([]*prod.UserInfo, 0)
		for i = 1; i < 6; i++ {
			req.Users = append(req.Users, &prod.UserInfo{UserId: i})
		}
		err := stream.Send(&req)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 1)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func getServerClientStreamResponse (conn *grpc.ClientConn) {
	client := prod.NewUserServiceClient(conn)
	var uid int32 = 1
	stream, err := client.GetUserScoreByTWStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for j := 1; j < 4; j++ {
		req := prod.UserScoreRequest{}
		req.Users = make([]*prod.UserInfo, 0)
		for i := 1; i < 6; i++ {
			req.Users = append(req.Users, &prod.UserInfo{UserId: uid})
			uid++
		}
		err := stream.Send(&req)
		if err != nil {
			fmt.Println(err.Error())
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(resp)
		time.Sleep(time.Second * 1)
	}
}
