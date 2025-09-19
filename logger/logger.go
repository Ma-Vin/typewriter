// This package provides the accessor [pkg/github.com/ma-vin/typewriter/logger.Log] to get the relevant logger,
// which provides all logging functions by implementing the interface [pkg/github.com/ma-vin/typewriter/logger.Logger]
package logger

import (
	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
)

type Logger interface {

	// Logs a message if debug level is enabled.
	// Arguments are handled in the manner of [fmt.Sprint].
	Debug(args ...any)

	// Logs a message together with a correlation id if debug level is enabled.
	// Arguments are handled in the manner of [fmt.Sprint].
	DebugWithCorrelation(correlationId string, args ...any)

	// Logs a message together with custom values if debug level is enabled.
	// Arguments are handled in the manner of [fmt.Sprint].
	DebugCustom(customValues map[string]any, args ...any)

	// Logs a message derived from format if debug level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	Debugf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if debug level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	DebugWithCorrelationf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if debug level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	DebugCustomf(customValues map[string]any, format string, args ...any)

	// Logs a message if information level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	Information(args ...any)

	// Logs a message together with a correlation id if information level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	InformationWithCorrelation(correlationId string, args ...any)

	// Logs a message together with custom values if information level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	InformationCustom(customValues map[string]any, args ...any)

	// Logs a message derived from format if information level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	Informationf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if information level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	InformationWithCorrelationf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if information level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	InformationCustomf(customValues map[string]any, format string, args ...any)

	// Logs a message if warning level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	Warning(args ...any)

	// Logs a message together with a correlation id if warning level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	WarningWithCorrelation(correlationId string, args ...any)

	// Logs a message together with custom values if warning level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	WarningCustom(customValues map[string]any, args ...any)

	// Logs a message derived from format if warning level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	Warningf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if warning level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	WarningWithCorrelationf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if warning level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	WarningCustomf(customValues map[string]any, format string, args ...any)

	// Logs a message if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	WarningWithPanic(args ...any)

	// Logs a message together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	WarningWithCorrelationAndPanic(correlationId string, args ...any)

	// Logs a message together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	WarningCustomWithPanic(customValues map[string]any, args ...any)

	// Logs a message derived from format if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	WarningWithPanicf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	WarningCustomWithPanicf(customValues map[string]any, format string, args ...any)

	// Logs a message if error level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	Error(args ...any)

	// Logs a message together with a correlation id if error level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	ErrorWithCorrelation(correlationId string, args ...any)

	// Logs a message together with custom values if error level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	ErrorCustom(customValues map[string]any, args ...any)

	// Logs a message derived from format if error level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	Errorf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if error level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	ErrorWithCorrelationf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if error level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	ErrorCustomf(customValues map[string]any, format string, args ...any)

	// Logs a message if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	ErrorWithPanic(args ...any)

	// Logs a message together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	ErrorWithCorrelationAndPanic(correlationId string, args ...any)

	// Logs a message together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	ErrorCustomWithPanic(customValues map[string]any, args ...any)

	// Logs a message derived from format if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	ErrorWithPanicf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any)

	// Logs a message if fatal level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	Fatal(args ...any)

	// Logs a message together with a correlation id if fatal level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalWithCorrelation(correlationId string, args ...any)

	// Logs a message together with custom values if fatal level is enabled
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalCustom(customValues map[string]any, args ...any)

	// Logs a message derived from format if fatal level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	Fatalf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if fatal level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalWithCorrelationf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if fatal level is enabled
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalCustomf(customValues map[string]any, format string, args ...any)

	// Logs a message if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalWithPanic(args ...any)

	// Logs a message together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalWithCorrelationAndPanic(correlationId string, args ...any)

	// Logs a message together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalCustomWithPanic(customValues map[string]any, args ...any)

	// Logs a message derived from format if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalWithPanicf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalCustomWithPanicf(customValues map[string]any, format string, args ...any)

	// Logs a message if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalWithExit(args ...any)

	// Logs a message together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalWithCorrelationAndExit(correlationId string, args ...any)

	// Logs a message together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
	// Arguments are handled in the manner of [fmt.Sprint].
	FatalCustomWithExit(customValues map[string]any, args ...any)

	// Logs a message derived from format if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalWithExitf(format string, args ...any)

	// Logs a message derived from format together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalWithCorrelationAndExitf(correlationId string, format string, args ...any)

	// Logs a message derived from format together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
	// Arguments are handled in the manner of [fmt.Sprintf].
	FatalCustomWithExitf(customValues map[string]any, format string, args ...any)

	// Indicator whether debug level is enabled or not
	IsDebugEnabled() bool

	// Indicator whether information level is enabled or not
	IsInformationEnabled() bool

	// Indicator whether warning level is enabled or not
	IsWarningEnabled() bool

	// Indicator whether error level is enabled or not
	IsErrorEnabled() bool

	// Indicator whether fatal level is enabled or not
	IsFatalEnabled() bool
}

