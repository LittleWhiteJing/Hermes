package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var client *http.Client

func main () {
	numberTasks := []string{"13456755448", " 13419385751", "13419317885", " 13434343439", "13438522395"}
	wg := &sync.WaitGroup{}
	beg := time.Now()
	for _, keyword := range numberTasks {
		wg.Add(1)
		go func(keyword string, wg *sync.WaitGroup) {
			body, err := NumberQueryRequest(keyword)
			if err != nil {
				fmt.Printf("error occurred in query keyword: %s, error: %s\n", keyword, err.Error())
			} else {
				fmt.Printf("search %s success, data size is %d\n, body is %s\n", keyword, len(body), string(body))
			}
			wg.Done()
		}(keyword, wg)
	}
	wg.Wait()
	fmt.Printf("time consumed: %fs", time.Now().Sub(beg).Seconds())
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
