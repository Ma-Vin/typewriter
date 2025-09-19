package appender

import (
	"fmt"
	"io"
	"os"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
)

// Writes message, formatted by a given formatter, to the standard output
type StandardOutputAppender struct {
	formatter *format.Formatter
	writer    io.Writer
	isClosed  *bool
}

// Creates a standard output appender with a given formatter
func CreateStandardOutputAppenderFromConfig(config *config.AppenderConfig, formatter *format.Formatter) (*Appender, error) {
	var closed = false
	var result Appender = StandardOutputAppender{formatter, os.Stdout, &closed}
	return &result, nil
}

// Writes the given logValues to the standard output
func (s StandardOutputAppender) Write(logValues *common.LogValues) {
	if *s.isClosed {
		return
	}
	fmt.Fprintln(s.writer, (*s.formatter).Format(logValues))
}

// Does nothing, but has to be declared to fulfill the interface
func (s StandardOutputAppender) Close() {
	*s.isClosed = true
	// Nothing to do. Closing os.Stdout may cause errors elsewhere: See documentation of os.Stdout
}
