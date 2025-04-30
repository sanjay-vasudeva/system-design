package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup

	go func() {
		Receive()
		wg.Done()
	}()
	go func() {
		Send()
		wg.Done()
	}()
	consumers := 2
	for i := range consumers {
		wg.Add(1)
		go func() {
			Consume(i)
			wg.Done()
		}()
	}
	wg.Add(1)

	go func() {
		CreateTasks()
		wg.Done()
	}()
	wg.Wait()
}
