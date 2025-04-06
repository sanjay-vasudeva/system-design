package main

import (
	"fmt"
)

func AddInts(a, b int) int {
	return a + b
}

func AddFloats(a, b float64) float64 {
	return a + b
}

func AddNumbers[T int | float64](a, b T) T {
	return a + b
}

type Number interface {
	int | float64
}

func Add[T Number](a, b T) T {
	return a + b
}
func main() {
	fmt.Println(AddInts(1, 2))
	fmt.Println(AddFloats(1.1, 2.2))
	fmt.Println(AddNumbers(1, 2))
	fmt.Println(AddNumbers(1.1, 2.2))
	fmt.Println(Add(1, 2))
	fmt.Println(Add(1.1, 2.2))
}
