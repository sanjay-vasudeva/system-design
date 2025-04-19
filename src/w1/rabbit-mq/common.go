package mq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Setup() (conn *amqp.Connection, channel *amqp.Channel, q amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	channel, err = conn.Channel()
	FailOnError(err, "Failed to open a channel")

	q, err = channel.QueueDeclare("hello", false, false, false, false, nil)
	FailOnError(err, "Failed to declare a queue")

	return conn, channel, q
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
