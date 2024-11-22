package appender

import (
	"fmt"
	"os"
	"sync"

	"github.com/ma-vin/typewriter/format"
)

// Writes message, formated by a given formatter, to the defined output file
type FileAppender struct {
	pathToLogFile string
	formatter     *format.Formatter
	writer        *os.File
	mu            *sync.Mutex
}

// For test usage only! Indicator whether to skip creation of the target output file
var SkipFileCreationForTest = false

// Creates a file appender to the file at “pathToLogFile” with a given formatter
func CreateFileAppender(pathToLogFile string, formatter *format.Formatter) Appender {
	var file *os.File = nil
	var err error = nil
	if !SkipFileCreationForTest {
		file, err = os.OpenFile(pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	}
	if err != nil {
		fmt.Printf("Fail to create file appender, use stdout instead: %s", err)
		fmt.Println()
		return CreateStandardOutputAppender(formatter)
	}
	return FileAppender{pathToLogFile, formatter, file, &sync.Mutex{}}
}

// Writes the given message to the defined output file
func (f FileAppender) Write(severity int, message string) {
	f.writeRecord((*f.formatter).Format(severity, message))
}

// Writes the given message and correlation id to the defined output file
func (f FileAppender) WriteWithCorrelation(severity int, correlationId string, message string) {
	f.writeRecord((*f.formatter).FormatWithCorrelation(severity, correlationId, message))
}

// Writes the given message and a map of custom values to the defined output file
func (f FileAppender) WriteCustom(severity int, message string, customValues map[string]any) {
	f.writeRecord((*f.formatter).FormatCustom(severity, message, customValues))
}

// Closes writer of the output file
func (f FileAppender) Close() {
	err := f.writer.Close()
	if err != nil {
		fmt.Printf("Fail to close writer of %s: %s", f.pathToLogFile, err)
		fmt.Println()
	}
}

func (f *FileAppender) writeRecord(formattedRecord string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fmt.Fprintln(f.writer, formattedRecord)
}
