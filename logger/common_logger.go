package logger

import (
	"os"
	"strings"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/constants"
)

const (
	DEFAULT_LOG_LEVEL_ENV_NAME = "TYPEWRITER_DEFAULT_LOG_LEVEL"
)

type CommonLogger struct {
	debugEnabled       bool
	informationEnabled bool
	warningEnabled     bool
	errorEnabled       bool
	fatalEnabled       bool
	appender           appender.Appender
}

func CreateCommonLogger(appender appender.Appender) CommonLogger {
	result := CommonLogger{appender: appender}
	determineSeverityFromEnv(&result)
	return result
}

// determines the default severity level from environment variable and activates different log levels
func determineSeverityFromEnv(l *CommonLogger) {
	switch strings.ToUpper(strings.TrimSpace(os.Getenv(DEFAULT_LOG_LEVEL_ENV_NAME))) {
	case "DEBUG":
		determineSeverityByLevel(l, constants.DEBUG_SEVERITY)
	case "INFO":
		determineSeverityByLevel(l, constants.INFORMATION_SEVERITY)
	case "INFORMATION":
		determineSeverityByLevel(l, constants.INFORMATION_SEVERITY)
	case "WARN":
		determineSeverityByLevel(l, constants.WARNING_SEVERITY)
	case "WARNING":
		determineSeverityByLevel(l, constants.WARNING_SEVERITY)
	case "ERROR":
		determineSeverityByLevel(l, constants.ERROR_SEVERITY)
	case "FATAL":
		determineSeverityByLevel(l, constants.FATAL_SEVERITY)
	default:
		determineSeverityByLevel(l, constants.OFF_SEVERITY)
	}
}

// activates different log levels
func determineSeverityByLevel(l *CommonLogger, severity int) {
	l.debugEnabled = constants.DEBUG_SEVERITY <= severity
	l.informationEnabled = constants.INFORMATION_SEVERITY <= severity
	l.warningEnabled = constants.WARNING_SEVERITY <= severity
	l.errorEnabled = constants.ERROR_SEVERITY <= severity
	l.fatalEnabled = constants.FATAL_SEVERITY <= severity
}

func (l CommonLogger) IsDebugEnabled() bool {
	return l.debugEnabled
}

func (l CommonLogger) IsInformationEnabled() bool {
	return l.informationEnabled
}

func (l CommonLogger) IsWarningEnabled() bool {
	return l.warningEnabled
}

func (l CommonLogger) IsErrorEnabled() bool {
	return l.errorEnabled
}

func (l CommonLogger) IsFatalEnabled() bool {
	return l.fatalEnabled
}
