package notification

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"os"
	"time"
)

var notificationEnabled = true

func Toggle() {
	notificationEnabled = !notificationEnabled
}

func IsEnabled() bool {
	return notificationEnabled
}

func Sound() {
	file, err := os.Open("assets/sound/notification.mp3")
	if err != nil {
		panic(err)
	}
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		panic(err)
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		panic(err)
	}
	speaker.Play(streamer)
}

func GetIcon() fyne.Resource {
	if IsEnabled() {
		return theme.VisibilityIcon()
	}
	return theme.VisibilityOffIcon()
}
