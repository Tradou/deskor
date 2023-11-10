package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func New() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		Compress:   true,
		MaxBackups: 10,
		MaxAge:     1,
	})
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
