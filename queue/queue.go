package queue

import (
)

type Msg struct{
	MType string
	Data interface{}
}

type Queue interface {
	NewTopic(string) error
	Put(string, *Msg) (bool, error)
	Get(string) (*Msg)
	Len(string) int
}
