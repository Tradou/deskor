package screen

import (
	"deskor/log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Auth(w fyne.Window) *fyne.Container {
	clientLog := logger.Get()
	defer clientLog.Close()
	clientLog.Write("Auth screen render")

	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Username")

	addrWidget := widget.NewEntry()
	addrWidget.SetPlaceHolder("IP:PORT")

	submitWidget := widget.NewButton("Submit", func() {
		w.SetContent(Chat(usernameWidget.Text))
	})

	return container.NewVBox(
		usernameWidget,
		addrWidget,
		submitWidget,
	)
}
