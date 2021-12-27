package queue

import (
	"fmt"
	"sync"
)

type ChannelQueue struct {
	sync.RWMutex
	Name2Queue map[string]chan *Msg
}

func NewChannelQueue() Queue {
	chanQueue := &ChannelQueue{
		Name2Queue: map[string]chan *Msg{},
	}
	return chanQueue
}

func(c *ChannelQueue) NewTopic(topic string) error {
	c.Lock()
	defer c.Unlock()

	c.Name2Queue[topic] = make(chan *Msg, 10000)
	return nil
}

func(c *ChannelQueue) Put(topic string, msg *Msg) (bool, error) {
	c.RLock()
	chQueue, ok := c.Name2Queue[topic]
	c.RUnlock()
	if !ok {
		return false, fmt.Errorf("topic not exists")
	}

	select{
	case chQueue <- msg:
		return true, nil
	default:
		return false, nil
	}
}

func(c *ChannelQueue) Get(topic string) (*Msg) {
	c.RLock()
	chQueue, ok := c.Name2Queue[topic]
	c.RUnlock()
	if ok {
		data := <- chQueue
		return data
	}
	return nil
}

func(c *ChannelQueue) Len(topic string) int {
	c.RLock()
	defer c.RUnlock()

	if q, ok := c.Name2Queue[topic]; ok {
		return len(q)
	}
	return 0
}