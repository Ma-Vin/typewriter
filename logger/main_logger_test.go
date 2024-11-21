package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/testutil"
)

var testMainCommonLoggerAppender appender.Appender = TestAppender{content: &[]string{}}
var testMainPackageLoggerAppender appender.Appender = TestAppender{content: &[]string{}}

var testMainCommonLogger CommonLogger
var testMainPackageLogger CommonLogger
var mainLogger MainLogger

func clearMainLoggerTestEnv() {
	os.Unsetenv(config.DEFAULT_LOG_LEVEL_ENV_NAME)
	os.Unsetenv(config.DEFAULT_LOG_LEVEL_ENV_NAME + "_LOGGER")
}

func initMainLoggerTest(envCommonLogLevel string, envPackageLogLevel string, packageName string) {
	clearMainLoggerTestEnv()

	*testMainCommonLoggerAppender.(TestAppender).content = []string{}
	*testMainPackageLoggerAppender.(TestAppender).content = []string{}

	testMainCommonLogger = CreateCommonLogger(&testMainCommonLoggerAppender, config.SeverityLevelMap[envCommonLogLevel])
	testMainPackageLogger = CreateCommonLogger(&testMainPackageLoggerAppender, config.SeverityLevelMap[envPackageLogLevel])

	mockPanicAndExitAtCommonLogger = true
	panicMockActivated = false
	exitMockAcitvated = false
	testCommonLoggerCounterAppenderClosed = 0
	testCommonLoggerCounterAppenderClosedExpected = 2

	mainLogger = MainLogger{
		commonLogger:       &testMainCommonLogger,
		existPackageLogger: true,
		packageLoggers:     map[string]*CommonLogger{strings.ToUpper(packageName): &testMainPackageLogger},
	}
}

func initMainLoggerViaPackageTest(envCommonLogLevel string, envPackageLogLevel string) {
	initMainLoggerTest(envCommonLogLevel, envPackageLogLevel, "LOGGER")
}

func initMainLoggerViaCommonTest(envCommonLogLevel string, envPackageLogLevel string) {
	initMainLoggerTest(envCommonLogLevel, envPackageLogLevel, "OTHER")
}

func initMainLoggerOnlyCommonTest(envCommonLogLevel string) {
	initMainLoggerViaCommonTest(envCommonLogLevel, envCommonLogLevel)
	testCommonLoggerCounterAppenderClosedExpected = 1
	mainLogger = MainLogger{
		commonLogger:       &testMainCommonLogger,
		existPackageLogger: false,
		packageLoggers:     make(map[string]*CommonLogger),
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

func TestMainLoggerOnlyCommonWarningCustomf(t *testing.T) {
	initMainLoggerOnlyCommonTest("WARN")

	mainLogger.WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCustomf", "3 map[test:123]warn test message")
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
	assertMessage(t, methodName, &testMainPackageLoggerAppender, &testMainCommonLoggerAppender, "package", "common", message)
	assertPanicAndExitMockNotActivated(t)
}

func assertMessageViaCommon(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainCommonLoggerAppender, &testMainPackageLoggerAppender, "common", "package", message)
	assertPanicAndExitMockNotActivated(t)
}

func assertMessageWithPanic(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainCommonLoggerAppender, &testMainPackageLoggerAppender, "common", "package", message)
	assertPanicMockAcitvated(t)
}

func assertMessageWithExit(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainCommonLoggerAppender, &testMainPackageLoggerAppender, "common", "package", message)
	assertExitMockAcitvated(t)
}

func assertMessageViaPackageWithExit(t *testing.T, methodName string, message string) {
	assertMessage(t, methodName, &testMainPackageLoggerAppender, &testMainCommonLoggerAppender, "package", "common", message)
	assertExitMockAcitvated(t)
}

func assertMessage(t *testing.T, methodName string, appenderWithMessage *appender.Appender, appenderWithoutMessage *appender.Appender, withMessageName string, withoutMessageName string, message string) {
	testutil.AssertEquals(1, len(*(*appenderWithMessage).(TestAppender).content), t, withMessageName+" "+methodName+": len(content)")
	testutil.AssertEquals(0, len(*(*appenderWithoutMessage).(TestAppender).content), t, withoutMessageName+" "+methodName+": len(content)")
	testutil.AssertEquals(message, (*(*appenderWithMessage).(TestAppender).content)[0], t, withMessageName+" "+methodName+": content[0]")
}

func assertNoMessageWithPanic(t *testing.T, methodName string) {
	assertNoMessage(t, methodName)
	assertPanicMockAcitvated(t)
}

func assertNoMessageWithExit(t *testing.T, methodName string) {
	assertNoMessage(t, methodName)
	assertExitMockAcitvated(t)
}

func assertNoMessage(t *testing.T, methodName string) {
	testutil.AssertEquals(0, len(*testMainCommonLoggerAppender.(TestAppender).content), t, "common "+methodName+": len(content)")
	testutil.AssertEquals(0, len(*testMainPackageLoggerAppender.(TestAppender).content), t, "package  "+methodName+": len(content)")
}

func assertEnabled(t *testing.T, isDebug bool, isInfo bool, isWarn bool, isError bool, isFatal bool) {
	testutil.AssertEquals(isDebug, mainLogger.IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertEquals(isInfo, mainLogger.IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertEquals(isWarn, mainLogger.IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertEquals(isError, mainLogger.IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertEquals(isFatal, mainLogger.IsFatalEnabled(), t, "IsFatalEnabled")
}
