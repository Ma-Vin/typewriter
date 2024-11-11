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
func determineLogger(l *MainLogger) *CommonLogger {
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
	determineLogger(&l).Debug(args...)
}

func (l MainLogger) DebugWithCorrelation(correlationId string, args ...any) {
	determineLogger(&l).DebugWithCorrelation(correlationId, args...)
}

func (l MainLogger) DebugCustom(customValues map[string]any, args ...any) {
	determineLogger(&l).DebugCustom(customValues, args...)
}

func (l MainLogger) Debugf(format string, args ...any) {
	determineLogger(&l).Debugf(format, args...)
}

func (l MainLogger) DebugWithCorrelationf(correlationId string, format string, args ...any) {
	determineLogger(&l).DebugWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) DebugCustomf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).DebugCustomf(customValues, format, args...)
}

func (l MainLogger) Information(args ...any) {
	determineLogger(&l).Information(args...)
}

func (l MainLogger) InformationWithCorrelation(correlationId string, args ...any) {
	determineLogger(&l).InformationWithCorrelation(correlationId, args...)
}

func (l MainLogger) InformationCustom(customValues map[string]any, args ...any) {
	determineLogger(&l).InformationCustom(customValues, args...)
}

func (l MainLogger) Informationf(format string, args ...any) {
	determineLogger(&l).Informationf(format, args...)
}

func (l MainLogger) InformationWithCorrelationf(correlationId string, format string, args ...any) {
	determineLogger(&l).InformationWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) InformationCustomf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).InformationCustomf(customValues, format, args...)
}

func (l MainLogger) Warning(args ...any) {
	determineLogger(&l).Warning(args...)
}

func (l MainLogger) WarningWithCorrelation(correlationId string, args ...any) {
	determineLogger(&l).WarningWithCorrelation(correlationId, args...)
}

func (l MainLogger) WarningCustom(customValues map[string]any, args ...any) {
	determineLogger(&l).WarningCustom(customValues, args...)
}

func (l MainLogger) Warningf(format string, args ...any) {
	determineLogger(&l).Warningf(format, args...)
}

func (l MainLogger) WarningWithCorrelationf(correlationId string, format string, args ...any) {
	determineLogger(&l).WarningWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) WarningCustomf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).WarningCustomf(customValues, format, args...)
}

func (l MainLogger) WarningWithPanic(args ...any) {
	determineLogger(&l).WarningWithPanic(args...)
}

func (l MainLogger) WarningWithCorrelationAndPanic(correlationId string, args ...any) {
	determineLogger(&l).WarningWithCorrelationAndPanic(correlationId, args...)
}

func (l MainLogger) WarningCustomWithPanic(customValues map[string]any, args ...any) {
	determineLogger(&l).WarningCustomWithPanic(customValues, args...)
}

func (l MainLogger) WarningWithPanicf(format string, args ...any) {
	determineLogger(&l).WarningWithPanicf(format, args...)
}

func (l MainLogger) WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	determineLogger(&l).WarningWithCorrelationAndPanicf(correlationId, format, args...)
}

func (l MainLogger) WarningCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).WarningCustomWithPanicf(customValues, format, args...)
}

func (l MainLogger) Error(args ...any) {
	determineLogger(&l).Error(args...)
}

func (l MainLogger) ErrorWithCorrelation(correlationId string, args ...any) {
	determineLogger(&l).ErrorWithCorrelation(correlationId, args...)
}

func (l MainLogger) ErrorCustom(customValues map[string]any, args ...any) {
	determineLogger(&l).ErrorCustom(customValues, args...)
}

func (l MainLogger) Errorf(format string, args ...any) {
	determineLogger(&l).Errorf(format, args...)
}

func (l MainLogger) ErrorWithCorrelationf(correlationId string, format string, args ...any) {
	determineLogger(&l).ErrorWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) ErrorCustomf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).ErrorCustomf(customValues, format, args...)
}

func (l MainLogger) ErrorWithPanic(args ...any) {
	determineLogger(&l).ErrorWithPanic(args...)
}

