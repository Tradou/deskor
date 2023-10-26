package main

import (
	"deskor/chat"
	"deskor/graphic"
	"deskor/notification"
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
	var notificationWidget *widget.Button
	notificationWidget = widget.NewButtonWithIcon("", notification.GetIcon(), func() {
		notification.Toggle()
		notificationWidget.SetIcon(notification.GetIcon())
	})

	topContainer := graphic.NewAdaptiveGridWithRatios([]float32{0.95, 0.05},
		usernameWidget,
		notificationWidget,
	)

	usernameWidget.SetPlaceHolder("Enter your username")
	messageWidget.SetPlaceHolder("Type your message and press Enter")

	messageWidget.OnSubmitted = func(text string) {
		sender := usernameWidget.Text
		chater := &chat.Client{}

		message, err := chater.EncodeMessage(sender, text)
		if err != nil {
			fmt.Println("Error while encoding message")
			close(exit)
		} else {
			err = chater.SendMessage(conn, message)
			if err != nil {
				fmt.Println("Error while sending message")
				close(exit)
			}
		}
		messageWidget.SetText("")
	}

	go func() {
		for {
			chater := &chat.Client{}

			message, err := chater.ReceiveMessage(conn)
			if err != nil {
				fmt.Println("Error while reading messages:", err)
				close(exit)
				break
			}

			var receivedMessage chat.Message

			if decodedMessage, err := chater.DecodeMessage(message); err == nil {
				receivedMessage = decodedMessage
			} else {
				fmt.Println("Error while decoding message:", err)
			}

			chatWidget.SetText(chatWidget.Text + "\n" + receivedMessage.Sender + ": " + receivedMessage.Text)
			chatScroller.ScrollToBottom()
			if notification.IsEnabled() && usernameWidget.Text != receivedMessage.Sender {
				notification.Sound()
			}
		}
	}()

	content := container.NewBorder(
		topContainer,
		messageWidget,
		nil,
		nil,
		chatScroller,
	)

	window.Resize(fyne.NewSize(600, 500))
	window.SetContent(content)
	window.ShowAndRun()
}
