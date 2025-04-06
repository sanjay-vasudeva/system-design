package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type SampleJson struct {
	Text string `json:"text"`
}

func DeferExample() {
	//File open
	f, err := os.Open("sample.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var s *SampleJson
	json.NewDecoder(f).Decode(&s)
	fmt.Println(s.Text)
}

func RecoverExample() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic", r)
		}
	}()
	panic("Panic message")
}

func ArrayExample() [5]int {
	arr := [...]int{1, 2, 3, 4, 5}
	return arr
}

func WordCount(s string) map[string]int {
	var mp = make(map[string]int)
	words := strings.Fields(s)
	for _, w := range words {
		if v, ok := mp[w]; ok {
			mp[w] = v + 1
			continue
		}
		mp[w] = 1
	}
	return mp
}

func funcTest(f func(string) string) {
	val := f("Hello")
	fmt.Println(val)
}

type Vertex struct {
	X int
	Y int
}

func (v Vertex) Print() {
	fmt.Println(v.X, v.Y)
}

func (v Vertex) String() string {
	return fmt.Sprintf("X: %v, Y: %v", v.X, v.Y)
}

func (v Vertex) ScaleValueReceiver(f int) Vertex {
	v.X = v.X * f
	v.Y = v.Y * f
	return v
}

func (v *Vertex) ScalePointerReceiver(f int) {
	v.X = v.X * f
	v.Y = v.Y * f
}

type MyFloat float64

func (f MyFloat) Add(v float64) float64 {
	return float64(f) + v
}

type Abser interface {
	Abs() float64
}

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func (v *Vertex) Abs() float64 {
	if v == nil {
		fmt.Println("<NIL> Vertex")
		return 0
	}
	return float64(v.X*v.X + v.Y*v.Y)
}

func describe(i interface{}) {
	fmt.Printf("%v %T\n", i, i)
}

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.

func (ip IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", ip[0], ip[1], ip[2], ip[3])
}

type MyError struct {
	When time.Time
	What string
}

func main() {
	// DeferExample()
	// RecoverExample()
	// funcTest(func(s string) string {
	// 	return s + "Hello"
	// })
	// f := MyFloat(1.2)
	// fmt.Println(f.Add(2.2))
	// v := Vertex{X: 1, Y: 2}
	// v.Print()

	// v = v.ScaleValueReceiver(5)
	// v.Print()

	// v.ScalePointerReceiver(2)
	// v.Print()

	// var a Abser
	// a = MyFloat(3)
	// fmt.Println(a.Abs())

	// a = &Vertex{1, 2}
	// fmt.Println(a.Abs())

	// switch a.(type) {
	// case MyFloat:
	// 	fmt.Println("MyFloat type")
	// case *Vertex:
	// 	fmt.Println("*Vertex type")
	// default:
	// 	fmt.Println("Unknown type")
	// }

	// var v *Vertex
	// a = v
	// fmt.Println(a.Abs())

	// var i interface{}
	// i = 52
	// describe(i)
	// i = "Hello"
	// describe(i)

	// var i1 interface{}
	// t, ok := i1.(int64)
	// if ok {
	// 	fmt.Print("success. ", t)
	// }
	// t = i1.(int64)
	// fmt.Print(t)
	// v1 := Vertex{1, 2}
	// fmt.Println(v1)
	// var e error = errors.New("")

	r := strings.NewReader("Hello, Reader!")

	for {
		b := make([]byte, 6)
		n, err := r.Read(b)
		fmt.Printf("n: %v, err: %v, Read: %s\n", n, err, b[:n])
		if err == io.EOF {
			break
		}
	}
	ch := make(chan int)
	wg.Add(2)
	go produce(ch)
	go consume(ch)

	wg.Wait()
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

var wg sync.WaitGroup = sync.WaitGroup{}

func produce(ch chan int) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("producer done")
}

func consume(ch chan int) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		v := <-ch
		fmt.Printf("Value: %v\n", v)
	}
	fmt.Println("consumer done")
}
