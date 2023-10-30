package screen

import (
	"deskor/log"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"net"
)

func Auth(w fyne.Window) *fyne.Container {
	clientLog := logger.Get()
	defer clientLog.Close()
	clientLog.Write("Auth screen render")

	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Username")

	addrWidget := widget.NewEntry()
	addrWidget.SetPlaceHolder("IP:PORT")

	passwordWidget := widget.NewPasswordEntry()
	passwordWidget.SetPlaceHolder("Password")

	submitWidget := widget.NewButton("Submit", func() {
		conn, err := net.Dial("tcp", addrWidget.Text)
		if err != nil {
			dialog.NewError(fmt.Errorf("error while connecting to server: %v", err), w).Show()
			return
		}

		success, connErr := SendAuthenticationRequest(conn, passwordWidget.Text)
		if !success {
			dialog.NewError(fmt.Errorf(connErr), w).Show()
			conn.Close()
			return
		}

		w.SetContent(Chat(usernameWidget.Text))
	})

	return container.NewVBox(
		usernameWidget,
		addrWidget,
		passwordWidget,
		submitWidget,
	)
}

func SendAuthenticationRequest(conn net.Conn, password string) (bool, string) {
	_, err := conn.Write([]byte(password))
	if err != nil {
		return false, "Error while sending request to server"
	}

	response := make([]byte, 64)

	_, err = conn.Read(response)

	if err != nil {
		return false, "Bad authentication"
	}

	return true, ""
}
