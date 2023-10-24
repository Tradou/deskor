package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"os/signal"
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

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	myApp := app.New()
	window := myApp.NewWindow("Deskor")
	usernameWidget := widget.NewEntry()
	messageWidget := widget.NewEntry()
	chatWidget := widget.NewLabel("Chat will appear here")

	chatScroller := container.NewVScroll(chatWidget)

	usernameWidget.SetPlaceHolder("Enter your username")
	messageWidget.SetPlaceHolder("Type your message and press Enter")

	messageWidget.OnSubmitted = func(text string) {
		message := Message{
			Sender: usernameWidget.Text,
			Text:   text,
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error while sending message")
			close(exit)
		} else {
			_, err = conn.Write(append(messageJSON, '\n'))
			if err != nil {
				fmt.Println("Error while sending message")
				close(exit)
			}
		}
		messageWidget.SetText("")
	}

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
				chatWidget.SetText(chatWidget.Text + "\n" + receivedMessage.Sender + ": " + receivedMessage.Text)
			} else {
				// TODO: idk in which case it's possible, and how to handle it
				chatWidget.SetText(chatWidget.Text + "\n" + receivedMessage.Sender + ": " + receivedMessage.Text)
			}
			chatScroller.ScrollToBottom()
		}
	}()

	content := container.NewBorder(
		usernameWidget,
		messageWidget,
		nil,
		nil,
		chatScroller,
	)

	window.Resize(fyne.NewSize(600, 500))
	window.SetContent(content)
	window.ShowAndRun()
}
