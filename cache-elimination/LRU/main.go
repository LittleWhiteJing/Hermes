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

type CallBack func(key interface{}, value interface{})

type LRU struct {
	capacity	int
	list		*list.List
	cache		map[interface{}] *list.Element
	callback	CallBack
	mu			*sync.Mutex
}

func New(capacity int, callback CallBack) *LRU {
	return &LRU{
		capacity: 	capacity,
		list:	 	list.New(),
		cache:		make(map[interface{}] *list.Element),
		callback: 	callback,
		mu:			new(sync.Mutex),
	}
}

func (l *LRU) Set(key interface{}, val interface{}) error {
	if l.list == nil {
		return	errors.New("internal list not init")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.cache[key]; ok {
		e.Value.(*Node).Val = val
		l.list.MoveToFront(e)
		return nil
	}
	ele := l.list.PushBack(Node{
		Key: key,
		Val: val,
	})
	l.cache[key] = ele
	if l.capacity != 0 && l.list.Len() > l.capacity {
		if e := l.list.Back(); e != nil {
			l.list.Remove(e)
			node := e.Value.(*Node)
			delete(l.cache, node.Key)
			if l.callback != nil {
				l.callback(node.Key, node.Val)
			}
		}
	}
	return nil
}

func (l *LRU) Get(key interface{}) (val interface{}, err error) {
	if l.list == nil {
		return nil, errors.New("internal list not init")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.cache[key]; ok {
		l.list.MoveToFront(e)
		return e.Value, nil
	}
	return nil, errors.New("key does not exist")
}





