package consul

import (
	_"fmt"
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	addr := "192.168.208.214:8500"
	token := ""
	_, err := NewLock(addr, token)
	assert.Equal(t, err, nil, "err must be nil")

	lockKey := "hello_lock"
	lock, err := GetLock(lockKey, "10s", time.Duration(2 * time.Second))
	assert.Equal(t, err, nil, "err must be nil")

	stopCh := make(chan struct{}, 1)
	ret, err := lock.Lock(stopCh)
	assert.Equal(t, err, nil, "err must be nil")
	assert.NotEqual(t, ret, nil, "lock must be not nil")
}
