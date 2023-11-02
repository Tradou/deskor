package logger

import (
	"log"
	"sync"
)

type Instance struct {
	logger *FileLogger
	once   sync.Once
}

var instance Instance

func (l *Instance) init() {
	var err error
	l.logger, err = NewFileLogger()

	if err != nil {
		log.Fatalf("Error while instantiated logger: %v", err)
	} else {
		l.logger.Write("Logger has just been instantiated")
	}
}

func New() {
	instance.once.Do(instance.init)
}

func Get() *FileLogger {
	return instance.logger
}
