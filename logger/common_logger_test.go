package logger

import (
	"fmt"
	"testing"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/testutil"
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

func (s TestAppender) Close() {
	testCommonLoggerCounterAppenderClosed++
}

var testCommonLoggerAppender appender.Appender = TestAppender{content: &[]string{}}
var testCommonLogger = CreateCommonLogger(&testCommonLoggerAppender, common.OFF_SEVERITY)
var testCommonLoggerCounterAppenderClosed = 0
var testCommonLoggerCounterAppenderClosedExpected = 1

func initTestCommonLogger(envLogLevel string) {
	*testCommonLoggerAppender.(TestAppender).content = []string{}
	determineSeverityByLevel(&testCommonLogger, config.SeverityLevelMap[envLogLevel])
	mockPanicAndExitAtCommonLogger = true
	panicMockActivated = false
	exitMockAcitvated = false
	testCommonLoggerCounterAppenderClosed = 0
	testCommonLoggerCounterAppenderClosedExpected = 1
}

func TestEnableDebugSeverityCommonLogger(t *testing.T) {
	determineSeverityByLevel(&testCommonLogger, common.DEBUG_SEVERITY)

	testutil.AssertTrue(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableInformationSeverityCommonLogger(t *testing.T) {
	determineSeverityByLevel(&testCommonLogger, common.INFORMATION_SEVERITY)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableWarningSeverityCommonLogger(t *testing.T) {
	determineSeverityByLevel(&testCommonLogger, common.WARNING_SEVERITY)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableErrorSeverityCommonLogger(t *testing.T) {
	determineSeverityByLevel(&testCommonLogger, common.ERROR_SEVERITY)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableFatalSeverityCommonLogger(t *testing.T) {
	determineSeverityByLevel(&testCommonLogger, common.FATAL_SEVERITY)

	testutil.AssertFalse(testCommonLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testCommonLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testCommonLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(testCommonLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testCommonLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableOffSeverityCommonLogger(t *testing.T) {
	determineSeverityByLevel(&testCommonLogger, common.OFF_SEVERITY)

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
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Debug("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	testutil.AssertEquals("51234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "DebugWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.Debugf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Debugf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	testutil.AssertEquals("51234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "DebugWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("DEBUG")

	testCommonLogger.DebugCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.DebugCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Information("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Information("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	testutil.AssertEquals("41234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "InformationWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.Informationf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Informationf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	testutil.AssertEquals("41234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "InformationWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("INFO")

	testCommonLogger.InformationCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.InformationCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Warning("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Warning("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "WarningWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.Warningf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Warningf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	testutil.AssertEquals("31234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "WarningWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	testutil.AssertEquals("3TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockAcitvated(t)
}

func TestWarningWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestWarningWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockAcitvated(t)
}

func TestWarningWithCorrelationAndPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestWarningCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockAcitvated(t)
}

func TestWarningCustomWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestWarningWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	testutil.AssertEquals("3Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockAcitvated(t)
}

func TestWarningWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestWarningWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("31234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockAcitvated(t)
}

func TestWarningWithCorrelationAndPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestWarningCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockAcitvated(t)
}

func TestWarningCustomWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestErrorCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Error("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Error("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "ErrorWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Errorf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Errorf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	testutil.AssertEquals("21234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "ErrorWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
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
	assertPanicMockAcitvated(t)
}

func TestErrorWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestErrorWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockAcitvated(t)
}

func TestErrorWithCorrelationAndPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestErrorCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockAcitvated(t)
}

func TestErrorCustomWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestErrorWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	testutil.AssertEquals("2Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockAcitvated(t)
}

func TestErrorWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestErrorWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("21234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockAcitvated(t)
}

func TestErrorWithCorrelationAndPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestErrorCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockAcitvated(t)
}

func TestErrorCustomWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Fatal("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.Fatal("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "FatalWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Fatalf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.Fatalf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfWithCorrelationCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "FatalWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfWithCorrelationInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfCustomCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfCustomInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockAcitvated(t)
}

func TestFatalWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockAcitvated(t)
}

func TestFatalWithCorrelationAndPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockAcitvated(t)
}

func TestFatalCustomWithPanicInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockAcitvated(t)
}

func TestFatalWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockAcitvated(t)
}

func TestFatalWithCorrelationAndPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockAcitvated(t)
}

func TestFatalCustomWithPanicfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	assertPanicMockAcitvated(t)
}

func TestFatalWithExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockAcitvated(t)
}

func TestFatalWithExitInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	assertExitMockAcitvated(t)
}

func TestFatalWithCorrelationAndExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockAcitvated(t)
}

func TestFatalWithCorrelationAndExitInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	assertExitMockAcitvated(t)
}

func TestFatalCustomWithExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockAcitvated(t)
}

func TestFatalCustomWithExitInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	assertExitMockAcitvated(t)
}

func TestFatalWithExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockAcitvated(t)
}

func TestFatalWithExitfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	assertExitMockAcitvated(t)
}

func TestFatalWithCorrelationAndExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockAcitvated(t)
}

func TestFatalWithCorrelationAndExitfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	assertExitMockAcitvated(t)
}

func TestFatalCustomWithExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockAcitvated(t)
}

func TestFatalCustomWithExitfInaktiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	assertExitMockAcitvated(t)
}

func assertPanicAndExitMockNotActivated(t *testing.T) {
	testutil.AssertFalse(panicMockActivated, t, "panic")
	testutil.AssertFalse(exitMockAcitvated, t, "exit")
	testutil.AssertEquals(0, testCommonLoggerCounterAppenderClosed, t, "appenderClosed")
}

func assertPanicMockAcitvated(t *testing.T) {
	testutil.AssertTrue(panicMockActivated, t, "panic")
	testutil.AssertFalse(exitMockAcitvated, t, "exit")
	testutil.AssertEquals(0, testCommonLoggerCounterAppenderClosed, t, "appenderClosed")
}

func assertExitMockAcitvated(t *testing.T) {
	testutil.AssertFalse(panicMockActivated, t, "panic")
	testutil.AssertTrue(exitMockAcitvated, t, "exit")
	testutil.AssertEquals(testCommonLoggerCounterAppenderClosedExpected, testCommonLoggerCounterAppenderClosed, t, "appenderClosed")
}
