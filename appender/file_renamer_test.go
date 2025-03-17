package appender

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

func getCronFileSchedulerTestLogFile(testCase string) string {
	SkipFileCreationForTest = false
	_, filename, _, _ := runtime.Caller(0)
	result := strings.Replace(filename, ".go", "_"+testCase+"_scratch.log", 1)
	os.Create(result)
	return result
}

func TestCheckFileNoRename(t *testing.T) {
	logFileName := getCronFileSchedulerTestLogFile("noRename")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	fileRenamerCrontab.NextTime = timePtr(fileRenamerTestTime.Add(time.Second))
	renamer := CreateCronFileRenamer(logFileName, file, fileRenamerCrontab, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	logFileName = filepath.Base(logFileName)
	os.Remove(logFileName)

	entriesBefore, err := os.ReadDir("./")
	testutil.AssertNil(err, t, "read dir before")

	renamer.CheckFile(logValues)

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

func TestCheckFileRename(t *testing.T) {
	logFileName := getCronFileSchedulerTestLogFile("rename")
	os.Remove(logFileName)
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	testutil.AssertNil(err, t, "create log file")

	fileRenamerCrontab.NextTime = timePtr(fileRenamerTestTime.Add(-time.Second))
	renamer := CreateCronFileRenamer(logFileName, file, fileRenamerCrontab, &fileRenamerMu)
	renamer.timeFileNameGenerator.referenceTime = &fileRenamerTestTime

	fmt.Fprintln(file, "first test entry")

	logFileName = filepath.Base(logFileName)
	indexOfFileEnding := strings.LastIndex(logFileName, ".")
	expectedNewFileName := logFileName[:indexOfFileEnding] + "_20250314_200100" + logFileName[indexOfFileEnding:]
	os.Remove(expectedNewFileName)

	entriesBefore, err := os.ReadDir("./")
	testutil.AssertNil(err, t, "read dir before")

	renamer.CheckFile(logValues)

	entriesAfter, err := os.ReadDir("./")
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
