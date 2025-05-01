package main

import (
	"fmt"
	"order/svc"
	"sync"
	"time"

	"github.com/sanjay-vasudeva/ioutil"
)

func main() {
	// reset()
	var wg sync.WaitGroup
	start := time.Now()
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Simulate some work
			order, err := svc.PlaceOrder()
			if err != nil {
				fmt.Println("Error placing order:", err)
				return
			}
			fmt.Println("Order placed successfully:", order.ID)
		}()
	}
	wg.Wait()
	fmt.Printf("Took %f seconds", time.Since(start).Seconds())
}

func reset() {
	conn := ioutil.NewConn("3308", "root", "password", "delivery")
	tx, err := conn.Begin()
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec("UPDATE stock SET is_reserved = 0, order_id = NULL")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec("UPDATE agent SET is_reserved = 0, order_id = NULL")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
