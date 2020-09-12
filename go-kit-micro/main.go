package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go-kit-micro/Services"
	"go-kit-micro/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	userService := Services.UserService{}
	endPoint := Services.GenUserEndpoint(userService)
	serverHandler := httptransport.NewServer(endPoint, Services.DecodeUserRequest, Services.EncodeUserResponse)

	r := mux.NewRouter()
	r.Methods("GET", "DELETE").Path(`/user/{userid:\d+}`).Handler(serverHandler)

	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	errC := make(chan error)

	go func() {
		util.RegisterService()
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			errC <- err
		}
	}()

	go func() {
		signC := make(chan os.Signal)
		signal.Notify(signC, syscall.SIGINT, syscall.SIGTERM)
		errC <- fmt.Errorf("%s", <-signC)
	}()

	getErr := <-errC
	util.UnRegisterService()
	fmt.Println(getErr)
}