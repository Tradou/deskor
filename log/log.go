package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger interface {
	Open() error
	Close() error
	Exists() bool
	Create() error
	Write(message string)
}

type FileLogger struct {
	logFile  *os.File
	filePath string
}

func NewFileLogger() (*FileLogger, error) {
	today := time.Now().Format(time.DateOnly)
	logFilePath := fmt.Sprintf("./logs/%s.txt", today)
	logger := &FileLogger{filePath: logFilePath}

	if err := logger.Open(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (f *FileLogger) Exists() bool {
	_, err := os.Stat(f.filePath)
	return !os.IsNotExist(err)
}

func (f *FileLogger) Create() error {
	file, err := os.Create(f.filePath)
	if err != nil {
		return err
	}
	f.logFile = file
	return nil
}

func (f *FileLogger) Open() error {
	if f.logFile != nil {
		f.logFile.Close()
	}
	if f.Exists() {
		file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		f.logFile = file
		return nil
	} else {
		err := f.Create()
		return err
	}
}

func (f *FileLogger) Close() error {
	if f.logFile != nil {
		err := f.logFile.Close()
		f.logFile = nil
		return err
	}
	return nil
}

func (f *FileLogger) OpenOrCreate() error {
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		file, err := os.Create(f.filePath)
		if err != nil {
			return err
		}
		f.logFile = file
	} else {
		file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		f.logFile = file
	}

	return nil
}

func (f *FileLogger) Write(message string) {
	if f.logFile == nil {
		log.Fatalf("Can't open log file : %s", f.filePath)
		return
	}

	logMessage := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.DateTime), message)
	log.Print(logMessage)
	f.logFile.WriteString(logMessage)
}
