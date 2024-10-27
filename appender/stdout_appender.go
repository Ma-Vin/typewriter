package appender

import (
	"fmt"
	"io"
	"os"
	"time"

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

func (s StandardOutputAppender) Write(timestamp time.Time, severity int, message string) {
	fmt.Fprintln(s.writer, (*s.formatter).Format(timestamp, severity, message))
}

func (s StandardOutputAppender) WriteWithCorrelation(timestamp time.Time, severity int, correlationId string, message string) {
	fmt.Fprintln(s.writer, (*s.formatter).FormatWithCorrelation(timestamp, severity, correlationId, message))
}

func (s StandardOutputAppender) WriteCustom(timestamp time.Time, severity int, message string, customValues map[string]any) {
	fmt.Fprintln(s.writer, (*s.formatter).FormatCustom(timestamp, severity, message, customValues))
}
