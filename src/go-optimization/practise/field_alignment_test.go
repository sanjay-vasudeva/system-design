package practise

import (
	"sync"
	"testing"
)

type PoorAligned struct {
	i  int32
	i8 int64
	b  bool
}

type WellAligned struct {
	i8 int64
	i  int32
	b  bool
}

func BenchmarkPoorAlignment(b *testing.B) {
	for b.Loop() {
		data := make([]PoorAligned, 100_000)
		for i := range len(data) {
			data[i].i8 = int64(i)
		}
	}
}

func BenchmarkGoodAlignment(b *testing.B) {
	for b.Loop() {
		data := make([]WellAligned, 100_000)
		for i := range len(data) {
			data[i].i8 = int64(i)
		}
	}
}

type SharedCounterBad struct {
	a int64
	b int64
}

type SharedCounterGood struct {
	a int64
	_ [56]byte // padding to avoid sharing cache line
	b int64
}

func BenchmarkFalseSharing(b *testing.B) {
	var c SharedCounterBad
	var wg sync.WaitGroup
	for b.Loop() {
		wg.Add(2)
		go func() {
			for range 1_000_000 {
				c.a++
			}
			wg.Done()
		}()
		go func() {
			for range 1_000_000 {
				c.b++
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

func BenchmarkGoodSharing(b *testing.B) {
	var c SharedCounterGood
	var wg sync.WaitGroup
	for b.Loop() {
		wg.Add(2)
		go func() {
			for range 1_000_000 {
				c.a++
			}
			wg.Done()
		}()
		go func() {
			for range 1_000_000 {
				c.b++
			}
			wg.Done()
		}()
		wg.Wait()
	}
}
