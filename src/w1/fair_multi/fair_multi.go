package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var numberOfPrimes int32
var MAX_INT int = 100000000

func main() {
	//basic_single()
	//basic_multi(10)
	// fair_multi(10)
}

func checkPrime(i int) {
	if i < 2 {
		return
	}
	for j := 2; j*j <= i; j++ {
		if i%j == 0 {
			return
		}
	}
	atomic.AddInt32(&numberOfPrimes, 1)
}

func basic_single() {
	now := time.Now()
	for i := 0; i < MAX_INT; i++ {
		checkPrime(i)
	}
	fmt.Printf("Thread completed in %v\n", time.Since(now))
	fmt.Printf("Total number of primes: %d\n", numberOfPrimes)
}

func basic_multi(threads int) {
	wg := sync.WaitGroup{}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			now := time.Now()
			for i := start; i < end; i++ {
				checkPrime(i)
			}
			fmt.Printf("Thread %d completed start:%d, end:%d in %v\n", i, start, end, time.Since(now))
		}((i*(MAX_INT/threads))+1, (i+1)*(MAX_INT/threads))
	}
	wg.Wait()
	fmt.Printf("Total number of primes: %d\n", numberOfPrimes)
}

func fair_multi(threads int) {
	wg := sync.WaitGroup{}
	var curr_num atomic.Int32 = atomic.Int32{}
	curr_num.Store(-1)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			now := time.Now()
			for {
				num := curr_num.Add(1)
				if int(num) > MAX_INT {
					break
				}
				checkPrime(int(num))
			}
			fmt.Printf("Thread %d completed in %v\n", i, time.Since(now))
		}()
	}
	wg.Wait()
	fmt.Printf("Total number of primes: %d\n", numberOfPrimes)
}
