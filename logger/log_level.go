package logger

import "fmt"

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	FATAL
)

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
