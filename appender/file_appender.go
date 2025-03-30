package appender

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
)

// Writes message, formatted by a given formatter, to the defined output file
type FileAppender struct {
	pathToLogFile string
	formatter     *format.Formatter
	cronRenamer   *CronFileRenamer
	sizeRenamer   *SizeFileRenamer
	writer        *os.File
	isClosed      *bool
	mu            *sync.Mutex
}

type fileDeduction struct {
	pathToLogFile *string
	mu            *sync.Mutex
	cronRenamer   *CronFileRenamer
	sizeRenamer   *SizeFileRenamer
}

// For test usage only! Indicator whether to skip creation of the target output file
var SkipFileCreationForTest = false
var fileAppenderMu = sync.Mutex{}
var fileDeductions = []fileDeduction{}

// Removes all registered mutex, renamer for log file writing
func CleanFileDeductions() {
	fileDeductions = []fileDeduction{}
}

// Creates a file appender to the file at “pathToLogFile” with a given formatter
func CreateFileAppender(pathToLogFile string, formatter *format.Formatter, cronExpression string, limitByteSize string) Appender {
	var file *os.File = nil
	var err error = nil
	var closed = false
	if !SkipFileCreationForTest {
		fileAppenderMu.Lock()
		defer fileAppenderMu.Unlock()
		file, err = os.OpenFile(pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	}
	if err != nil {
		fmt.Printf("Fail to create file appender, use stdout instead: %s", err)
		fmt.Println()
		return CreateStandardOutputAppender(formatter)
	}

	mu, cronRenamer, sizeRenamer := getOrCreateDeductionForFile(&pathToLogFile, file, cronExpression, limitByteSize)

	return FileAppender{pathToLogFile, formatter, cronRenamer, sizeRenamer, file, &closed, mu}
}

// gets an existing struct of deduced elements from pathToLogFile or creates a new one
func getOrCreateDeductionForFile(pathToLogFile *string, file *os.File, cronExpression string, limitByteSize string) (*sync.Mutex, *CronFileRenamer, *SizeFileRenamer) {
	for _, fd := range fileDeductions {
		if *pathToLogFile == *fd.pathToLogFile {
			return fd.mu, fd.cronRenamer, fd.sizeRenamer
		}
	}
	mu := sync.Mutex{}
	cronRenamer := createCrontabAndRenamer(*pathToLogFile, file, cronExpression, &mu)
	sizeRenamer := createSizeRenamer(*pathToLogFile, file, limitByteSize, &mu, cronRenamer)
	fileDeductions = append(fileDeductions, fileDeduction{pathToLogFile, &mu, cronRenamer, sizeRenamer})
	return &mu, cronRenamer, sizeRenamer
}

// creates a new cron renamer and its crontab
func createCrontabAndRenamer(pathToLogFile string, file *os.File, cronExpression string, mu *sync.Mutex) *CronFileRenamer {
	if cronExpression == "" {
		return nil
	}
	crontab := common.CreateCrontab(cronExpression)
	return CreateCronFileRenamer(pathToLogFile, file, crontab, mu)
}

// creates a new size renamer
func createSizeRenamer(pathToLogFile string, file *os.File, limitByteSizeText string, mu *sync.Mutex, cronRenamer *CronFileRenamer) *SizeFileRenamer {
	if cronRenamer != nil || limitByteSizeText == "" {
		return nil
	}

	sizeFactor := 1
	if strings.HasSuffix(limitByteSizeText, "kb") || strings.HasSuffix(limitByteSizeText, "KB") {
		limitByteSizeText = strings.TrimSpace(limitByteSizeText[:len(limitByteSizeText)-2])
		sizeFactor = 1000
	} else if strings.HasSuffix(limitByteSizeText, "mb") || strings.HasSuffix(limitByteSizeText, "MB") {
		limitByteSizeText = strings.TrimSpace(limitByteSizeText[:len(limitByteSizeText)-2])
		sizeFactor = 1000000
	}

	limitByteSize, err := strconv.Atoi(limitByteSizeText)
	if err != nil {
		fmt.Printf("Fail to parse byte size limit for log file renaming %s: %s", limitByteSizeText, err)
		fmt.Println()
		return nil
	}
	return CreateSizeFileRenamer(pathToLogFile, file, int64(limitByteSize*sizeFactor), mu)
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
	if f.sizeRenamer != nil {
		f.sizeRenamer.CheckFile(formattedRecord)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	fmt.Fprintln(f.writer, formattedRecord)
}
