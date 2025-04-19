package practise

import (
	"sync"
	"testing"
)

type Buffer struct {
	Data [10240]int
}

var bufferPool = sync.Pool{
	New: func() any {
		return &Buffer{}
	},
}

func BenchmarkNonPool(b *testing.B) {
	for b.Loop() {
		obj := &Buffer{}
		for i := range len(obj.Data) {
			obj.Data[i] = 1
		}
	}
}

func BenchmarkPool(b *testing.B) {
	for b.Loop() {
		obj := bufferPool.Get().(*Buffer)
		for i := range len(obj.Data) {
			obj.Data[i] = 1
		}
		bufferPool.Put(obj)
	}
}
