package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// clients is a map that holds connected clients
var (
	clients    = make(map[net.Conn]struct{})
	clientsMtx sync.Mutex
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// main is the entry point of the chat server
func main() {
	// Create a context that can be used to gracefully shutdown the server
	ctx, cancel := context.WithCancel(context.Background())

	// Listen for incoming connections
	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Listening on " + connHost + ":" + connPort)

	// Start a goroutine to handle incoming signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Received signal, shutting down...")
		cancel()
	}()

	// Start a goroutine to handle incoming connections
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		wg.Add(1)
		go handleConnection(ctx, &wg, conn)
	}

	// Wait for all client goroutines to finish
	wg.Wait()
}

// handleConnection manages the lifecycle of a client's connection to the chat server.
// It sends a request for the client's username, reads incoming messages from the client,
// and broadcasts them to other clients. The function also monitors the context
// for cancellation, which indicates that the server is shutting down.
func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	clientsMtx.Lock()
	clients[conn] = struct{}{}
	clientsMtx.Unlock()

	defer func() {
		removeClient(conn)
		wg.Done()
	}()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// Send a message to the client asking for their username
	encoder.Encode(Message{Type: "request_username"})

	usernameMsg := Message{}
	err := decoder.Decode(&usernameMsg)
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	username := usernameMsg.Username

	broadcast(Message{Type: "user_joined", Username: username}, conn)

	for {
		message := Message{}
		err := decoder.Decode(&message)
		if err != nil {
			broadcast(Message{Type: "user_left", Username: username}, conn)
			return
		}
		broadcast(Message{Type: "message", Username: username, Message: message.Message}, conn)

		// Check if the context has been cancelled
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// broadcast sends a message to all connected clients except the sender.
// The message is a struct containing the message type, sender's username, and the message text.
func broadcast(message Message, sender net.Conn) {
	clientsMtx.Lock()
	defer clientsMtx.Unlock()
	for conn := range clients {
		if conn != sender {
			encoder := json.NewEncoder(conn)
			encoder.Encode(message)
		}
	}
}

// removeClient removes a client from the clients map and closes their connection.
// It acquires a lock on the clientsMtx mutex to ensure safe concurrent access.
func removeClient(conn net.Conn) {
	clientsMtx.Lock()
	delete(clients, conn)
	clientsMtx.Unlock()
	conn.Close()
}
