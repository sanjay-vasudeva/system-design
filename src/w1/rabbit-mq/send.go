package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Send() {
	conn, chn, q := Setup()
	defer conn.Close()
	defer chn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	for range 10 {
		err := chn.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
			Body:        []byte(body),
			ContentType: "text/plain",
		})
		FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
		time.Sleep(100 * time.Millisecond)
	}
}
