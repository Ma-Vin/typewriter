package logger

import (
	"os"
	"testing"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
)

var testMainGeneralLoggerAppender appender.Appender = TestAppender{content: &[]string{}}
var testMainPackageLoggerAppender appender.Appender = TestAppender{content: &[]string{}}

var testMainGeneralLogger GeneralLogger
var testMainPackageLogger GeneralLogger
var mainLogger MainLogger

func clearMainLoggerTestEnv() {
	os.Unsetenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME)
	os.Unsetenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME + "_LOGGER")
}

func initMainLoggerTest(envCommonLogLevel string, envPackageLogLevel string, packageName string, isFullQualified bool) {
	clearMainLoggerTestEnv()

	*testMainGeneralLoggerAppender.(TestAppender).content = []string{}
	*testMainPackageLoggerAppender.(TestAppender).content = []string{}

	appenders = []appender.Appender{testMainGeneralLoggerAppender, testMainPackageLoggerAppender}

	testMainGeneralLogger = CreateGeneralLoggerForTest(&testMainGeneralLoggerAppender, config.SeverityLevelMap[envCommonLogLevel], false)
	testMainPackageLogger = CreateGeneralLoggerForTest(&testMainPackageLoggerAppender, config.SeverityLevelMap[envPackageLogLevel], false)

	mockPanicAndExitAtGeneralLogger = true
	panicMockActivated = false
	exitMockActivated = false
	testGeneralLoggerCounterAppenderClosed = 0
	testGeneralLoggerCounterAppenderClosedExpected = 2

	mainLogger = MainLogger{
		generalLogger:           &testMainGeneralLogger,
		existPackageLogger:      true,
		useFullQualifiedPackage: isFullQualified,
		packageLoggers:          map[string]*GeneralLogger{packageName: &testMainPackageLogger},
	}
}

func initMainLoggerViaPackageTest(envCommonLogLevel string, envPackageLogLevel string) {
	initMainLoggerTest(envCommonLogLevel, envPackageLogLevel, "logger", false)
}

func initMainLoggerViaFullQualifiedPackageTest(envCommonLogLevel string, envPackageLogLevel string) {
	initMainLoggerTest(envCommonLogLevel, envPackageLogLevel, "github.com/ma-vin/typewriter/logger", true)
}

func initMainLoggerViaCommonTest(envCommonLogLevel string, envPackageLogLevel string) {
	initMainLoggerTest(envCommonLogLevel, envPackageLogLevel, "other", false)
}

func initMainLoggerViaCommonFullQualifiedTest(envCommonLogLevel string, envPackageLogLevel string) {
	initMainLoggerTest(envCommonLogLevel, envPackageLogLevel, "github.com/ma-vin/typewriter/other", true)
}

func initMainLoggerOnlyCommonTest(envCommonLogLevel string) {
	initMainLoggerViaCommonTest(envCommonLogLevel, envCommonLogLevel)
	testGeneralLoggerCounterAppenderClosedExpected = 1
	mainLogger = MainLogger{
		generalLogger:      &testMainGeneralLogger,
		existPackageLogger: false,
		packageLoggers:     make(map[string]*GeneralLogger),
	}
}

// -------------------
//
// Is Enabled Via Package Block
//
// -------------------

func TestEnableDebugSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "DEBUG")

	assertEnabled(t, true, true, true, true, true)
}

func TestEnableInformationSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "INFORMATION")

	assertEnabled(t, false, true, true, true, true)
}

func TestEnableInfoSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "INFO")

	assertEnabled(t, false, true, true, true, true)
}

func TestEnableWarningSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "WARNING")

	assertEnabled(t, false, false, true, true, true)
}

func TestEnableWarnSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "WARN")

	assertEnabled(t, false, false, true, true, true)
}

func TestEnableErrorSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "ERROR")

	assertEnabled(t, false, false, false, true, true)
}

func TestEnableFatalSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "FATAL")

	assertEnabled(t, false, false, false, false, true)
}

func TestEnableOffSeverityMainLoggerViaPackage(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	assertEnabled(t, false, false, false, false, false)
}

// -------------------
//
// Is Enabled Via Common Block
//
// -------------------

func TestEnableDebugSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "OFF")

	assertEnabled(t, true, true, true, true, true)
}

func TestEnableInformationSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("INFORMATION", "OFF")

	assertEnabled(t, false, true, true, true, true)
}

func TestEnableInfoSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "OFF")

	assertEnabled(t, false, true, true, true, true)
}

func TestEnableWarningSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("WARNING", "OFF")

	assertEnabled(t, false, false, true, true, true)
}

func TestEnableWarnSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "OFF")

	assertEnabled(t, false, false, true, true, true)
}

func TestEnableErrorSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "OFF")

	assertEnabled(t, false, false, false, true, true)
}

func TestEnableFatalSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "OFF")

	assertEnabled(t, false, false, false, false, true)
}

func TestEnableOffSeverityMainLoggerViaCommon(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	assertEnabled(t, false, false, false, false, false)
}

// -------------------
//
// Is Enabled Only Common Block
//
// -------------------

func TestEnableDebugSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	assertEnabled(t, true, true, true, true, true)
}

func TestEnableInformationSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFORMATION")

	assertEnabled(t, false, true, true, true, true)
}

func TestEnableInfoSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	assertEnabled(t, false, true, true, true, true)
}

func TestEnableWarningSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARNING")

	assertEnabled(t, false, false, true, true, true)
}

func TestEnableWarnSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	assertEnabled(t, false, false, true, true, true)
}

func TestEnableErrorSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	assertEnabled(t, false, false, false, true, true)
}

func TestEnableFatalSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	assertEnabled(t, false, false, false, false, true)
}

func TestEnableOffSeverityMainLoggerOnlyCommon(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	assertEnabled(t, false, false, false, false, false)
}

// -------------------
//
// Debug Package Block
//
// -------------------

func TestMainLoggerViaPackageDebug(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.Debug("debug test message")

	assertMessageViaPackage(t, "Debug", "5debug test message")
}

func TestMainLoggerViaPackageInactiveDebug(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestMainLoggerViaPackageDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertMessageViaPackage(t, "DebugWithCorrelation", "51234debug test message")
}

func TestMainLoggerViaPackageInactiveDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestMainLoggerViaPackageDebugCustom(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaPackage(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestMainLoggerViaPackageInactiveDebugCustom(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestMainLoggerViaPackageDebugCtx(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertMessageViaPackage(t, "DebugCtx", "51234debug test message")
}

func TestMainLoggerViaPackageInactiveDebugCtx(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}
func TestMainLoggerViaPackageDebugf(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.Debugf("debug test %s", "message")

	assertMessageViaPackage(t, "Debugf", "5debug test message")
}

func TestMainLoggerViaPackageInactiveDebugf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestMainLoggerViaPackageDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaPackage(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestMainLoggerViaPackageInactiveDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestMainLoggerViaPackageDebugCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaPackage(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestMainLoggerViaPackageInactiveDebugCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestMainLoggerViaPackageDebugCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaPackage(t, "DebugCtxf", "51234debug test message")
}

func TestMainLoggerViaPackageInactiveDebugCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Debug Full Qualified Package Block
//
// -------------------

func TestMainLoggerViaFullQualifiedPackageDebug(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.Debug("debug test message")

	assertMessageViaPackage(t, "Debug", "5debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebug(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestMainLoggerViaFullQualifiedPackageDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertMessageViaPackage(t, "DebugWithCorrelation", "51234debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestMainLoggerViaFullQualifiedPackageDebugCustom(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaPackage(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugCustom(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestMainLoggerViaFullQualifiedPackageDebugCtx(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertMessageViaPackage(t, "DebugCtx", "51234debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugCtx(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}

func TestMainLoggerViaFullQualifiedPackageDebugf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.Debugf("debug test %s", "message")

	assertMessageViaPackage(t, "Debugf", "5debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestMainLoggerViaFullQualifiedPackageDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaPackage(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestMainLoggerViaFullQualifiedPackageDebugCustomf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaPackage(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugCustomf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestMainLoggerViaFullQualifiedPackageDebugCtxf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("DEBUG", "DEBUG")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaPackage(t, "DebugCtxf", "51234debug test message")
}

func TestMainLoggerViaFullQualifiedPackageInactiveDebugCtxf(t *testing.T) {
	initMainLoggerViaFullQualifiedPackageTest("INFO", "INFO")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Debug Common Block
//
// -------------------

func TestMainLoggerViaCommonDebug(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.Debug("debug test message")

	assertMessageViaCommon(t, "Debug", "5debug test message")
}

func TestMainLoggerViaCommonInactiveDebug(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestMainLoggerViaCommonDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertMessageViaCommon(t, "DebugWithCorrelation", "51234debug test message")
}

func TestMainLoggerViaCommonInactiveDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestMainLoggerViaCommonDebugCustom(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaCommon(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestMainLoggerViaCommonInactiveDebugCustom(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestMainLoggerViaCommonDebugCtx(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertMessageViaCommon(t, "DebugCtx", "51234debug test message")
}

func TestMainLoggerViaCommonInactiveDebugCtx(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}

func TestMainLoggerViaCommonDebugf(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.Debugf("debug test %s", "message")

	assertMessageViaCommon(t, "Debugf", "5debug test message")
}

func TestMainLoggerViaCommonInactiveDebugf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestMainLoggerViaCommonDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaCommon(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestMainLoggerViaCommonInactiveDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestMainLoggerViaCommonDebugCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestMainLoggerViaCommonInactiveDebugCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestMainLoggerViaCommonDebugCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("DEBUG", "DEBUG")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCtxf", "51234debug test message")
}

func TestMainLoggerViaCommonInactiveDebugCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Debug Common Block and given full qualified package
//
// -------------------

func TestMainLoggerViaCommonFullQualifiedDebug(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.Debug("debug test message")

	assertMessageViaCommon(t, "Debug", "5debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebug(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestMainLoggerViaCommonFullQualifiedDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertMessageViaCommon(t, "DebugWithCorrelation", "51234debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestMainLoggerViaCommonFullQualifiedDebugCustom(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaCommon(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugCustom(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestMainLoggerViaCommonFullQualifiedDebugCtx(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertMessageViaCommon(t, "DebugCtx", "51234debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugCtx(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}

func TestMainLoggerViaCommonFullQualifiedDebugf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.Debugf("debug test %s", "message")

	assertMessageViaCommon(t, "Debugf", "5debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestMainLoggerViaCommonFullQualifiedDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaCommon(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestMainLoggerViaCommonFullQualifiedDebugCustomf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugCustomf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestMainLoggerViaCommonFullQualifiedDebugCtxf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("DEBUG", "DEBUG")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCtxf", "51234debug test message")
}

func TestMainLoggerViaCommonFullQualifiedInactiveDebugCtxf(t *testing.T) {
	initMainLoggerViaCommonFullQualifiedTest("INFO", "INFO")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Debug Only Common Block
//
// -------------------

func TestMainLoggerOnlyCommonDebug(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.Debug("debug test message")

	assertMessageViaCommon(t, "Debug", "5debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebug(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestMainLoggerOnlyCommonDebugWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertMessageViaCommon(t, "DebugWithCorrelation", "51234debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestMainLoggerOnlyCommonDebugCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaCommon(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestMainLoggerOnlyCommonDebugCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertMessageViaCommon(t, "DebugCtx", "51234debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}

func TestMainLoggerOnlyCommonDebugf(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.Debugf("debug test %s", "message")

	assertMessageViaCommon(t, "Debugf", "5debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestMainLoggerOnlyCommonDebugWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaCommon(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestMainLoggerOnlyCommonDebugCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestMainLoggerOnlyCommonDebugCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("DEBUG")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCtxf", "51234debug test message")
}

func TestMainLoggerOnlyCommonInactiveDebugCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Information Package Block
//
// -------------------

func TestMainLoggerViaPackageInformation(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.Information("info test message")

	assertMessageViaPackage(t, "Information", "4info test message")
}

func TestMainLoggerViaPackageInactiveInformation(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.Information("info test message")

	assertNoMessage(t, "Information")
}

func TestMainLoggerViaPackageInformationWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.InformationWithCorrelation("1234", "info test message")

	assertMessageViaPackage(t, "InformationWithCorrelation", "41234info test message")
}

func TestMainLoggerViaPackageInactiveInformationWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.InformationWithCorrelation("1234", "info test message")

	assertNoMessage(t, "InformationWithCorrelation")
}

func TestMainLoggerViaPackageInformationCustom(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.InformationCustom(map[string]any{"test": 123}, "info test message")

	assertMessageViaPackage(t, "InformationCustom", "4 map[test:123]info test message")
}

func TestMainLoggerViaPackageInactiveInformationCustom(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.InformationCustom(map[string]any{"test": 123}, "info test message")

	assertNoMessage(t, "InformationCustom")
}

func TestMainLoggerViaPackageInformationCtx(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.InformationCtx(testDummyContext, "info test message")

	assertMessageViaPackage(t, "InformationCtx", "41234info test message")
}

func TestMainLoggerViaPackageInactiveInformationCtx(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.InformationCtx(testDummyContext, "info test message")

	assertNoMessage(t, "InformationCtx")
}

func TestMainLoggerViaPackageInformationf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.Informationf("info test %s", "message")

	assertMessageViaPackage(t, "Informationf", "4info test message")
}

func TestMainLoggerViaPackageInactiveInformationf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.Informationf("info test %s", "message")

	assertNoMessage(t, "Informationf")
}

func TestMainLoggerViaPackageInformationWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.InformationWithCorrelationf("1234", "info test %s", "message")

	assertMessageViaPackage(t, "InformationWithCorrelationf", "41234info test message")
}

func TestMainLoggerViaPackageInactiveInformationWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.InformationWithCorrelationf("1234", "info test %s", "message")

	assertNoMessage(t, "InformationWithCorrelationf")
}

func TestMainLoggerViaPackageInformationCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertMessageViaPackage(t, "InformationCustomf", "4 map[test:123]info test message")
}

func TestMainLoggerViaPackageInactiveInformationCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertNoMessage(t, "InformationCustomf")
}

func TestMainLoggerViaPackageInformationCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("INFO", "INFO")

	mainLogger.InformationCtxf(testDummyContext, "info test %s", "message")

	assertMessageViaPackage(t, "InformationCtxf", "41234info test message")
}

func TestMainLoggerViaPackageInactiveInformationCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.InformationCtxf(testDummyContext, "info test %s", "message")

	assertNoMessage(t, "InformationCtxf")
}

// -------------------
//
// Information Common Block
//
// -------------------

func TestMainLoggerViaCommonInformation(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.Information("info test message")

	assertMessageViaCommon(t, "Information", "4info test message")
}

func TestMainLoggerViaCommonInactiveInformation(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.Information("info test message")

	assertNoMessage(t, "Information")
}

func TestMainLoggerViaCommonInformationWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.InformationWithCorrelation("1234", "info test message")

	assertMessageViaCommon(t, "InformationWithCorrelation", "41234info test message")
}

func TestMainLoggerViaCommonInactiveInformationWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.InformationWithCorrelation("1234", "info test message")

	assertNoMessage(t, "InformationWithCorrelation")
}

func TestMainLoggerViaCommonInformationCustom(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.InformationCustom(map[string]any{"test": 123}, "info test message")

	assertMessageViaCommon(t, "InformationCustom", "4 map[test:123]info test message")
}

func TestMainLoggerViaCommonInactiveInformationCustom(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.InformationCustom(map[string]any{"test": 123}, "info test message")

	assertNoMessage(t, "InformationCustom")
}

func TestMainLoggerViaCommonInformationCtx(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.InformationCtx(testDummyContext, "info test message")

	assertMessageViaCommon(t, "InformationCtx", "41234info test message")
}

func TestMainLoggerViaCommonInactiveInformationCtx(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.InformationCtx(testDummyContext, "info test message")

	assertNoMessage(t, "InformationCtx")
}

func TestMainLoggerViaCommonInformationf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.Informationf("info test %s", "message")

	assertMessageViaCommon(t, "Informationf", "4info test message")
}

func TestMainLoggerViaCommonInactiveInformationf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.Informationf("info test %s", "message")

	assertNoMessage(t, "Informationf")
}

func TestMainLoggerViaCommonInformationWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.InformationWithCorrelationf("1234", "info test %s", "message")

	assertMessageViaCommon(t, "InformationWithCorrelationf", "41234info test message")
}

func TestMainLoggerViaCommonInactiveInformationWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.InformationWithCorrelationf("1234", "info test %s", "message")

	assertNoMessage(t, "InformationWithCorrelationf")
}

func TestMainLoggerViaCommonInformationCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCustomf", "4 map[test:123]info test message")
}

func TestMainLoggerViaCommonInactiveInformationCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertNoMessage(t, "InformationCustomf")
}

func TestMainLoggerViaCommonInformationCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("INFO", "INFO")

	mainLogger.InformationCtxf(testDummyContext, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCtxf", "41234info test message")
}

func TestMainLoggerViaCommonInactiveInformationCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.InformationCtxf(testDummyContext, "info test %s", "message")

	assertNoMessage(t, "InformationCtxf")
}

// -------------------
//
// Information Only Common Block
//
// -------------------

func TestMainLoggerOnlyCommonInformation(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.Information("info test message")

	assertMessageViaCommon(t, "Information", "4info test message")
}

func TestMainLoggerOnlyCommonInactiveInformation(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.Information("info test message")

	assertNoMessage(t, "Information")
}

func TestMainLoggerOnlyCommonInformationWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.InformationWithCorrelation("1234", "info test message")

	assertMessageViaCommon(t, "InformationWithCorrelation", "41234info test message")
}

func TestMainLoggerOnlyCommonInactiveInformationWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.InformationWithCorrelation("1234", "info test message")

	assertNoMessage(t, "InformationWithCorrelation")
}

func TestMainLoggerOnlyCommonInformationCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.InformationCustom(map[string]any{"test": 123}, "info test message")

	assertMessageViaCommon(t, "InformationCustom", "4 map[test:123]info test message")
}

func TestMainLoggerOnlyCommonInactiveInformationCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.InformationCustom(map[string]any{"test": 123}, "info test message")

	assertNoMessage(t, "InformationCustom")
}

func TestMainLoggerOnlyCommonInformationCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.InformationCtx(testDummyContext, "info test message")

	assertMessageViaCommon(t, "InformationCtx", "41234info test message")
}

func TestMainLoggerOnlyCommonInactiveInformationCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.InformationCtx(testDummyContext, "info test message")

	assertNoMessage(t, "InformationCtx")
}

func TestMainLoggerOnlyCommonInformationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.Informationf("info test %s", "message")

	assertMessageViaCommon(t, "Informationf", "4info test message")
}

