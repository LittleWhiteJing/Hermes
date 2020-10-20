package main

import (
	"fmt"
	"time"
)

type MethodRequest int

const (
	ADD    MethodRequest =  iota + 1
	Lessen
)

type Service struct {
	queue chan MethodRequest
	value int
}

func NewService(buffer int) *Service {
	srv := &Service{
		queue: make(chan MethodRequest, buffer),
	}
	go srv.schedule()
	return srv
}

func (s *Service) schedule() {
	for {
		select {
			case req := <- s.queue:
				switch req {
					case ADD:
						fmt.Println("action: value++")
						s.value++
					case Lessen:
						fmt.Println("action: value--")
						s.value--
				}
			case <-time.After(time.Second * 5):
				fmt.Println("5s没有写入")
		}
	}
}

func (s *Service) Add() {
	fmt.Println("method request: add")
	s.queue <- ADD
}

func (s *Service) Lessen() {
	fmt.Println("method request: lessen")
	s.queue <- Lessen
}