package screen

import (
	"crypto/tls"
	"deskor/chat"
	"deskor/graphic"
	"deskor/notification"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"os/signal"
)

func Chat(username string, conn *tls.Conn, app fyne.App) *fyne.Container {
	client := &chat.Client{
		Conn: conn,
	}

	usernameWidget := widget.NewEntry()
	usernameWidget.SetText(username)
	usernameWidget.Disable()

	messageWidget := widget.NewEntry()
	messageWidget.SetPlaceHolder("Type your message and press Enter")

	chatWidget := widget.NewLabel("Chat will appear here")

	connectedWidget := widget.NewEntry()
	connectedWidget.SetText("Loading connected people")
	connectedWidget.Disable()

	chatScroller := container.NewVScroll(chatWidget)
	var notificationWidget *widget.Button
	notificationWidget = widget.NewButtonWithIcon("", notification.GetIcon(), func() {
		notification.Toggle()
		notificationWidget.SetIcon(notification.GetIcon())
	})

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	messageWidget.OnSubmitted = func(text string) {
		sender := usernameWidget.Text

		message, err := client.EncodeMessage(sender, text)
		if err != nil {
			log.Printf("Error while encoding message: %s", err)
			close(exit)
		} else {
			err = client.SendMessage(message)
			if err != nil {
				log.Printf("Error while sending message: %s", err)
				close(exit)
			}
		}
		messageWidget.SetText("")
	}

	go func() {
		for {
			message, err := client.ReceiveMessage()
			if err != nil {
				log.Printf("Error while receiving message: %s", err)
				close(exit)
				break
			}

			var receivedMessage chat.Message
			decodedMessage, err := client.DecodeMessage(message)

			if err == nil {
				receivedMessage = decodedMessage

				chatWidget.SetText(chatWidget.Text + "\n" + receivedMessage.Sender + ": " + receivedMessage.Text)
				chatScroller.ScrollToBottom()
				if notification.IsEnabled() && usernameWidget.Text != receivedMessage.Sender {
					notification.Sound()
					notification.Popup(app, receivedMessage.Sender)
				}
				connectedWidget.SetText(fmt.Sprintf("Connected people: %d", receivedMessage.Connected))
			} else {
				log.Println("Error while reading message", err)
			}
		}
	}()

	topContainer := graphic.NewAdaptiveGridWithRatios([]float32{0.60, 0.35, 0.05},
		usernameWidget,
		connectedWidget,
		notificationWidget,
	)

	return container.NewBorder(
		topContainer,
		messageWidget,
		nil,
		nil,
		chatScroller,
	)
}
