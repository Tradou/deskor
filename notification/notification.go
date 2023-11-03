package notification

import (
	"bytes"
	"deskor/assets/bundle"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"io"
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
	soundBytes := bundle.ResourceAssetsSoundNotificationMp3.StaticContent

	soundReader := bytes.NewReader(soundBytes)

	readCloser := io.NopCloser(soundReader)

	streamer, format, err := mp3.Decode(readCloser)
	if err != nil {
		panic(err)
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		panic(err)
	}
	speaker.Play(streamer)
}

func Popup(app fyne.App, username string) {
	app.SendNotification(fyne.NewNotification("New message", fmt.Sprintf("%s just send a new message", username)))
}

func GetIcon() fyne.Resource {
	if IsEnabled() {
		return theme.VisibilityIcon()
	}
	return theme.VisibilityOffIcon()
}
