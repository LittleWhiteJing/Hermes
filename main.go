package main

import (
	"fmt"
	"net/http"
	"github.com/TyrellJing/Hermes/routers"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World !\n")
}

func sayGoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good Bye !\n")
}

func main() {
	r := routers.NewRouter()
	r.GET("/helloworld", helloWorldHandler)
	r.GET("/saygoodbye", sayGoodbyeHandler)
	http.ListenAndServe(":8000", r)
}