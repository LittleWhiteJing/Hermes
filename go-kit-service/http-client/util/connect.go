package util

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"io"
	"log"
	"micro-client/Services"
	"net/url"
	"os"
	"time"
)

func DirectConnect() (string, error) {
	tg, err := url.Parse("http://localhost:8080")
	if err != nil {
		return "", err
	}
	client := httptransport.NewClient("GET", tg, Services.GetUserInfoRequest, Services.GetUserInfoResponse)
	getUserInfo := client.Endpoint()
	res, err := getUserInfo(context.Background(), Services.UserRequest{Uid: 101})
	if err != nil {
		return "", err
	}
	userInfo := res.(Services.UserResponse)
	return userInfo.Result, nil
}

func ConsulConnect() (string, error) {
	cfg := consulapi.DefaultConfig()
	cfg.Address = "http://localhost:8500"
	api_client, err := consulapi.NewClient(cfg)
	if err != nil {
		return "", err
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

	//endpointlb := lb.NewRoundRobin(endpointer)
	endpointlb := lb.NewRandom(endpointer, time.Now().UnixNano())

	getUserInfo, err := endpointlb.Endpoint()
	if err != nil {
		return "", err
	}
	res, err := getUserInfo(context.Background(), Services.UserRequest{Uid: 102})
	if err != nil {
		return "", err
	}
	userInfo := res.(Services.UserResponse)
	return userInfo.Result, nil
}
