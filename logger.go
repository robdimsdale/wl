package wundergo

import (
	"log"
)

// Logger provides the ability to log a message.
type Logger interface {
	Println(message string)
}

// PrintlnLogger is an implementation of Logger.
type defaultLogger struct {
}

// Println writes the message to std out
func (l defaultLogger) Println(message string) {
	log.Println(message)
}
