package main

import (
	"container/list"
	"errors"
	"sync"
)

type Node struct {
	Key		interface{}
	Val		interface{}
	Freq	int
}

type LFU struct {
	capacity	int
	counter		int
	minFreq		int
	keyMap		map[interface{}] *list.Element
	freqMap		map[int] *list.List
	mu			*sync.Mutex
}

func New(capacity int) *LFU {
	return &LFU{
		capacity: capacity,
		counter: 0,
		minFreq: 1,
		keyMap: make(map[interface{}] *list.Element),
		freqMap: make(map[int] *list.List),
		mu: new(sync.Mutex),
	}
}

func (l *LFU) Set(key interface{}, val interface{}) error {
	if l.capacity <= 0 {
		return	errors.New("no cache space available")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.keyMap[key]; ok {
		e.Value.(*Node).Val = val
		l.freqAdd(e)
	} else {
		l.removeEle()
		ele := &list.Element{
			Value: &Node{
				Key: key,
				Val: val,
				Freq: 1,
			},
		}
		l.keyMap[key] = ele
		if _, ok := l.freqMap[1]; ok {
			l.freqMap[1].PushBack(ele)
		} else {
			l.freqMap[1] = list.New()
			l.freqMap[1].PushBack(ele)
		}
		l.minFreq = 1
		l.counter++
	}
	return nil
}

func (l *LFU) Get(key interface{}) (val interface{}) {
	if e, ok := l.keyMap[key]; ok {
		l.freqAdd(e)
		return e.Value.(*Node).Val
	}
	return nil
}

func (l *LFU) freqAdd(e *list.Element) {
	freq := e.Value.(*Node).Freq
	l.freqMap[freq].Remove(e)

	if l.minFreq == freq && l.freqMap[freq].Len() == 0 {
		l.minFreq++
	}
	freq++
	e.Value.(*Node).Freq++
	if _, ok := l.freqMap[freq]; ok {
		l.freqMap[freq].PushBack(e)
	} else {
		l.freqMap[freq] = list.New()
		l.freqMap[freq].PushBack(e)
	}
}

func (l *LFU) removeEle() {
	if l.counter < l.capacity {
		return
	}
	minList := l.freqMap[l.minFreq]
	key := minList.Back().Value.(*Node).Key
	minList.Remove(minList.Back())
	delete(l.keyMap, key)
	l.counter--
}

