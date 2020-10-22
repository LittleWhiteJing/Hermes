package distributed_lock

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func TestRedLock (t *testing.T) {
	redisClients := []redis.Client{*redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		}),*redis.NewClient(&redis.Options{
			Addr: "localhost:6380",
			Password: "",
			DB: 0,
		}),*redis.NewClient(&redis.Options{
			Addr: "localhost:6381",
			Password: "",
			DB: 0,
		}),
	}
	for _, client := range redisClients {
		_, err := client.Ping().Result()
		if err != nil {
			t.Fatal("redis-service unavailable")
		}
	}

	redisPool := NewRedisPool(redisClients...)
	redLock := redisPool.NewRedLock("redLock-test")
	if err := redLock.Lock(); err != nil {
		t.Fatal(err)
	}
	fmt.Println("get lock success")
	if _, err := redLock.Unlock(); err != nil {
		t.Fatal(err)
	}
}
