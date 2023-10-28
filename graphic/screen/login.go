package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Auth(w fyne.Window) *fyne.Container {
	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Username")

	addrWidget := widget.NewEntry()
	addrWidget.SetPlaceHolder("IP:PORT")

	passwordWidget := widget.NewPasswordEntry()
	passwordWidget.SetPlaceHolder("Password")

	submitWidget := widget.NewButton("Submit", func() {
		w.SetContent(Chat())
	})

	return container.NewVBox(
		usernameWidget,
		addrWidget,
		passwordWidget,
		submitWidget,
	)
}
