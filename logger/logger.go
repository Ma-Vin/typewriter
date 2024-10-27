package logger

type Logger interface {
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
