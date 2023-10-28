package screen

import (
	logger "deskor/log"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"net"
)

func Auth(w fyne.Window) *fyne.Container {
	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Username")

	addrWidget := widget.NewEntry()
	addrWidget.SetPlaceHolder("IP:PORT")

	passwordWidget := widget.NewPasswordEntry()
	passwordWidget.SetPlaceHolder("Password")

	submitWidget := widget.NewButton("Submit", func() {
		_, err := SendAuthenticationRequest(passwordWidget.Text, addrWidget.Text)
		if err != "" {
			dialog.NewError(fmt.Errorf(err), w).Show()
		} else {
			w.SetContent(Chat(usernameWidget.Text))
		}
	})

	return container.NewVBox(
		usernameWidget,
		addrWidget,
		passwordWidget,
		submitWidget,
	)
}

func SendAuthenticationRequest(password, addr string) (bool, string) {
	l, lErr := logger.NewFileLogger()
	if lErr != nil {
		log.Fatalf("Erreur while instantiating logger : %v", lErr)
	}
	defer l.Close()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false, fmt.Sprintf("Error while connectin to server", err)
	}

	_, err = conn.Write([]byte(password))
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
