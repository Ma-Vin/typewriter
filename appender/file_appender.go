package appender

import (
	"fmt"
	"os"
	"sync"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
)

// Writes message, formatted by a given formatter, to the defined output file
type FileAppender struct {
	pathToLogFile string
	formatter     *format.Formatter
	cronRenamer   *CronFileRenamer
	writer        *os.File
	isClosed      *bool
	mu            *sync.Mutex
}

type fileDeduction struct {
	pathToLogFile *string
	mu            *sync.Mutex
	cronRenamer   *CronFileRenamer
}

// For test usage only! Indicator whether to skip creation of the target output file
var SkipFileCreationForTest = false
var fileDeductions = []fileDeduction{}

// Removes all registered mutex, renamer for log file writing
func CleanFileDeductions() {
	fileDeductions = []fileDeduction{}
}

// Creates a file appender to the file at “pathToLogFile” with a given formatter
func CreateFileAppender(pathToLogFile string, formatter *format.Formatter, cronExpression string) Appender {
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

	mu, cronRenamer := getOrCreateDeductionForFile(&pathToLogFile, file, cronExpression)

	return FileAppender{pathToLogFile, formatter, cronRenamer, file, &closed, mu}
}

// gets an existing struct of deduced elements from pathToLogFile or creates a new one
func getOrCreateDeductionForFile(pathToLogFile *string, file *os.File, cronExpression string) (*sync.Mutex, *CronFileRenamer) {
	for _, fd := range fileDeductions {
		if *pathToLogFile == *fd.pathToLogFile {
			return fd.mu, fd.cronRenamer
		}
	}
	mu := sync.Mutex{}
	renamer := createCrontabAndRenamer(*pathToLogFile, file, cronExpression, &mu)
	fileDeductions = append(fileDeductions, fileDeduction{pathToLogFile, &mu, renamer})
	return &mu, renamer
}

// creates a new renamer and its crontab
func createCrontabAndRenamer(pathToLogFile string, file *os.File, cronExpression string, mu *sync.Mutex) *CronFileRenamer {
	if cronExpression == "" {
		return nil
	}
	crontab := common.CreateCrontab(cronExpression)
	return CreateCronFileRenamer(pathToLogFile, file, crontab, mu)
}

// Writes the given logValues to the defined output file and checks whether to rename existing log file or not
func (f FileAppender) Write(logValues *common.LogValues) {
	if f.cronRenamer != nil {
		f.cronRenamer.CheckFile(logValues)
	}
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
