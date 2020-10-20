package main

import (
	"fmt"
	"time"
)

func main () {
	srv := NewService(1)
	for i := 0; i < 20; i++ {
		srv.Add()
		fmt.Printf("value: %d\n", srv.value)
	}
	fmt.Printf("value: %d\n", srv.value)
	time.Sleep(time.Second * 3)
	fmt.Printf("final: %d\n", srv.value)
}
