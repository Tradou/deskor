package screen

import (
	"crypto/tls"
	"deskor/chat"
	"deskor/encrypt"
	"deskor/graphic"
	"deskor/notification"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"os/signal"
)

func Chat(username string, conn *tls.Conn) *fyne.Container {
	encrypter := encrypt.EncrypterImpl{}
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
		chater := &chat.Client{}

		message, err := chater.EncodeMessage(sender, text)
		if err != nil {
			fmt.Print("Error while encoding message")
			close(exit)
		} else {
			err = chater.SendMessage(conn, message)
			if err != nil {
				fmt.Print("Error while sending message")
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
				fmt.Print("Error while receiving message")
				close(exit)
				break
			}

			var receivedMessage chat.Message
			decodedMessage, err := chater.DecodeMessage(message)

			if err == nil {
				receivedMessage = decodedMessage

				chatWidget.SetText(chatWidget.Text + "\n" + receivedMessage.Sender + ": " + receivedMessage.Text)
				chatScroller.ScrollToBottom()
				if notification.IsEnabled() && usernameWidget.Text != receivedMessage.Sender {
					notification.Sound()
				}
				connectedWidget.SetText(fmt.Sprintf("Connected people: %d", receivedMessage.Connected))
			} else {
				fmt.Println("Error while reading message", err)
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
