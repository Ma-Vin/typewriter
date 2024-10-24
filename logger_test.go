package typewriter

import (
	"os"
	"testing"

	testutil "github.com/ma-vin/typewriter/util"
)

func TestEnableDebugSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "DEBUG")

	InitLogConfig()

	testutil.AssertTrue(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableInformationSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "INFORMATION")

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableInfoSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "INFO")

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableWarningSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "WARNING")

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableWarnSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "WARN")

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableErrorSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "ERROR")

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableFatalSeverity(t *testing.T) {
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, "FATAL")

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableOffSeverity(t *testing.T) {
	os.Unsetenv(DEFAULT_LOG_LEVEL_ENV_NAME)

	InitLogConfig()

	testutil.AssertFalse(IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertFalse(IsFatalEnabled(), t, "IsFatalEnabled")
}
