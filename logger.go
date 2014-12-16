package wundergo

import (
	"log"
)

type Logger interface {
	LogLine(message string)
}

type PrintlnLogger struct {
}

func NewPrintlnLogger() *PrintlnLogger {
	return &PrintlnLogger{}
}

func (l PrintlnLogger) LogLine(message string) {
	log.Println(message)
}