// Returns a pointer to the main logger which provides all methods of the Logger interface
func Log() Logger {
	return getLoggers()
}

// Logs a message if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func Debug(args ...any) {
	getLoggers().Debug(args...)
}

// Logs a message together with a correlation id if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func DebugWithCorrelation(correlationId string, args ...any) {
	getLoggers().DebugWithCorrelation(correlationId, args...)
}

// Logs a message together with custom values if debug level is enabled.
// Arguments are handled in the manner of [fmt.Sprint].
func DebugCustom(customValues map[string]any, args ...any) {
	getLoggers().DebugCustom(customValues, args...)
}

// Logs a message derived from format if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func Debugf(format string, args ...any) {
	getLoggers().Debugf(format, args...)
}

// Logs a message derived from format together with a correlation id if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func DebugWithCorrelationf(correlationId string, format string, args ...any) {
	getLoggers().DebugWithCorrelationf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if debug level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func DebugCustomf(customValues map[string]any, format string, args ...any) {
	getLoggers().DebugCustomf(customValues, format, args...)
}

// Logs a message if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func Information(args ...any) {
	getLoggers().Information(args...)
}

// Logs a message together with a correlation id if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func InformationWithCorrelation(correlationId string, args ...any) {
	getLoggers().InformationWithCorrelation(correlationId, args...)
}

// Logs a message together with custom values if information level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func InformationCustom(customValues map[string]any, args ...any) {
	getLoggers().InformationCustom(customValues, args...)
}

// Logs a message derived from format if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func Informationf(format string, args ...any) {
	getLoggers().Informationf(format, args...)
}

// Logs a message derived from format together with a correlation id if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func InformationWithCorrelationf(correlationId string, format string, args ...any) {
	getLoggers().InformationWithCorrelationf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if information level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func InformationCustomf(customValues map[string]any, format string, args ...any) {
	getLoggers().InformationCustomf(customValues, format, args...)
}

// Logs a message if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func Warning(args ...any) {
	getLoggers().Warning(args...)
}

// Logs a message together with a correlation id if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func WarningWithCorrelation(correlationId string, args ...any) {
	getLoggers().WarningWithCorrelation(correlationId, args...)
}

// Logs a message together with custom values if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func WarningCustom(customValues map[string]any, args ...any) {
	getLoggers().WarningCustom(customValues, args...)
}

// Logs a message derived from format if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func Warningf(format string, args ...any) {
	getLoggers().Warningf(format, args...)
}

// Logs a message derived from format together with a correlation id if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func WarningWithCorrelationf(correlationId string, format string, args ...any) {
	getLoggers().WarningWithCorrelationf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if warning level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func WarningCustomf(customValues map[string]any, format string, args ...any) {
	getLoggers().WarningCustomf(customValues, format, args...)
}

// Logs a message if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func WarningWithPanic(args ...any) {
	getLoggers().WarningWithPanic(args...)
}

// Logs a message together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func WarningWithCorrelationAndPanic(correlationId string, args ...any) {
	getLoggers().WarningWithCorrelationAndPanic(correlationId, args...)
}

// Logs a message together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func WarningCustomWithPanic(customValues map[string]any, args ...any) {
	getLoggers().WarningCustomWithPanic(customValues, args...)
}

// Logs a message derived from format if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func WarningWithPanicf(format string, args ...any) {
	getLoggers().WarningWithPanicf(format, args...)
}

// Logs a message derived from format together with a correlation id if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func WarningWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	getLoggers().WarningWithCorrelationAndPanicf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if warning level is enabled and calls built-in function panic to stop current goroutine (independent if warning level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func WarningCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	getLoggers().WarningCustomWithPanicf(customValues, format, args...)
}

