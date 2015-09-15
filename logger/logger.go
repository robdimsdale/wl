package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Error(string, error, ...interface{})
}

type logger struct {
	sinks []Sink
}

func NewLogger(minLogLevel LogLevel) Logger {
	sink := writerSink{
		writer:      os.Stdout,
		minLogLevel: minLogLevel,
	}
	return &logger{
		sinks: []Sink{sink},
	}
}

func NewTestLogger(writer io.Writer) Logger {
	sink := writerSink{
		writer:      writer,
		minLogLevel: DEBUG,
	}
	return &logger{
		sinks: []Sink{sink},
	}
}

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

func (l logger) Debug(message string, data ...interface{}) {
	for _, sink := range l.sinks {
		sink.Log(DEBUG, l.toByteArray(message, data...))
	}
}

func (l logger) Error(message string, err error, data ...interface{}) {
	combined := []interface{}{err}
	combined = append(combined, data...)

	for _, sink := range l.sinks {
		sink.Log(ERROR, l.toByteArray(message, combined...))
	}
}
