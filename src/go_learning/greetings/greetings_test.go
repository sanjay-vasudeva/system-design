package greetings

import (
	"regexp"
	"testing"
)

func TestHelloname(t *testing.T) {
	name := "Sanjay"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello("Sanjay")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Sanjay") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
