package test

import (
	"encoding/json"
	"fmt"
	"github.com/TyrellJing/Hermes/router"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func helloWorldGETHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World By GET Method!\n")
}

func helloWorldPOSTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World By POST Method!\n")
}

func helloWorldPUTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World By PUT Method!\n")
}

func helloWorldDELETEHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World By DELETE Method!\n")
}

func TestRouterRegister(t *testing.T) {
	r := radix_tree.NewRouter()

	r.GET("/method/hello/get",   		helloWorldGETHandler)
	r.POST("/method/hello/post", 		helloWorldPOSTHandler)
	r.PUT("/method/hello/put",   		helloWorldPUTHandler)
	r.DELETE("/method/hello/delete",  helloWorldDELETEHandler)

	http.ListenAndServe("127.0.0.1:8000", r)
}

func TestRouterGetHandler(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8000/method/hello/get")
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	c, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(c) + "\n")
}

func TestRouterPostHandler(t *testing.T) {
	body := map[string]string{"server1": "127.0.0.1", "server2": "127.0.0.2"}
	b, _ := json.Marshal(body)
	resp, err := http.Post("http://127.0.0.1:8000/method/hello/post", "application/json", strings.NewReader("heel="+string(b)))
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	c, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(c) + "\n")
}

func TestRouterPutHandler(t *testing.T) {
	body := map[string]string{"server1": "127.0.0.1", "server2": "127.0.0.2"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/method/hello/put", strings.NewReader("heel="+string(b)))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	c, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(c) + "\n")
}

func TestRouterDeleteHandler(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8000/method/hello/delete", nil)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	c, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(c) + "\n")
}