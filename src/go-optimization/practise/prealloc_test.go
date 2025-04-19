package practise

import "testing"

func BenchmarkNopreAlloc(b *testing.B) {
	for b.Loop() {

		s := make([]int, 0)
		for range 10000 {
			s = append(s, 5)
		}
	}
}

func BenchmarkPreAlloc(b *testing.B) {
	for b.Loop() {
		s := make([]int, 0, 10000)
		for range 10000 {
			s = append(s, 5)
		}
	}
}
