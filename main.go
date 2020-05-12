package main

import (
	"fmt"
	"net/http"
	"github.com/TyrellJing/Hermes/router"
	"github.com/TyrellJing/Hermes/server"
	"time"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World !\n")
}

func sayGoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good Bye !\n")
}

func main() {
	r := router.NewRouter()
	r.GET("/helloworld", helloWorldHandler)
	r.GET("/saygoodbye", sayGoodbyeHandler)
	s := server.NewServer(*r, "127.0.0.1:8000", 300*time.Second, 300*time.Second)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}