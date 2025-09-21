package logger

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	testGeneralLoggerCounterAppenderClosed++
}

type testContext struct {
	correlationId any
}

func (c testContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c testContext) Done() <-chan struct{} {
	return nil
}

func (c testContext) Err() error {
	return nil
}

func (c testContext) Value(key any) any {
	if key == "correlationId" {
		return c.correlationId
	}
	return nil
}

func CreateGeneralLoggerForTest(appender *appender.Appender, severity int, isCallerToSet bool) GeneralLogger {
	commonConfig := config.CommonLoggerConfig{CorrelationIdKey: "correlationId"}
	return CreateGeneralLoggerFromConfig(config.GeneralLoggerConfig{
		Common:        &commonConfig,
		Severity:      severity,
		IsCallerToSet: isCallerToSet,
	}, appender)
}

var testGeneralLoggerAppender appender.Appender = TestAppender{content: &[]string{}}
var testGeneralLogger = CreateGeneralLoggerForTest(&testGeneralLoggerAppender, common.OFF_SEVERITY, false)
var testGeneralLoggerCounterAppenderClosed = 0
var testGeneralLoggerCounterAppenderClosedExpected = 1
var testDummyContext context.Context = testContext{correlationId: "1234"}

func initTestGeneralLogger(envLogLevel string) {
	*testGeneralLoggerAppender.(TestAppender).content = []string{}
	determineSeverityByLevel(&testGeneralLogger, config.SeverityLevelMap[envLogLevel])
	mockPanicAndExitAtGeneralLogger = true
	panicMockActivated = false
	exitMockActivated = false
	testGeneralLoggerCounterAppenderClosed = 0
	testGeneralLoggerCounterAppenderClosedExpected = 1
}

