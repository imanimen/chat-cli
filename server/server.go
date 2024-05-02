package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

//
type Client struct {
	Name 	string
	Conn    net.Conn
	Messages chan string
}

var clients = make(map[net.Conn]Client)
var mutex sync.Mutex

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer listener.Close()
	log.Println("Server Sstarted. Listening on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("Client connected: ", conn.RemoteAddr())

	client := Client{
		Conn: conn,
		Messages: make(chan string),
	}
	// prompt the username
	_, _ = conn.Write([]byte("Enter your name: "))
	name, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		log.Println("Failed to read name:", err)
		return
	}

	client.Name = strings.TrimSpace(name)
	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	// broadcast users that joined
	broadcastMessage(fmt.Sprintf("User %s joined the chat", client.Name), nil) //
	// broadcast
	go broadcastClientMessage(client)
	// [iman] salam khobi?
	for {
		select {
			case  message := <- client.Messages:
				broadcastMessage(fmt.Sprintf("[%s]: %s ", client.Name, message), conn)
		}
	}

}

func broadcastClientMessage(client Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Failed to read message:", err)
			return
		}
		message = strings.TrimSpace(message)
		if message != "" {
			client.Messages <- message
		}
	}
}

func broadcastMessage(message string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn, _ := range clients {
		if conn != sender {
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Failed to write message:", err)
			}
		}
	}
}

