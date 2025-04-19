package mq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateTasks() {
	conn, chn, q := Setup()
	defer conn.Close()
	defer chn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := range 10 {
		body := fmt.Sprintf("msg %d", i)
		err := chn.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
		})
		FailOnError(err, "Failed to publish a message")
		// log.Printf(" [x] Sent %s", body)
		time.Sleep(100 * time.Millisecond)
	}
}