func TestEnableDebugSeverityGeneralLogger(t *testing.T) {
	determineSeverityByLevel(&testGeneralLogger, common.DEBUG_SEVERITY)

	testutil.AssertTrue(testGeneralLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testGeneralLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testGeneralLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testGeneralLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testGeneralLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableInformationSeverityGeneralLogger(t *testing.T) {
	determineSeverityByLevel(&testGeneralLogger, common.INFORMATION_SEVERITY)

	testutil.AssertFalse(testGeneralLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertTrue(testGeneralLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testGeneralLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testGeneralLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testGeneralLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableWarningSeverityGeneralLogger(t *testing.T) {
	determineSeverityByLevel(&testGeneralLogger, common.WARNING_SEVERITY)

	testutil.AssertFalse(testGeneralLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testGeneralLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertTrue(testGeneralLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testGeneralLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testGeneralLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableErrorSeverityGeneralLogger(t *testing.T) {
	determineSeverityByLevel(&testGeneralLogger, common.ERROR_SEVERITY)

	testutil.AssertFalse(testGeneralLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testGeneralLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testGeneralLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertTrue(testGeneralLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testGeneralLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableFatalSeverityGeneralLogger(t *testing.T) {
	determineSeverityByLevel(&testGeneralLogger, common.FATAL_SEVERITY)

	testutil.AssertFalse(testGeneralLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testGeneralLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testGeneralLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(testGeneralLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertTrue(testGeneralLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestEnableOffSeverityGeneralLogger(t *testing.T) {
	determineSeverityByLevel(&testGeneralLogger, common.OFF_SEVERITY)

	testutil.AssertFalse(testGeneralLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertFalse(testGeneralLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertFalse(testGeneralLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertFalse(testGeneralLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertFalse(testGeneralLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}

func TestDebugGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.Debug("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.Debug("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.DebugWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	testutil.AssertEquals("51234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "DebugWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.DebugWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.DebugCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.DebugCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCtxGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.DebugCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugCtx: len(content)")
	testutil.AssertEquals("51234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "DebugCtx: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCtxInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.DebugCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugCtx: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.Debugf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.Debugf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.DebugWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	testutil.AssertEquals("51234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "DebugWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.DebugWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.DebugCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	testutil.AssertEquals("5 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "debug: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugfCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.DebugCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Debug: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCtxfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.DebugCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugCtxf: len(content)")
	testutil.AssertEquals("51234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "DebugCtxf: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestDebugCtxfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.DebugCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugCtxf: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.Information("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.Information("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.InformationWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	testutil.AssertEquals("41234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "InformationWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.InformationWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.InformationCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.InformationCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCtxGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.InformationCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationCtx: len(content)")
	testutil.AssertEquals("41234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "InformationCtx: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCtxInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.InformationCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationCtx: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.Informationf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.Informationf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.InformationWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	testutil.AssertEquals("41234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "InformationWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.InformationWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.InformationCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	testutil.AssertEquals("4 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "info: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationfCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.InformationCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Information: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCtxfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("INFO")

	testGeneralLogger.InformationCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationCtxf: len(content)")
	testutil.AssertEquals("41234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "InformationCtxf: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestInformationCtxfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.InformationCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "InformationCtxf: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.Warning("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.Warning("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "WarningWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCtxGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtx: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "WarningCtx: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCtxInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtx: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.Warningf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.Warningf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	testutil.AssertEquals("31234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "WarningWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningfCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Warning: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCtxfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtxf: len(content)")
	testutil.AssertEquals("31234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "WarningCtxf: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningCtxfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtxf: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestWarningWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	testutil.AssertEquals("3TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	testutil.AssertEquals("3 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningCtxWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCtxWithPanic(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtxWithPanic: len(content)")
	testutil.AssertEquals("31234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "WarningCtxWithPanic: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningCtxWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCtxWithPanic(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtxWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	testutil.AssertEquals("3Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("31234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningWithCorrelationAndPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningWithCorrelationAndPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	testutil.AssertEquals("3 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "warn: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningCustomWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCustomWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestWarningCtxWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("WARN")

	testGeneralLogger.WarningCtxWithPanicf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtxWithPanicf: len(content)")
	testutil.AssertEquals("31234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "WarningCtxWithPanicf: content[0]")
	assertPanicMockActivated(t)
}

func TestWarningCtxWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.WarningCtxWithPanicf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "WarningCtxWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.Error("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.Error("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "ErrorWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCtxGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtx: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "ErrorCtx: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCtxInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtx: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.Errorf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.Errorf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	testutil.AssertEquals("21234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "ErrorWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorfCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Error: len(content)")
}

func TestErrorCtxfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtxf: len(content)")
	testutil.AssertEquals("21234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "ErrorCtxf: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorCtxfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtxf: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestErrorWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	testutil.AssertEquals("2TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	testutil.AssertEquals("2 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCtxWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCtxWithPanic(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtxWithPanic: len(content)")
	testutil.AssertEquals("21234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "ErrorCtxWithPanic: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorCtxWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCtxWithPanic(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtxWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	testutil.AssertEquals("2Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("21234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorWithCorrelationAndPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorWithCorrelationAndPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	testutil.AssertEquals("2 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "error: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorCustomWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCustomWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestErrorCtxWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("ERROR")

	testGeneralLogger.ErrorCtxWithPanicf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtxWithPanicf: len(content)")
	testutil.AssertEquals("21234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "ErrorCtxWithPanicf: content[0]")
	assertPanicMockActivated(t)
}

func TestErrorCtxWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.ErrorCtxWithPanicf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "ErrorCtxWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.Fatal("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.Fatal("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithCorrelation("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCustom(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCtxGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtx: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalCtx: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCtxInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtx: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.Fatalf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.Fatalf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfWithCorrelationGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	testutil.AssertEquals("11234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalWithCorrelation: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfWithCorrelationInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithCorrelationf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelation: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfCustomGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalfCustomInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCustomf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "Fatal: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCtxfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalCtxf: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalCtxfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCtxf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxf: len(content)")
	assertPanicAndExitMockNotActivated(t)
}

func TestFatalWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	testutil.AssertEquals("1TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithPanic("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithCorrelationAndPanic("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCtxWithPanicGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCtxWithPanic(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithPanic: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalCtxWithPanic: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalCtxWithPanicInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCtxWithPanic(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithPanic: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	testutil.AssertEquals("1Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithPanicf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalWithCorrelationAndPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithCorrelationAndPanicf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalCustomWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalCtxWithPanicfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCtxWithPanicf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithPanicf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalCtxWithPanicf: content[0]")
	assertPanicMockActivated(t)
}

func TestFatalCtxWithPanicfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCtxWithPanicf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithPanicf: len(content)")
	assertPanicMockActivated(t)
}

func TestFatalWithExitGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	testutil.AssertEquals("1TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithExitInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithExit("Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithCorrelationAndExit("1234", "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	testutil.AssertEquals("1 map[test:123]TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCustomWithExit(map[string]any{"test": 123}, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalCtxWithExitGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCtxWithExit(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithExit: len(content)")
	testutil.AssertEquals("11234TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalCtxWithExit: content[0]")
	assertExitMockActivated(t)
}

func TestFatalCtxWithExitInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCtxWithExit(testDummyContext, "Test", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithExit: len(content)")
	assertExitMockActivated(t)
}

func TestFatalWithExitfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	testutil.AssertEquals("1Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithExitfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithExitf("Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithExitf: len(content)")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalWithCorrelationAndExitfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalWithCorrelationAndExitf("1234", "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalWithCorrelationAndExitf: len(content)")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	testutil.AssertEquals("1 map[test:123]Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "fatal: content[0]")
	assertExitMockActivated(t)
}

func TestFatalCustomWithExitfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCustomWithExitf: len(content)")
	assertExitMockActivated(t)
}

func TestFatalCtxWithExitfGeneralLogger(t *testing.T) {
	initTestGeneralLogger("FATAL")

	testGeneralLogger.FatalCtxWithExitf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithExitf: len(content)")
	testutil.AssertEquals("11234Test Message", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "FatalCtxWithExitf: content[0]")
	assertExitMockActivated(t)
}

func TestFatalCtxWithExitfInactiveGeneralLogger(t *testing.T) {
	initTestGeneralLogger("OFF")

	testGeneralLogger.FatalCtxWithExitf(testDummyContext, "Test %s", "Message")

	testutil.AssertEquals(0, len(*testGeneralLoggerAppender.(TestAppender).content), t, "FatalCtxWithExitf: len(content)")
	assertExitMockActivated(t)
}

func TestNilCorrelationIdPropertyCtxGeneralLogger(t *testing.T) {
	initTestGeneralLogger("DEBUG")

	testGeneralLogger.correlationIdKey = "anyThing"
	testGeneralLogger.DebugCtx(testDummyContext, "Test", "Message")

	testutil.AssertEquals(1, len(*testGeneralLoggerAppender.(TestAppender).content), t, "DebugCtx: len(content)")
	testutil.AssertEquals("5TestMessage", (*testGeneralLoggerAppender.(TestAppender).content)[0], t, "DebugCtx: content[0]")
	assertPanicAndExitMockNotActivated(t)
}

func assertPanicAndExitMockNotActivated(t *testing.T) {
	testutil.AssertFalse(panicMockActivated, t, "panic")
	testutil.AssertFalse(exitMockActivated, t, "exit")
	testutil.AssertEquals(0, testGeneralLoggerCounterAppenderClosed, t, "appenderClosed")
}

func assertPanicMockActivated(t *testing.T) {
	testutil.AssertTrue(panicMockActivated, t, "panic")
	testutil.AssertFalse(exitMockActivated, t, "exit")
	testutil.AssertEquals(0, testGeneralLoggerCounterAppenderClosed, t, "appenderClosed")
}

func assertExitMockActivated(t *testing.T) {
	testutil.AssertFalse(panicMockActivated, t, "panic")
	testutil.AssertTrue(exitMockActivated, t, "exit")
	testutil.AssertEquals(testGeneralLoggerCounterAppenderClosedExpected, testGeneralLoggerCounterAppenderClosed, t, "appenderClosed")
}
