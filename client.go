package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

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

			fmt.Printf("\033[1A\033[K")
			fmt.Print(message)
		}
	}()

	go func() {
		for {
			text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			message := username + ": " + text
			_, err := conn.Write([]byte(message))
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
