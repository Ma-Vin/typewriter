package appender

import (
	"fmt"
	"io"
	"os"

	"github.com/ma-vin/typewriter/common"
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

// Writes the given logValues to the standard output
func (s StandardOutputAppender) Write(logValues *common.LogValues) {
	fmt.Fprintln(s.writer, (*s.formatter).Format(logValues))
}

// Does nothing, but has to be declared to fullfill the interface
func (s StandardOutputAppender) Close() {
	// Nothing to do. Closing os.Stdout may cause errors elsewhere: See documentation of os.Stdout
}
