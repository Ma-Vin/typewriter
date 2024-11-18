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

func CreateStandardOutputAppender(formatter *format.Formatter) Appender {
	return StandardOutputAppender{formatter, os.Stdout}
}

func (s StandardOutputAppender) Write(severity int, message string) {
	fmt.Fprintln(s.writer, (*s.formatter).Format(severity, message))
}

func (s StandardOutputAppender) WriteWithCorrelation(severity int, correlationId string, message string) {
	fmt.Fprintln(s.writer, (*s.formatter).FormatWithCorrelation(severity, correlationId, message))
}

func (s StandardOutputAppender) WriteCustom(severity int, message string, customValues map[string]any) {
	fmt.Fprintln(s.writer, (*s.formatter).FormatCustom(severity, message, customValues))
}

func (s StandardOutputAppender) Close() {
	// Nothing to do. closing os.Stdout may cause errors elsewhere: See documentation of os.Stdout
}
