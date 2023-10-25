package main

import (
	"deskor/chat"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
)

var clients = make(map[chat.Client]bool)
var join = make(chan chat.Client)
var leave = make(chan chat.Disconnect)
var messages = make(chan chat.Message)

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

	go broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := chat.Client{
			Conn:     conn,
			Messages: make(chan chat.Message),
		}

		join <- client
	}
}

func broadcast() {
	for {
		select {
		case client := <-join:
			fmt.Println("New client joined:", client.Conn.RemoteAddr())
			go handleClient(client)
		case disconnect := <-leave:
			fmt.Println("Client left:", disconnect.Client.Conn.RemoteAddr())
		case message := <-messages:
			fmt.Println("Message received:", message.Text)
			for client := range clients {
				go func(c chat.Client, msg chat.Message) {
					select {
					case c.Messages <- msg:
					default:
						fmt.Println("Message not sent to clients:", c.Conn.RemoteAddr())
					}
				}(client, message)
			}
		}
	}
}

func handleClient(client chat.Client) {
	clients[client] = true
	defer func() {
		delete(clients, client)
		leave <- chat.Disconnect{client}
		client.Conn.Close()
	}()

	go func() {
		for message := range client.Messages {
			msg := chat.Message{
				Sender:   message.Sender,
				SenderIp: client.Conn.RemoteAddr(),
				Text:     message.Text,
			}

			messageJSON, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				return
			}
			_, err = client.Conn.Write(messageJSON)
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				return
			}
		}
	}()

	for {
		message := make([]byte, 512)
		n, err := client.Conn.Read(message)
		if err != nil {
			break
		}
		msg := string(message[:n])

		var chatMsg chat.Message
		if err := json.Unmarshal([]byte(msg), &chatMsg); err == nil {
			messages <- chatMsg
		} else {
			fmt.Println("Received invalid message:", msg)
		}
	}
}
