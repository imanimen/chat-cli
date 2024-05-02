package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your username: ")
	scanner.Scan()

	name := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println("Error in reading name:", err)
		return
	}

	_, err = conn.Write([]byte(name + "\n"))
	if err != nil {
		log.Println("Failed to send name:", err)
		return
	}

	// receive messages
	go receiveMessages(conn)

	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("Failed to send message:", err)
			break
		}
	}
	
	if scanner.Err() != nil {
		log.Println("Error in reading from Stdin:", scanner.Err())
	}
	
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		fmt.Println(message)
	}
}