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
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

type Message struct {
	msg  string
	wait chan bool
}

type Result string
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep((time.Duration(rand.Intn(100)) * time.Millisecond))
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

func Google(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	timeout := time.After(80 * time.Millisecond)
	for range 3 {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}
	return
}
func main() {
	// fmt.Println("I'm listening...")
	// c := boring("boring")
	// for i := 0; i < 5; i++ {
	// 	fmt.Printf("You say %q\n", <-c)
	// }
	// fmt.Println("You're boring.. I'm leaving")
	// c := fanIn(boring("joe"), boring("ann"))
	// for range 10 {
	// 	msg1 := <-c
	// 	fmt.Println(msg1.msg)
	// 	msg2 := <-c
	// 	fmt.Println(msg2.msg)
	// 	msg1.wait <- true
	// 	msg2.wait <- true
	// }
	// fmt.Println("You're both boring.. I'm leaving")

	//Timeout using select
	// c := boring("joe")
	// timeout := time.After(3 * time.Second)
	// for {
	// 	select {
	// 	case s := <-c:
	// 		fmt.Println(s.msg)
	// 		s.wait <- true
	// 	case <-timeout:
	// 		fmt.Println("You're too slow..")
	// 		return
	// 	}
	// }

	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Printf("Search took %s\n", elapsed)
}