func TestMainLoggerOnlyCommonInactiveInformationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.Informationf("info test %s", "message")

	assertNoMessage(t, "Informationf")
}

func TestMainLoggerOnlyCommonInformationWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.InformationWithCorrelationf("1234", "info test %s", "message")

	assertMessageViaCommon(t, "InformationWithCorrelationf", "41234info test message")
}

func TestMainLoggerOnlyCommonInactiveInformationWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.InformationWithCorrelationf("1234", "info test %s", "message")

	assertNoMessage(t, "InformationWithCorrelationf")
}

func TestMainLoggerOnlyCommonInformationCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCustomf", "4 map[test:123]info test message")
}

func TestMainLoggerOnlyCommonInformationCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("INFO")

	mainLogger.InformationCtxf(testDummyContext, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCtxf", "41234info test message")
}

func TestMainLoggerOnlyCommonInactiveInformationCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.InformationCtxf(testDummyContext, "info test %s", "message")

	assertNoMessage(t, "InformationCtxf")
}

// -------------------
//
// Warning Package Block
//
// -------------------

func TestMainLoggerViaPackageWarning(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.Warning("warn test message")

	assertMessageViaPackage(t, "Warning", "3warn test message")
}

func TestMainLoggerViaPackageInactiveWarning(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.Warning("warn test message")

	assertNoMessage(t, "Warning")
}

func TestMainLoggerViaPackageWarningWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.WarningWithCorrelation("1234", "warn test message")

	assertMessageViaPackage(t, "WarningWithCorrelation", "31234warn test message")
}

func TestMainLoggerViaPackageInactiveWarningWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.WarningWithCorrelation("1234", "warn test message")

	assertNoMessage(t, "WarningWithCorrelation")
}

func TestMainLoggerViaPackageWarningCustom(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertMessageViaPackage(t, "WarningCustom", "3 map[test:123]warn test message")
}

func TestMainLoggerViaPackageInactiveWarningCustom(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertNoMessage(t, "WarningCustom")
}

func TestMainLoggerViaPackageWarningCtx(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.WarningCtx(testDummyContext, "warn test message")

	assertMessageViaPackage(t, "WarningCtx", "31234warn test message")
}

func TestMainLoggerViaPackageInactiveWarningCtx(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.WarningCtx(testDummyContext, "warn test message")

	assertNoMessage(t, "WarningCtx")
}

