package main

import (
	"fmt"
	"time"
)

type Ball struct {
	hits int
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Printf("%s hit the ball. Hit count: %d\n", name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func main() {
	// Ping pong game simulation
	table := make(chan *Ball)
	go player("Sanjay", table)
	go player("Aarthi", table)
	//toss the ball
	table <- new(Ball)
	time.Sleep(5 * time.Second)
	<-table
}
