package logger

import (
	"runtime"
	"strings"
)

type MainLogger struct {
	commonLogger       *CommonLogger
	existPackageLogger bool
	packageLoggers     map[string]*CommonLogger
}

// Determines which logger is relevant. If der exists a logger for a package equal to the callers package, this logger will be return, else the commonlogger.
func (l *MainLogger) determineLogger() *CommonLogger {
	if l.existPackageLogger {
		pc, _, _, ok := runtime.Caller(2)
		if !ok {
			return l.commonLogger
		}

		pl, found := l.packageLoggers[determinePackageName(runtime.FuncForPC(pc).Name())]
		if found {
			return pl
		}
	}
	return l.commonLogger
}

// extracts the packename from a given function name. E.g. the result with parameter "github.com/ma-vin/typewriter/logger.determinePackageName" would be "logger"
func determinePackageName(functionName string) string {
	packageBegin := strings.LastIndex(functionName, "/") + 1
	var functionNameSuffix string
	if packageBegin > 0 {
		functionNameSuffix = functionName[packageBegin:]
	} else {
		functionNameSuffix = functionName
	}
	packageEnd := strings.Index(functionNameSuffix, ".")
	if packageEnd > -1 {
		return strings.ToUpper(functionNameSuffix[:packageEnd])
	}
	return strings.ToUpper(functionNameSuffix)

}

func (l MainLogger) Debug(args ...any) {
	l.determineLogger().Debug(args...)
}

func (l MainLogger) DebugWithCorrelation(correlationId string, args ...any) {
	l.determineLogger().DebugWithCorrelation(correlationId, args...)
}

func (l MainLogger) DebugCustom(customValues map[string]any, args ...any) {
	l.determineLogger().DebugCustom(customValues, args...)
}

func (l MainLogger) Debugf(format string, args ...any) {
	l.determineLogger().Debugf(format, args...)
}

func (l MainLogger) DebugWithCorrelationf(correlationId string, format string, args ...any) {
	l.determineLogger().DebugWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) DebugCustomf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().DebugCustomf(customValues, format, args...)
}

func (l MainLogger) Information(args ...any) {
	l.determineLogger().Information(args...)
}

func (l MainLogger) InformationWithCorrelation(correlationId string, args ...any) {
	l.determineLogger().InformationWithCorrelation(correlationId, args...)
}

func (l MainLogger) InformationCustom(customValues map[string]any, args ...any) {
	l.determineLogger().InformationCustom(customValues, args...)
}

func (l MainLogger) Informationf(format string, args ...any) {
	l.determineLogger().Informationf(format, args...)
}

func (l MainLogger) InformationWithCorrelationf(correlationId string, format string, args ...any) {
	l.determineLogger().InformationWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) InformationCustomf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().InformationCustomf(customValues, format, args...)
}

func (l MainLogger) Warning(args ...any) {
	l.determineLogger().Warning(args...)
}

func (l MainLogger) WarningWithCorrelation(correlationId string, args ...any) {
	l.determineLogger().WarningWithCorrelation(correlationId, args...)
}

func (l MainLogger) WarningCustom(customValues map[string]any, args ...any) {
	l.determineLogger().WarningCustom(customValues, args...)
}

func (l MainLogger) Warningf(format string, args ...any) {
	l.determineLogger().Warningf(format, args...)
}

func (l MainLogger) WarningWithCorrelationf(correlationId string, format string, args ...any) {
	l.determineLogger().WarningWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) WarningCustomf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().WarningCustomf(customValues, format, args...)
}

func (l MainLogger) WarningWithPanic(args ...any) {
	l.determineLogger().WarningWithPanic(args...)
}

func (l MainLogger) WarningWithCorrelationAndPanic(correlationId string, args ...any) {
	l.determineLogger().WarningWithCorrelationAndPanic(correlationId, args...)
}

func (l MainLogger) WarningCustomWithPanic(customValues map[string]any, args ...any) {
	l.determineLogger().WarningCustomWithPanic(customValues, args...)
}

func (l MainLogger) WarningWithPanicf(format string, args ...any) {
	l.determineLogger().WarningWithPanicf(format, args...)
}

func (l MainLogger) WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.determineLogger().WarningWithCorrelationAndPanicf(correlationId, format, args...)
}

func (l MainLogger) WarningCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().WarningCustomWithPanicf(customValues, format, args...)
}

func (l MainLogger) Error(args ...any) {
	l.determineLogger().Error(args...)
}

func (l MainLogger) ErrorWithCorrelation(correlationId string, args ...any) {
	l.determineLogger().ErrorWithCorrelation(correlationId, args...)
}

func (l MainLogger) ErrorCustom(customValues map[string]any, args ...any) {
	l.determineLogger().ErrorCustom(customValues, args...)
}

func (l MainLogger) Errorf(format string, args ...any) {
	l.determineLogger().Errorf(format, args...)
}

