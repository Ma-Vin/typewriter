package appender

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/ma-vin/typewriter/common"
)

type CronFileRenamer struct {
	pathToLogFile         string
	writer                *os.File
	crontab               *common.Crontab
	timeFileNameGenerator *TimeFileNameGenerator
	mu                    *sync.Mutex
}

type SizeFileRenamer struct {
	pathToLogFile         string
	writer                *os.File
	limitByteSize         int64
	currentByteSize       int64
	timeFileNameGenerator *TimeFileNameGenerator
	mu                    *sync.Mutex
}

type TimeFileNameGenerator struct {
	basePath      string
	fileEnding    string
	referenceTime *time.Time
}

var newLineUtf8Size int = utf8.RuneCountInString(fmt.Sprintln())

// Creates a new CronFileNamer for a given path and crontab
func CreateCronFileRenamer(pathToLogFile string, writer *os.File, crontab *common.Crontab, mu *sync.Mutex) *CronFileRenamer {
	return &CronFileRenamer{pathToLogFile, writer, crontab, CreateTimeFileNameGenerator(pathToLogFile), mu}
}

// Checks whether the next time of crontab is reached or not. In positive case the current file will be renamed to a name given by filename generator.
func (c *CronFileRenamer) CheckFile(logValues *common.LogValues) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if logValues.Time.Before(*c.crontab.NextTime) {
		return
	}
	renameLogFile(&c.pathToLogFile, c.writer, c.timeFileNameGenerator, c.prepareNextInterval)
}

func (c *CronFileRenamer) prepareNextInterval() {
	c.timeFileNameGenerator.referenceTime = c.crontab.NextTime
	c.crontab.CalculateNextTime()
}

// Creates a new SizeFileRenamer for a given path and size limit
func CreateSizeFileRenamer(pathToLogFile string, writer *os.File, limitByteSize int64, mu *sync.Mutex) *SizeFileRenamer {
	stat, err := os.Stat(pathToLogFile)
	var currentSize int64 = 0
	if err == nil {
		currentSize = stat.Size()
	}
	return &SizeFileRenamer{pathToLogFile, writer, limitByteSize, currentSize, CreateTimeFileNameGenerator(pathToLogFile), mu}
}

// Checks whether the size limit is reached or not. In positive case the current file will be renamed to a name given by filename generator.
func (c *SizeFileRenamer) CheckFile(formattedRecord string) {
	sizeToAdd := int64(utf8.RuneCountInString(formattedRecord) + newLineUtf8Size)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currentByteSize+sizeToAdd < c.limitByteSize {
		c.currentByteSize += sizeToAdd
		return
	}

	renameLogFile(&c.pathToLogFile, c.writer, c.timeFileNameGenerator, c.prepareNextInterval)
	c.currentByteSize = sizeToAdd
}

func (c *SizeFileRenamer) prepareNextInterval() {
	referenceTime := common.GetNow()
	c.timeFileNameGenerator.referenceTime = &referenceTime
	c.currentByteSize = 0
}

func renameLogFile(pathToLogFile *string, writer *os.File, timeFileNameGenerator *TimeFileNameGenerator, prepareNextInterval func()) {

	if _, err := os.Stat(*pathToLogFile); err != nil {
		return
	}

	newPath := timeFileNameGenerator.determineNextPathToLogFile()

	err := writer.Close()
	if err != nil {
		fmt.Println("Failed to close log file before renaming from", pathToLogFile, "to", newPath)
		prepareNextInterval()
		return
	}

	err = os.Rename(*pathToLogFile, newPath)
	if err != nil {
		fmt.Println("Failed to rename log file from", pathToLogFile, "to", newPath)
		prepareNextInterval()
		return
	}

	if !SkipFileCreationForTest {
		file, err := os.OpenFile(*pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			*writer = *file
		} else {
			fmt.Println("Failed to create new log file from", pathToLogFile, " after renaming to", newPath)
		}
	}
	prepareNextInterval()
}

// Creates a new TimeFileNameGenerator
func CreateTimeFileNameGenerator(pathToLogFile string) *TimeFileNameGenerator {
	indexOfFileEnding := strings.LastIndex(pathToLogFile, ".")
	refTime := common.GetNow()
	return &TimeFileNameGenerator{pathToLogFile[:indexOfFileEnding], pathToLogFile[indexOfFileEnding+1:], &refTime}
}

// creates the next path of log file. If the file exists a count from 1 to 10 will be tried to append
func (t *TimeFileNameGenerator) determineNextPathToLogFile() string {
	withoutFileEnding := fmt.Sprintf("%s_%d%s%s_%s%s%s", t.basePath,
		t.referenceTime.Year(), determineTwoDigits(int(t.referenceTime.Month())), determineTwoDigits(t.referenceTime.Day()),
		determineTwoDigits(t.referenceTime.Hour()), determineTwoDigits(t.referenceTime.Minute()), determineTwoDigits(t.referenceTime.Second()))

	result := withoutFileEnding + "." + t.fileEnding

	for i := range 10 {
		if _, err := os.Stat(result); err != nil {
			break
		}
		result = fmt.Sprintf("%s_%d.%s", withoutFileEnding, i+1, t.fileEnding)
	}
	return result
}

// puts a zero at front if the number is lower than ten
func determineTwoDigits(number int) string {
	if number < 10 {
		return "0" + strconv.Itoa(number)
	}
	return strconv.Itoa(number)
}
