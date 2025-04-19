package main

import (
	"sync"
	mq "w1/rabbit-mq"
)

func main() {
	//basic_single()
	//basic_multi(10)
	// fair_multi(10)

	//cp := NewConnectionPool(5)
	// benchmarkNonPool(150)
	// benchmarkPool(10000)

	//sse
	// http.HandleFunc("/events", sse.EventHandler)
	// http.ListenAndServe(":8080", nil)

	//rabbitmq
	var wg sync.WaitGroup

	// go func() {
	// 	mq.Receive()
	// 	wg.Done()
	// }()
	// go func() {
	// 	mq.Send()
	// 	wg.Done()
	// }()
	consumers := 2
	for i := range consumers {
		wg.Add(1)
		go func() {
			mq.Consume(i)
			wg.Done()
		}()
	}
	wg.Add(1)

	go func() {
		mq.CreateTasks()
		wg.Done()
	}()
	wg.Wait()
}
