package test

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/logger"
	"github.com/ma-vin/typewriter/testutil"
)

const (
	minutesToRun   int = 2
	goRoutineCount int = 30
)

func TestFileAppenderCronRenameLongRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	logFilePath := testutil.GetTestCaseFilePath("longRun", true)

	os.Clearenv()
	logger.Reset()
	os.Setenv("TYPEWRITER_LOG_LEVEL", "DEBUG")
	os.Setenv("TYPEWRITER_LOG_APPENDER_TYPE", "FILE")
	os.Setenv("TYPEWRITER_LOG_APPENDER_FILE", logFilePath)
	os.Setenv("TYPEWRITER_LOG_APPENDER_CRON_RENAMING", "* * * * *")

	c := make(chan []int, goRoutineCount)

	waitForStartTime()
	fmt.Println("start go routines")
	for i := range goRoutineCount {
		go logForFileAppenderCronRename(i, c)
	}

	var logEntryCount int
	for i := range goRoutineCount {
		r := <-c
		fmt.Println(i, "thread", r[0], "with", r[1], "log entries done")
		logEntryCount += r[1]
	}

	logFilePaths := testutil.GetExistingTestCaseFilePaths("longRun")

	testutil.AssertEquals(minutesToRun+1, len(logFilePaths), t, "len(logFilePaths)")

	lineCount := 0
	for _, filePath := range logFilePaths {
		file, err := os.Open(filePath)
		testutil.AssertNil(err, t, "os.Open(filePath)")
		fileScanner := bufio.NewScanner(file)
		lineCountPerFile := 0
		for fileScanner.Scan() {
			lineCountPerFile++
		}
		testutil.AssertTrue(lineCountPerFile > 0, t, "lineCount positive at "+filePath)
		lineCount += lineCountPerFile
	}

	testutil.AssertEquals(logEntryCount, lineCount, t, "lineCount")
}

func waitForStartTime() {
	second := time.Now().Second()
	if second > 0 && second < 10 {
		return
	}
	fmt.Println("Wait", 70-second, "seconds for starting point of ten seconds after each minute")
	time.Sleep(time.Duration(70-second) * time.Second)
}

func logForFileAppenderCronRename(thread int, c chan []int) {
	millisPerIteration := 101 + rand.IntN(900)
	iterations := calcMillisToRun() / millisPerIteration

	fmt.Println("thread", thread, "with", iterations, "iterations every", millisPerIteration, "millis")

	for i := range iterations {
		level := rand.IntN(9) + 1
		switch level {
		case common.DEBUG_SEVERITY:
			logger.Debug("some debug message", thread, i)
		case common.INFORMATION_SEVERITY:
			logger.Information("some info message", thread, i)
		case common.WARNING_SEVERITY:
			logger.Warning("some warn message", thread, i)
		case common.ERROR_SEVERITY:
			logger.Error("some error message", thread, i)
		case common.FATAL_SEVERITY:
			logger.Fatal("some fatal message", thread, i)
		case common.DEBUG_SEVERITY + 5:
			logger.DebugWithCorrelation(strconv.Itoa(thread), "some debug message with correlation", i)
		case common.INFORMATION_SEVERITY + 5:
			logger.InformationWithCorrelation(strconv.Itoa(thread), "some info message with correlation", i)
		case common.WARNING_SEVERITY + 5:
			logger.WarningWithCorrelation(strconv.Itoa(thread), "some warn message with correlation", i)
		case common.ERROR_SEVERITY + 5:
			logger.ErrorWithCorrelation(strconv.Itoa(thread), "some error message with correlation", i)
		case common.FATAL_SEVERITY + 5:
			logger.FatalWithCorrelation(strconv.Itoa(thread), "some fatal message with correlation", i)
		}

		time.Sleep(time.Duration(millisPerIteration) * time.Millisecond)
	}

	c <- []int{thread, iterations}
}

func calcMillisToRun() int {
	return minutesToRun * 60 * 1000
}
