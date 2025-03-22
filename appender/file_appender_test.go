package appender

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
	"github.com/ma-vin/typewriter/testutil"
)

const testCronExpression = "* * * * *"

var testJsonFormatter = format.CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339, "caller", "file", "line", false)
var jsonFormatTestTime = time.Date(2024, time.November, 18, 16, 00, 0, 0, time.UTC)
var jsonFormatTestTimeText = jsonFormatTestTime.Format(time.RFC3339Nano)

func getAppenderTestLogFile(testCase string) string {
	SkipFileCreationForTest = false
	result := testutil.GetTestCaseFilePath(testCase, true)
	os.Create(result)
	return result
}

func TestCreateFileAppenderDifferentLogFilePaths(t *testing.T) {
	SkipFileCreationForTest = true
	CleanFileDeductions()
	appender1 := CreateFileAppender("Path1.log", &testJsonFormatter, testCronExpression, "").(FileAppender)
	appender2 := CreateFileAppender("Path2.log", &testJsonFormatter, testCronExpression, "").(FileAppender)

	testutil.AssertEquals(2, len(fileDeductions), t, "len(fileToMutex)")
	testutil.AssertNotEquals(appender1.mu, appender2.mu, t, "mu")
	testutil.AssertNotEquals(appender1.cronRenamer, appender2.cronRenamer, t, "cronRenamer")
}

func TestCreateFileAppenderEqualLogFilePaths(t *testing.T) {
	SkipFileCreationForTest = true
	CleanFileDeductions()
	appender1 := CreateFileAppender("PathEqual.log", &testJsonFormatter, testCronExpression, "").(FileAppender)
	appender2 := CreateFileAppender("PathEqual.log", &testJsonFormatter, testCronExpression, "").(FileAppender)

	testutil.AssertEquals(1, len(fileDeductions), t, "len(fileToMutex)")
	testutil.AssertEquals(appender1.mu, appender2.mu, t, "mu")
	testutil.AssertEquals(appender1.cronRenamer, appender2.cronRenamer, t, "cronRenamer")
}

