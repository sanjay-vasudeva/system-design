package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   int
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte // Channel to send messages to the client
}

func (c *Client) ReadMessages() {
	defer func() {
		c.Hub.Unregister <- c // Unregister client on read completion
		c.Conn.Close()        // Close the connection
	}()
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(s string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		mType, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s - Error reading message from client %d: %v\n", time.Now(), c.Id, err.Error())
			c.Hub.Unregister <- c // Unregister client on read error
			return
		}
		if mType != websocket.TextMessage {
			fmt.Printf("Unexpected message type from client %d: %d\n", c.Id, mType)
			continue // Ignore non-text messages
		}
		message = bytes.Replace(message, []byte("\n"), []byte(" "), -1) // Remove newlines from message
		fmt.Printf("Message received from client %d: %s\n", c.Id, message)
		c.Hub.Broadcast <- Message{Body: message, ClientId: c.Id} // Broadcast message to the hub
	}
}

var pingPeriod time.Duration = 2 * time.Second
var writeWait time.Duration = 5 * time.Second
var pongWait time.Duration = 5 * time.Second

func (c *Client) WriteMessages() {
	defer func() {
		c.Hub.Unregister <- c // Unregister client on write completion
		c.Conn.Close()        // Close the connection
	}()

	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Hub.Unregister <- c // Unregister client on write error
				return
			}
			fmt.Printf("Writing message to client %d: %s\n", c.Id, msg)
			writer.Write(msg)

			n := len(c.Send)
			for range n {
				writer.Write([]byte("\n"))
				writer.Write(<-c.Send)
			}
			writer.Close()
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Printf("Failed to write ping to client %d: %v\n", c.Id, err)
				c.Hub.Unregister <- c // Unregister client on ping error
				return
			}
		}
	}
}
