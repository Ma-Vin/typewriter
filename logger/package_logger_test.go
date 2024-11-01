package logger

import (
	"os"
	"testing"

	"github.com/ma-vin/typewriter/testutil"
)

func TestCreatePackageLoggers(t *testing.T) {
	os.Clearenv()

	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "DEBUG")
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME+"_"+"first", "info")
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME+"_"+"SECOND", "WARN")

	loggers := CreatePackageLoggers(&testCommonLoggerAppender)

	testutil.AssertNotNil(loggers, t, "logger")
	testutil.AssertEquals(2, len(loggers), t, "logger")

	firstEntry, foundFirstEntry := loggers["first"]
	testutil.AssertTrue(foundFirstEntry, t, "foundFirstEntry")
	testutil.AssertFalse(firstEntry.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(firstEntry.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(firstEntry.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(firstEntry.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(firstEntry.IsFatalEnabled(), t, "IsFatalEnabled")

	secondEntry, foundSecondEntry := loggers["second"]
	testutil.AssertTrue(foundSecondEntry, t, "foundFirstEntry")
	testutil.AssertFalse(secondEntry.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(secondEntry.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(secondEntry.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(secondEntry.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(secondEntry.IsFatalEnabled(), t, "IsFatalEnabled")
}
