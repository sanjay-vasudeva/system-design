package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
)

type NegativeSqrtError float64

func (e NegativeSqrtError) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number. Value: %f", float64(e))
}

func Sqrt(v float64) (float64, error) {
	if v < 0 {
		return 0, NegativeSqrtError(v)
	}
	//To do: implement Sqrt
	return v, nil
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *Person) String() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func readJson() *Person {
	personJson := `{"name": "John", "age": 30`
	var person Person
	err := json.Unmarshal([]byte(personJson), &person)
	if err != nil {
		panic(fmt.Sprintf("Failed to read json with error %q", err.Error()))
	}
	return &person
}
func main() {
	// v, err := Sqrt(-2)
	// if err != nil {
	// 	// if nerr, ok := err.(NegativeSqrtError); ok {
	// 	// 	fmt.Println("NegativeSqrtError: ", nerr.Error())
	// 	// }
	// 	switch err.(type) {
	// 	case NegativeSqrtError:
	// 		fmt.Println("NegativeSqrtError")
	// 	default:
	// 		fmt.Println("Unknown error")
	// 	}
	// } else {
	// 	fmt.Println(v)
	// }
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered in main: ", err, string(debug.Stack()))
		}
	}()
	person := readJson()
	fmt.Println(person)
}
