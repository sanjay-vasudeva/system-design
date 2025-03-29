package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan Message {
	c := make(chan Message)
	go func() {
		waitForIt := make(chan bool)
		for i := 0; ; i++ {
			c <- Message{msg: fmt.Sprintf("%s %d", msg, i), wait: waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c
}

func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

type Message struct {
	msg  string
	wait chan bool
}

func main() {
	// fmt.Println("I'm listening...")
	// c := boring("boring")
	// for i := 0; i < 5; i++ {
	// 	fmt.Printf("You say %q\n", <-c)
	// }
	// fmt.Println("You're boring.. I'm leaving")
	c := fanIn(boring("joe"), boring("ann"))
	for range 10 {
		msg1 := <-c
		fmt.Println(msg1.msg)
		msg2 := <-c
		fmt.Println(msg2.msg)
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're both boring.. I'm leaving")

}
