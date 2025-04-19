package practise

import (
	"io"
	"os"
	"testing"

	"golang.org/x/exp/mmap"
)

func BenchmarkCopy(b *testing.B) {
	data := make([]byte, 64*1024)
	for b.Loop() {
		buf := make([]byte, len(data))
		copy(buf, data)
	}
}

func BenchmarkSlice(b *testing.B) {
	data := make([]byte, 64*1024)
	for b.Loop() {
		_ = data[:]
	}
}

const BUFFER_SIZE int = 2 * 1024 * 1024

func BenchmarkReadWithCopy(b *testing.B) {
	for b.Loop() {
		f, err := os.Open("large_file.txt")
		if err != nil {
			b.Fatal(err)
		}
		defer f.Close()

		buf := make([]byte, BUFFER_SIZE)
		_, err = f.ReadAt(buf, 0)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadWithMmap(b *testing.B) {
	for b.Loop() {
		r, err := mmap.Open("large_file.txt")
		if err != nil {
			b.Fatal(err)
		}
		defer r.Close()
		buf := make([]byte, r.Len())
		_, err = r.ReadAt(buf, 0)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}
