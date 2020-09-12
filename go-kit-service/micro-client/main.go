package main

import (
	"fmt"
	"io"
	"net/url"
	"log"
	"context"
	"os"
	"micro-client/Services"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	kitlog "github.com/go-kit/kit/log"
)

func main() {
	//directConnect()
	consulConnect()
}

func directConnect() {
	tg, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	client := httptransport.NewClient("GET", tg, Services.GetUserInfoRequest, Services.GetUserInfoResponse)
	getUserInfo := client.Endpoint()
	res, err := getUserInfo(context.Background(), Services.UserRequest{Uid: 101})
	if err != nil {
		log.Fatal(err)
	}
	userInfo := res.(Services.UserResponse)
	fmt.Println(userInfo.Result)
}

func consulConnect() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = "http://localhost:8500"
	api_client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	client := consul.NewClient(api_client)
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stdout)

	tags := []string{"primary"}
	instance := consul.NewInstancer(client, logger, "userservice", tags, true)

	f := func(instance string) (endpoint.Endpoint, io.Closer, error) {
		tg, err := url.Parse("http://"+instance)
		if err != nil {
			log.Fatal(err)
		}
		return httptransport.NewClient("GET", tg, Services.GetUserInfoRequest, Services.GetUserInfoResponse).Endpoint(), nil, nil
	}

	endpointer := sd.NewEndpointer(instance, f, logger)
	endpoints, err := endpointer.Endpoints()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("发现", len(endpoints), "个服务")

	getUserInfo := endpoints[0]
	res, err := getUserInfo(context.Background(), Services.UserRequest{Uid: 102})
	if err != nil {
		log.Fatal(err)
	}
	userInfo := res.(Services.UserResponse)
	fmt.Println(userInfo.Result)

}
