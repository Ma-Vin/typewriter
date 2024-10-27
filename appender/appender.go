package appender

type Appender interface {
	// Writes the given message
	Write(severity int, message string)
	// Writes the given message and correlation id
	WriteWithCorrelation(severity int, correlationId string, message string)
	// Writes the given message and a map of custom values
	WriteCustom(severity int, message string, customValues map[string]any)
}