func TestMainLoggerViaPackageWarningf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.Warningf("warn test %s", "message")

	assertMessageViaPackage(t, "Warningf", "3warn test message")
}

func TestMainLoggerViaPackageInactiveWarningf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.Warningf("warn test %s", "message")

	assertNoMessage(t, "Warningf")
}

func TestMainLoggerViaPackageWarningWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.WarningWithCorrelationf("1234", "warn test %s", "message")

	assertMessageViaPackage(t, "WarningWithCorrelationf", "31234warn test message")
}

func TestMainLoggerViaPackageInactiveWarningWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.WarningWithCorrelationf("1234", "warn test %s", "message")

	assertNoMessage(t, "WarningWithCorrelationf")
}

func TestMainLoggerViaPackageWarningCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaPackage(t, "WarningCustomf", "3 map[test:123]warn test message")
}

func TestMainLoggerViaPackageInactiveWarningCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessage(t, "WarningCustomf")
}

func TestMainLoggerViaPackageWarningCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("WARN", "WARN")

	mainLogger.WarningCtxf(testDummyContext, "warn test %s", "message")

	assertMessageViaPackage(t, "WarningCtxf", "31234warn test message")
}

func TestMainLoggerViaPackageInactiveWarningCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.WarningCtxf(testDummyContext, "warn test %s", "message")

	assertNoMessage(t, "WarningCtxf")
}

// -------------------
//
// Warning Common Block
//
// -------------------

func TestMainLoggerViaCommonWarning(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.Warning("warn test message")

	assertMessageViaCommon(t, "Warning", "3warn test message")
}

func TestMainLoggerViaCommonInactiveWarning(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.Warning("warn test message")

	assertNoMessage(t, "Warning")
}

func TestMainLoggerViaCommonWarningWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningWithCorrelation("1234", "warn test message")

	assertMessageViaCommon(t, "WarningWithCorrelation", "31234warn test message")
}

func TestMainLoggerViaCommonInactiveWarningWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningWithCorrelation("1234", "warn test message")

	assertNoMessage(t, "WarningWithCorrelation")
}

func TestMainLoggerViaCommonWarningCustom(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertMessageViaCommon(t, "WarningCustom", "3 map[test:123]warn test message")
}

func TestMainLoggerViaCommonInactiveWarningCustom(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertNoMessage(t, "WarningCustom")
}

func TestMainLoggerViaCommonWarningCtx(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCtx(testDummyContext, "warn test message")

	assertMessageViaCommon(t, "WarningCtx", "31234warn test message")
}

func TestMainLoggerViaCommonInactiveWarningCtx(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCtx(testDummyContext, "warn test message")

	assertNoMessage(t, "WarningCtx")
}

func TestMainLoggerViaCommonWarningf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.Warningf("warn test %s", "message")

	assertMessageViaCommon(t, "Warningf", "3warn test message")
}

func TestMainLoggerViaCommonInactiveWarningf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.Warningf("warn test %s", "message")

	assertNoMessage(t, "Warningf")
}

func TestMainLoggerViaCommonWarningWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningWithCorrelationf("1234", "warn test %s", "message")

	assertMessageViaCommon(t, "WarningWithCorrelationf", "31234warn test message")
}

func TestMainLoggerViaCommonInactiveWarningWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningWithCorrelationf("1234", "warn test %s", "message")

	assertNoMessage(t, "WarningWithCorrelationf")
}

func TestMainLoggerViaCommonWarningCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCustomf", "3 map[test:123]warn test message")
}

func TestMainLoggerViaCommonInactiveWarningCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessage(t, "WarningCustomf")
}

func TestMainLoggerViaCommonWarningCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCtxf(testDummyContext, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCtxf", "31234warn test message")
}

func TestMainLoggerViaCommonInactiveWarningCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCtxf(testDummyContext, "warn test %s", "message")

	assertNoMessage(t, "WarningCtxf")
}

// -------------------
//
// Warning Only Common Block
//
// -------------------

func TestMainLoggerOnlyCommonWarning(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.Warning("warn test message")

	assertMessageViaCommon(t, "Warning", "3warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarning(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.Warning("warn test message")

	assertNoMessage(t, "Warning")
}

func TestMainLoggerOnlyCommonWarningWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningWithCorrelation("1234", "warn test message")

	assertMessageViaCommon(t, "WarningWithCorrelation", "31234warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarningWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.WarningWithCorrelation("1234", "warn test message")

	assertNoMessage(t, "WarningWithCorrelation")
}

func TestMainLoggerOnlyCommonWarningCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertMessageViaCommon(t, "WarningCustom", "3 map[test:123]warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarningCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertNoMessage(t, "WarningCustom")
}

func TestMainLoggerOnlyCommonWarningCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningCtx(testDummyContext, "warn test message")

	assertMessageViaCommon(t, "WarningCtx", "31234warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarningCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.WarningCtx(testDummyContext, "warn test message")

	assertNoMessage(t, "WarningCtx")
}

func TestMainLoggerOnlyCommonWarningf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.Warningf("warn test %s", "message")

	assertMessageViaCommon(t, "Warningf", "3warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarningf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.Warningf("warn test %s", "message")

	assertNoMessage(t, "Warningf")
}

func TestMainLoggerOnlyCommonWarningWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningWithCorrelationf("1234", "warn test %s", "message")

	assertMessageViaCommon(t, "WarningWithCorrelationf", "31234warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarningWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.WarningWithCorrelationf("1234", "warn test %s", "message")

	assertNoMessage(t, "WarningWithCorrelationf")
}

func TestMainLoggerOnlyCommonWarningCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCustomf", "3 map[test:123]warn test message")
}
func TestMainLoggerOnlyCommonInactiveWarningCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessage(t, "WarningCustom")
}

func TestMainLoggerOnlyCommonWarningCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningCtxf(testDummyContext, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCtxf", "31234warn test message")
}

func TestMainLoggerOnlyCommonInactiveWarningCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.WarningCtxf(testDummyContext, "warn test %s", "message")

	assertNoMessage(t, "WarningCtxf")
}

// -------------------
//
// Warning With Panic Block
//
// -------------------

func TestMainLoggerWarningWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningWithPanic("warn test message")

	assertMessageWithPanic(t, "WarningWithPanic", "3warn test message")
}

func TestMainLoggerInactiveWarningWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningWithPanic("warn test message")

	assertNoMessageWithPanic(t, "WarningWithPanic")
}

func TestMainLoggerWarningWithCorrelationAndPanic(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningWithCorrelationAndPanic("1234", "warn test message")

	assertMessageWithPanic(t, "WarningWithCorrelationAndPanic", "31234warn test message")
}

func TestMainLoggerInactiveWarningWithCorrelationAndPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningWithCorrelationAndPanic("1234", "warn test message")

	assertNoMessageWithPanic(t, "WarningWithCorrelationAndPanic")
}

func TestMainLoggerWarningCustomWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "warn test message")

	assertMessageWithPanic(t, "WarningCustomWithPanic", "3 map[test:123]warn test message")
}

func TestMainLoggerInactiveWarningCustomWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCustomWithPanic(map[string]any{"test": 123}, "warn test message")

	assertNoMessageWithPanic(t, "WarningCustomWithPanic")
}

func TestMainLoggerWarningCtxWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCtxWithPanic(testDummyContext, "warn test message")

	assertMessageWithPanic(t, "WarningCtxWithPanic", "31234warn test message")
}

func TestMainLoggerInactiveWarningCtxWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCtxWithPanic(testDummyContext, "warn test message")

	assertNoMessageWithPanic(t, "WarningCtxWithPanic")
}

func TestMainLoggerWarningWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningWithPanicf("warn test %s", "message")

	assertMessageWithPanic(t, "WarningWithPanicf", "3warn test message")
}

func TestMainLoggerInactiveWarningWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningWithPanicf("warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningWithPanicf")
}

func TestMainLoggerWarningWithCorrelationAndPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningWithCorrelationAndPanicf("1234", "warn test %s", "message")

	assertMessageWithPanic(t, "WarningWithCorrelationAndPanicf", "31234warn test message")
}

func TestMainLoggerInactiveWarningWithCorrelationAndPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningWithCorrelationAndPanicf("1234", "warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningWithCorrelationAndPanicf")
}

func TestMainLoggerWarningCustomWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageWithPanic(t, "WarningCustomWithPanicf", "3 map[test:123]warn test message")
}

func TestMainLoggerInactiveWarningCustomWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCustomWithPanicf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningCustomWithPanicf")
}

func TestMainLoggerWarningCtxWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("WARN", "WARN")

	mainLogger.WarningCtxWithPanicf(testDummyContext, "warn test %s", "message")

	assertMessageWithPanic(t, "WarningCtxWithPanicf", "31234warn test message")
}

func TestMainLoggerInactiveWarningCtxWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.WarningCtxWithPanicf(testDummyContext, "warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningCtxWithPanicf")
}

// -------------------
//
// Error Package Block
//
// -------------------

func TestMainLoggerViaPackageError(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.Error("error test message")

	assertMessageViaPackage(t, "Error", "2error test message")
}

func TestMainLoggerViaPackageInactiveError(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.Error("error test message")

	assertNoMessage(t, "Error")
}

func TestMainLoggerViaPackageErrorWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.ErrorWithCorrelation("1234", "error test message")

	assertMessageViaPackage(t, "ErrorWithCorrelation", "21234error test message")
}

func TestMainLoggerViaPackageInactiveErrorWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.ErrorWithCorrelation("1234", "error test message")

	assertNoMessage(t, "ErrorWithCorrelation")
}

func TestMainLoggerViaPackageErrorCustom(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertMessageViaPackage(t, "ErrorCustom", "2 map[test:123]error test message")
}

func TestMainLoggerViaPackageInactiveErrorCustom(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertNoMessage(t, "ErrorCustom")
}

func TestMainLoggerViaPackageErrorCtx(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.ErrorCtx(testDummyContext, "error test message")

	assertMessageViaPackage(t, "ErrorCtx", "21234error test message")
}

func TestMainLoggerViaPackageInactiveErrorCtx(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.ErrorCtx(testDummyContext, "error test message")

	assertNoMessage(t, "ErrorCtx")
}

func TestMainLoggerViaPackageErrorf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.Errorf("error test %s", "message")

	assertMessageViaPackage(t, "Errorf", "2error test message")
}

func TestMainLoggerViaPackageInactiveErrorf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.Errorf("error test %s", "message")

	assertNoMessage(t, "Errorf")
}

func TestMainLoggerViaPackageErrorWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.ErrorWithCorrelationf("1234", "error test %s", "message")

	assertMessageViaPackage(t, "ErrorWithCorrelationf", "21234error test message")
}

func TestMainLoggerViaPackageInactiveErrorWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.ErrorWithCorrelationf("1234", "error test %s", "message")

	assertNoMessage(t, "ErrorWithCorrelationf")
}

func TestMainLoggerViaPackageErrorCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageViaPackage(t, "ErrorCustomf", "2 map[test:123]error test message")
}

func TestMainLoggerViaPackageInactiveErrorCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessage(t, "ErrorCustomf")
}

func TestMainLoggerViaPackageErrorCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("ERROR", "ERROR")

	mainLogger.ErrorCtxf(testDummyContext, "error test %s", "message")

	assertMessageViaPackage(t, "ErrorCtxf", "21234error test message")
}

func TestMainLoggerViaPackageInactiveErrorCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.ErrorCtxf(testDummyContext, "error test %s", "message")

	assertNoMessage(t, "ErrorCtxf")
}

// -------------------
//
// Error Common Block
//
// -------------------

func TestMainLoggerViaCommonError(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.Error("error test message")

	assertMessageViaCommon(t, "Error", "2error test message")
}

func TestMainLoggerViaCommonInactiveError(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.Error("error test message")

	assertNoMessage(t, "Error")
}

func TestMainLoggerViaCommonErrorWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorWithCorrelation("1234", "error test message")

	assertMessageViaCommon(t, "ErrorWithCorrelation", "21234error test message")
}

func TestMainLoggerViaCommonInactiveErrorWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorWithCorrelation("1234", "error test message")

	assertNoMessage(t, "ErrorWithCorrelation")
}

func TestMainLoggerViaCommonErrorCustom(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertMessageViaCommon(t, "ErrorCustom", "2 map[test:123]error test message")
}

func TestMainLoggerViaCommonInactiveErrorCustom(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertNoMessage(t, "ErrorCustom")
}

func TestMainLoggerViaCommonErrorCtx(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCtx(testDummyContext, "error test message")

	assertMessageViaCommon(t, "ErrorCtx", "21234error test message")
}

func TestMainLoggerViaCommonInactiveErrorCtx(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCtx(testDummyContext, "error test message")

	assertNoMessage(t, "ErrorCtx")
}

func TestMainLoggerViaCommonErrorf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.Errorf("error test %s", "message")

	assertMessageViaCommon(t, "Errorf", "2error test message")
}

func TestMainLoggerViaCommonInactiveErrorf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.Errorf("error test %s", "message")

	assertNoMessage(t, "Errorf")
}

func TestMainLoggerViaCommonErrorWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorWithCorrelationf("1234", "error test %s", "message")

	assertMessageViaCommon(t, "ErrorWithCorrelationf", "21234error test message")
}

func TestMainLoggerViaCommonInactiveErrorWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorWithCorrelationf("1234", "error test %s", "message")

	assertNoMessage(t, "ErrorWithCorrelationf")
}

func TestMainLoggerViaCommonErrorCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCustomf", "2 map[test:123]error test message")
}

func TestMainLoggerViaCommonInactiveErrorCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessage(t, "ErrorCustomf")
}

func TestMainLoggerViaCommonErrorCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCtxf(testDummyContext, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCtxf", "21234error test message")
}

func TestMainLoggerViaCommonInactiveErrorCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCtxf(testDummyContext, "error test %s", "message")

	assertNoMessage(t, "ErrorCtxf")
}

// -------------------
//
// Error Only Common Block
//
// -------------------

