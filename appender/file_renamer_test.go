package appender

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/common"
)

const testResourceTarget = "genTestResources"

var fileRenamerCrontab = common.CreateCrontab("* * * * *")
var fileRenamerTestTime = time.Date(2025, time.March, 14, 20, 1, 0, 0, time.UTC)
var fileRenamerMu = sync.Mutex{}
var logValues = &common.LogValues{Time: fileRenamerTestTime}

func getFileRenamerTestLogFile(testCase string) string {
	SkipFileCreationForTest = false
	result := testutil.DetermineTestCaseFilePathAt(testCase, "log", true, true, testResourceTarget)
	os.Create(result)
	return result
}

func TestCheckCronFileNoRename(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("CronNoRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	fileRenamerCrontab.NextTime = timePtr(fileRenamerTestTime.Add(time.Second))
	renamer := CreateCronFileRenamer(logFileName, file, fileRenamerCrontab, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	checkNoRenaming(logFileName, t, func() {
		fmt.Fprintln(file, "first test entry")
		renamer.CheckFile(logValues)
	})
}

func TestCheckCronFileRename(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("CronRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	fileRenamerCrontab.NextTime = timePtr(fileRenamerTestTime.Add(-time.Second))
	renamer := CreateCronFileRenamer(logFileName, file, fileRenamerCrontab, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	indexOfFileEnding := strings.LastIndex(logFileName, ".")
	expectedNewFileName := logFileName[:indexOfFileEnding] + "_20250314_200100" + logFileName[indexOfFileEnding:]

	checkRenaming(file, logFileName, "CronRename", expectedNewFileName, t, func() {
		fmt.Fprintln(file, "first test entry")
		renamer.CheckFile(logValues)
	})
}

func TestCheckCronFileRenameButAlreadyExists(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("ExistingCronRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	fileRenamerCrontab.NextTime = timePtr(fileRenamerTestTime.Add(-time.Second))
	renamer := CreateCronFileRenamer(logFileName, file, fileRenamerCrontab, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	indexOfFileEnding := strings.LastIndex(logFileName, ".")
	existingFileName := logFileName[:indexOfFileEnding] + "_20250314_200100" + logFileName[indexOfFileEnding:]
	expectedNewFileName := logFileName[:indexOfFileEnding] + "_20250314_200100_1" + logFileName[indexOfFileEnding:]

	os.Create(existingFileName)

	checkRenaming(file, logFileName, "ExistingCronRename", expectedNewFileName, t, func() {
		fmt.Fprintln(file, "first test entry")
		renamer.CheckFile(logValues)
	})
}

func TestCheckSizeFileNoRename(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("SizeNoRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	firstEntry := "first test entry"
	secondEntry := "second test entry"
	limitByteSize := len(firstEntry) + len(secondEntry) + 2*len(fmt.Sprintln()) + 1

	renamer := CreateSizeFileRenamer(logFileName, file, int64(limitByteSize), &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	checkNoRenaming(logFileName, t, func() {
		fmt.Fprintln(file, firstEntry)
		renamer.CheckFile(secondEntry)
	})
}

func TestCheckSizeFileRename(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("SizeRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	firstEntry := "first test entry"
	secondEntry := "second test entry"
	limitByteSize := len(firstEntry) + len(secondEntry) + 2*len(fmt.Sprintln())

	renamer := CreateSizeFileRenamer(logFileName, file, int64(limitByteSize), &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	indexOfFileEnding := strings.LastIndex(logFileName, ".")
	expectedNewFileName := logFileName[:indexOfFileEnding] + "_20250314_200100" + logFileName[indexOfFileEnding:]

	checkRenaming(file, logFileName, "SizeRename", expectedNewFileName, t, func() {
		fmt.Fprintln(file, firstEntry)
		stat, err := os.Stat(logFileName)
		testutil.AssertNil(err, t, "os.Stat(logFileName)")
		renamer.currentByteSize = stat.Size()
		renamer.CheckFile(secondEntry)
	})
}

func TestCheckSizeFileRenameButAlreadyExists(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("ExistingSizeRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	firstEntry := "first test entry"
	secondEntry := "second test entry"
	limitByteSize := len(firstEntry) + len(secondEntry) + 2*len(fmt.Sprintln())

	renamer := CreateSizeFileRenamer(logFileName, file, int64(limitByteSize), &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	indexOfFileEnding := strings.LastIndex(logFileName, ".")
	existingFileName := logFileName[:indexOfFileEnding] + "_20250314_200100" + logFileName[indexOfFileEnding:]
	expectedNewFileName := logFileName[:indexOfFileEnding] + "_20250314_200100_1" + logFileName[indexOfFileEnding:]

	os.Create(existingFileName)

	checkRenaming(file, logFileName, "ExistingSizeRename", expectedNewFileName, t, func() {
		fmt.Fprintln(file, firstEntry)
		stat, err := os.Stat(logFileName)
		testutil.AssertNil(err, t, "os.Stat(logFileName)")
		renamer.currentByteSize = stat.Size()
		renamer.CheckFile(secondEntry)
	})
}

func checkNoRenaming(testCase string, t *testing.T, toBeExecutedForTest func()) {
	entriesBefore := testutil.DetermineExistingTestCaseFilePathsAt(testCase, testResourceTarget)

	toBeExecutedForTest()

	entriesAfter := testutil.DetermineExistingTestCaseFilePathsAt(testCase, testResourceTarget)

	testutil.AssertEquals(len(entriesBefore), len(entriesAfter), t, "compare dirs")
	for _, dirEntryBefore := range entriesBefore {
		contains := false
		for _, dirEntryAfter := range entriesAfter {
			contains = contains || dirEntryBefore == dirEntryAfter
		}
		testutil.AssertTrue(contains, t, "missing "+path.Base(dirEntryBefore)+" at entriesAfter")
	}
}

func checkRenaming(file *os.File, logFileName string, testCase string, expectedNewFileName string, t *testing.T, toBeExecutedForTest func()) {
	entriesBefore := testutil.DetermineExistingTestCaseFilePathsAt(testCase, testResourceTarget)

	toBeExecutedForTest()

	entriesAfter := testutil.DetermineExistingTestCaseFilePathsAt(testCase, testResourceTarget)

	testutil.AssertEquals(len(entriesBefore)+1, len(entriesAfter), t, "compare dirs")

	for _, dirEntryBefore := range entriesBefore {
		contains := false
		for _, dirEntryAfter := range entriesAfter {
			contains = contains || dirEntryBefore == dirEntryAfter
		}
		testutil.AssertTrue(contains, t, "missing "+path.Base(dirEntryBefore)+" at entriesAfter")
	}

	contains := false
	for _, dirEntryAfter := range entriesAfter {
		contains = contains || expectedNewFileName == dirEntryAfter
	}
	testutil.AssertTrue(contains, t, "missing "+path.Base(expectedNewFileName)+" at entriesAfter")

	_, err := fmt.Fprintln(file, "second test entry")
	testutil.AssertNil(err, t, "write second entry")

	dat, err := os.ReadFile(expectedNewFileName)
	testutil.AssertNil(err, t, "read renamed log file")
	testutil.AssertEquals("first test entry", strings.TrimSpace(string(dat)), t, "wrong first entry")

	dat, err = os.ReadFile(logFileName)
	testutil.AssertNil(err, t, "read original named log file")
	testutil.AssertEquals("second test entry", strings.TrimSpace(string(dat)), t, "wrong second entry")
}

func timePtr(t time.Time) *time.Time {
	return &t
}
