package appender

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
	"github.com/ma-vin/typewriter/testutil"
)

var testJsonFormatter = format.CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339, "caller", "file", "line", false)
var jsonFormatTestTime = time.Date(2024, time.November, 18, 16, 00, 0, 0, time.UTC)
var jsonFormatTestTimeText = jsonFormatTestTime.Format(time.RFC3339Nano)

func getAppenderTestLogFile(testCase string) string {
	SkipFileCreationForTest = false
	_, filename, _, _ := runtime.Caller(0)
	result := strings.Replace(filename, ".go", "_"+testCase+"_scratch.log", 1)
	os.Create(result)
	return result
}

func TestCreateFileAppenderDifferentLogFilePaths(t *testing.T) {
	SkipFileCreationForTest = true
	CleanFileToMutex()
	appender1 := CreateFileAppender("Path1", &testJsonFormatter).(FileAppender)
	appender2 := CreateFileAppender("Path2", &testJsonFormatter).(FileAppender)

	testutil.AssertEquals(2, len(fileToMutex), t, "len(fileToMutex)")
	testutil.AssertNotEquals(appender1.mu, appender2.mu, t, "mu")
}

func TestCreateFileAppenderEqualLogFilePaths(t *testing.T) {
	SkipFileCreationForTest = true
	CleanFileToMutex()
	appender1 := CreateFileAppender("PathEqual", &testJsonFormatter).(FileAppender)
	appender2 := CreateFileAppender("PathEqual", &testJsonFormatter).(FileAppender)

	testutil.AssertEquals(1, len(fileToMutex), t, "len(fileToMutex)")
	testutil.AssertEquals(appender1.mu, appender2.mu, t, "mu")
}

func TestFileAppenderWrite(t *testing.T) {
	logFilePath := getAppenderTestLogFile("write")
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	appender := CreateFileAppender(logFilePath, &testJsonFormatter).(FileAppender)

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	testutil.AssertFalse(*appender.isClosed, t, "isNotClosed")
	appender.Write(&logValuesToFormat)
	appender.Close()
	testutil.AssertTrue(*appender.isClosed, t, "isClosed")
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteWithCorrelation(t *testing.T) {
	logFilePath := getAppenderTestLogFile("correlation")
	common.SetLogValuesMockTime(&jsonFormatTestTime)
	correlation := "someCorrelationId"

	appender := CreateFileAppender(logFilePath, &testJsonFormatter).(FileAppender)

	logValuesToFormat := common.CreateLogValuesWithCorrelation(common.INFORMATION_SEVERITY, &correlation, "Testmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteCustom(t *testing.T) {
	logFilePath := getAppenderTestLogFile("custom")
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	appender := CreateFileAppender(logFilePath, &testJsonFormatter).(FileAppender)

	customProperties := map[string]any{
		"first": "abc",
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)
	appender.Write(&logValuesToFormat)
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func checkLogFileEntry(logFilePath string, entry string, t *testing.T) {
	logFile, err := os.Open(logFilePath)
	testutil.AssertNil(err, t, "open logFile err")
	defer logFile.Close()
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		testutil.AssertEquals(entry, scanner.Text(), t, "logFile line")
	}
}
