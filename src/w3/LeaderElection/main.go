package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	serverName = flag.String("server-name", "", "server name")
	port       = flag.String("port", "", "server port")
)

const (
	follower = iota
	candidate
	leader
)

var (
	healthTimeout      = 100
	electionTimeoutMin = 150
	electionTimeoutMax = 300
)

type Server struct {
	peers         []string
	serverState   int
	currentTerm   int
	leaderId      string
	votedFor      string
	electionTimer *time.Timer
}

func RequestVoteFromPeer(peer string, wg *sync.WaitGroup, server *Server, vote *int32) {
	host := strings.Split(peer, ":")[0]
	port, err := strconv.Atoi(strings.Split(peer, ":")[1])
	if err != nil {
		log.Fatalf("Invalid port in peer address: %v", err)
		wg.Done()
		return
	}

	c, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(host), Port: port})
	if err != nil {
		log.Printf("Failed to connect to peer: %v", err)
		wg.Done()
		return
	}
	c.Write([]byte("request vote:" + strconv.Itoa(server.currentTerm) + ":" + *serverName + "\n"))
	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Printf("Failed to read vote response: %v", err)
		c.Close()
		wg.Done()
		return
	}
	if strings.TrimSpace(message) == "granted" {
		atomic.AddInt32(vote, 1)
	}
	c.Close()
	wg.Done()
}

func SendHeartbeatToPeers(server *Server) {
	for {
		time.Sleep(time.Duration(healthTimeout) * time.Millisecond)
		for _, peer := range server.peers {
			go SendHeartbeatToPeer(peer, server)
		}
	}
}

func SendHeartbeatToPeer(peer string, server *Server) {
	host := strings.Split(peer, ":")[0]
	port, err := strconv.Atoi(strings.Split(peer, ":")[1])
	if err != nil {
		log.Fatalf("Invalid port in peer address: %v", err)
		return
	}

	c, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(host), Port: port})
	if err != nil {
		// log.Printf("Failed to connect to peer: %v", err)
		return
	}
	c.Write(fmt.Appendf(nil, "health:%d:%s\n", server.currentTerm, *serverName))
	c.Close()
}

func StartElectionIfTimerExpired(server *Server) {
	for {
		duration := rand.Intn(electionTimeoutMax-electionTimeoutMin+1) + electionTimeoutMin
		server.electionTimer.Reset(time.Duration(duration) * time.Millisecond)
		electionExpired := <-server.electionTimer.C
		log.Printf("Election timeout expired at %v, starting election", electionExpired)
		server.serverState = candidate
		server.currentTerm++
		server.votedFor = *serverName

		var votes int32 = 1
		var wg sync.WaitGroup
		wg.Add(len(server.peers))
		for _, peer := range server.peers {
			go RequestVoteFromPeer(peer, &wg, server, &votes)
		}
		wg.Wait()

		if votes > int32(len(server.peers)/2) {
			log.Printf("Server %s became leader for term %d with %d votes", *serverName, server.currentTerm, votes)
			server.serverState = leader
			server.leaderId = *serverName
			go SendHeartbeatToPeers(server)
			return
		} else {
			log.Printf("Server %s failed to become leader for term %d with only %d votes", *serverName, server.currentTerm, votes)
			server.serverState = follower
		}
	}
}

func main() {
	parseFlags()
	time.Sleep(time.Duration(5000) * time.Millisecond) // Stagger startup times
	l, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}
	defer l.Close()
	peers := []string{"localhost:8001", "localhost:8002", "localhost:8003", "localhost:8004", "localhost:8005"}
	for i, peer := range peers {
		if strings.HasSuffix(peer, *port) {
			peers = slices.Delete(peers, i, i+1)
			break
		}
	}
	server := &Server{
		peers:         peers,
		serverState:   0, //follower
		currentTerm:   0,
		leaderId:      "",
		votedFor:      "",
		electionTimer: time.NewTimer(time.Duration(electionTimeoutMin) * time.Millisecond),
	}

	go StartElectionIfTimerExpired(server)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		HandleConnection(c, server)
	}
}

func HandleConnection(c net.Conn, server *Server) {
	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Printf("Failed to read message: %v", err)
		c.Close()
		return
	}
	if strings.HasPrefix(message, "health") {
		HandleHealthRequest(server, message, c)
	} else if strings.HasPrefix(message, "request vote") {
		HandleRequestVoteRequest(message, server, c)
	}
}

func HandleRequestVoteRequest(message string, server *Server, c net.Conn) {
	// request vote:<term>:server
	log.Printf("Received vote request: %s", message)
	term := strings.Split(message, ":")[1]
	termInt, err := strconv.Atoi(term)
	candidateId := strings.Split(message, ":")[2]
	if err != nil {
		log.Printf("Invalid term in vote request: %v", err)
		return
	}
	if server.currentTerm < termInt &&
		(server.votedFor == "" || server.votedFor == candidateId) {
		server.currentTerm = termInt
		server.votedFor = candidateId
		server.serverState = follower
		c.Write([]byte("granted\n"))
		log.Printf("Voted for %s in term %d", candidateId, termInt)
		go StartElectionIfTimerExpired(server)
	} else {
		log.Printf("Denied vote for %s in term %d", candidateId, termInt)
		c.Write([]byte("denied\n"))
	}
}

func HandleHealthRequest(server *Server, message string, c net.Conn) {
	if server.serverState == candidate {
		// If a candidate receives a heartbeat, it steps down to follower
		server.serverState = follower
	}

	parts := strings.Split(message, ":")
	term, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Invalid term in health message: %v", err)
		return
	}
	leader := parts[2]
	if term >= server.currentTerm {
		server.currentTerm = term
		server.leaderId = leader
		server.serverState = follower
	}
	c.Write([]byte("OK\n"))
	// log.Printf("Received heartbeat from leader %s for term %d", leader, term)
	server.electionTimer.Reset(time.Duration(rand.Intn(electionTimeoutMax-electionTimeoutMin+1)+electionTimeoutMin) * time.Millisecond)
	c.Close()
}

func parseFlags() {
	flag.Parse()
	if *serverName == "" {
		log.Fatalf("server-name is required")
	}
	if *port == "" {
		log.Fatalf("port is required")
	}
}
