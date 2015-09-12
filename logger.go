package wundergo

import (
	"log"
)

// Logger provides the ability to log a message.
type Logger interface {
	Println(message string)
}

// DefaultLogger is an implementation of Logger.
type DefaultLogger struct {
}

// Println writes the message to std out
func (l DefaultLogger) Println(message string) {
	log.Println(message)
}
