package util

import (
	"log"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	consoleClient *consulapi.Client
	serviceID	string
	serviceName string
	servicePort int
)

func init() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	consoleClient = client
	serviceID = "userservice" + uuid.New().String()
}

func SetServicePortAndName(name string, port int) {
	serviceName = name
	servicePort = port
}

func RegisterService () {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = serviceID
	reg.Name = serviceName
	reg.Port = servicePort
	reg.Address = "192.168.1.104"
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
	consoleClient.Agent().ServiceDeregister(serviceID)
}