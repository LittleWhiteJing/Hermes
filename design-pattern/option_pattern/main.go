package main

import "fmt"

func main() {
	cfg := NewOption(
		SetIpAddr("127.0.0.1"),
		SetPort(8080),
		SetUsername("YaSong"),
		SetPassword("ZnideD"),
	)
	fmt.Println(cfg)
}