// Logs a message if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func Error(args ...any) {
	getLoggers().Error(args...)
}

// Logs a message together with a correlation id if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func ErrorWithCorrelation(correlationId string, args ...any) {
	getLoggers().ErrorWithCorrelation(correlationId, args...)
}

// Logs a message together with custom values if error level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func ErrorCustom(customValues map[string]any, args ...any) {
	getLoggers().ErrorCustom(customValues, args...)
}

// Logs a message derived from format if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func Errorf(format string, args ...any) {
	getLoggers().Errorf(format, args...)
}

// Logs a message derived from format together with a correlation id if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func ErrorWithCorrelationf(correlationId string, format string, args ...any) {
	getLoggers().ErrorWithCorrelationf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if error level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func ErrorCustomf(customValues map[string]any, format string, args ...any) {
	getLoggers().ErrorCustomf(customValues, format, args...)
}

// Logs a message if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func ErrorWithPanic(args ...any) {
	getLoggers().ErrorWithPanic(args...)
}

// Logs a message together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func ErrorWithCorrelationAndPanic(correlationId string, args ...any) {
	getLoggers().ErrorWithCorrelationAndPanic(correlationId, args...)
}

// Logs a message together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func ErrorCustomWithPanic(customValues map[string]any, args ...any) {
	getLoggers().ErrorCustomWithPanic(customValues, args...)
}

// Logs a message derived from format if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func ErrorWithPanicf(format string, args ...any) {
	getLoggers().ErrorWithPanicf(format, args...)
}

// Logs a message derived from format together with a correlation id if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func ErrorWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	getLoggers().ErrorWithCorrelationAndPanicf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if error level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func ErrorCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	getLoggers().ErrorCustomWithPanicf(customValues, format, args...)
}

// Logs a message if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func Fatal(args ...any) {
	getLoggers().Fatal(args...)
}

// Logs a message together with a correlation id if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func FatalWithCorrelation(correlationId string, args ...any) {
	getLoggers().FatalWithCorrelation(correlationId, args...)
}

// Logs a message together with custom values if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprint].
func FatalCustom(customValues map[string]any, args ...any) {
	getLoggers().FatalCustom(customValues, args...)
}

// Logs a message derived from format if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func Fatalf(format string, args ...any) {
	getLoggers().Fatalf(format, args...)
}

// Logs a message derived from format together with a correlation id if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalWithCorrelationf(correlationId string, format string, args ...any) {
	getLoggers().FatalWithCorrelationf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if fatal level is enabled
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalCustomf(customValues map[string]any, format string, args ...any) {
	getLoggers().FatalCustomf(customValues, format, args...)
}

// Logs a message if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func FatalWithPanic(args ...any) {
	getLoggers().FatalWithPanic(args...)
}

// Logs a message together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func FatalWithCorrelationAndPanic(correlationId string, args ...any) {
	getLoggers().FatalWithCorrelationAndPanic(correlationId, args...)
}

// Logs a message together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func FatalCustomWithPanic(customValues map[string]any, args ...any) {
	getLoggers().FatalCustomWithPanic(customValues, args...)
}

// Logs a message derived from format if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalWithPanicf(format string, args ...any) {
	getLoggers().FatalWithPanicf(format, args...)
}

// Logs a message derived from format together with a correlation id if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalWithCorrelationAndPanicf(correlationId string, format string, args ...any) {
	getLoggers().FatalWithCorrelationAndPanicf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if fatal level is enabled and calls built-in function panic to stop current goroutine (independent if error level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalCustomWithPanicf(customValues map[string]any, format string, args ...any) {
	getLoggers().FatalCustomWithPanicf(customValues, format, args...)
}

// Logs a message if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func FatalWithExit(args ...any) {
	getLoggers().FatalWithExit(args...)
}

// Logs a message together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func FatalWithCorrelationAndExit(correlationId string, args ...any) {
	getLoggers().FatalWithCorrelationAndExit(correlationId, args...)
}

// Logs a message together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprint].
func FatalCustomWithExit(customValues map[string]any, args ...any) {
	getLoggers().FatalCustomWithExit(customValues, args...)
}

// Logs a message derived from format if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalWithExitf(format string, args ...any) {
	getLoggers().FatalWithExitf(format, args...)
}

