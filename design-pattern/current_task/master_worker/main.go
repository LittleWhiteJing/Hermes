package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	TaskTimeout	 = 5
	MaxGoroutine = 10
)

var client *http.Client

func main() {
	numberTasks := []string{"13456755448", " 13419385751", "13419317885", " 13434343439", "13438522395"}
	beg := time.Now()
	results := make([]chan string, len(numberTasks))
	limiter := make(chan bool, MaxGoroutine)
	//worker
	for i, numberTask := range numberTasks {
		results[i] = make(chan string, 1)
		limiter <- true
		go worker(numberTask, TaskTimeout, limiter, results[i])
	}
	//receiver
	for _, result := range results {
		fmt.Println("result:", <-result)
	}
	fmt.Printf("time consumed: %fs", time.Now().Sub(beg).Seconds())
}

func worker(task string, timeout int, limiter chan bool, results chan string) {
	receiver := make(chan string, 1)
	go exec(task, limiter, receiver)
	select {
		case rs := <-receiver:
			results <- rs
		case <-time.After(time.Duration(timeout) * time.Second):
			results <- "task run out of time"
	}
}

func exec(task string, limiter chan bool, results chan string) {
	respBody, err := NumberQueryRequest(task)
	if err != nil{
		fmt.Printf("error occurred in NumberQueryRequest: %s\n", task)
		results <- err.Error()
	}else{
		results <- string(respBody)
	}
	<-limiter
}

func NumberQueryRequest(keyword string)(body []byte, err error){
	url := fmt.Sprintf("https://api.binstd.com/shouji/query?appkey=df2720f76a0991fa&shouji=%s", keyword)
	resp, err := client.Get(url)
	if err != nil{
		return nil, err
	}
	if resp.StatusCode != http.StatusOK{
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("response status code is not OK, response code is %d, body:%s", resp.StatusCode, string(data))
	}
	if resp != nil && resp.Body != nil{
		defer resp.Body.Close()
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	return body, nil
}