func TestFileAppenderWrite(t *testing.T) {
	logFilePath := getAppenderTestLogFile("write")
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	appender := CreateFileAppender(logFilePath, &testJsonFormatter, "", "").(FileAppender)

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	testutil.AssertFalse(*appender.isClosed, t, "isNotClosed")
	appender.Write(&logValuesToFormat)
	appender.Close()
	testutil.AssertTrue(*appender.isClosed, t, "isClosed")
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteCronRenameFile(t *testing.T) {
	logFilePath := getAppenderTestLogFile("writeCronRename")
	checkAppenderWriteCronRename(logFilePath, t, func() FileAppender {
		return CreateFileAppender(logFilePath, &testJsonFormatter, testCronExpression, "").(FileAppender)
	})
}

func TestFileAppenderWriteCronAndSizeRenameFile(t *testing.T) {
	logFilePath := getAppenderTestLogFile("writeCronAndSizeRename")
	checkAppenderWriteCronRename(logFilePath, t, func() FileAppender {
		return CreateFileAppender(logFilePath, &testJsonFormatter, testCronExpression, "76").(FileAppender)
	})
}

func checkAppenderWriteCronRename(logFilePath string, t *testing.T, createFileAppender func() FileAppender) {
	indexOfFileEnding := strings.LastIndex(logFilePath, ".")
	expectedNewFileName := logFilePath[:indexOfFileEnding] + "_20241118_160000" + logFilePath[indexOfFileEnding:]

	common.SetLogValuesMockTime(&jsonFormatTestTime)

	appender := createFileAppender()

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	appender.Write(&logValuesToFormat)

	modifiedTestTime := jsonFormatTestTime.Add(time.Hour)
	common.SetLogValuesMockTime(&modifiedTestTime)

	logValuesToFormat = common.CreateLogValues(common.INFORMATION_SEVERITY, "OtherTestmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()
	testutil.AssertTrue(*appender.isClosed, t, "isClosed")
	appender.Close()

	checkLogFileEntry(expectedNewFileName, "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
	checkLogFileEntry(logFilePath, "{\"message\":\"OtherTestmessage\",\"severity\":\"INFO\",\"time\":\""+modifiedTestTime.Format(time.RFC3339Nano)+"\"}", t)
}

func TestFileAppenderWriteSizeRenameFile(t *testing.T) {
	logFilePath := getAppenderTestLogFile("writeSizeRename")
	indexOfFileEnding := strings.LastIndex(logFilePath, ".")
	expectedNewFileName := logFilePath[:indexOfFileEnding] + "_20241118_160000" + logFilePath[indexOfFileEnding:]

	common.SetLogValuesMockTime(&jsonFormatTestTime)

	expectedFirstLogEntry := "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}"
	appender := CreateFileAppender(logFilePath, &testJsonFormatter, "", strconv.Itoa(len(expectedFirstLogEntry)+2)).(FileAppender)

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	appender.Write(&logValuesToFormat)
	logValuesToFormat = common.CreateLogValues(common.INFORMATION_SEVERITY, "OtherTestmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()
	testutil.AssertTrue(*appender.isClosed, t, "isClosed")
	appender.Close()

	checkLogFileEntry(expectedNewFileName, expectedFirstLogEntry, t)
	checkLogFileEntry(logFilePath, "{\"message\":\"OtherTestmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteSizeRenameFileInvalid(t *testing.T) {
	logFilePath := getAppenderTestLogFile("writeInvalidSizeRename")
	indexOfFileEnding := strings.LastIndex(logFilePath, ".")
	expectedNewFileName := logFilePath[:indexOfFileEnding] + "_20241118_160000" + logFilePath[indexOfFileEnding:]

	common.SetLogValuesMockTime(&jsonFormatTestTime)

	appender := CreateFileAppender(logFilePath, &testJsonFormatter, "", "invalidNumber").(FileAppender)

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	appender.Write(&logValuesToFormat)
	logValuesToFormat = common.CreateLogValues(common.INFORMATION_SEVERITY, "OtherTestmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()
	testutil.AssertTrue(*appender.isClosed, t, "isClosed")
	appender.Close()

	checkLogFileEntries(logFilePath,
		[]string{"{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
			"{\"message\":\"OtherTestmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}"},
		t)
	_, err := os.Stat(expectedNewFileName)
	testutil.AssertNotNil(err, t, "expectedNewFileName should not exist")
}

func TestFileAppenderWriteWithCorrelation(t *testing.T) {
	logFilePath := getAppenderTestLogFile("correlation")
	common.SetLogValuesMockTime(&jsonFormatTestTime)
	correlation := "someCorrelationId"

	appender := CreateFileAppender(logFilePath, &testJsonFormatter, testCronExpression, "").(FileAppender)

	logValuesToFormat := common.CreateLogValuesWithCorrelation(common.INFORMATION_SEVERITY, &correlation, "Testmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func TestFileAppenderWriteCustom(t *testing.T) {
	logFilePath := getAppenderTestLogFile("custom")
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	appender := CreateFileAppender(logFilePath, &testJsonFormatter, testCronExpression, "").(FileAppender)

	customProperties := map[string]any{
		"first": "abc",
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)
	appender.Write(&logValuesToFormat)
	appender.Close()

	checkLogFileEntry(logFilePath, "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\""+jsonFormatTestTimeText+"\"}", t)
}

func checkLogFileEntry(logFilePath string, entry string, t *testing.T) {
	checkLogFileEntries(logFilePath, []string{entry}, t)
}

func checkLogFileEntries(logFilePath string, entries []string, t *testing.T) {
	logFile, err := os.Open(logFilePath)
	testutil.AssertNil(err, t, "open logFile err")
	defer logFile.Close()
	scanner := bufio.NewScanner(logFile)
	i := 0
	for scanner.Scan() {
		if len(entries) > i {
			testutil.AssertEquals(entries[i], scanner.Text(), t, "logFile line")
		}
		i++
	}
}
