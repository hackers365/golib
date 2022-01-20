package consul

import (
	"time"
	"github.com/hashicorp/consul/api"
)

type ConsulLock struct {
	client 		*api.Client
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
		client: client,
	}
	return c, nil
}

func C() *ConsulLock {
	return c
}

func GetLock(lockKey string, ttl string, lockDelay time.Duration) (*api.Lock, error) {
	opts := &api.LockOptions{
		Key: lockKey,
		SessionTTL: ttl,
		LockDelay: lockDelay,
	}
	return C().client.LockOpts(opts)
}
