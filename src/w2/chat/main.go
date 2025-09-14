package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
}

var rClient *redis.Client
var pubsub *redis.PubSub

func main() {
	flag.Parse()
	hub := NewHub()
	go hub.Run() // Start the hub to manage clients and messages
	rClient = createRedisClient()
	if res := rClient.Ping(context.Background()); res.Err() != nil {
		fmt.Printf("Error connecting to Redis: %v\n", res.Err())
		return
	}
	pubsub = rClient.Subscribe(context.Background())

	go func() {
		for msg := range pubsub.Channel() {
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				fmt.Printf("Error unmarshalling message from Redis: %v\n", err)
				continue
			}
			for c := range hub.Clients {
				if c.Id == message.ClientId {
					continue // Skip sending the message back to the sender
				}
				select {
				case c.Send <- message.Body:
				default:
					fmt.Println("Failed to send message to client, closing connection")
					close(c.Send)          // Close the send channel to prevent further writes
					delete(hub.Clients, c) // Remove client if sending fails
				}
			}
			// hub.Broadcast <- Message{Body: []byte(message.Body), ClientId: message.ClientId} // Broadcast messages received from Redis
		}
	}()
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
			return
		}
		client := &Client{Conn: conn, Hub: hub, Send: make(chan []byte, 10), Id: rand.Intn(1000)}
		hub.Register <- client

		//subscribe to channels that the client is part of
		pubsub.Subscribe(context.Background(), "ch:0")
		go client.WriteMessages() // Start writing messages to the client
		go client.ReadMessages()  // Start reading messages from the client
	})

	r.Run(fmt.Sprintf(":%s", flag.Arg(0))) // Start the server on port 8080
}
