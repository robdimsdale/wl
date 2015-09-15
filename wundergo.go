/*
Package wundergo provides a client to the Wunderlist API.

For more information on the API, see https://developer.wunderlist.com/documentation
*/
package wundergo

const (
	// APIURL is the default URL for Wunderlist API.
	APIURL = "https://a.wunderlist.com/api/v1"
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
