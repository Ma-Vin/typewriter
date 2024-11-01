package logger

import (
	"fmt"
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
	appender           *appender.Appender
}

var mockPanicAndExitAtCommonLogger = false
var panicMockActivated = false
var exitMockAcitvated = false

func CreateCommonLogger(appender *appender.Appender) CommonLogger {
	result := CommonLogger{appender: appender}
	determineSeverityFromEnv(&result)
	return result
}

// determines the default severity level from environment variable and activates different log levels
func determineSeverityFromEnv(l *CommonLogger) {
	determineSeverity(strings.ToUpper(strings.TrimSpace(os.Getenv(DEFAULT_LOG_LEVEL_ENV_NAME))), l)
}

// determines the severity level from a given variable and activates different log levels
func determineSeverity(severityLevel string, l *CommonLogger) {
	switch severityLevel {
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

func (l CommonLogger) Debug(args ...any) {
	if l.debugEnabled {
		(*l.appender).Write(constants.DEBUG_SEVERITY, fmt.Sprint(args...))
	}
}

func (l CommonLogger) DebugWithCorrelation(correlationId string, args ...any) {
	if l.debugEnabled {
		(*l.appender).WriteWithCorrelation(constants.DEBUG_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

func (l CommonLogger) DebugCustom(customValues map[string]any, args ...any) {
	if l.debugEnabled {
		(*l.appender).WriteCustom(constants.DEBUG_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

func (l CommonLogger) Debugf(format string, args ...any) {
	if l.debugEnabled {
		(*l.appender).Write(constants.DEBUG_SEVERITY, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) DebugWithCorrelationf(correlationId string, format string, args ...any) {
	if l.debugEnabled {
		(*l.appender).WriteWithCorrelation(constants.DEBUG_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) DebugCustomf(customValues map[string]any, format string, args ...any) {
	if l.debugEnabled {
		(*l.appender).WriteCustom(constants.DEBUG_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

func (l CommonLogger) Information(args ...any) {
	if l.informationEnabled {
		(*l.appender).Write(constants.INFORMATION_SEVERITY, fmt.Sprint(args...))
	}
}

func (l CommonLogger) InformationWithCorrelation(correlationId string, args ...any) {
	if l.informationEnabled {
		(*l.appender).WriteWithCorrelation(constants.INFORMATION_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

func (l CommonLogger) InformationCustom(customValues map[string]any, args ...any) {
	if l.informationEnabled {
		(*l.appender).WriteCustom(constants.INFORMATION_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

func (l CommonLogger) Informationf(format string, args ...any) {
	if l.informationEnabled {
		(*l.appender).Write(constants.INFORMATION_SEVERITY, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) InformationWithCorrelationf(correlationId string, format string, args ...any) {
	if l.informationEnabled {
		(*l.appender).WriteWithCorrelation(constants.INFORMATION_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) InformationCustomf(customValues map[string]any, format string, args ...any) {
	if l.informationEnabled {
		(*l.appender).WriteCustom(constants.INFORMATION_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

func (l CommonLogger) Warning(args ...any) {
	if l.warningEnabled {
		(*l.appender).Write(constants.WARNING_SEVERITY, fmt.Sprint(args...))
	}
}

func (l CommonLogger) WarningWithCorrelation(correlationId string, args ...any) {
	if l.warningEnabled {
		(*l.appender).WriteWithCorrelation(constants.WARNING_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

func (l CommonLogger) WarningCustom(customValues map[string]any, args ...any) {
	if l.warningEnabled {
		(*l.appender).WriteCustom(constants.WARNING_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

func (l CommonLogger) Warningf(format string, args ...any) {
	if l.warningEnabled {
		(*l.appender).Write(constants.WARNING_SEVERITY, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) WarningWithCorrelationf(correlationId string, format string, args ...any) {
	if l.warningEnabled {
		(*l.appender).WriteWithCorrelation(constants.WARNING_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) WarningCustomf(customValues map[string]any, format string, args ...any) {
	if l.warningEnabled {
		(*l.appender).WriteCustom(constants.WARNING_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

func (l CommonLogger) WarningWithPanic(args ...any) {
	l.Warning(args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) WarningWithCorrelationAndPanic(correlationId string, args ...any) {
	l.WarningWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) WarningCustomWithPanic(customValues map[string]any, args ...any) {
	l.WarningCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) WarningWithPanicf(format string, args ...any) {
	l.Warningf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.WarningWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) WarningCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.WarningCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) Error(args ...any) {
	if l.errorEnabled {
		(*l.appender).Write(constants.ERROR_SEVERITY, fmt.Sprint(args...))
	}
}

func (l CommonLogger) ErrorWithCorrelation(correlationId string, args ...any) {
	if l.errorEnabled {
		(*l.appender).WriteWithCorrelation(constants.ERROR_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

func (l CommonLogger) ErrorCustom(customValues map[string]any, args ...any) {
	if l.errorEnabled {
		(*l.appender).WriteCustom(constants.ERROR_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

func (l CommonLogger) Errorf(format string, args ...any) {
	if l.errorEnabled {
		(*l.appender).Write(constants.ERROR_SEVERITY, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) ErrorWithCorrelationf(correlationId string, format string, args ...any) {
	if l.errorEnabled {
		(*l.appender).WriteWithCorrelation(constants.ERROR_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) ErrorCustomf(customValues map[string]any, format string, args ...any) {
	if l.errorEnabled {
		(*l.appender).WriteCustom(constants.ERROR_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

func (l CommonLogger) ErrorWithPanic(args ...any) {
	l.Error(args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) ErrorWithCorrelationAndPanic(correlationId string, args ...any) {
	l.ErrorWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) ErrorCustomWithPanic(customValues map[string]any, args ...any) {
	l.ErrorCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) ErrorWithPanicf(format string, args ...any) {
	l.Errorf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.ErrorWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.ErrorCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) Fatal(args ...any) {
	if l.fatalEnabled {
		(*l.appender).Write(constants.FATAL_SEVERITY, fmt.Sprint(args...))
	}
}

func (l CommonLogger) FatalWithCorrelation(correlationId string, args ...any) {
	if l.fatalEnabled {
		(*l.appender).WriteWithCorrelation(constants.FATAL_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

func (l CommonLogger) FatalCustom(customValues map[string]any, args ...any) {
	if l.fatalEnabled {
		(*l.appender).WriteCustom(constants.FATAL_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

func (l CommonLogger) Fatalf(format string, args ...any) {
	if l.fatalEnabled {
		(*l.appender).Write(constants.FATAL_SEVERITY, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) FatalWithCorrelationf(correlationId string, format string, args ...any) {
	if l.fatalEnabled {
		(*l.appender).WriteWithCorrelation(constants.FATAL_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

func (l CommonLogger) FatalCustomf(customValues map[string]any, format string, args ...any) {
	if l.fatalEnabled {
		(*l.appender).WriteCustom(constants.FATAL_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

func (l CommonLogger) FatalWithPanic(args ...any) {
	l.Fatal(args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) FatalWithCorrelationAndPanic(correlationId string, args ...any) {
	l.FatalWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) FatalCustomWithPanic(customValues map[string]any, args ...any) {
	l.FatalCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) FatalWithPanicf(format string, args ...any) {
	l.Fatalf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.FatalWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) FatalCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.FatalCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

func (l CommonLogger) FatalWithExit(args ...any) {
	l.Fatal(args...)
	exitOrMock(1)
}

func (l CommonLogger) FatalWithCorrelationAndExit(correlationId string, args ...any) {
	l.FatalWithCorrelation(correlationId, args...)
	exitOrMock(1)
}

func (l CommonLogger) FatalCustomWithExit(customValues map[string]any, args ...any) {
	l.FatalCustom(customValues, args...)
	exitOrMock(1)
}

func (l CommonLogger) FatalWithExitf(format string, args ...any) {
	l.Fatalf(format, args...)
	exitOrMock(1)
}

func (l CommonLogger) FatalWithCorrelationAndExitf(correlationId string, format string, args ...any) {
	l.FatalWithCorrelationf(correlationId, format, args...)
	exitOrMock(1)
}

func (l CommonLogger) FatalCustomWithExitf(customValues map[string]any, format string, args ...any) {
	l.FatalCustomf(customValues, format, args...)
	exitOrMock(1)
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

func panicOrMock(message string) {
	if mockPanicAndExitAtCommonLogger {
		panicMockActivated = true
		return
	}
	panic(message)
}

func exitOrMock(code int) {
	if mockPanicAndExitAtCommonLogger {
		exitMockAcitvated = true
		return
	}
	os.Exit(code)
}
