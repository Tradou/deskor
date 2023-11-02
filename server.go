package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"deskor/chat"
	"deskor/log"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var clients = make(map[chat.Client]bool)
var join = make(chan chat.Client)
var leave = make(chan chat.Disconnect)
var messages = make(chan chat.Message)
var l *logger.FileLogger

func main() {
	logger.New()
	l = logger.Get()
	defer l.Close()

	l.Write("Server has just started")

	err := godotenv.Load(".env.server")
	if err != nil {
		l.Write("Error loading env var")
	}
	port := os.Getenv("PORT")

	cert, err := tls.LoadX509KeyPair("./cert/server.pem", "./cert/server.key")
	if err != nil {
		l.Write(fmt.Sprintf("Error while loading pair certificate: %s", err))

	}

	caCert, err := os.ReadFile("./cert/ca.crt")
	if err != nil {
		l.Write(fmt.Sprintf("Error while reading ca certificate: %s", err))
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port), &config)
	if err != nil {
		l.Write(fmt.Sprintf("Error starting the server: %s", err))
		return
	}
	defer listener.Close()

	go broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			l.Write(fmt.Sprintf("Error accepting connection: %s", err))
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
			l.Write(fmt.Sprintf("New client has joined: %s", client.Conn.RemoteAddr()))
			go func(client chat.Client) {
				welcomeMessage := chat.Message{
					Sender:   "Server",
					SenderIp: "",
					Text:     "Someone has arrived",
				}
				welcomeMessageJSON, _ := json.Marshal(welcomeMessage)
				_, err := client.Conn.Write(welcomeMessageJSON)
				if err != nil {
					l.Write(fmt.Sprintf("Error while sending welcome message to %s : %s", client.Conn.RemoteAddr(), err))
				}

				go handleClient(client)
			}(client)
		case disconnect := <-leave:
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
		leave <- chat.Disconnect{client}
		client.Conn.Close()
	}()

	go func() {
		for message := range client.Messages {
			msg := chat.Message{
				Sender:   message.Sender,
				SenderIp: client.Conn.RemoteAddr().String(),
				Text:     message.Text,
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
