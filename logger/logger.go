package logger

import (
	"fmt"
	"os"

	"github.com/pivotal-golang/lager"
)

// LogLevel is a typedef for a string.
type LogLevel string

// LogLevels are provided as constants below.
const (
	LogLevelInvalid LogLevel = ""
	LogLevelDebug   LogLevel = "debug"
	LogLevelInfo    LogLevel = "info"
	LogLevelError   LogLevel = "error"
	LogLevelFatal   LogLevel = "fatal"
)

// InitializeLogger is used to create a lager.Logger with a ReconfigurableSink
func InitializeLogger(minLogLevel LogLevel) (lager.Logger, *lager.ReconfigurableSink, error) {
	var minLagerLogLevel lager.LogLevel
	switch minLogLevel {
	case LogLevelDebug:
		minLagerLogLevel = lager.DEBUG
	case LogLevelInfo:
		minLagerLogLevel = lager.INFO
	case LogLevelError:
		minLagerLogLevel = lager.ERROR
	case LogLevelFatal:
		minLagerLogLevel = lager.FATAL
	default:
		return nil, nil, fmt.Errorf("unknown log level: %s", minLogLevel)
	}

	logger := lager.NewLogger("wundergo")

	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), minLagerLogLevel)
	logger.RegisterSink(sink)

	return logger, sink, nil
}
