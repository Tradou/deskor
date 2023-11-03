package main

import (
	"deskor/graphic/screen"
	logger "deskor/log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	logger.New()
	clientLog := logger.Get()
	defer clientLog.Close()
	clientLog.Write("Start running client")

	myApp := app.New()
	window := myApp.NewWindow("Deskor")

	window.Resize(fyne.NewSize(600, 500))
	window.SetContent(screen.Auth(myApp, window))
	window.ShowAndRun()
}
