package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Message struct {
	Body     []byte
	ClientId int
}

type Hub struct {
	Clients    map[*Client]bool // Connected clients
	Broadcast  chan Message     // Channel to broadcast messages to clients
	Register   chan *Client     // Register new clients
	Unregister chan *Client     // Unregister clients
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			fmt.Printf("Client registered: %d\n", client.Id)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				client.Conn.Close()
			}
		case msg := <-h.Broadcast:
			jmsg, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error marshalling message: %v\n", err)
				continue
			}
			if res := rClient.Publish(context.Background(), "ch:0", jmsg); res.Err() != nil {
				fmt.Printf("Error publishing message to Redis: %v\n", res.Err())
			}

			for c := range h.Clients {
				if c.Id == msg.ClientId {
					continue // Skip sending the message back to the sender
				}
				select {
				case c.Send <- msg.Body:
				default:
					fmt.Println("Failed to send message to client, closing connection")
					close(c.Send)        // Close the send channel to prevent further writes
					delete(h.Clients, c) // Remove client if sending fails
				}
			}
		}
	}
}
