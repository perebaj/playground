package semaphore

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type semaphore struct {
	sem chan struct{}
}

func newSemaphore(maxConcurrency int) *semaphore {
	return &semaphore{make(chan struct{}, maxConcurrency)}
}

func (s *semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Release releases the semaphore. Must be defered after a successful process
func (s *semaphore) Release() {
	<-s.sem
}

func main() {
	var nSlice = flag.Int("n", 5, "number of slice elements")
	var maxConcurrency = flag.Int("c", 5, "max concurrency")

	flag.Parse()

	var list []int
	for i := 1; i <= *nSlice; i++ {
		list = append(list, i)
	}
	fmt.Println("Number of slice elements: ", len(list))
	sem := newSemaphore(*maxConcurrency)

	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)
		sem.Acquire()
		go func(v int) {
			defer wg.Done()
			defer sem.Release()
			// time sleep using event processing time

			time.Sleep(2 * time.Second)

		}(v)
	}
	wg.Wait()

}

// How to use:
// Time go run semaphore.go -n 10 -c 5
