package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

func main() {
	err := godotenv.Load(".env.client")
	if err != nil {
		log.Fatal("Error loading env var")
	}

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	serverAddr := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Print("Type your username: ")
	username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("Connected")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error while reading messages:")
				close(exit)
				break
			}
			message := string(buffer[:n])

			var receivedMessage Message
			if err := json.Unmarshal([]byte(message), &receivedMessage); err == nil {
				fmt.Printf("\033[1A\033[K")
				fmt.Printf("%s: %s\n", receivedMessage.Sender, receivedMessage.Text)
			} else {
				// TODO: idk in which case it's possible, and how to handle it
				fmt.Printf("\033[1A\033[K")
				fmt.Print(message)
			}
		}
	}()

	go func() {
		for {
			text, _ := bufio.NewReader(os.Stdin).ReadString('\n')

			message := Message{
				Sender: username,
				Text:   text,
			}

			messageJSON, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error while sending message")
				close(exit)
				break
			}
			_, err = conn.Write(append(messageJSON, '\n'))
			if err != nil {
				fmt.Println("Error while sending message")
				close(exit)
				break
			}
		}
	}()

	<-exit
	fmt.Println("Disconnected")
}
