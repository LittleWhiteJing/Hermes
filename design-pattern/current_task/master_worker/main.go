package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const RoutineCount = 5

var client *http.Client

func main() {
	numberTasks := []string{"13456755448", " 13419385751", "13419317885", " 13434343439", "13438522395"}
	wg := &sync.WaitGroup{}
	beg := time.Now()
	tasks := make(chan string)
	results := make(chan string)
	//receiver
	for i := 0; i < len(numberTasks); i++ {
		result := <-results
		fmt.Println("result:", result)
	}
	//worker
	for i := 0; i < RoutineCount; i++ {
		wg.Add(1)
		go worker(wg, tasks, results)
	}
	//caller
	for _, task := range numberTasks {
		tasks <- task
	}
	wg.Wait()
	fmt.Printf("time consumed: %fs", time.Now().Sub(beg).Seconds())
}

func worker(wg *sync.WaitGroup, tasks chan string, results chan string) {
	task := <-tasks
	respBody, err := NumberQueryRequest(task)
	if err != nil{
		fmt.Printf("error occurred in NumberQueryRequest: %s\n", task)
		results <- err.Error()
	}else{
		results <- string(respBody)
	}
	wg.Done()
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
