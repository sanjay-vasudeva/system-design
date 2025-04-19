package blocking_queue

import (
	"fmt"
	"testing"
	"time"
)

// Testing
func TestBlockingQueue(t *testing.T) {
	q := NewBlockingQueue[int](1)
	c := make(chan bool)

	go func() {
		q.Put(1)
		time.Sleep(100 * time.Millisecond)
		q.Put(2)
		time.Sleep(100 * time.Millisecond)
	}()

	go func(c chan bool) {
		item := q.Take()
		fmt.Printf("Got %v\n", item)

		item = q.Take()
		fmt.Printf("Got %v\n", item)
		c <- true
	}(c)

	<-c
}
