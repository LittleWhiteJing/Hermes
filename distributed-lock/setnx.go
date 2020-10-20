package distributed_lock

/**
 * 基于redis的单机分布式锁实现
 * 1.基于 SETNX 命令 保证 kv 和 expire 设置的原子性
 * 2.锁超时机制，避免持有锁的线程死掉导致锁不可用
 * 3.基于客户端提供的uuid进行释放锁身份验证
 * 4.支持客户端主动续约机制，避免锁超时业务未处理完
 */

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

type DisLockRedis struct {
	conn 	*redis.Client
	ttl		time.Duration
	key  	string
	count   int
	cancel 	context.CancelFunc
}

func NewDisLockRedis (conn *redis.Client, key string) *DisLockRedis {
	dlr := &DisLockRedis{
		conn: conn,
		ttl:  30,
		key:  key,
	}
	return dlr
}

func (dlr *DisLockRedis) TryLock (uuid string) (bool, error) {
	if dlr.conn.Get(dlr.key).String() == uuid {
		dlr.count++
		return true, nil
	}
	res, err := dlr.conn.SetNX(dlr.key, uuid, dlr.ttl).Result()
	ctx, cancelFunc := context.WithCancel(context.Background())
	dlr.cancel = cancelFunc
	dlr.count  = 1
	//自动续约
	dlr.renew(ctx)
	return res, err
}

func (dlr *DisLockRedis) UnLock (uuid string) error {
	if dlr.conn.Get(dlr.key).String() == uuid {
		dlr.count--
		if dlr.count == 0 {
			//关闭续约协程
			dlr.cancel()
			return dlr.conn.Del(dlr.key).Err()
		}
	}
	return errors.New("lock is not own")
}

func (dlr *DisLockRedis) renew (ctx context.Context) {
	go func() {
		for {
			select {
				case <-ctx.Done():
					return
				default:
					dlr.conn.Expire(dlr.key, dlr.ttl)
			}
			time.Sleep((dlr.ttl / 3) * time.Second)
		}
	}()
}








