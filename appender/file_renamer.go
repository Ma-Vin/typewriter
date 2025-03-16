package appender

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ma-vin/typewriter/common"
)

type CronFileNamer struct {
	pathToLogFile         string
	writer                *os.File
	crontab               *common.Crontab
	timeFileNameGenerator *TimeFileNameGenerator
	mu                    *sync.Mutex
}

type TimeFileNameGenerator struct {
	basePath      string
	fileEnding    string
	referenceTime *time.Time
}

// Creates a new CronFileNamer for a given path and crontab
func CreateCronFileRenamer(pathToLogFile string, writer *os.File, crontab *common.Crontab, mu *sync.Mutex) *CronFileNamer {
	indexOfFileEnding := strings.LastIndex(pathToLogFile, ".")
	refTime := time.Now()
	fileNameCreator := TimeFileNameGenerator{pathToLogFile[:indexOfFileEnding], pathToLogFile[indexOfFileEnding+1:], &refTime}
	return &CronFileNamer{pathToLogFile, writer, crontab, &fileNameCreator, mu}
}

// Checks whether the next time of crontab is reached or not. In positive case the current file will be renamed to a name given by filename generator.
func (c *CronFileNamer) CheckFile(logValues *common.LogValues) {
	if logValues.Time.Before(*c.crontab.NextTime) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(c.pathToLogFile); err != nil {
		return
	}

	newPath := c.timeFileNameGenerator.determineNextPathToLogFile()
	err := os.Rename(c.pathToLogFile, newPath)
	if err != nil {
		fmt.Println("Failed to rename log file from", c.pathToLogFile, "to", newPath)
	} else if !SkipFileCreationForTest {
		file, err := os.OpenFile(c.pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			*c.writer = *file
		}
	}

	c.timeFileNameGenerator.referenceTime = c.crontab.NextTime
	c.crontab.CalculateNextTime()
}

// creates the next path of log file
func (t *TimeFileNameGenerator) determineNextPathToLogFile() string {
	return fmt.Sprintf("%s_%d%s%s_%s%s%s.%s", t.basePath,
		t.referenceTime.Year(), determineTwoDigits(int(t.referenceTime.Month())), determineTwoDigits(t.referenceTime.Day()),
		determineTwoDigits(t.referenceTime.Hour()), determineTwoDigits(t.referenceTime.Minute()), determineTwoDigits(t.referenceTime.Second()),
		t.fileEnding)
}

// puts a zero at front if the number is lower than ten
func determineTwoDigits(number int) string {
	if number < 10 {
		return "0" + strconv.Itoa(number)
	}
	return strconv.Itoa(number)
}