func TestMainLoggerOnlyCommonError(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.Error("error test message")

	assertMessageViaCommon(t, "Error", "2error test message")
}

func TestMainLoggerOnlyCommonInactiveError(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.Error("error test message")

	assertNoMessage(t, "Error")
}

func TestMainLoggerOnlyCommonErrorWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.ErrorWithCorrelation("1234", "error test message")

	assertMessageViaCommon(t, "ErrorWithCorrelation", "21234error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.ErrorWithCorrelation("1234", "error test message")

	assertNoMessage(t, "ErrorWithCorrelation")
}

func TestMainLoggerOnlyCommonErrorCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertMessageViaCommon(t, "ErrorCustom", "2 map[test:123]error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertNoMessage(t, "ErrorCustom")
}

func TestMainLoggerOnlyCommonErrorCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.ErrorCtx(testDummyContext, "error test message")

	assertMessageViaCommon(t, "ErrorCtx", "21234error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.ErrorCtx(testDummyContext, "error test message")

	assertNoMessage(t, "ErrorWithCorErrorCtxelation")
}

func TestMainLoggerOnlyCommonErrorf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.Errorf("error test %s", "message")

	assertMessageViaCommon(t, "Errorf", "2error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.Errorf("error test %s", "message")

	assertNoMessage(t, "Errorf")
}

func TestMainLoggerOnlyCommonErrorWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.ErrorWithCorrelationf("1234", "error test %s", "message")

	assertMessageViaCommon(t, "ErrorWithCorrelationf", "21234error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.ErrorWithCorrelationf("1234", "error test %s", "message")

	assertNoMessage(t, "ErrorWithCorrelationf")
}

func TestMainLoggerOnlyCommonErrorCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCustomf", "2 map[test:123]error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessage(t, "ErrorCustomf")
}

func TestMainLoggerOnlyCommonErrorCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("ERROR")

	mainLogger.ErrorCtxf(testDummyContext, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCtxf", "21234error test message")
}

func TestMainLoggerOnlyCommonInactiveErrorCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.ErrorCtxf(testDummyContext, "error test %s", "message")

	assertNoMessage(t, "ErrorCtxf")
}

// -------------------
//
// Error With Panic Block
//
// -------------------

func TestMainLoggerErrorWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorWithPanic("error test message")

	assertMessageWithPanic(t, "ErrorWithPanic", "2error test message")
}

func TestMainLoggerInactiveErrorWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorWithPanic("error test message")

	assertNoMessageWithPanic(t, "ErrorWithPanic")
}

func TestMainLoggerErrorWithCorrelationAndPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorWithCorrelationAndPanic("1234", "error test message")

	assertMessageWithPanic(t, "ErrorWithCorrelationAndPanic", "21234error test message")
}

func TestMainLoggerInactiveErrorWithCorrelationAndPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorWithCorrelationAndPanic("1234", "error test message")

	assertNoMessageWithPanic(t, "ErrorWithCorrelationAndPanic")
}

func TestMainLoggerErrorCustomWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "error test message")

	assertMessageWithPanic(t, "ErrorCustomWithPanic", "2 map[test:123]error test message")
}

func TestMainLoggerInactiveErrorCustomWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCustomWithPanic(map[string]any{"test": 123}, "error test message")

	assertNoMessageWithPanic(t, "ErrorCustomWithPanic")
}

func TestMainLoggerErrorCtxWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCtxWithPanic(testDummyContext, "error test message")

	assertMessageWithPanic(t, "ErrorCtxWithPanic", "21234error test message")
}

func TestMainLoggerInactiveErrorCtxWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCtxWithPanic(testDummyContext, "error test message")

	assertNoMessageWithPanic(t, "ErrorCtxWithPanic")
}

func TestMainLoggerErrorWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorWithPanicf("error test %s", "message")

	assertMessageWithPanic(t, "ErrorWithPanicf", "2error test message")
}

func TestMainLoggerInactiveErrorWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorWithPanicf("error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorWithPanicf")
}

func TestMainLoggerErrorWithCorrelationAndPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorWithCorrelationAndPanicf("1234", "error test %s", "message")

	assertMessageWithPanic(t, "ErrorWithCorrelationAndPanicf", "21234error test message")
}

func TestMainLoggerInactiveErrorWithCorrelationAndPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorWithCorrelationAndPanicf("1234", "error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorWithCorrelationAndPanicf")
}

func TestMainLoggerErrorCustomWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageWithPanic(t, "ErrorCustomWithPanicf", "2 map[test:123]error test message")
}

func TestMainLoggerInactiveErrorCustomWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCustomWithPanicf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorCustomWithPanicf")
}

func TestMainLoggerErrorCtxWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("ERROR", "ERROR")

	mainLogger.ErrorCtxWithPanicf(testDummyContext, "error test %s", "message")

	assertMessageWithPanic(t, "ErrorCtxWithPanicf", "21234error test message")
}

func TestMainLoggerInactiveErrorCtxWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.ErrorCtxWithPanicf(testDummyContext, "error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorCtxWithPanicf")
}

// -------------------
//
// Fatal Package Block
//
// -------------------

func TestMainLoggerViaPackageFatal(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.Fatal("fatal test message")

	assertMessageViaPackage(t, "Fatal", "1fatal test message")
}

func TestMainLoggerViaPackageInactiveFatal(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.Fatal("fatal test message")

	assertNoMessage(t, "Fatal")
}

func TestMainLoggerViaPackageFatalWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelation("1234", "fatal test message")

	assertMessageViaPackage(t, "FatalWithCorrelation", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalWithCorrelation(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalWithCorrelation("1234", "fatal test message")

	assertNoMessage(t, "FatalWithCorrelation")
}

func TestMainLoggerViaPackageFatalCustom(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaPackage(t, "FatalCustom", "1 map[test:123]fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCustom(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertNoMessage(t, "FatalCustom")
}

func TestMainLoggerViaPackageFatalCtx(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCtx(testDummyContext, "fatal test message")

	assertMessageViaPackage(t, "FatalCtx", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCtx(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCtx(testDummyContext, "fatal test message")

	assertNoMessage(t, "FatalCtx")
}

func TestMainLoggerViaPackageFatalf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.Fatalf("fatal test %s", "message")

	assertMessageViaPackage(t, "Fatalf", "1fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.Fatalf("fatal test %s", "message")

	assertNoMessage(t, "Fatalf")
}

func TestMainLoggerViaPackageFatalWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertMessageViaPackage(t, "FatalWithCorrelationf", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalWithCorrelationf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertNoMessage(t, "FatalWithCorrelationf")
}

func TestMainLoggerViaPackageFatalCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaPackage(t, "FatalCustomf", "1 map[test:123]fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCustomf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessage(t, "FatalCustomf")
}

func TestMainLoggerViaPackageFatalCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertMessageViaPackage(t, "FatalCtxf", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCtxf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertNoMessage(t, "FatalCtxf")
}

// -------------------
//
// Fatal Common Block
//
// -------------------

func TestMainLoggerViaCommonFatal(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.Fatal("fatal test message")

	assertMessageViaCommon(t, "Fatal", "1fatal test message")
}

func TestMainLoggerViaCommonInactiveFatal(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.Fatal("fatal test message")

	assertNoMessage(t, "Fatal")
}

func TestMainLoggerViaCommonFatalWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelation("1234", "fatal test message")

	assertMessageViaCommon(t, "FatalWithCorrelation", "11234fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalWithCorrelation(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithCorrelation("1234", "fatal test message")

	assertNoMessage(t, "FatalWithCorrelation")
}

func TestMainLoggerViaCommonFatalCustom(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaCommon(t, "FatalCustom", "1 map[test:123]fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalCustom(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertNoMessage(t, "FatalCustom")
}

func TestMainLoggerViaCommonFatalCtx(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCtx(testDummyContext, "fatal test message")

	assertMessageViaCommon(t, "FatalCtx", "11234fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalCtx(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCtx(testDummyContext, "fatal test message")

	assertNoMessage(t, "FatalCtx")
}

func TestMainLoggerViaCommonFatalf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.Fatalf("fatal test %s", "message")

	assertMessageViaCommon(t, "Fatalf", "1fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.Fatalf("fatal test %s", "message")

	assertNoMessage(t, "Fatalf")
}

func TestMainLoggerViaCommonFatalWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalWithCorrelationf", "11234fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalWithCorrelationf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertNoMessage(t, "FatalWithCorrelationf")
}

func TestMainLoggerViaCommonFatalCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCustomf", "1 map[test:123]fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalCustomf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessage(t, "FatalCustomf")
}

func TestMainLoggerViaCommonFatalCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCtxf", "11234fatal test message")
}

func TestMainLoggerViaCommonInactiveFatalCtxf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertNoMessage(t, "FatalCtxf")
}

// -------------------
//
// Fatal Only Common Block
//
// -------------------

func TestMainLoggerOnlyCommonFatal(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.Fatal("fatal test message")

	assertMessageViaCommon(t, "Fatal", "1fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatal(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.Fatal("fatal test message")

	assertNoMessage(t, "Fatal")
}

func TestMainLoggerOnlyCommonFatalWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalWithCorrelation("1234", "fatal test message")

	assertMessageViaCommon(t, "FatalWithCorrelation", "11234fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalWithCorrelation(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.FatalWithCorrelation("1234", "fatal test message")

	assertNoMessage(t, "FatalWithCorrelation")
}

func TestMainLoggerOnlyCommonFatalCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaCommon(t, "FatalCustom", "1 map[test:123]fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalCustom(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertNoMessage(t, "FatalCustom")
}

func TestMainLoggerOnlyCommonFatalCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalCtx(testDummyContext, "fatal test message")

	assertMessageViaCommon(t, "FatalCtx", "11234fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalCtx(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.FatalCtx(testDummyContext, "fatal test message")

	assertNoMessage(t, "FatalCtx")
}

func TestMainLoggerOnlyCommonFatalf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.Fatalf("fatal test %s", "message")

	assertMessageViaCommon(t, "Fatalf", "1fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalf(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.Fatalf("fatal test %s", "message")

	assertNoMessage(t, "Fatalf")
}

func TestMainLoggerOnlyCommonFatalWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalWithCorrelationf", "11234fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalWithCorrelationf(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertNoMessage(t, "FatalWithCorrelationf")
}

func TestMainLoggerOnlyCommonFatalCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCustomf", "1 map[test:123]fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessage(t, "FatalCustomf")
}

func TestMainLoggerOnlyCommonFatalCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCtxf", "11234fatal test message")
}

func TestMainLoggerOnlyCommonInactiveFatalCtxf(t *testing.T) {
	initMainLoggerOnlyCommonTest("OFF")

	mainLogger.FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertNoMessage(t, "FatalCtxf")
}

// -------------------
//
// Fatal With Panic Block
//
// -------------------

func TestMainLoggerFatalWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithPanic("fatal test message")

	assertMessageWithPanic(t, "FatalWithPanic", "1fatal test message")
}

func TestMainLoggerInactiveFatalWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithPanic("fatal test message")

	assertNoMessageWithPanic(t, "FatalWithPanic")
}

func TestMainLoggerFatalWithCorrelationAndPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationAndPanic("1234", "fatal test message")

	assertMessageWithPanic(t, "FatalWithCorrelationAndPanic", "11234fatal test message")
}

func TestMainLoggerInactiveFatalWithCorrelationAndPanic(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationAndPanic("1234", "fatal test message")

	assertNoMessageWithPanic(t, "FatalWithCorrelationAndPanic")
}

func TestMainLoggerFatalCustomWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "fatal test message")

	assertMessageWithPanic(t, "FatalCustomWithPanic", "1 map[test:123]fatal test message")
}

func TestMainLoggerInactiveFatalCustomWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCustomWithPanic(map[string]any{"test": 123}, "fatal test message")

	assertNoMessageWithPanic(t, "FatalCustomWithPanic")
}

func TestMainLoggerFatalCtxWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCtxWithPanic(testDummyContext, "fatal test message")

	assertMessageWithPanic(t, "FatalCtxWithPanic", "11234fatal test message")
}

func TestMainLoggerInactiveFatalCtxWithPanic(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCtxWithPanic(testDummyContext, "fatal test message")

	assertNoMessageWithPanic(t, "FatalCtxWithPanic")
}

func TestMainLoggerFatalWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithPanicf("fatal test %s", "message")

	assertMessageWithPanic(t, "FatalWithPanicf", "1fatal test message")
}

func TestMainLoggerInactiveFatalWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithPanicf("fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalWithPanicf")
}

func TestMainLoggerFatalWithCorrelationAndPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationAndPanicf("1234", "fatal test %s", "message")

	assertMessageWithPanic(t, "FatalWithCorrelationAndPanicf", "11234fatal test message")
}

func TestMainLoggerInactiveFatalWithCorrelationAndPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationAndPanicf("1234", "fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalWithCorrelationAndPanicf")
}

func TestMainLoggerFatalCustomWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageWithPanic(t, "FatalCustomWithPanicf", "1 map[test:123]fatal test message")
}

func TestMainLoggerInactiveFatalCustomWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCustomWithPanicf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalCustomWithPanicf")
}

func TestMainLoggerFatalCtxWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCtxWithPanicf(testDummyContext, "fatal test %s", "message")

	assertMessageWithPanic(t, "FatalCtxWithPanicf", "11234fatal test message")
}

func TestMainLoggerInactiveFatalCtxWithPanicf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCtxWithPanicf(testDummyContext, "fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalCtxWithPanicf")
}

// -------------------
//
// Fatal With Exit Block
//
// -------------------

func TestMainLoggerFatalWithExit(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithExit("fatal test message")

	assertMessageWithExit(t, "FatalWithExit", "1fatal test message")
}

func TestMainLoggerInactiveFatalWithExit(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithExit("fatal test message")

	assertNoMessageWithExit(t, "FatalWithExit")
}

func TestMainLoggerFatalWithCorrelationAndExit(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationAndExit("1234", "fatal test message")

	assertMessageWithExit(t, "FatalWithCorrelationAndExit", "11234fatal test message")
}

