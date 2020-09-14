package registry

type Service struct {
	ServiceID 	string
	ServiceName string
	ServicePort int
	ServiceAddr	string
	ServiceTags []string
}

type IServiceRegistry interface {
	RegisterService(service Service) error
	UnregisterService(serviceID string) error
}
