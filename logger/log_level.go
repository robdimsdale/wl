package logger

import "fmt"

// LogLevel represents the levels at which logs are written and rendered.
type LogLevel int

// DEBUG < INFO < ERROR < FATAL
const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	FATAL
)

// LogLevelFromString returns a LogLevel for the provided logLevel string.
// It will panic if an unknown logLevel is provided.
func LogLevelFromString(logLevel string) LogLevel {
	switch logLevel {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		panic(fmt.Errorf("unknown log level: %s", logLevel))
	}

}
