package distributed_lock

import (
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"testing"
	"time"
)

func TestSetnx(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		t.Fatal("redis-service unavailable")
	}

	dLockR := NewDisLockRedis(rdb, "setNxLock-test")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(seq int) {
			for {
				res, _ := dLockR.TryLock()
				if res == true {
					fmt.Printf("goroutine:%d get lock success\n", seq)
					time.Sleep(3 * time.Second)
					dLockR.UnLock()
				} else {
					fmt.Printf("goroutine:%d get lock failed\n", seq)
				}
				time.Sleep(1 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
