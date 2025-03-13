package main

import (
	"fmt"
	"log"

	"go_learn.com/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(1)
	const name = "Sanjay"
	message, err := greetings.Hello(name)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(message)

	names := []string{"Sanjay", "Gabi", "Kitty"}
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(messages)
}
