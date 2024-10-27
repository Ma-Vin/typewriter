package logger

import (
	"os"
	"testing"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/format"
	testutil "github.com/ma-vin/typewriter/util"
)

var testDelimiterFormatter = format.CreateDelimiterFormatter(" - ")
var testCommonLogger = CreateCommonLogger(appender.CreateStandardOutputAppender(&testDelimiterFormatter))

func TestEnableDebugSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "DEBUG")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertTrue(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableInformationSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "INFORMATION")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableInfoSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "INFO")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableWarningSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "WARNING")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableWarnSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "WARN")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableErrorSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "ERROR")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableFatalSeverityCommonLogger(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "FATAL")

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableOffSeverityCommonLogger(t *testing.T) {
	os.Unsetenv(DEFAULT_LOG_LEVEL_ENV_NAME)

	determineSeverityFromEnv(&testCommonLogger)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertFalse(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}
