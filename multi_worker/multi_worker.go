package multi_worker

import (
	"sync"
)

func Run(tasks chan interface{}, workerNum int, workFunc func(interface{})) {
	var wg sync.WaitGroup

	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go func() {
			for task := range tasks {
				workFunc(task)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
