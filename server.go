package main

import (
	"deskor/chat"
	"deskor/connect"
	"deskor/log"
	"encoding/json"
	"fmt"
)

var clients = make(map[chat.Client]bool)
var join = make(chan chat.Client)
var leave = make(chan chat.Disconnect)
var messages = make(chan chat.Message)
var l *logger.FileLogger
var connected int

func main() {
	logger.New()
	l = logger.Get()
	defer l.Close()

	l.Write("Server has just started")

	listener, err := connect.Setup()
	if err != nil {
		l.Write(fmt.Sprint(err))
	}
	defer listener.Close()

	go broadcast()

	for {
		client, err := connect.Accept(listener)
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
		}
		join <- client
	}
}

func broadcast() {
	for {
		select {
		case client := <-join:
			l.Write(fmt.Sprintf("New client has joined: %s", client.Conn.RemoteAddr()))
			go func(client chat.Client) {
				welcomeMessage := chat.Message{
					Sender:    "SERVER",
					SenderIp:  "",
					Text:      "Someone has arrived",
					Connected: connected,
				}
				welcomeMessageJSON, _ := json.Marshal(welcomeMessage)
				_, err := client.Conn.Write(welcomeMessageJSON)
				if err != nil {
					l.Write(fmt.Sprintf("Error while sending welcome message to %s : %s", client.Conn.RemoteAddr(), err))
				}
				connected++
				go handleClient(client)
			}(client)
		case disconnect := <-leave:
			connected--
			l.Write(fmt.Sprintf("Client left: %s", disconnect.Client.Conn.RemoteAddr()))
		case message := <-messages:
			for client := range clients {
				go func(c chat.Client, msg chat.Message) {
					select {
					case c.Messages <- msg:
					default:
						l.Write("Message not sent to clients")
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
		leave <- chat.Disconnect{Client: client}
		client.Conn.Close()
	}()

	go func() {
		for message := range client.Messages {
			msg := chat.Message{
				Sender:    message.Sender,
				SenderIp:  client.Conn.RemoteAddr().String(),
				Text:      message.Text,
				Connected: connected,
			}

			if chat.IsCommand(msg) {
				msg = chat.Dispatch(msg)
			}

			messageJSON, err := json.Marshal(msg)
			if err != nil {
				l.Write(fmt.Sprintf("Error sending message to client: %s", err))
				return
			}
			_, err = client.Conn.Write(messageJSON)
			if err != nil {
				l.Write(fmt.Sprintf("Error sending message to client: %s", err))
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
			l.Write(fmt.Sprintf("Received invalid message: %s", err))
		}
	}
}
