package queue

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

var chQueue Queue

func init() {
	chQueue = NewChannelQueue()
}

func TestPut(t *testing.T) {
	topic := "test"

	err := chQueue.NewTopic(topic)
	assert.Equal(t, err, nil, "err must be nil")

	bol, err := chQueue.Put(topic, &Msg{MType: "hello", Data: "data"})
	assert.Equal(t, err, nil, "err must be nil")
	assert.Equal(t, bol, true, "bol must be true")

	bol, err = chQueue.Put("tests", &Msg{MType: "hello", Data: "data"})
	assert.NotEqual(t, err, nil, "err must be nil")
	assert.Equal(t, bol, false, "bol must be false")

}

func TestGet(t *testing.T) {
	topic := "test"

	msg := chQueue.Get(topic)
	assert.NotEqual(t, msg, nil, "msg must not be nil")
}
