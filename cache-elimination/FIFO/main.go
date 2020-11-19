package main

import (
	"container/list"
	"errors"
	"sync"
)

type Node struct {
	Key		interface{}
	Val		interface{}
}

type FIFO struct {
	capacity	int
	list		*list.List
	keyMap		map[interface{}] *list.Element
	mu			*sync.Mutex
}

func New(capacity int) *FIFO {
	return &FIFO{
		capacity: capacity,
		list: list.New(),
		keyMap: make(map[interface{}] *list.Element),
		mu: new(sync.Mutex),
	}
}

func (f *FIFO) Set(key interface{}, val interface{}) error {
	if f.list == nil || f.capacity == 0 {
		return errors.New("cache has not been init")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.keyMap[key]; ok {
		f.keyMap[key].Value.(*Node).Val = val
		return nil
	}
	ele := &list.Element{
		Value: &Node{
			Key: key,
			Val: val,
		},
	}
	f.keyMap[key] = ele
	f.list.PushBack(ele)
	if f.list.Len() > f.capacity {
		front := f.list.Front()
		key := front.Value.(*Node).Key
		f.list.Remove(front)
		delete(f.keyMap, key)
	}
	return nil
}

func (f *FIFO) Get(key interface{}) (val interface{}) {
	if e, ok := f.keyMap[key]; ok {
		return e.Value.(*Node).Val
	}
	return nil
}




