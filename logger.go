package wundergo

import (
	"log"
)

// Logger provides the ability to log a message.
type Logger interface {
	LogLine(message string)
}

// PrintlnLogger is an implementation of Logger.
type PrintlnLogger struct {
}

// LogLine writes the message to std out
func (l PrintlnLogger) LogLine(message string) {
	log.Println(message)
}
