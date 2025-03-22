package appender

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/testutil"
)

var fileRenamerCrontab = common.CreateCrontab("* * * * *")
var fileRenamerTestTime = time.Date(2025, time.March, 14, 20, 1, 0, 0, time.UTC)
var fileRenamerMu = sync.Mutex{}
var logValues = &common.LogValues{Time: fileRenamerTestTime}

func getFileRenamerTestLogFile(testCase string) string {
	SkipFileCreationForTest = false
	result := testutil.GetTestCaseFilePath(testCase, true)
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

	logFileName = filepath.Base(logFileName)
	indexOfFileEnding := strings.LastIndex(logFileName, ".")
	expectedNewFileName := logFileName[:indexOfFileEnding] + "_20250314_200100" + logFileName[indexOfFileEnding:]

	checkRenaming(file, logFileName, expectedNewFileName, t, func() {
		fmt.Fprintln(file, "first test entry")
		renamer.CheckFile(logValues)
	})
}

func TestCheckSizeFileNoRename(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("SizeNoRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	renamer := CreateSizeFileRenamer(logFileName, file, 64, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	checkNoRenaming(logFileName, t, func() {
		fmt.Fprintln(file, "first test entry")
		renamer.CheckFile("second test entry")
	})
}

func TestCheckSizeFileRename(t *testing.T) {
	logFileName := getFileRenamerTestLogFile("SizeRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	renamer := CreateSizeFileRenamer(logFileName, file, 20, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	logFileName = filepath.Base(logFileName)
	expectedNewFileName := logFileName[:len(logFileName)-4] + "_20250314_200100.log"

	checkRenaming(file, logFileName, expectedNewFileName, t, func() {
		fmt.Fprintln(file, "first test entry")
		stat, err := os.Stat(logFileName)
		testutil.AssertNil(err, t, "os.Stat(logFileName)")
		renamer.currentByteSize = stat.Size()
		renamer.CheckFile("second test entry")
	})
}

func checkNoRenaming(logFileName string, t *testing.T, toBeExecutedForTest func()) {
	entriesBefore, err := os.ReadDir(path.Dir(logFileName))
	testutil.AssertNil(err, t, "read dir before")

	toBeExecutedForTest()

	entriesAfter, err := os.ReadDir("./")
	testutil.AssertNil(err, t, "read dir after")

	testutil.AssertEquals(len(entriesBefore), len(entriesAfter), t, "compare dirs")
	for _, dirEntryBefore := range entriesBefore {
		contains := false
		for _, dirEntryAfter := range entriesAfter {
			contains = contains || dirEntryBefore.Name() == dirEntryAfter.Name()
		}
		testutil.AssertTrue(contains, t, "missing "+dirEntryBefore.Name()+" at entriesAfter")
	}
}

func checkRenaming(file *os.File, logFileName string, expectedNewFileName string, t *testing.T, toBeExecutedForTest func()) {
	entriesBefore, err := os.ReadDir(path.Dir(logFileName))
	testutil.AssertNil(err, t, "read dir before")

	toBeExecutedForTest()

	entriesAfter, err := os.ReadDir(path.Dir(logFileName))
	testutil.AssertNil(err, t, "read dir after")

	testutil.AssertEquals(len(entriesBefore)+1, len(entriesAfter), t, "compare dirs")

	for _, dirEntryBefore := range entriesBefore {
		contains := false
		for _, dirEntryAfter := range entriesAfter {
			contains = contains || dirEntryBefore.Name() == dirEntryAfter.Name() || logFileName == dirEntryBefore.Name()
		}
		testutil.AssertTrue(contains, t, "missing "+dirEntryBefore.Name()+" at entriesAfter")
	}

	contains := false
	for _, dirEntryAfter := range entriesAfter {
		contains = contains || expectedNewFileName == dirEntryAfter.Name()
	}
	testutil.AssertTrue(contains, t, "missing "+expectedNewFileName+" at entriesAfter")

	_, err = fmt.Fprintln(file, "second test entry")
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
