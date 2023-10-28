package screen

import (
	"deskor/graphic"
	"deskor/notification"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Chat() *fyne.Container {
	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Enter your username")

	messageWidget := widget.NewEntry()
	messageWidget.SetPlaceHolder("Type your message and press Enter")

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

	return container.NewBorder(
		topContainer,
		messageWidget,
		nil,
		nil,
		chatScroller,
	)
}
