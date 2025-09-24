package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

// A general logger which delegates messages directly to the appender if log level is enabled.
type GeneralLogger struct {
	debugEnabled       bool
	informationEnabled bool
	warningEnabled     bool
	errorEnabled       bool
	fatalEnabled       bool
	isCallerToSet      bool
	appender           *appender.Appender
	correlationIdKey   string
}

var mockPanicAndExitAtGeneralLogger = false
var panicMockActivated = false
var exitMockActivated = false

// Creates a general logger which delegates messages to the given appender if the log level is enabled by given severity
func CreateGeneralLoggerFromConfig(generalLoggerConfig config.GeneralLoggerConfig, appender *appender.Appender) GeneralLogger {
	result := GeneralLogger{
		appender:         appender,
		isCallerToSet:    generalLoggerConfig.IsCallerToSet,
		correlationIdKey: generalLoggerConfig.Common.CorrelationIdKey,
	}
	determineSeverityByLevel(&result, generalLoggerConfig.Severity)
	return result
}

// activates different log levels
func determineSeverityByLevel(l *GeneralLogger, severity int) {
	l.debugEnabled = common.DEBUG_SEVERITY <= severity
	l.informationEnabled = common.INFORMATION_SEVERITY <= severity
	l.warningEnabled = common.WARNING_SEVERITY <= severity
	l.errorEnabled = common.ERROR_SEVERITY <= severity
	l.fatalEnabled = common.FATAL_SEVERITY <= severity
}

