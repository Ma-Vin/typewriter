package appender

import (
	"time"
)

type Appender interface {
	// Writes the given message
	Write(timestamp time.Time, severity int, message string)
	// Writes the given message and correlation id
	WriteWithCorrelation(timestamp time.Time, severity int, correlationId string, message string)
	// Writes the given message and a map of custom values
	WriteCustom(timestamp time.Time, severity int, message string, customValues map[string]any)
}
