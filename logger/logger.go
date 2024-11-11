package logger

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
	FatalWithCorrelationAndExitf(fcorrelationId string, ormat string, args ...any)

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
