package main

import (
	"testing"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"The quick brown fox jumped over the lazy dog", "god yzal eht revo depmuj xof nworb kciuq ehT"},
		{"Hello, sanjay", "yajnas ,olleH"},
		{"", ""},
	}
	for _, testcase := range testcases {
		got := reverse(testcase.in)
		if got != testcase.want {
			t.Errorf("reverse(%q) = %q; want %q", testcase.in, got, testcase.want)
		}
	}
}
