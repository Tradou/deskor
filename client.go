package main

import (
	"deskor/graphic/screen"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	fmt.Print("Start running client")

	myApp := app.New()
	window := myApp.NewWindow("Deskor")

	window.Resize(fyne.NewSize(600, 500))
	window.SetContent(screen.Auth(myApp, window))
	window.ShowAndRun()
}
