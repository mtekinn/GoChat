package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/mtekinn/go-chat/server"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// Read and handle incoming messages from the server
	go func() {
		for {
			message := server.Message{}
			err := decoder.Decode(&message)
			if err != nil {
				fmt.Println("Error decoding message from the server:", err)
				return
			}
			fmt.Printf("%s: %s\n", message.Username, message.Message)
		}
	}()

	// Send the username to the server
	fmt.Print("Enter your username: ")
	username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	username = strings.TrimSpace(username)
	encoder.Encode(server.Message{Type: "request_username", Username: username})

	// Read and send messages to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		encoder.Encode(server.Message{Type: "message", Username: username, Message: message})
	}
}
