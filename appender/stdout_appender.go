package appender

import (
	"fmt"
	"io"
	"os"

	"github.com/ma-vin/typewriter/format"
)

// Writes message, formated by a given formatter, to the standard output
type StandardOutputAppender struct {
	formatter *format.Formatter
	writer    io.Writer
}

// Creates a standard output appender with a given formatter
func CreateStandardOutputAppender(formatter *format.Formatter) Appender {
	return StandardOutputAppender{formatter, os.Stdout}
}

// Writes the given message to the standard output
func (s StandardOutputAppender) Write(severity int, message string) {
	fmt.Fprintln(s.writer, (*s.formatter).Format(severity, message))
}

// Writes the given message and correlation id to the standard output
func (s StandardOutputAppender) WriteWithCorrelation(severity int, correlationId string, message string) {
	fmt.Fprintln(s.writer, (*s.formatter).FormatWithCorrelation(severity, correlationId, message))
}

// Writes the given message and a map of custom values to the standard output
func (s StandardOutputAppender) WriteCustom(severity int, message string, customValues map[string]any) {
	fmt.Fprintln(s.writer, (*s.formatter).FormatCustom(severity, message, customValues))
}

// Does nothing, but has to be declared to fullfill the interface
func (s StandardOutputAppender) Close() {
	// Nothing to do. Closing os.Stdout may cause errors elsewhere: See documentation of os.Stdout
}
