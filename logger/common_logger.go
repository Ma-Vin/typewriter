package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/common"
)

// A common logger which delegates messages directly to the appender if log level is enabled.
type CommonLogger struct {
	debugEnabled       bool
	informationEnabled bool
	warningEnabled     bool
	errorEnabled       bool
	fatalEnabled       bool
	isCallerToSet      bool
	appender           *appender.Appender
}

var mockPanicAndExitAtCommonLogger = false
var panicMockActivated = false
var exitMockActivated = false

// Creates a common logger which delegates messages to the given appender if the log level is enabled by given severity
func CreateCommonLogger(appender *appender.Appender, severity int, isCallerToSet bool) CommonLogger {
	result := CommonLogger{appender: appender, isCallerToSet: isCallerToSet}
	determineSeverityByLevel(&result, severity)
	return result
}

// activates different log levels
func determineSeverityByLevel(l *CommonLogger, severity int) {
	l.debugEnabled = common.DEBUG_SEVERITY <= severity
	l.informationEnabled = common.INFORMATION_SEVERITY <= severity
	l.warningEnabled = common.WARNING_SEVERITY <= severity
	l.errorEnabled = common.ERROR_SEVERITY <= severity
	l.fatalEnabled = common.FATAL_SEVERITY <= severity
}

// Logs a message if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) Debug(args ...any) {
	if l.debugEnabled {
		l.write(common.DEBUG_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) DebugWithCorrelation(correlationId string, args ...any) {
	if l.debugEnabled {
		l.writeWithCorrelation(common.DEBUG_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) DebugCustom(customValues map[string]any, args ...any) {
	if l.debugEnabled {
		l.writeCustom(common.DEBUG_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message derived from format if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) Debugf(format string, args ...any) {
	if l.debugEnabled {
		l.write(common.DEBUG_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) DebugWithCorrelationf(correlationId string, format string, args ...any) {
	if l.debugEnabled {
		l.writeWithCorrelation(common.DEBUG_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) DebugCustomf(customValues map[string]any, format string, args ...any) {
	if l.debugEnabled {
		l.writeCustom(common.DEBUG_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) Information(args ...any) {
	if l.informationEnabled {
		l.write(common.INFORMATION_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) InformationWithCorrelation(correlationId string, args ...any) {
	if l.informationEnabled {
		l.writeWithCorrelation(common.INFORMATION_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) InformationCustom(customValues map[string]any, args ...any) {
	if l.informationEnabled {
		l.writeCustom(common.INFORMATION_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message derived from format if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) Informationf(format string, args ...any) {
	if l.informationEnabled {
		l.write(common.INFORMATION_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) InformationWithCorrelationf(correlationId string, format string, args ...any) {
	if l.informationEnabled {
		l.writeWithCorrelation(common.INFORMATION_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) InformationCustomf(customValues map[string]any, format string, args ...any) {
	if l.informationEnabled {
		l.writeCustom(common.INFORMATION_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) Warning(args ...any) {
	if l.warningEnabled {
		l.write(common.WARNING_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) WarningWithCorrelation(correlationId string, args ...any) {
	if l.warningEnabled {
		l.writeWithCorrelation(common.WARNING_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) WarningCustom(customValues map[string]any, args ...any) {
	if l.warningEnabled {
		l.writeCustom(common.WARNING_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message derived from format if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) Warningf(format string, args ...any) {
	if l.warningEnabled {
		l.write(common.WARNING_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) WarningWithCorrelationf(correlationId string, format string, args ...any) {
	if l.warningEnabled {
		l.writeWithCorrelation(common.WARNING_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) WarningCustomf(customValues map[string]any, format string, args ...any) {
	if l.warningEnabled {
		l.writeCustom(common.WARNING_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) WarningWithPanic(args ...any) {
	l.Warning(args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) WarningWithCorrelationAndPanic(correlationId string, args ...any) {
	l.WarningWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) WarningCustomWithPanic(customValues map[string]any, args ...any) {
	l.WarningCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) WarningWithPanicf(format string, args ...any) {
	l.Warningf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.WarningWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) WarningCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.WarningCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) Error(args ...any) {
	if l.errorEnabled {
		l.write(common.ERROR_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) ErrorWithCorrelation(correlationId string, args ...any) {
	if l.errorEnabled {
		l.writeWithCorrelation(common.ERROR_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) ErrorCustom(customValues map[string]any, args ...any) {
	if l.errorEnabled {
		l.writeCustom(common.ERROR_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message derived from format if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) Errorf(format string, args ...any) {
	if l.errorEnabled {
		l.write(common.ERROR_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) ErrorWithCorrelationf(correlationId string, format string, args ...any) {
	if l.errorEnabled {
		l.writeWithCorrelation(common.ERROR_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) ErrorCustomf(customValues map[string]any, format string, args ...any) {
	if l.errorEnabled {
		l.writeCustom(common.ERROR_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) ErrorWithPanic(args ...any) {
	l.Error(args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) ErrorWithCorrelationAndPanic(correlationId string, args ...any) {
	l.ErrorWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) ErrorCustomWithPanic(customValues map[string]any, args ...any) {
	l.ErrorCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) ErrorWithPanicf(format string, args ...any) {
	l.Errorf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.ErrorWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.ErrorCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) Fatal(args ...any) {
	if l.fatalEnabled {
		l.write(common.FATAL_SEVERITY, fmt.Sprint(args...))
	}
}

// Logs a message together with a correlation id if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalWithCorrelation(correlationId string, args ...any) {
	if l.fatalEnabled {
		l.writeWithCorrelation(common.FATAL_SEVERITY, correlationId, fmt.Sprint(args...))
	}
}

// Logs a message together with custom values if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalCustom(customValues map[string]any, args ...any) {
	if l.fatalEnabled {
		l.writeCustom(common.FATAL_SEVERITY, fmt.Sprint(args...), customValues)
	}
}

// Logs a message derived from format if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) Fatalf(format string, args ...any) {
	if l.fatalEnabled {
		l.write(common.FATAL_SEVERITY, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with a correlation id if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalWithCorrelationf(correlationId string, format string, args ...any) {
	if l.fatalEnabled {
		l.writeWithCorrelation(common.FATAL_SEVERITY, correlationId, fmt.Sprintf(format, args...))
	}
}

// Logs a message derived from format together with custom values if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalCustomf(customValues map[string]any, format string, args ...any) {
	if l.fatalEnabled {
		l.writeCustom(common.FATAL_SEVERITY, fmt.Sprintf(format, args...), customValues)
	}
}

// Logs a message if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalWithPanic(args ...any) {
	l.Fatal(args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalWithCorrelationAndPanic(correlationId string, args ...any) {
	l.FatalWithCorrelation(correlationId, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalCustomWithPanic(customValues map[string]any, args ...any) {
	l.FatalCustom(customValues, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalWithPanicf(format string, args ...any) {
	l.Fatalf(format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.FatalWithCorrelationf(correlationId, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message derived from format together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.FatalCustomf(customValues, format, args...)
	panicOrMock(fmt.Sprint(args...))
}

// Logs a message if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalWithExit(args ...any) {
	l.Fatal(args...)
	l.exitOrMock(1)
}

// Logs a message together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalWithCorrelationAndExit(correlationId string, args ...any) {
	l.FatalWithCorrelation(correlationId, args...)
	l.exitOrMock(1)
}

// Logs a message together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func (l CommonLogger) FatalCustomWithExit(customValues map[string]any, args ...any) {
	l.FatalCustom(customValues, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalWithExitf(format string, args ...any) {
	l.Fatalf(format, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalWithCorrelationAndExitf(correlationId string, format string, args ...any) {
	l.FatalWithCorrelationf(correlationId, format, args...)
	l.exitOrMock(1)
}

// Logs a message derived from format together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func (l CommonLogger) FatalCustomWithExitf(customValues map[string]any, format string, args ...any) {
	l.FatalCustomf(customValues, format, args...)
	l.exitOrMock(1)
}

// Indicator whether debug level is enabled or not
func (l CommonLogger) IsDebugEnabled() bool {
	return l.debugEnabled
}

// Indicator whether information level is enabled or not
func (l CommonLogger) IsInformationEnabled() bool {
	return l.informationEnabled
}

// Indicator whether warning level is enabled or not
func (l CommonLogger) IsWarningEnabled() bool {
	return l.warningEnabled
}

// Indicator whether error level is enabled or not
func (l CommonLogger) IsErrorEnabled() bool {
	return l.errorEnabled
}

// Indicator whether fatal level is enabled or not
func (l CommonLogger) IsFatalEnabled() bool {
	return l.fatalEnabled
}

func (l *CommonLogger) closeAppender() {
	(*l.appender).Close()
}

func panicOrMock(message string) {
	if mockPanicAndExitAtCommonLogger {
		panicMockActivated = true
		return
	}
	panic(message)
}

func (l *CommonLogger) exitOrMock(code int) {
	l.closeAppender()

	if mockPanicAndExitAtCommonLogger {
		exitMockActivated = true
		return
	}
	os.Exit(code)
}

func (l *CommonLogger) write(severity int, message string) {
	logValuesToWrite := common.CreateLogValues(severity, message)
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *CommonLogger) writeWithCorrelation(severity int, correlationId string, message string) {
	logValuesToWrite := common.CreateLogValuesWithCorrelation(severity, &correlationId, message)
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *CommonLogger) writeCustom(severity int, message string, customValues map[string]any) {
	logValuesToWrite := common.CreateLogValuesCustom(severity, message, &customValues)
	l.setCallerValues(&logValuesToWrite)
	(*l.appender).Write(&logValuesToWrite)
}

func (l *CommonLogger) setCallerValues(logValuesToWrite *common.LogValues) {
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
