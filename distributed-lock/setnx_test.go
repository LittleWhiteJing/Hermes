package distributed_lock

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
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

	dLockR := NewDisLockRedis(rdb, "distributed-lock")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(seq int) {
			uuid := uuid.New()
			uuidString := uuid.String()
			for {
				res, _ := dLockR.TryLock(uuidString)
				if res == true {
					fmt.Printf("goroutine:%d uuid:%s get lock success\n", seq, uuidString)
					time.Sleep(3 * time.Second)
					dLockR.UnLock(uuidString)
				} else {
					fmt.Printf("goroutine:%d uuid:%s get lock failed\n", seq, uuidString)
				}
				time.Sleep(1 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
