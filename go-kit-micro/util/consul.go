package util

import (
	consulapi "github.com/hashicorp/consul/api"
	"log"
)

var (
	consoleClient *consulapi.Client
)

func init() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	consoleClient = client
}

func RegisterService () {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = "userservice"
	reg.Name = "userservice"
	reg.Address = "192.168.1.104"
	reg.Port = 8080
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = "http://192.168.1.104:8080/health"
	reg.Check = &check

	err := consoleClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func UnRegisterService () {
	consoleClient.Agent().ServiceDeregister("userservice")
}