package typewriter

import (
	"os"
	"strings"
)

const (
	DEFAULT_LOG_LEVEL_ENV_NAME = "TYPEWRITER_DEFAULT_LOG_LEVEL"

	OFF_SEVERITY = iota
	FATAL_SEVERITY
	ERROR_SEVERITY
	WARNING_SEVERITY
	INFORMATION_SEVERITY
	DEBUG_SEVERITY
)

var debugEnabled = false
var informationEnabled = false
var warningEnabled = false
var errorEnabled = true
var fatalEnabled = true

// Initalize logging by setting severity
func InitLogConfig() {
	determineSeverityFromEnv()
}

// determines the default severity level from environment variable and activates different log levels
func determineSeverityFromEnv() {
	switch strings.ToUpper(strings.TrimSpace(os.Getenv(DEFAULT_LOG_LEVEL_ENV_NAME))) {
	case "DEBUG":
		determineSeverityByLevel(DEBUG_SEVERITY)
	case "INFO":
		determineSeverityByLevel(INFORMATION_SEVERITY)
	case "INFORMATION":
		determineSeverityByLevel(INFORMATION_SEVERITY)
	case "WARN":
		determineSeverityByLevel(WARNING_SEVERITY)
	case "WARNING":
		determineSeverityByLevel(WARNING_SEVERITY)
	case "ERROR":
		determineSeverityByLevel(ERROR_SEVERITY)
	case "FATAL":
		determineSeverityByLevel(FATAL_SEVERITY)
	default:
		determineSeverityByLevel(OFF_SEVERITY)
	}
}

// activates different log levels
func determineSeverityByLevel(severity int) {
	debugEnabled = DEBUG_SEVERITY <= severity
	informationEnabled = INFORMATION_SEVERITY <= severity
	warningEnabled = WARNING_SEVERITY <= severity
	errorEnabled = ERROR_SEVERITY <= severity
	fatalEnabled = FATAL_SEVERITY <= severity
}

func IsDebugEnabled() bool {
	return debugEnabled
}

func IsInformationEnabled() bool {
	return informationEnabled
}

func IsWarningEnabled() bool {
	return warningEnabled
}

func IsErrorEnabled() bool {
	return errorEnabled
}

func IsFatalEnabled() bool {
	return fatalEnabled
}
