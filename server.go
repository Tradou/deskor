package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
)

type Client struct {
	conn     net.Conn
	messages chan string
}

type Disconnect struct {
	client Client
}

var clients = make(map[Client]bool)
var join = make(chan Client)
var leave = make(chan Disconnect)
var messages = make(chan string)

func main() {
	err := godotenv.Load(".env.server")
	if err != nil {
		log.Fatal("Error loading env var")
	}
	port := os.Getenv("PORT")

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port 8080")

	go handleMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := Client{
			conn:     conn,
			messages: make(chan string),
		}

		join <- client
		go handleClient(client)
	}
}

func handleClient(client Client) {
	clients[client] = true

	defer func() {
		delete(clients, client)
		leave <- Disconnect{client}
		client.conn.Close()
	}()

	go func() {
		for msg := range client.messages {
			_, err := client.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				return
			}
		}
	}()

	for {
		message := make([]byte, 512)
		_, err := client.conn.Read(message)
		if err != nil {
			break
		}
		msg := string(message)
		messages <- msg
	}
}

func handleMessages() {
	for {
		select {
		case client := <-join:
			fmt.Println("New client joined:", client.conn.RemoteAddr())
		case disconnect := <-leave:
			fmt.Println("Client left:", disconnect.client.conn.RemoteAddr())
		case message := <-messages:
			fmt.Println("Message received:", message)
			for client := range clients {
				go func(c Client, msg string) {
					select {
					case c.messages <- msg:
					default:
						fmt.Println("Message not sent to clients:", c.conn.RemoteAddr())
					}
				}(client, message)
			}
		}
	}
}
