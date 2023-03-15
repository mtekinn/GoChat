package server

import (
	"encoding/json"
	"fmt"
	"net"
)

/*
// Message is a struct to hold messages to be broadcasted.
type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
*/

// Connection is a struct to hold connection information for clients.
type Connection struct {
	Username string
	Address  net.Addr
	Encoder  *json.Encoder
}

// ClientList is a struct to hold a list of connected clients.
type ClientList struct {
	List []*Connection
}

// AddClient adds a client to the list of clients.
func (cl *ClientList) AddClient(c *Connection) {
	cl.List = append(cl.List, c)
}

// RemoveClient removes a client from the list of clients.
// It iterates through the list and removes the client if its address matches.
func (cl *ClientList) RemoveClient(c *Connection) {
	for i, client := range cl.List {
		if client.Address == c.Address {
			cl.List = append(cl.List[:i], cl.List[i+1:]...)
			break
		}
	}
}

// BroadcastMessage broadcasts a message to all clients.
// It iterates through the list of clients and sends the message using each client's Encoder.
func (cl *ClientList) BroadcastMessage(m Message) {
	for _, client := range cl.List {
		err := client.Encoder.Encode(m)
		if err != nil {
			fmt.Println("Error broadcasting message:", err)
		}
	}
}
