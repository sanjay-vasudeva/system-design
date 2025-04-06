package main

import (
	"fmt"
)

func reverse(a string) string {
	b := []byte(a)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	a := "The quick brown fox jumped over the lazy dog"
	rev := reverse(a)
	double_rev := reverse(rev)
	fmt.Println(rev)
	fmt.Println(double_rev)
}
