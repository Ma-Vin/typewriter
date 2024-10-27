package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/ma-vin/typewriter/appender"
	testutil "github.com/ma-vin/typewriter/util"
)

type TestAppender struct {
	content *[]string
}

func (s TestAppender) Write(severity int, message string) {
	*s.content = append(*s.content, fmt.Sprint(severity, message))
}

func (s TestAppender) WriteWithCorrelation(severity int, correlationId string, message string) {
	*s.content = append(*s.content, fmt.Sprint(severity, correlationId, message))
}

func (s TestAppender) WriteCustom(severity int, message string, customValues map[string]any) {
	*s.content = append(*s.content, fmt.Sprint(severity, customValues, message))
}

var testCommonLoggerAppender appender.Appender = TestAppender{content: &[]string{}}
var testCommonLogger = CreateCommonLogger(&testCommonLoggerAppender)

func initTestCommonLogger(envLogLevel string) {
	*testCommonLoggerAppender.(TestAppender).content = []string{}
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, envLogLevel)
	determineSeverityFromEnv(&testCommonLogger)
	mockPanicAndExitAtCommonLogger = true
	panicMockActivated = false
	exitMockAcitvated = false
}

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

func TestDebugCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.Debug("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
}

func TestDebugInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Debug("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
}

func TestDebugWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	testutil.AssertEquals("51234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "DebugWithCorrelation: content[0]")
}

func TestDebugWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
}

func TestDebugCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
}

func TestDebugCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
}

func TestDebugfCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.Debugf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
}

func TestDebugfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Debugf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
}

func TestDebugfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	testutil.AssertEquals("51234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "DebugWithCorrelation: content[0]")
}

func TestDebugfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
}

func TestDebugfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
}

func TestDebugfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
}

func TestInformationCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Information("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
}

func TestInformationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Information("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
}

func TestInformationWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	testutil.AssertEquals("41234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "InformationWithCorrelation: content[0]")
}

func TestInformationWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
}

func TestInformationCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
}

func TestInformationCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
}

func TestInformationfCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Informationf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
}

func TestInformationfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Informationf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
}

func TestInformationfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	testutil.AssertEquals("41234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "InformationWithCorrelation: content[0]")
}

func TestInformationfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
}

func TestInformationfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
}

func TestInformationfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
}

func TestWarningCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Warning("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
}

func TestWarningInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Warning("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
}

func TestWarningWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "WarningWithCorrelation: content[0]")
}

func TestWarningWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
}

func TestWarningCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
}

func TestWarningCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
}

func TestWarningfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Warningf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
}

func TestWarningfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Warningf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
}

func TestWarningfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	testutil.AssertEquals("31234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "WarningWithCorrelation: content[0]")
}

func TestWarningfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
}

func TestWarningfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
}

func TestWarningfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
}

func TestWarningWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	testutil.AssertEquals("3TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithCorrelationAndPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningCustomWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	testutil.AssertEquals("3Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("31234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningWithCorrelationAndPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestWarningCustomWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Error("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
}

func TestErrorInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Error("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
}

func TestErrorWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "ErrorWithCorrelation: content[0]")
}

func TestErrorWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
}

func TestErrorCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
}

func TestErrorCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
}

func TestErrorfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Errorf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
}

func TestErrorfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Errorf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
}

func TestErrorfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	testutil.AssertEquals("21234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "ErrorWithCorrelation: content[0]")
}

func TestErrorfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
}

func TestErrorfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
}

func TestErrorfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
}

func TestErrorWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	testutil.AssertEquals("2TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithCorrelationAndPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorCustomWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	testutil.AssertEquals("2Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("21234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorWithCorrelationAndPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestErrorCustomWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Fatal("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
}

func TestFatalInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.Fatal("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
}

func TestFatalWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "FatalWithCorrelation: content[0]")
}

func TestFatalWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
}

func TestFatalCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
}

func TestFatalCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
}

func TestFatalfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Fatalf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
}

func TestFatalfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.Fatalf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
}

func TestFatalfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "FatalWithCorrelation: content[0]")
}

func TestFatalfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
}

func TestFatalfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
}

func TestFatalfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
}

func TestFatalWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithCorrelationAndPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalCustomWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithCorrelationAndPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalCustomWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	testutil.AssertTrue(panicMockActivated, t, "panic")
}

func TestFatalWithExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithExitInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithCorrelationAndExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithCorrelationAndExitInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalCustomWithExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalCustomWithExitInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithExitfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithCorrelationAndExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalWithCorrelationAndExitfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalCustomWithExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}

func TestFatalCustomWithExitfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	testutil.AssertTrue(exitMockAcitvated, t, "panic")
}
