package appender

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/format"
	"github.com/ma-vin/typewriter/testutil"
)

var testJsonFormatter = format.CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339, false)
var jsonFormatTestTime = time.Date(2024, time.November, 18, 16, 00, 0, 0, time.UTC)
var jsonFormatTestTimeText = jsonFormatTestTime.Format(time.RFC3339Nano)

func getAppenderTestLogFile(testCase string) string {
	SkipFileCreationForTest = false
	_, filename, _, _ := runtime.Caller(0)
	result := strings.Replace(filename, ".go", "_"+testCase+"_scratch.log", 1)
	os.Create(result)
	return result
}

func TestFileAppenderWrite(t *testing.T) {
	logFilePath := getAppenderTestLogFile("write")
	format.SetFormatterMockTime(&jsonFormatTestTime)
	appender := CreateFileAppender(logFilePath, &testJsonFormatter).(FileAppender)

	appender.Write(constants.INFORMATION_SEVERITY, "Testmessage")
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteWithCorrelation(t *testing.T) {
	logFilePath := getAppenderTestLogFile("correlation")
	format.SetFormatterMockTime(&jsonFormatTestTime)
	appender := CreateFileAppender(logFilePath, &testJsonFormatter).(FileAppender)

	appender.WriteWithCorrelation(constants.INFORMATION_SEVERITY, "someCorrelationId", "Testmessage")
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteCustom(t *testing.T) {
	logFilePath := getAppenderTestLogFile("custom")
	format.SetFormatterMockTime(&jsonFormatTestTime)
	appender := CreateFileAppender(logFilePath, &testJsonFormatter).(FileAppender)

	customProperties := map[string]any{
		"first": "abc",
	}

	appender.WriteCustom(constants.INFORMATION_SEVERITY, "Testmessage", customProperties)
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
