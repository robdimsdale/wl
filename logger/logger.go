package logger

import (
	"fmt"
	"io"
	"os"
)

// Logger supports writing messages and arbitrary data
// at different log levels.
type Logger interface {
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Error(string, error, ...interface{})
}

type logger struct {
	sinks []Sink
}

// NewLogger returns a logger writing to stdout whose level is controlled
// by the provided minLogLevel
func NewLogger(minLogLevel LogLevel) Logger {
	sink := writerSink{
		writer:      os.Stdout,
		minLogLevel: minLogLevel,
	}
	return &logger{
		sinks: []Sink{sink},
	}
}

// NewTestLogger returns a logger writing to the provided writer.
// Its level is fixed at DEBUG
// It is primarily used in testing to write to e.g. the GinkgoWriter
func NewTestLogger(writer io.Writer) Logger {
	sink := writerSink{
		writer:      writer,
		minLogLevel: DEBUG,
	}
	return &logger{
		sinks: []Sink{sink},
	}
}

// Info logs the message and any provided data at Info level.
func (l logger) Info(message string, data ...interface{}) {
	for _, sink := range l.sinks {
		sink.Log(INFO, l.toByteArray(message, data...))
	}
}

func (l logger) toByteArray(message string, data ...interface{}) []byte {
	for _, d := range data {
		message = fmt.Sprintf("%s %v", message, d)
	}
	return []byte(message)
}

// Info logs the message and any provided data at Debug level.
func (l logger) Debug(message string, data ...interface{}) {
	for _, sink := range l.sinks {
		sink.Log(DEBUG, l.toByteArray(message, data...))
	}
}

// Info logs the message and any provided data at Error level.
func (l logger) Error(message string, err error, data ...interface{}) {
	combined := []interface{}{err}
	combined = append(combined, data...)

	for _, sink := range l.sinks {
		sink.Log(ERROR, l.toByteArray(message, combined...))
	}
}
