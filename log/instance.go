package logger

import (
	"fmt"
	"log"
	"sync"
)

type Instance struct {
	logger Logger
	once   sync.Once
}

var instance Instance

func (l *Instance) init() {
	var err error
	l.logger, err = NewFileLogger()
	fmt.Println("Attempt to instantiate new logger")
	if err != nil {
		log.Fatalf("Error while instantiated logger: %v", err)
	} else {
		l.logger.Write("Logger has just been instantiated")
	}
}

func New() {
	fmt.Println("Start new fn")
	instance.once.Do(instance.init)
}

func Get() Logger {
	fmt.Println("Get instance")
	return instance.logger
}
