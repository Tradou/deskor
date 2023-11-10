package main

import (
	"deskor/graphic/screen"
	logger "deskor/log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
)

func main() {
	logger.New()
	log.Print("Start running client")

	myApp := app.New()
	window := myApp.NewWindow("Deskor")

	window.Resize(fyne.NewSize(600, 500))
	window.SetContent(screen.Auth(myApp, window))
	window.ShowAndRun()
}
