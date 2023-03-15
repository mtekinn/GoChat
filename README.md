# Go Chat

A simple chat server and client implementation in Go.

## Features

- TCP-based chat server and client
- Client-server communication using JSON messages
- Supports multiple concurrent clients
- Broadcasts messages to all connected clients
- Graceful server shutdown using signals

## Usage

### Running the server

1. Clone the repository: git clone https://github.com/mtekinn/go-chat.git
2. Change to the project directory: cd go-chat
3. Build the server: go build main.go client_list.go
4. Run the server: ./main

The server will listen for incoming connections on `localhost:8080`.

### Running the client

To connect to the chat server, use a terminal-based application like `nc` or `telnet`. For example: nc localhost 8080 or telnet localhost 8080

## Contributing

Contributions, issues, and feature requests are welcome!
