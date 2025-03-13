package main

import (
	"fmt"
)

type Animal interface {
	DoSomething()
}

type Cat struct {
	Name string
}

func (c Cat) DoSomething() {
	fmt.Println(c.Name, "is meowing!")
}

func helloWorld() {
	fmt.Println("Hello, World!")
}

func main() {
	var i int = 10
	var p *int = &i
	fmt.Println("Value of i: ", i)
	fmt.Println("Value of i through pointer: ", *p)
	fmt.Println("Address of i: ", p)

	//Function types
	var func1, func2 func()
	func1 = helloWorld
	func2 = helloWorld

	fmt.Println("Value of func1: ", func1)
	fmt.Println("Value of func2: ", func2)

	//Interface types
	var animal1, animal2 Animal
	cat1 := Cat{Name: "Kitty"}
	animal1 = cat1
	// animal1.Name = "Gabi" // cannot do this

	fmt.Println("\nInterface with non-pointer value:")
	animal1.DoSomething()
	cat1.Name = "Gabi"
	animal1.DoSomething()

	cat2 := Cat{Name: "Kitty"}
	animal2 = &cat2
	// animal2.Name = "Gabi" // cannot do this

	fmt.Println("\nInterface with pointer value:")
	animal2.DoSomething()
	cat2.Name = "Gabi"
	animal2.DoSomething()

	/*
	 * String types
	 */
	str1 := "Cello"
	str2 := []byte(str1)

	fmt.Println("\nValue of str1:", str1)
	fmt.Println("Value of str2:", string(str2))

	// str1[0] = 'H' // cannot do this
	str2[0] = byte('H')
	fmt.Println("Value of str2 after modification:", string(str2))
}
