package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Auth() *fyne.Container {
	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Username")

	ipWidget := widget.NewEntry()
	ipWidget.SetPlaceHolder("IP:PORT")

	passwordWidget := widget.NewPasswordEntry()
	passwordWidget.SetPlaceHolder("Password")

	submitWidget := widget.NewButton("Submit", func() {
	})

	return container.NewVBox(
		usernameWidget,
		ipWidget,
		passwordWidget,
		submitWidget,
	)
}
