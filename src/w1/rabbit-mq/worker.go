package main

import (
	"log"
	"math/rand/v2"
	"time"
)

func Consume(num int) {
	conn, ch, q := Setup()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("[%d] Received a message: %s", num, d.Body)
			time.Sleep(time.Duration(rand.Int64N(1000)) * time.Millisecond)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
