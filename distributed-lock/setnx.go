package distributed_lock

/**
 * 基于 redis 的单机分布式锁实现
 * 1.基于 SETNX 命令 保证 kv 和 expire 设置的原子性
 * 2.锁超时机制，避免持有锁的线程死掉导致锁不可用
 * 3.基于客户端提供的 uuid 进行锁操作身份验证
 * 4.支持客户端主动续约机制，避免锁超时业务未处理完
 * 5.支持锁重入，一个多次加锁内部计数，计数为0时释放
 */

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

type DisLockRedis struct {
	conn 	*redis.Client
	ttl		time.Duration
	key  	string
	count   int
	sign 	chan string
	exit 	chan bool
}

func NewDisLockRedis (conn *redis.Client, key string) *DisLockRedis {
	dlr := &DisLockRedis{
		conn: conn,
		ttl:  10,
		key:  key,
		sign: make(chan string, 1),
		exit: make(chan bool, 1),
	}
	return dlr
}

func (dlr *DisLockRedis) TryLock (uuid string) (bool, chan string) {
	value, _ := dlr.conn.Get(dlr.key).Result()
	if value == uuid {
		dlr.count++
		return true, dlr.sign
	}
	res, _ := dlr.conn.SetNX(dlr.key, uuid, dlr.ttl * time.Second).Result()
	if res == true {
		dlr.count  = 1
		//开启续约协程
		dlr.renew(dlr.sign, dlr.exit)
	}
	return res, dlr.sign
}

func (dlr *DisLockRedis) UnLock (uuid string) error {
	if dlr.conn.Get(dlr.key).String() == uuid {
		dlr.count--
		if dlr.count == 0 {
			//关闭续约协程
			dlr.exit <- true
			return dlr.conn.Del(dlr.key).Err()
		}
	}
	return errors.New("lock is not own")
}

func (dlr *DisLockRedis) renew (sign chan string, exit chan bool) {
	go func() {
		for {
			select {
				case uuid := <-sign:
					if uuid == dlr.conn.Get(dlr.key).String() {
						dlr.conn.Expire(dlr.key, dlr.ttl * time.Second)
					}
				case <-exit:
					return
			}
		}
	}()
}

