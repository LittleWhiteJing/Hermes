package util

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	consoleClient *consulapi.Client
	ServiceID	string
	ServiceName string
	ServicePort int
)

func init() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	consoleClient = client
	ServiceID = "userservice" + uuid.New().String()
}

func SetServicePortAndName(name string, port int) {
	ServiceName = name
	ServicePort = port
}

func RegisterService () {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = ServiceID
	reg.Name = ServiceName
	reg.Port = ServicePort
	reg.Address = "192.168.1.104"
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = fmt.Sprintf("http://%s:%d/health", reg.Address, ServicePort)
	reg.Check = &check

	err := consoleClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func UnRegisterService () {
	consoleClient.Agent().ServiceDeregister(ServiceID)
}