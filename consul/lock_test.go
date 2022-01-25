package consul

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	addr := "192.168.208.214:8500"
	token := ""
	_, err := NewLock(addr, token)
	assert.Equal(t, err, nil, "err must be nil")

	lockKey := "session_timeout_lock_key"
	for {
		closeChan, err := GetLock(lockKey, "10s", time.Duration(2*time.Second))
		assert.Equal(t, err, nil, "err must be nil")
		fmt.Println("lock success")
		select {
		case <-closeChan:

		default:
			fmt.Println("default")
		}

		closeChan, err = GetLock(lockKey, "10s", time.Duration(2*time.Second))
		assert.Equal(t, err, nil, "err must be nil")

		<-closeChan

		fmt.Println("lock close")
	}

	/*
		stopCh := make(chan struct{}, 1)
		ret, err := lock.Lock(stopCh)
		assert.Equal(t, err, nil, "err must be nil")
		assert.NotEqual(t, ret, nil, "lock must be not nil")

		fmt.Println("lock success")
		<-ret
		fmt.Println("lost lock")

		lock, err = GetLock(lockKey, "10s", time.Duration(2 * time.Second))
		assert.Equal(t, err, nil, "err must be nil")
		ret, err = lock.Lock(stopCh)
		assert.Equal(t, err, nil, "err must be nil")
		assert.NotEqual(t, ret, nil, "lock must be not nil")
	*/
	//time.Sleep(100 * time.Second)
}