func TestMainLoggerInactiveFatalWithCorrelationAndExit(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationAndExit("1234", "fatal test message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExit")
}

func TestMainLoggerFatalCustomWithExit(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertMessageWithExit(t, "FatalCustomWithExit", "1 map[test:123]fatal test message")
}

func TestMainLoggerInactiveFatalCustomWithExit(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertNoMessageWithExit(t, "FatalCustomWithExit")
}

func TestMainLoggerFatalCtxWithExit(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCtxWithExit(testDummyContext, "fatal test message")

	assertMessageWithExit(t, "FatalCtxWithExit", "11234fatal test message")
}

func TestMainLoggerInactiveFatalCtxWithExit(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCtxWithExit(testDummyContext, "fatal test message")

	assertNoMessageWithExit(t, "FatalCtxWithExit")
}

func TestMainLoggerFatalWithExitf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithExitf("fatal test %s", "message")

	assertMessageWithExit(t, "FatalWithExitf", "1fatal test message")
}

func TestMainLoggerInactiveFatalWithExitf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithExitf("fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithExitf")
}

func TestMainLoggerFatalWithCorrelationAndExitf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertMessageWithExit(t, "FatalWithCorrelationAndExitf", "11234fatal test message")
}

func TestMainLoggerInactiveFatalWithCorrelationAndExitf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExitf")
}

func TestMainLoggerFatalCustomWithExitf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageWithExit(t, "FatalCustomWithExitf", "1 map[test:123]fatal test message")
}

func TestMainLoggerInactiveFatalCustomWithExitf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCustomWithExitf")
}

func TestMainLoggerFatalCtxWithExitf(t *testing.T) {
	initMainLoggerViaCommonTest("FATAL", "FATAL")

	mainLogger.FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertMessageWithExit(t, "FatalCtxWithExitf", "11234fatal test message")
}

func TestMainLoggerInactiveFatalCtxWithExitf(t *testing.T) {
	initMainLoggerViaCommonTest("OFF", "OFF")

	mainLogger.FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCtxWithExitf")
}

// -------------------
//
// Fatal Package With Exit Block
//
// -------------------

func TestMainLoggerViaPackageFatalWithExit(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalWithExit("fatal test message")

	assertMessageViaPackageWithExit(t, "FatalWithExit", "1fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalWithExit(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalWithExit("fatal test message")

	assertNoMessageWithExit(t, "FatalWithExit")
}

func TestMainLoggerViaPackageFatalWithCorrelationAndExit(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationAndExit("1234", "fatal test message")

	assertMessageViaPackageWithExit(t, "FatalWithCorrelationAndExit", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalWithCorrelationAndExit(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationAndExit("1234", "fatal test message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExit")
}

func TestMainLoggerViaPackageFatalCustomWithExit(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaPackageWithExit(t, "FatalCustomWithExit", "1 map[test:123]fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCustomWithExit(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertNoMessageWithExit(t, "FatalCustomWithExit")
}

func TestMainLoggerViaPackageFatalCtxWithExit(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCtxWithExit(testDummyContext, "fatal test message")

	assertMessageViaPackageWithExit(t, "FatalCtxWithExit", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCtxWithExit(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCtxWithExit(testDummyContext, "fatal test message")

	assertNoMessageWithExit(t, "FatalCtxWithExit")
}

func TestMainLoggerViaPackageFatalWithExitf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalWithExitf("fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalWithExitf", "1fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalWithExitf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalWithExitf("fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithExitf")
}

func TestMainLoggerViaPackageFatalWithCorrelationAndExitf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalWithCorrelationAndExitf", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalWithCorrelationAndExitf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExitf")
}

func TestMainLoggerViaPackageFatalCustomWithExitf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalCustomWithExitf", "1 map[test:123]fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCustomWithExitf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCustomWithExitf")
}

func TestMainLoggerViaPackageFatalCtxWithExitf(t *testing.T) {
	initMainLoggerViaPackageTest("FATAL", "FATAL")

	mainLogger.FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalCtxWithExitf", "11234fatal test message")
}

func TestMainLoggerViaPackageInactiveFatalCtxWithExitf(t *testing.T) {
	initMainLoggerViaPackageTest("OFF", "OFF")

	mainLogger.FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCtxWithExitf")
}

// -------------------
//
// Fatal Only Common Block
//
// -------------------

func TestMainLoggerOnlyCommonFatalWithExit(t *testing.T) {
	initMainLoggerOnlyCommonTest("FATAL")

	mainLogger.FatalWithExit("fatal test message")

	assertMessageWithExit(t, "Fatal", "1fatal test message")
}

// -------------------
//
// Assert Block
//
// -------------------

func assertMessageViaPackage(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainPackageLoggerAppender, &testMainGeneralLoggerAppender, "package", "common", message)
	assertPanicAndExitMockNotActivated(t)
}

func assertMessageViaCommon(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainGeneralLoggerAppender, &testMainPackageLoggerAppender, "common", "package", message)
	assertPanicAndExitMockNotActivated(t)
}

func assertMessageWithPanic(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainGeneralLoggerAppender, &testMainPackageLoggerAppender, "common", "package", message)
	assertPanicMockActivated(t)
}

func assertMessageWithExit(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainGeneralLoggerAppender, &testMainPackageLoggerAppender, "common", "package", message)
	assertExitMockActivated(t)
}

func assertMessageViaPackageWithExit(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainPackageLoggerAppender, &testMainGeneralLoggerAppender, "package", "common", message)
	assertExitMockActivated(t)
}

func assertMessage(t *testing.T, methodName string, appenderWithMessage *appender.Appender, appenderWithoutMessage *appender.Appender, withMessageName string, withoutMessageName string, message string) {
	testutil.AssertEquals(1, len(*(*appenderWithMessage).(TestAppender).content), t, withMessageName+" "+methodName+": len(content)")
	testutil.AssertEquals(0, len(*(*appenderWithoutMessage).(TestAppender).content), t, withoutMessageName+" "+methodName+": len(content)")
	testutil.AssertEquals(message, (*(*appenderWithMessage).(TestAppender).content)[0], t, withMessageName+" "+methodName+": content[0]")
}

func assertNoMessageWithPanic(t *testing.T, methodName string) {
	assertNoMessage(t, methodName)
	assertPanicMockActivated(t)
}

func assertNoMessageWithExit(t *testing.T, methodName string) {
	assertNoMessage(t, methodName)
	assertExitMockActivated(t)
}

func assertNoMessage(t *testing.T, methodName string) {
	testutil.AssertEquals(0, len(*testMainGeneralLoggerAppender.(TestAppender).content), t, "common "+methodName+": len(content)")
	testutil.AssertEquals(0, len(*testMainPackageLoggerAppender.(TestAppender).content), t, "package  "+methodName+": len(content)")
}

func assertEnabled(t *testing.T, isDebug bool, isInfo bool, isWarn bool, isError bool, isFatal bool) {
	testutil.AssertEquals(isDebug, mainLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertEquals(isInfo, mainLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertEquals(isWarn, mainLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertEquals(isError, mainLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertEquals(isFatal, mainLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}
