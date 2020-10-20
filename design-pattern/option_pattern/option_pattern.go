package main

type config struct {
	ip			string
	port		int
	username	string
	password    string
}

type Option func(* config)

func NewOption(opts ...Option) *config {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func SetIpAddr(ipAddr string) Option {
	return func(c *config) {
		c.ip = ipAddr
	}
}

func SetPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

func SetUsername(username string) Option {
	return func(c *config) {
		c.username = username
	}
}

func SetPassword(password string) Option {
	return func(c *config) {
		c.password = password
	}
}

