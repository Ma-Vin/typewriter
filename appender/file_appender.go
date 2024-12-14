package appender

import (
	"fmt"
	"os"
	"sync"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
)

// Writes message, formated by a given formatter, to the defined output file
type FileAppender struct {
	pathToLogFile string
	formatter     *format.Formatter
	writer        *os.File
	isClosed      *bool
	mu            *sync.Mutex
}

type fileMutex struct {
	pathToLogFile *string
	mu            *sync.Mutex
}

// For test usage only! Indicator whether to skip creation of the target output file
var SkipFileCreationForTest = false
var fileToMutex = []fileMutex{}

// Removes all registered mutex for log file writing
func CleanFileToMutex() {
	fileToMutex = []fileMutex{}
}

// Creates a file appender to the file at “pathToLogFile” with a given formatter
func CreateFileAppender(pathToLogFile string, formatter *format.Formatter) Appender {
	var file *os.File = nil
	var err error = nil
	var closed = false
	if !SkipFileCreationForTest {
		file, err = os.OpenFile(pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	}
	if err != nil {
		fmt.Printf("Fail to create file appender, use stdout instead: %s", err)
		fmt.Println()
		return CreateStandardOutputAppender(formatter)
	}
	return FileAppender{pathToLogFile, formatter, file, &closed, getOrCreateMutexForFile(&pathToLogFile)}
}

func getOrCreateMutexForFile(pathToLogFile *string) *sync.Mutex {
	for _, fm := range fileToMutex {
		if *pathToLogFile == *fm.pathToLogFile {
			return fm.mu
		}
	}
	result := sync.Mutex{}
	fileToMutex = append(fileToMutex, fileMutex{pathToLogFile, &result})
	return &result
}

// Writes the given logValues to the defined output file
func (f FileAppender) Write(logValues *common.LogValues) {
	f.writeRecord((*f.formatter).Format(logValues))
}

// Closes writer of the output file
func (f FileAppender) Close() {
	if *f.isClosed {
		return
	}
	err := f.writer.Close()
	*f.isClosed = true
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