// Logs a message if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) Debug(args ...any) {
	if l.debugEnabled {
		l.write(common.DEBUG_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) DebugWithCorrelation(correlationId string, args ...any) {
	if l.debugEnabled {
		l.writeWithCorrelation(common.DEBUG_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) DebugCustom(customValues map[string]any, args ...any) {
	if l.debugEnabled {
		l.writeCustom(common.DEBUG_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message together with a correlation id from context if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) DebugCtx(context context.Context, args ...any) {
	if l.debugEnabled {
		l.writeWithCtx(common.DEBUG_SEVERITY, context, fmt.Sprint(args...))
	}
}

// Logs a message derived from format if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) Debugf(format string, args ...any) {
	if l.debugEnabled {
		l.write(common.DEBUG_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) DebugWithCorrelationf(correlationId string, format string, args ...any) {
	if l.debugEnabled {
		l.writeWithCorrelation(common.DEBUG_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) DebugCustomf(customValues map[string]any, format string, args ...any) {
	if l.debugEnabled {
		l.writeCustom(common.DEBUG_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message derived from format together with a correlation id from context if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) DebugCtxf(context context.Context, format string, args ...any) {
	if l.debugEnabled {
		l.writeWithCtx(common.DEBUG_SEVERITY, context, fmt.Sprintf(format, args...))
	}
}

// Logs a message if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) Information(args ...any) {
	if l.informationEnabled {
		l.write(common.INFORMATION_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) InformationWithCorrelation(correlationId string, args ...any) {
	if l.informationEnabled {
		l.writeWithCorrelation(common.INFORMATION_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) InformationCustom(customValues map[string]any, args ...any) {
	if l.informationEnabled {
		l.writeCustom(common.INFORMATION_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message together with a correlation id from context if information level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) InformationCtx(context context.Context, args ...any) {
	if l.informationEnabled {
		l.writeWithCtx(common.INFORMATION_SEVERITY, context, fmt.Sprint(args...))
	}
}

// Logs a message derived from format if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) Informationf(format string, args ...any) {
	if l.informationEnabled {
		l.write(common.INFORMATION_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) InformationWithCorrelationf(correlationId string, format string, args ...any) {
	if l.informationEnabled {
		l.writeWithCorrelation(common.INFORMATION_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) InformationCustomf(customValues map[string]any, format string, args ...any) {
	if l.informationEnabled {
		l.writeCustom(common.INFORMATION_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message derived from format together with a correlation id from context if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) InformationCtxf(context context.Context, format string, args ...any) {
	if l.informationEnabled {
		l.writeWithCtx(common.INFORMATION_SEVERITY, context, fmt.Sprintf(format, args...))
	}
}

// Logs a message if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) Warning(args ...any) {
	if l.warningEnabled {
		l.write(common.WARNING_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningWithCorrelation(correlationId string, args ...any) {
	if l.warningEnabled {
		l.writeWithCorrelation(common.WARNING_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningCustom(customValues map[string]any, args ...any) {
	if l.warningEnabled {
		l.writeCustom(common.WARNING_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message together with a correlation id from context if warning level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningCtx(context context.Context, args ...any) {
	if l.warningEnabled {
		l.writeWithCtx(common.WARNING_SEVERITY, context, fmt.Sprint(args...))
	}
}

// Logs a message derived from format if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) Warningf(format string, args ...any) {
	if l.warningEnabled {
		l.write(common.WARNING_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningWithCorrelationf(correlationId string, format string, args ...any) {
	if l.warningEnabled {
		l.writeWithCorrelation(common.WARNING_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningCustomf(customValues map[string]any, format string, args ...any) {
	if l.warningEnabled {
		l.writeCustom(common.WARNING_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message derived from format together with a correlation id from context if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningCtxf(context context.Context, format string, args ...any) {
	if l.warningEnabled {
		l.writeWithCtx(common.WARNING_SEVERITY, context, fmt.Sprintf(format, args...))
	}
}

// Logs a message if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningWithPanic(args ...any) {
	l.Warning(args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningWithCorrelationAndPanic(correlationId string, args ...any) {
	l.WarningWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningCustomWithPanic(customValues map[string]any, args ...any) {
	l.WarningCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id from context if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) WarningCtxWithPanic(context context.Context, args ...any) {
	l.WarningCtx(context, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningWithPanicf(format string, args ...any) {
	l.Warningf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.WarningWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.WarningCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id from context if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) WarningCtxWithPanicf(context context.Context, format string, args ...any) {
	l.WarningCtxf(context, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) Error(args ...any) {
	if l.errorEnabled {
		l.write(common.ERROR_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorWithCorrelation(correlationId string, args ...any) {
	if l.errorEnabled {
		l.writeWithCorrelation(common.ERROR_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorCustom(customValues map[string]any, args ...any) {
	if l.errorEnabled {
		l.writeCustom(common.ERROR_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message together with a correlation id from context if error level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorCtx(context context.Context, args ...any) {
	if l.errorEnabled {
		l.writeWithCtx(common.ERROR_SEVERITY, context, fmt.Sprint(args...))
	}
}

// Logs a message derived from format if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) Errorf(format string, args ...any) {
	if l.errorEnabled {
		l.write(common.ERROR_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorWithCorrelationf(correlationId string, format string, args ...any) {
	if l.errorEnabled {
		l.writeWithCorrelation(common.ERROR_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorCustomf(customValues map[string]any, format string, args ...any) {
	if l.errorEnabled {
		l.writeCustom(common.ERROR_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message derived from format together with a correlation id from context if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorCtxf(context context.Context, format string, args ...any) {
	if l.errorEnabled {
		l.writeWithCtx(common.ERROR_SEVERITY, context, fmt.Sprintf(format, args...))
	}
}

// Logs a message if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorWithPanic(args ...any) {
	l.Error(args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorWithCorrelationAndPanic(correlationId string, args ...any) {
	l.ErrorWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorCustomWithPanic(customValues map[string]any, args ...any) {
	l.ErrorCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id from context if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) ErrorCtxWithPanic(context context.Context, args ...any) {
	l.ErrorCtx(context, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorWithPanicf(format string, args ...any) {
	l.Errorf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.ErrorWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.ErrorCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id from context if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) ErrorCtxWithPanicf(context context.Context, format string, args ...any) {
	l.ErrorCtxf(context, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) Fatal(args ...any) {
	if l.fatalEnabled {
		l.write(common.FATAL_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalWithCorrelation(correlationId string, args ...any) {
	if l.fatalEnabled {
		l.writeWithCorrelation(common.FATAL_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalCustom(customValues map[string]any, args ...any) {
	if l.fatalEnabled {
		l.writeCustom(common.FATAL_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message together with a correlation id from context if fatal level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalCtx(context context.Context, args ...any) {
	if l.fatalEnabled {
		l.writeWithCtx(common.FATAL_SEVERITY, context, fmt.Sprint(args...))
	}
}

// Logs a message derived from format if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) Fatalf(format string, args ...any) {
	if l.fatalEnabled {
		l.write(common.FATAL_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalWithCorrelationf(correlationId string, format string, args ...any) {
	if l.fatalEnabled {
		l.writeWithCorrelation(common.FATAL_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalCustomf(customValues map[string]any, format string, args ...any) {
	if l.fatalEnabled {
		l.writeCustom(common.FATAL_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message derived from format together with a correlation id from context if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalCtxf(context context.Context, format string, args ...any) {
	if l.fatalEnabled {
		l.writeWithCtx(common.FATAL_SEVERITY, context, fmt.Sprintf(format, args...))
	}
}

// Logs a message if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalWithPanic(args ...any) {
	l.Fatal(args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalWithCorrelationAndPanic(correlationId string, args ...any) {
	l.FatalWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalCustomWithPanic(customValues map[string]any, args ...any) {
	l.FatalCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id from context if fatal level is enabled  and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalCtxWithPanic(context context.Context, args ...any) {
	l.FatalCtx(context, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalWithPanicf(format string, args ...any) {
	l.Fatalf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.FatalWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.FatalCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id from context if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalCtxWithPanicf(context context.Context, format string, args ...any) {
	l.FatalCtxf(context, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalWithExit(args ...any) {
	l.Fatal(args...)
	l.exitOrMock(1)
}

// Logs a message together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalWithCorrelationAndExit(correlationId string, args ...any) {
	l.FatalWithCorrelation(correlationId, args...)
	l.exitOrMock(1)
}

// Logs a message together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalCustomWithExit(customValues map[string]any, args ...any) {
	l.FatalCustom(customValues, args...)
	l.exitOrMock(1)
}

// Logs a message together with a correlation id from context if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l GeneralLogger) FatalCtxWithExit(context context.Context, args ...any) {
	l.FatalCtx(context, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalWithExitf(format string, args ...any) {
	l.Fatalf(format, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalWithCorrelationAndExitf(correlationId string, format string, args ...any) {
	l.FatalWithCorrelationf(correlationId, format, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalCustomWithExitf(customValues map[string]any, format string, args ...any) {
	l.FatalCustomf(customValues, format, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format together with a correlation id from context if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l GeneralLogger) FatalCtxWithExitf(context context.Context, format string, args ...any) {
	l.FatalCtxf(context, format, args...)
	l.exitOrMock(1)
}

// Indicator whether debug level is enabled or not
func (l GeneralLogger) IsDebugEnabled() bool {
	return l.debugEnabled
}

// Indicator whether information level is enabled or not
func (l GeneralLogger) IsInformationEnabled() bool {
	return l.informationEnabled
}

// Indicator whether warning level is enabled or not
func (l GeneralLogger) IsWarningEnabled() bool {
	return l.warningEnabled
}

// Indicator whether error level is enabled or not
func (l GeneralLogger) IsErrorEnabled() bool {
	return l.errorEnabled
}

// Indicator whether fatal level is enabled or not
func (l GeneralLogger) IsFatalEnabled() bool {
	return l.fatalEnabled
}

func (l *GeneralLogger) closeAppender() {
	(*l.appender).Close()
}

func panicOrMock(message string) {
	if mockPanicAndExitAtGeneralLogger {
		panicMockActivated = true
		return
	}
	panic(message)
}

func (l *GeneralLogger) exitOrMock(code int) {
	l.closeAppender()

	if mockPanicAndExitAtGeneralLogger {
		exitMockActivated = true
		return
	}
	os.Exit(code)
}

func (l *GeneralLogger) write(severity int, message string) {
	logValuesToWrite := common.CreateLogValues(severity, message)
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *GeneralLogger) writeWithCorrelation(severity int, correlationId string, message string) {
	logValuesToWrite := common.CreateLogValuesWithCorrelation(severity, &correlationId, message)
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *GeneralLogger) writeCustom(severity int, message string, customValues map[string]any) {
	logValuesToWrite := common.CreateLogValuesCustom(severity, message, &customValues)
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *GeneralLogger) writeWithCtx(severity int, context context.Context, message string) {
	correlationId, exists := context.Value(l.correlationIdKey).(string)
	var logValuesToWrite common.LogValues
	if exists {
		logValuesToWrite = common.CreateLogValuesWithCorrelation(severity, &correlationId, message)
	} else {
		logValuesToWrite = common.CreateLogValues(severity, message)
	}
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *GeneralLogger) setCallerValues(logValuesToWrite *common.LogValues) {
	if l.isCallerToSet {
		rpc := make([]uintptr, 2)
		callersCount := runtime.Callers(5, rpc)
		if callersCount < 1 {
			logValuesToWrite.IsCallerSet = false
			return
		}
		frames := runtime.CallersFrames(rpc)
		frame, more := frames.Next()
		if setCallerFromFrame(&frame, logValuesToWrite) {
			return
		}
		if more {
			frame, _ = frames.Next()
			if setCallerFromFrame(&frame, logValuesToWrite) {
				return
			}
		}
	}
	logValuesToWrite.IsCallerSet = false
}

func setCallerFromFrame(frame *runtime.Frame, logValuesToWrite *common.LogValues) bool {
	if isRelevantCaller(frame) {
		adoptFrameValues(frame, logValuesToWrite)
		return true
	}
	return false
}

func isRelevantCaller(frame *runtime.Frame) bool {
	return frame.PC != 0 && (!strings.HasPrefix(frame.Func.Name(), "github.com/ma-vin/typewriter/logger") || !strings.HasSuffix(frame.File, "logger.go"))
}

func adoptFrameValues(source *runtime.Frame, target *common.LogValues) {
	target.CallerFunction = source.Func.Name()
	target.CallerFile = source.File
	target.CallerFileLine = source.Line
	target.IsCallerSet = true
}
