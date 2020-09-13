package main

import (
	"fmt"
	"log"
	"micro-client/util"

	hystrix "github.com/afex/hystrix-go/hystrix"
)

func main() {
	cfg := hystrix.CommandConfig{
		Timeout: 2000,
		MaxConcurrentRequests: 5,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold: 20,
		SleepWindow: 100,
	}
	hystrix.ConfigureCommand("calluserservice", cfg)
	err := hystrix.Do("calluserservice", func() error {
		res, err := util.ConsulConnect()
		fmt.Println(res)
		return err
	}, func(err error) error {
		fmt.Println("降级用户")
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

}
