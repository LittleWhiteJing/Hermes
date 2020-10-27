package distributed_lock

/**
 * 基于 redis 的单机分布式锁实现
 * 1.基于 SETNX 命令 保证 kv 和 expire 设置的原子性
 * 2.锁超时机制，避免持有锁的线程死掉导致锁不可用
 * 3.基于客户端提供的 uuid 进行锁操作身份验证
 * 4.支持客户端主动续约机制，避免锁超时业务未处理完
 */

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"time"
)

type DisLockRedis struct {
	conn 	*redis.Client
	ttl		time.Duration
	key  	string
	val 	string
}

func NewDisLockRedis (conn *redis.Client, key string) *DisLockRedis {
	dlr := &DisLockRedis{
		conn: conn,
		ttl:  10,
		key:  key,
	}
	return dlr
}

func (dlr *DisLockRedis) TryLock () (bool, error) {
	value := uuid.New().String()
	res, err := dlr.conn.SetNX(dlr.key, value, dlr.ttl * time.Second).Result()
	if res == true {
		dlr.val    = value
	}
	return res, err
}

func (dlr *DisLockRedis) UnLock () error {
	if dlr.val == dlr.conn.Get(dlr.key).String() {
		return dlr.conn.Del(dlr.key).Err()
	}
	return errors.New("lock is not own")
}

func (dlr *DisLockRedis) Extend () (bool, error) {
	if dlr.val == dlr.conn.Get(dlr.key).String() {
		return dlr.conn.Expire(dlr.key, dlr.ttl * time.Second).Result()
	}
	return false, errors.New("lock is not own")
}

