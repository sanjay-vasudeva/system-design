package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var count atomic.Int32

func IncrementCounter() {
	oldValue := count.Load()
	newValue := oldValue + 1
	if !count.CompareAndSwap(oldValue, newValue) {
		fmt.Println("Increment failed")
	}
}

func main() {
	wg := sync.WaitGroup{}

	for range 1000 {
		wg.Add(1)
		go func() {
			IncrementCounter()
			wg.Done()
		}()
	}
	wg.Wait()

}
