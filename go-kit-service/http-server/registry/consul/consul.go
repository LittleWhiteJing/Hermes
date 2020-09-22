package consul

import (
	"micro-server/registry"
	consulapi "github.com/hashicorp/consul/api"
)

type options struct {
	ConsulClient *consulapi.Client
}

type Option func(o *options)

func NewConsulRegister(opts... Option) (options, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o, nil
}

func SetConsulClient(address string) Option {
	return func(o *options) {
		cfg := consulapi.DefaultConfig()
		cfg.Address = address
		client, _ := consulapi.NewClient(cfg)
		o.ConsulClient = client
	}
}

type ServiceCheck struct {
	CheckAddr 	string
	CheckIntval string
}

func (o *options) RegisterService (service registry.Service, serviceCheck ServiceCheck) error {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID      = service.ServiceID
	reg.Name    = service.ServiceName
	reg.Address = service.ServiceAddr
	reg.Port    = service.ServicePort
	reg.Tags    = service.ServiceTags

	var check consulapi.AgentServiceCheck
	check.Interval = serviceCheck.CheckIntval
	check.HTTP = serviceCheck.CheckAddr
	reg.Check = &check

	return o.ConsulClient.Agent().ServiceRegister(&reg)
}

func (o *options) UnRegisterService (serviceID string) error {
	return o.ConsulClient.Agent().ServiceDeregister(serviceID)
}