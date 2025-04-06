package main

import (
	"sync"
)

type BlockingQueue[T any] struct {
	mu       *sync.Mutex
	capacity int
	data     []T
	cond     *sync.Cond
}

func (queue *BlockingQueue[T]) Put(item T) {
	queue.cond.L.Lock()
	defer queue.cond.L.Unlock()

	for queue.IsFull() {
		queue.cond.Wait()
	}
	queue.data = append(queue.data, item)
	queue.cond.Signal()
}

func (queue *BlockingQueue[T]) IsFull() bool {
	return len(queue.data) >= queue.capacity
}

func (queue *BlockingQueue[T]) Take() T {
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

func (queue *BlockingQueue[T]) IsEmpty() bool {
	return len(queue.data) == 0
}

func NewBlockingQueue[T any](capacity int) *BlockingQueue[T] {
	queue := &BlockingQueue[T]{
		mu:       &sync.Mutex{},
		capacity: capacity,
		data:     make([]T, 0, capacity),
	}
	queue.cond = &sync.Cond{L: queue.mu}
	return queue
}
