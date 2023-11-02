package screen

import (
	"crypto/tls"
	logger "deskor/log"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Auth(w fyne.Window) *fyne.Container {
	logger.New()
	clientLog := logger.Get()
	defer clientLog.Close()
	clientLog.Write("Auth screen render")

	usernameWidget := widget.NewEntry()
	usernameWidget.SetPlaceHolder("Username")

	addrWidget := widget.NewEntry()
	addrWidget.SetPlaceHolder("IP:PORT")

	submitWidget := widget.NewButton("Submit", func() {
		cert, err := tls.LoadX509KeyPair("./cert/client.pem", "./cert/client.key")
		if err != nil {
			dialog.NewError(fmt.Errorf("error while loading keys: %v", err), w).Show()
			return
		}
		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		conn, err := tls.Dial("tcp", addrWidget.Text, &config)
		if err != nil {
			dialog.NewError(fmt.Errorf("error while connecting to server: %v", err), w).Show()
			conn.Close()
			return
		}

		w.SetContent(Chat(usernameWidget.Text, conn))
	})

	return container.NewVBox(
		usernameWidget,
		addrWidget,
		submitWidget,
	)
}
