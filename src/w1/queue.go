package main

import (
	"sync"
)

type BlockingQueue struct {
	mu       *sync.Mutex
	capacity int
	data     []any
	cond     *sync.Cond
}

func (queue *BlockingQueue) Put(item any) {
	queue.cond.L.Lock()
	defer queue.cond.L.Unlock()

	for queue.IsFull() {
		queue.cond.Wait()
	}
	queue.data = append(queue.data, item)
	queue.cond.Signal()
}

func (queue *BlockingQueue) IsFull() bool {
	return len(queue.data) >= queue.capacity
}

func (queue *BlockingQueue) Take() any {
	queue.cond.L.Lock()
	defer queue.cond.L.Unlock()

	for queue.IsEmpty() {
		queue.cond.Wait()
	}
	item := queue.data[0]
	queue.data = queue.data[1:]
	queue.cond.Signal()
	return item
}

func (queue *BlockingQueue) IsEmpty() bool {
	return len(queue.data) == 0
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	queue := &BlockingQueue{
		mu:       &sync.Mutex{},
		capacity: capacity,
		data:     make([]any, 0, capacity),
	}
	queue.cond = &sync.Cond{L: queue.mu}
	return queue
}
