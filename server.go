package main

import (
	"deskor/chat"
	"deskor/connect"
	logger "deskor/log"
	"encoding/json"
	"log"
)

var server *chat.Server

func main() {
	logger.New()
	log.Print("Server has just started")

	server := chat.NewServer()

	listener, err := connect.Setup()
	if err != nil {
		log.Print(err)
		return
	}
	defer listener.Close()

	go broadcast()

	for {
		client, err := connect.Accept(listener)
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
		}
		server.JoinClient(client)
	}
}

func broadcast() {
	for {
		select {
		case client := <-server.Join:
			handleJoin(client)
		case disconnect := <-server.Leave:
			handleLeave(disconnect.Client)
		case message := <-server.Messages:
			handleMessage(message)
		}
	}
}

func handleLeave(client chat.Client) {
	server.DecrementConnectedCount()
	log.Printf("Client left: %s", client.Conn.RemoteAddr())
}

func handleJoin(client chat.Client) {
	log.Printf("New client has joined: %s", client.Conn.RemoteAddr())
	go func(client chat.Client) {
		if err := server.SendWelcomeMessage(client); err != nil {
			log.Printf("Error while sending welcome message to %s : %s", client.Conn.RemoteAddr(), err)
		}
		server.IncrementConnectedCount()
		go handleClient(client)
	}(client)
}

func handleMessage(message chat.Message) {
	for client := range server.Clients {
		go func(c chat.Client, msg chat.Message) {
			select {
			case c.Messages <- msg:
			default:
				log.Print("Message not sent to clients")
			}
		}(client, message)
	}
}

func handleClient(client chat.Client) {
	server.Clients[client] = true
	defer server.Clean(client)

	go sendClientMessages(client)

	for {
		message := make([]byte, 512)
		n, err := client.Conn.Read(message)
		if err != nil {
			break
		}
		msg := string(message[:n])

		var chatMsg chat.Message
		if err := json.Unmarshal([]byte(msg), &chatMsg); err == nil {
			server.AddMessage(chatMsg)
		} else {
			log.Printf("Received invalid message: %s", err)
		}
	}
}

func sendClientMessages(client chat.Client) {
	for message := range client.Messages {
		msg := chat.Message{
			Sender:    message.Sender,
			SenderIp:  client.Conn.RemoteAddr().String(),
			Text:      message.Text,
			Connected: server.GetConnectedCount(),
		}

		if chat.IsCommand(msg) {
			msg = chat.Dispatch(msg)
		}

		messageJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error sending message to client: %s", err)
			return
		}
		_, err = client.Conn.Write(messageJSON)
		if err != nil {
			log.Printf("Error sending message to client: %s", err)
			return
		}
	}
}
