package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
)

var backends []Backend

func main() {
	// conn := CreateZooKeeperConn()
	// WatchInRepeat(conn, "/test")
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	backends = append(backends, Backend{1, "localhost", "8081", "backend1", true, 0})
	backends = append(backends, Backend{2, "localhost", "8082", "backend2", true, 0})
	backends = append(backends, Backend{3, "localhost", "8083", "backend3", true, 0})

	log.Printf("Load balancer started on port 8000 with %d backends\n", len(backends))
	var incoming IncomingChannel = make(chan IncomingRequest, 100)
	workers := 100
	for i := 0; i < workers; i++ {
		go handleConnection(incoming)
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			for i, _ := range backends {
				count := atomic.LoadInt32(&backends[i].NumRequests)
				fmt.Printf("Backend %s has handled %d requests\n", backends[i].name, count)
			}
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			break
		}
		reqId, _ := uuid.DefaultGenerator.NewV7()

		incoming <- IncomingRequest{conn, reqId.String()}
	}
}

func handleConnection(incoming IncomingChannel) {
	for req := range incoming {
		backend := &backends[rand.Intn(len(backends))]
		bConn, err := net.Dial("tcp", net.JoinHostPort(backend.host, backend.port))
		if err != nil {
			req.conn.Write([]byte("Error connecting to backend: " + backend.name + "\n"))
			req.conn.Close()
			return
		}
		atomic.AddInt32(&backend.NumRequests, 1)
		go io.Copy(bConn, req.conn)
		go io.Copy(req.conn, bConn)
	}
}

type IncomingChannel chan IncomingRequest

type IncomingRequest struct {
	conn  net.Conn
	reqId string
}

type Backend struct {
	id          int
	host        string
	port        string
	name        string
	isHealthy   bool
	NumRequests int32
}
