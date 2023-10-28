package main

import (
	"deskor/graphic/screen"
	"deskor/log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
)

func main() {
	l, lErr := logger.NewFileLogger()
	if lErr != nil {
		log.Fatalf("Erreur while instantiating logger : %v", lErr)
	}
	defer l.Close()
	l.Write("Start app")

	app := app.New()
	window := app.NewWindow("Deskor")

	window.Resize(fyne.NewSize(600, 500))
	window.SetContent(screen.Auth(window))
	window.ShowAndRun()
}
