package practise

import "testing"

type Worker interface {
	Work()
}

type Largejob struct {
	payload [4096]byte
}

func (Largejob) Work() {

}

const count int = 1000

func BenchmarkBoxing(b *testing.B) {
	jobs := make([]Worker, 0, count)
	for b.Loop() {
		jobs = jobs[:0]
		for range count {
			var job Largejob
			jobs = append(jobs, job)
		}
	}
}

func BenchmarkNonBoxing(b *testing.B) {
	jobs := make([]Worker, 0, count)
	for b.Loop() {
		jobs = jobs[:0]
		for range count {
			job := &Largejob{}
			jobs = append(jobs, job)
		}
	}
}

// bench-call-start
var sink Worker

func call(w Worker) {
	sink = w
}

func BenchmarkCallWithValue(b *testing.B) {
	for b.Loop() {
		var j Largejob
		call(j)
	}
}

func BenchmarkCallWithPointer(b *testing.B) {
	for b.Loop() {
		j := &Largejob{}
		call(j)
	}
}
