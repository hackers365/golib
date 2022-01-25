package consul

import (
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type ConsulLock struct {
	sync.RWMutex

	client  *api.Client
	lockMap map[string]*lockOpts
}

type lockOpts struct {
	lock      *api.Lock
	closeChan <-chan struct{}
}

var c *ConsulLock

func NewLock(consulAddr string, token string) (*ConsulLock, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	c = &ConsulLock{
		client:  client,
		lockMap: map[string]*lockOpts{},
	}
	return c, nil
}

func C() *ConsulLock {
	return c
}

func GetLock(lockKey string, ttl string, lockDelay time.Duration) (<-chan struct{}, error) {
	return C().GetLock(lockKey, ttl, lockDelay)
}

func (c *ConsulLock) GetLock(lockKey string, ttl string, lockDelay time.Duration) (<-chan struct{}, error) {
	c.Lock()
	defer c.Unlock()

	if lockInstance, ok := c.lockMap[lockKey]; ok {
		//检查closeChan是否已经失效
		select {
		case <-lockInstance.closeChan:

		default:
			return lockInstance.closeChan, nil
		}
	}

	opts := &api.LockOptions{
		Key:         lockKey,
		SessionTTL:  ttl,
		LockDelay:   lockDelay,
		LockTryOnce: true,
	}
	lockInstance, err := c.client.LockOpts(opts)
	if err != nil {
		return nil, err
	}

	stopCh := make(chan struct{}, 1)
	closeChan, err := lockInstance.Lock(stopCh)
	if err != nil {
		return nil, err
	}

	c.lockMap[lockKey] = &lockOpts{
		lock:      lockInstance,
		closeChan: closeChan,
	}

	return closeChan, nil
}