func (l MainLogger) ErrorWithCorrelationf(correlationId string, format string, args ...any) {
	l.determineLogger().ErrorWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) ErrorCustomf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().ErrorCustomf(customValues, format, args...)
}

func (l MainLogger) ErrorWithPanic(args ...any) {
	l.determineLogger().ErrorWithPanic(args...)
}

func (l MainLogger) ErrorWithCorrelationAndPanic(correlationId string, args ...any) {
	l.determineLogger().ErrorWithCorrelationAndPanic(correlationId, args...)
}

func (l MainLogger) ErrorCustomWithPanic(customValues map[string]any, args ...any) {
	l.determineLogger().ErrorCustomWithPanic(customValues, args...)
}

func (l MainLogger) ErrorWithPanicf(format string, args ...any) {
	l.determineLogger().ErrorWithPanicf(format, args...)
}

func (l MainLogger) ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.determineLogger().ErrorWithCorrelationAndPanicf(correlationId, format, args...)
}

func (l MainLogger) ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().ErrorCustomWithPanicf(customValues, format, args...)
}

func (l MainLogger) Fatal(args ...any) {
	l.determineLogger().Fatal(args...)
}

func (l MainLogger) FatalWithCorrelation(correlationId string, args ...any) {
	l.determineLogger().FatalWithCorrelation(correlationId, args...)
}

func (l MainLogger) FatalCustom(customValues map[string]any, args ...any) {
	l.determineLogger().FatalCustom(customValues, args...)
}

func (l MainLogger) Fatalf(format string, args ...any) {
	l.determineLogger().Fatalf(format, args...)
}

func (l MainLogger) FatalWithCorrelationf(correlationId string, format string, args ...any) {
	l.determineLogger().FatalWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) FatalCustomf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().FatalCustomf(customValues, format, args...)
}

func (l MainLogger) FatalWithPanic(args ...any) {
	l.determineLogger().FatalWithPanic(args...)
}

func (l MainLogger) FatalWithCorrelationAndPanic(correlationId string, args ...any) {
	l.determineLogger().FatalWithCorrelationAndPanic(correlationId, args...)
}

func (l MainLogger) FatalCustomWithPanic(customValues map[string]any, args ...any) {
	l.determineLogger().FatalCustomWithPanic(customValues, args...)
}

func (l MainLogger) FatalWithPanicf(format string, args ...any) {
	l.determineLogger().FatalWithPanicf(format, args...)
}

func (l MainLogger) FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	l.determineLogger().FatalWithCorrelationAndPanicf(correlationId, format, args...)
}

func (l MainLogger) FatalCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	l.determineLogger().FatalCustomWithPanicf(customValues, format, args...)
}

func (l MainLogger) FatalWithExit(args ...any) {
	l.determineLoggerAndCloseOthersAppender().FatalWithExit(args...)
}

func (l MainLogger) FatalWithCorrelationAndExit(correlationId string, args ...any) {
	l.determineLoggerAndCloseOthersAppender().FatalWithCorrelationAndExit(correlationId, args...)
}

func (l MainLogger) FatalCustomWithExit(customValues map[string]any, args ...any) {
	l.determineLoggerAndCloseOthersAppender().FatalCustomWithExit(customValues, args...)
}

func (l MainLogger) FatalWithExitf(format string, args ...any) {
	l.determineLoggerAndCloseOthersAppender().FatalWithExitf(format, args...)
}

func (l MainLogger) FatalWithCorrelationAndExitf(correlationId string, format string, args ...any) {
	l.determineLoggerAndCloseOthersAppender().FatalWithCorrelationAndExitf(correlationId, format, args...)
}

func (l MainLogger) FatalCustomWithExitf(customValues map[string]any, format string, args ...any) {
	l.determineLoggerAndCloseOthersAppender().FatalCustomWithExitf(customValues, format, args...)
}

func (l MainLogger) IsDebugEnabled() bool {
	return l.determineLogger().IsDebugEnabled()
}

func (l MainLogger) IsInformationEnabled() bool {
	return l.determineLogger().IsInformationEnabled()
}

func (l MainLogger) IsWarningEnabled() bool {
	return l.determineLogger().IsWarningEnabled()
}

func (l MainLogger) IsErrorEnabled() bool {
	return l.determineLogger().IsErrorEnabled()
}

func (l MainLogger) IsFatalEnabled() bool {
	return l.determineLogger().IsFatalEnabled()
}

func (l *MainLogger) determineLoggerAndCloseOthersAppender() *CommonLogger {
	relevantLogger := l.determineLogger()
	l.closeAppender(relevantLogger)
	return relevantLogger
}

func (l *MainLogger) closeAppender(loggerToSkip *CommonLogger) {
	if l.commonLogger != loggerToSkip {
		l.commonLogger.closeAppender()
	}
	if !l.existPackageLogger {
		// There exists only commonlogger: nothing to do
		return
	}
	for _, pLog := range l.packageLoggers {
		if pLog != loggerToSkip {
			pLog.closeAppender()
		}
	}
}
