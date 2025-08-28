package logger

import (
	"fmt"
	"testing"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

type TestAppender struct {
	content *[]string
}

func (s TestAppender) Write(logValues *common.LogValues) {
	if logValues.CorrelationId != nil {
		*s.content = append(*s.content, fmt.Sprint(logValues.Severity, *logValues.CorrelationId, logValues.Message))
		return
	}
	if logValues.CustomValues != nil {
		*s.content = append(*s.content, fmt.Sprint(logValues.Severity, *logValues.CustomValues, logValues.Message))
		return
	}
	*s.content = append(*s.content, fmt.Sprint(logValues.Severity, logValues.Message))
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

func CreateCommonLoggerForTest(appender *appender.Appender, severity int, isCallerToSet bool) CommonLogger {
	commonConfig := config.CommonLoggerConfig{}
	return CreateCommonLoggerFromConfig(config.GeneralLoggerConfig{
		Common:        &commonConfig,
		Severity:      severity,
		IsCallerToSet: isCallerToSet,
	}, appender)
}

var testCommonLoggerAppender appender.Appender = TestAppender{content: &[]string{}}
var testCommonLogger = CreateCommonLoggerForTest(&testCommonLoggerAppender, common.OFF_SEVERITY, false)
var testCommonLoggerCounterAppenderClosed = 0
var testCommonLoggerCounterAppenderClosedExpected = 1

func initTestCommonLogger(envLogLevel string) {
	*testCommonLoggerAppender.(TestAppender).content = []string{}
	determineSeverityByLevel(&testCommonLogger, config.SeverityLevelMap[envLogLevel])
	mockPanicAndExitAtCommonLogger = true
	panicMockActivated = false
	exitMockActivated = false
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

func TestDebugInactiveCommonLogger(t *testing.T) {
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

func TestDebugWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestDebugCustomInactiveCommonLogger(t *testing.T) {
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

func TestDebugfInactiveCommonLogger(t *testing.T) {
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

func TestDebugfWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestDebugfCustomInactiveCommonLogger(t *testing.T) {
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

func TestInformationInactiveCommonLogger(t *testing.T) {
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

func TestInformationWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestInformationCustomInactiveCommonLogger(t *testing.T) {
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

func TestInformationfInactiveCommonLogger(t *testing.T) {
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

func TestInformationfWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestInformationfCustomInactiveCommonLogger(t *testing.T) {
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

func TestWarningInactiveCommonLogger(t *testing.T) {
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

func TestWarningWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestWarningCustomInactiveCommonLogger(t *testing.T) {
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

func TestWarningfInactiveCommonLogger(t *testing.T) {
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

func TestWarningfWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestWarningfCustomInactiveCommonLogger(t *testing.T) {
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
	assertPanicMockActivated(t)
}

func TestWarningWithPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	testutil.AssertEquals("3Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("31234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("WARN")

	testCommonLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.Error("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorInactiveCommonLogger(t *testing.T) {
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

func TestErrorWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestErrorCustomInactiveCommonLogger(t *testing.T) {
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

func TestErrorfInactiveCommonLogger(t *testing.T) {
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

func TestErrorfWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestErrorfCustomInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "Error: len(content)")
}

func TestErrorWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	testutil.AssertEquals("2TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	testutil.AssertEquals("2Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("21234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("ERROR")

	testCommonLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.Fatal("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalInactiveCommonLogger(t *testing.T) {
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

func TestFatalWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestFatalCustomInactiveCommonLogger(t *testing.T) {
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

func TestFatalfInactiveCommonLogger(t *testing.T) {
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

func TestFatalfWithCorrelationInactiveCommonLogger(t *testing.T) {
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

func TestFatalfCustomInactiveCommonLogger(t *testing.T) {
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
	assertPanicMockActivated(t)
}

func TestFatalWithPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	testutil.AssertEquals("1TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithExitInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalWithExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	testutil.AssertEquals("1Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithExitfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitfCommonLogger(t *testing.T) {
	initTestCommonLogger("FATAL")

	testCommonLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testCommonLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitfInactiveCommonLogger(t *testing.T) {
	initTestCommonLogger("OFF")

	testCommonLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testCommonLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	assertExitMockActivated(t)
}

func assertPanicAndExitMockNotActivated(t *testing.T) {
	testutil.AssertFalse(panicMockActivated, t, "panic")
	testutil.AssertFalse(exitMockActivated, t, "exit")
	testutil.AssertEquals(0, testCommonLoggerCounterAppenderClosed, t, "appenderClosed")
}

func assertPanicMockActivated(t *testing.T) {
	testutil.AssertTrue(panicMockActivated, t, "panic")
	testutil.AssertFalse(exitMockActivated, t, "exit")
	testutil.AssertEquals(0, testCommonLoggerCounterAppenderClosed, t, "appenderClosed")
}

func assertExitMockActivated(t *testing.T) {
	testutil.AssertFalse(panicMockActivated, t, "panic")
	testutil.AssertTrue(exitMockActivated, t, "exit")
	testutil.AssertEquals(testCommonLoggerCounterAppenderClosedExpected, testCommonLoggerCounterAppenderClosed, t, "appenderClosed")
}
