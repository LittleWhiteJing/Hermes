package distributed_lock

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/go-redis/redis"
	"github.com/hashicorp/go-multierror"
	"time"
)

type redisPool struct {
	instances []redis.Client
}

func NewRedisPool(instances ...redis.Client) *redisPool {
	return &redisPool{
		instances: instances,
	}
}

func (rp *redisPool) NewRedLock(name string, opts ...SetRedLock) *RedLock {
	rl := &RedLock{
		name: 		name,
		expire: 	8 * time.Second,
		tries: 		32,
		delayFunc:  genDelay,
		getValFunc: genValue,
		factor: 	0.01,
		quorum: 	len(rp.instances)/2 + 1,
		instances: 	rp.instances,
	}
	for _, opt := range opts {
		opt(rl)
	}
	return rl
}

type DelayFunc func(tries int) time.Duration

type GetValFunc func() (string, error)

func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func genDelay(tries int) time.Duration {
	return 500 * time.Millisecond
}

type RedLock struct {
	name 		string
	expire		time.Duration

	tries		int
	delayFunc	DelayFunc

	factor		float64

	quorum		int

	getValFunc 	GetValFunc
	value		string
	until		time.Time

	instances 	[]redis.Client
}

type SetRedLock func (rl *RedLock)

func SetRedLockExpire(expire time.Duration) SetRedLock {
	return func(rl *RedLock) {
		rl.expire = expire
	}
}

func SetRetriesTimes(tries int) SetRedLock {
	return func(rl *RedLock) {
		rl.tries = tries
	}
}

func SetDelayFunc(fn func(tries int) time.Duration) SetRedLock {
	return func(rl *RedLock) {
		rl.delayFunc = fn
	}
}

func SetGetValFunc(fn func() (string, error)) SetRedLock {
	return func(rl *RedLock) {
		rl.getValFunc = fn
	}
}

func (rl *RedLock) Lock() error {
	value, err := rl.getValFunc()
	if err != nil {
		return err
	}

	for i := 0; i < rl.tries; i++ {
		if i != 0 {
			time.Sleep(rl.delayFunc(i))
		}

		start := time.Now()
		n, err := rl.actOnClientsAsync(func(instance redis.Client) (bool, error) {
			return rl.acquire(instance, value)
		})
		if n == 0 && err != nil {
			return err
		}

		now := time.Now()
		until := now.Add(rl.expire - now.Sub(start) - time.Duration(int64(float64(rl.expire) * rl.factor)))
		if n >= rl.quorum && now.Before(until) {
			rl.value = value
			rl.until = until
			return nil
		}
		_, _ = rl.actOnClientsAsync(func(instance redis.Client) (bool, error) {
			return rl.release(instance, value)
		})
	}
	return errors.New("get redLock failed")
}

func (rl *RedLock) Unlock() (bool, error) {
	n, err := rl.actOnClientsAsync(func(instance redis.Client) (bool, error) {
		return rl.release(instance, rl.value)
	})
	if n < rl.quorum {
		return false, err
	}
	return true, nil
}

func (rl *RedLock) Extend() (bool, error) {
	n, err := rl.actOnClientsAsync(func(instance redis.Client) (bool, error) {
		return rl.touch(instance, rl.value, int(rl.expire/time.Millisecond))
	})
	if n < rl.quorum {
		return false, err
	}
	return true, nil
}

func (rl *RedLock) acquire(client redis.Client, value string) (bool, error) {
	reply, err := client.SetNX(rl.name, value, rl.expire).Result()
	if err != nil {
		return false, err
	}
	return reply, nil
}

var deleteScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`

func (rl *RedLock) release(client redis.Client, value string) (bool, error) {
	status, err := client.Eval(deleteScript, []string{rl.name, value}).Result()
	if err != nil {
		return false, err
	}
	return status != 0, nil
}

var touchScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`

func (rl *RedLock) touch(client redis.Client, value string, extend int) (bool, error) {
	status, err := client.Eval(touchScript, []string{rl.name, value, string(extend)}).Result()
	if err != nil {
		return false, err
	}
	return status != "ERR", nil
}

func (rl *RedLock) actOnClientsAsync(actFn func(redis.Client) (bool, error)) (int, error) {
	type result struct {
		Status 	bool
		Err		error
	}

	ch := make(chan result)
	for _, instance := range rl.instances {
		go func(instance redis.Client) {
			r := result{}
			r.Status, r.Err = actFn(instance)
			ch <- r
		}(instance)
	}
	n := 0
	var err error
	for range rl.instances {
		r := <-ch
		if r.Status {
			n++
		} else if r.Err != nil {
			err = multierror.Append(err, r.Err)
		}
	}
	return n, err

}