func (l MainLogger) ErrorWithCorrelationAndPanic(correlationId string, args ...any) {
	determineLogger(&l).ErrorWithCorrelationAndPanic(correlationId, args...)
}

func (l MainLogger) ErrorCustomWithPanic(customValues map[string]any, args ...any) {
	determineLogger(&l).ErrorCustomWithPanic(customValues, args...)
}

func (l MainLogger) ErrorWithPanicf(format string, args ...any) {
	determineLogger(&l).ErrorWithPanicf(format, args...)
}

func (l MainLogger) ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	determineLogger(&l).ErrorWithCorrelationAndPanicf(correlationId, format, args...)
}

func (l MainLogger) ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).ErrorCustomWithPanicf(customValues, format, args...)
}

func (l MainLogger) Fatal(args ...any) {
	determineLogger(&l).Fatal(args...)
}

func (l MainLogger) FatalWithCorrelation(correlationId string, args ...any) {
	determineLogger(&l).FatalWithCorrelation(correlationId, args...)
}

func (l MainLogger) FatalCustom(customValues map[string]any, args ...any) {
	determineLogger(&l).FatalCustom(customValues, args...)
}

func (l MainLogger) Fatalf(format string, args ...any) {
	determineLogger(&l).Fatalf(format, args...)
}

func (l MainLogger) FatalWithCorrelationf(correlationId string, format string, args ...any) {
	determineLogger(&l).FatalWithCorrelationf(correlationId, format, args...)
}

func (l MainLogger) FatalCustomf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).FatalCustomf(customValues, format, args...)
}

func (l MainLogger) FatalWithPanic(args ...any) {
	determineLogger(&l).FatalWithPanic(args...)
}

func (l MainLogger) FatalWithCorrelationAndPanic(correlationId string, args ...any) {
	determineLogger(&l).FatalWithCorrelationAndPanic(correlationId, args...)
}

func (l MainLogger) FatalCustomWithPanic(customValues map[string]any, args ...any) {
	determineLogger(&l).FatalCustomWithPanic(customValues, args...)
}

func (l MainLogger) FatalWithPanicf(format string, args ...any) {
	determineLogger(&l).FatalWithPanicf(format, args...)
}

func (l MainLogger) FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	determineLogger(&l).FatalWithCorrelationAndPanicf(correlationId, format, args...)
}

func (l MainLogger) FatalCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).FatalCustomWithPanicf(customValues, format, args...)
}

func (l MainLogger) FatalWithExit(args ...any) {
	determineLogger(&l).FatalWithExit(args...)
}

func (l MainLogger) FatalWithCorrelationAndExit(correlationId string, args ...any) {
	determineLogger(&l).FatalWithCorrelationAndExit(correlationId, args...)
}

func (l MainLogger) FatalCustomWithExit(customValues map[string]any, args ...any) {
	determineLogger(&l).FatalCustomWithExit(customValues, args...)
}

func (l MainLogger) FatalWithExitf(format string, args ...any) {
	determineLogger(&l).FatalWithExitf(format, args...)
}

func (l MainLogger) FatalWithCorrelationAndExitf(correlationId string, format string, args ...any) {
	determineLogger(&l).FatalWithCorrelationAndExitf(correlationId, format, args...)
}

func (l MainLogger) FatalCustomWithExitf(customValues map[string]any, format string, args ...any) {
	determineLogger(&l).FatalCustomWithExitf(customValues, format, args...)
}

func (l MainLogger) IsDebugEnabled() bool {
	return determineLogger(&l).IsDebugEnabled()
}

func (l MainLogger) IsInformationEnabled() bool {
	return determineLogger(&l).IsInformationEnabled()
}

func (l MainLogger) IsWarningEnabled() bool {
	return determineLogger(&l).IsWarningEnabled()
}

func (l MainLogger) IsErrorEnabled() bool {
	return determineLogger(&l).IsErrorEnabled()
}

func (l MainLogger) IsFatalEnabled() bool {
	return determineLogger(&l).IsFatalEnabled()
}
