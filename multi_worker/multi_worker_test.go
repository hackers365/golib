package multi_worker

import (
	"fmt"
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	tasks := make(chan interface{}, 10)
	for i := 0; i < 10; i++ {
		tasks <- i
	}
	close(tasks)

	l := []int{}
	var lock sync.RWMutex
	Run(tasks, 10, func(task interface{}) {
		i := task.(int)
		fmt.Println(i)
		lock.Lock()
		l = append(l, i)
		lock.Unlock()
	})

	if len(l) != 10 {
		t.Errorf("%v", l)
		return
	}
	fmt.Println(l)
}
