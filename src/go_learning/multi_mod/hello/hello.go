package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	fmt.Println(reverse.String("Hello, World!"))
	fmt.Println(reverse.Int(569))
}