// Logs a message derived from format together with a correlation id if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalWithCorrelationAndExitf(correlationId string, format string, args ...any) {
	getLoggers().FatalWithCorrelationAndExitf(correlationId, format, args...)
}

// Logs a message derived from format together with custom values if fatal level is enabled and calls [os.Exit](1) (independent if fatal level is enabled)
// Arguments are handled in the manner of [fmt.Sprintf].
func FatalCustomWithExitf(customValues map[string]any, format string, args ...any) {
	getLoggers().FatalCustomWithExitf(customValues, format, args...)
}

// Indicator whether debug level is enabled or not
func IsDebugEnabled() bool {
	return getLoggers().IsDebugEnabled()
}

// Indicator whether information level is enabled or not
func IsInformationEnabled() bool {
	return getLoggers().IsInformationEnabled()
}

// Indicator whether warning level is enabled or not
func IsWarningEnabled() bool {
	return getLoggers().IsWarningEnabled()
}

// Indicator whether error level is enabled or not
func IsErrorEnabled() bool {
	return getLoggers().IsErrorEnabled()
}

// Indicator whether fatal level is enabled or not
func IsFatalEnabled() bool {
	return getLoggers().IsFatalEnabled()
}

// Resets loggers. Configuration will be loaded and loggers will be created again.
// The registered custom appenders, formatters and their configurations will be reset also
func Reset() {
	ResetRegisteredAppenderAndFormatters()
	config.ResetRegisteredAppenderAndFormatterConfigs()
}

// Registers an appender and its configuration constructor functions: appenderCreator and appenderConfigCreator.
// The 'appenderType' contains the name which can be referenced by the configuration entry TYPEWRITER_LOG_APPENDER_TYPE or the package variant at environment or file.
// To provide all relevant configuration entries from environment or file their key prefix has to be set by 'keyPrefixes'
func RegisterAppenderWithConfig(appenderType string, keyPrefixes []string,
	appenderCreator func(appenderConfig *config.AppenderConfig, formatter *format.Formatter) (*appender.Appender, error),
	appenderConfigCreator func(relevantKeyValues *map[string]string, commonConfig *config.CommonAppenderConfig) (*config.AppenderConfig, error)) error {

	err := config.RegisterAppenderConfig(appenderType, keyPrefixes, appenderConfigCreator)
	if err != nil {
		return err
	}

	err = RegisterAppender(appenderType, appenderCreator)
	if err != nil {
		config.DeregisterAppenderConfig(appenderType)
		return err
	}
	return nil
}

// Removes an appender and its configuration constructor functions from registration
func DeregisterAppenderTogetherWithConfig(appenderType string) error {
	err := config.DeregisterAppenderConfig(appenderType)
	if err != nil {
		return err
	}
	return DeregisterAppender(appenderType)
}

// Registers a formatter and its configuration constructor functions: formatterCreator and formatterConfigCreator.
// The 'formatterType' contains the name which can be referenced by the configuration entry TYPEWRITER_LOG_FORMATTER_TYPE or the package variant at environment or file.
// To provide all relevant configuration entries from environment or file their key prefix has to be set by 'keyPrefixes'
func RegisterFormatterWithConfig(formatterType string, keyPrefixes []string,
	formatterCreator func(formatterConfig *config.FormatterConfig) (*format.Formatter, error),
	formatterConfigCreator func(relevantKeyValues *map[string]string, commonConfig *config.CommonFormatterConfig) (*config.FormatterConfig, error)) error {

	err := config.RegisterFormatterConfig(formatterType, keyPrefixes, formatterConfigCreator)
	if err != nil {
		return err
	}

	err = RegisterFormatter(formatterType, formatterCreator)
	if err != nil {
		config.DeregisterFormatterConfig(formatterType)
		return err
	}
	return nil
}

// Removes an formatter and its configuration constructor functions from registration
func DeregisterFormatterTogetherWithConfig(formatterType string) error {
	err := config.DeregisterFormatterConfig(formatterType)
	if err != nil {
		return err
	}
	return DeregisterFormatter(formatterType)
}

// Closes all appenders and marks the loggers and configuration as uninitialized
func Close() {
	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	loggersInitialized = false
	config.ClearConfig()

	for _, a := range appenders {
		a.Close()
	}
}
