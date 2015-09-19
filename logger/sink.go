package logger

import "io"

// A Sink represents a write destination for a Logger.
// Based off https://github.com/pivotal-golang/lager
type Sink interface {
	//Log to the sink.  Best effort -- no need to worry about errors.
	Log(level LogLevel, payload []byte)
}

type writerSink struct {
	writer      io.Writer
	minLogLevel LogLevel
}

// Log writes the provided log to any sinks whose level is greater
// than the provided LogLevel
func (sink writerSink) Log(level LogLevel, log []byte) {
	if level < sink.minLogLevel {
		return
	}

	sink.writer.Write(log)
	if len(log) > 0 {
		sink.writer.Write([]byte("\n"))
	}
}
