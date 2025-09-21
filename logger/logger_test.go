package logger

import (
	"os"
	"reflect"
	"testing"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
)

func initLoggerViaPackageTest(envCommonLogLevel string, envPackageLogLevel string) {
	// dummy config initialization, otherwise loggersInitialized=true will be without effect
	config.GetConfig()
	initMainLoggerViaPackageTest(envCommonLogLevel, envPackageLogLevel)
	mLogger = mainLogger
	loggersInitialized = true
}

func initLoggerViaCommonTest(envCommonLogLevel string, envPackageLogLevel string) {
	// dummy config initialization, otherwise loggersInitialized=true will be without effect
	config.GetConfig()
	initMainLoggerViaCommonTest(envCommonLogLevel, envPackageLogLevel)
	mLogger = mainLogger
	loggersInitialized = true
}

func initLoggerOnlyCommonTest(envCommonLogLevel string) {
	// dummy config initialization, otherwise loggersInitialized=true will be without effect
	config.GetConfig()
	initMainLoggerOnlyCommonTest(envCommonLogLevel)
	mLogger = mainLogger
	loggersInitialized = true
}

// -------------------
//
// Is Enabled Via Package Block
//
// -------------------

func TestEnableDebugSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "DEBUG")

	assertLoggerEnabled(t, true, true, true, true, true)
}

func TestEnableInformationSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "INFORMATION")

	assertLoggerEnabled(t, false, true, true, true, true)
}

func TestEnableInfoSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "INFO")

	assertLoggerEnabled(t, false, true, true, true, true)
}

func TestEnableWarningSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "WARNING")

	assertLoggerEnabled(t, false, false, true, true, true)
}

func TestEnableWarnSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "WARN")

	assertLoggerEnabled(t, false, false, true, true, true)
}

func TestEnableErrorSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "ERROR")

	assertLoggerEnabled(t, false, false, false, true, true)
}

func TestEnableFatalSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "FATAL")

	assertLoggerEnabled(t, false, false, false, false, true)
}

func TestEnableOffSeverityLoggerViaPackage(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	assertLoggerEnabled(t, false, false, false, false, false)
}

// -------------------
//
// Is Enabled Via Common Block
//
// -------------------

func TestEnableDebugSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "OFF")

	assertLoggerEnabled(t, true, true, true, true, true)
}

func TestEnableInformationSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("INFORMATION", "OFF")

	assertLoggerEnabled(t, false, true, true, true, true)
}

func TestEnableInfoSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("INFO", "OFF")

	assertLoggerEnabled(t, false, true, true, true, true)
}

func TestEnableWarningSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("WARNING", "OFF")

	assertLoggerEnabled(t, false, false, true, true, true)
}

func TestEnableWarnSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("WARN", "OFF")

	assertLoggerEnabled(t, false, false, true, true, true)
}

func TestEnableErrorSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "OFF")

	assertLoggerEnabled(t, false, false, false, true, true)
}

func TestEnableFatalSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "OFF")

	assertLoggerEnabled(t, false, false, false, false, true)
}

func TestEnableOffSeverityLoggerViaCommon(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	assertLoggerEnabled(t, false, false, false, false, false)
}

// -------------------
//
// Is Enabled Only Common Block
//
// -------------------

func TestEnableDebugSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	assertLoggerEnabled(t, true, true, true, true, true)
}

func TestEnableInformationSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("INFORMATION")

	assertLoggerEnabled(t, false, true, true, true, true)
}

func TestEnableInfoSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	assertLoggerEnabled(t, false, true, true, true, true)
}

func TestEnableWarningSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("WARNING")

	assertLoggerEnabled(t, false, false, true, true, true)
}

func TestEnableWarnSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	assertLoggerEnabled(t, false, false, true, true, true)
}

func TestEnableErrorSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	assertLoggerEnabled(t, false, false, false, true, true)
}

func TestEnableFatalSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	assertLoggerEnabled(t, false, false, false, false, true)
}

func TestEnableOffSeverityLoggerOnlyCommon(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	assertLoggerEnabled(t, false, false, false, false, false)
}

func assertLoggerEnabled(t *testing.T, isDebug bool, isInfo bool, isWarn bool, isError bool, isFatal bool) {
	testutil.AssertEquals(isDebug, IsDebugEnabled(), t, "IsDebugEnabled")
	testutil.AssertEquals(isInfo, IsInformationEnabled(), t, "IsInformationEnabled")
	testutil.AssertEquals(isWarn, IsWarningEnabled(), t, "IsWarningEnabled")
	testutil.AssertEquals(isError, IsErrorEnabled(), t, "IsErrorEnabled")
	testutil.AssertEquals(isFatal, IsFatalEnabled(), t, "IsFatalEnabled")
}

// -------------------
//
// Debug Package Block
//
// -------------------

func TestLoggerViaPackageDebug(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	Debug("debug test message")

	assertMessageViaPackage(t, "Debug", "5debug test message")
}

func TestLoggerViaPackageInactiveDebug(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestLoggerViaPackageDebugWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	DebugWithCorrelation("1234", "debug test message")

	assertMessageViaPackage(t, "DebugWithCorrelation", "51234debug test message")
}

func TestLoggerViaPackageInactiveDebugWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestLoggerViaPackageDebugCustom(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaPackage(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestLoggerViaPackageInactiveDebugCustom(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestLoggerViaPackageDebugCtx(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	DebugCtx(testDummyContext, "debug test message")

	assertMessageViaPackage(t, "DebugWithCorrelation", "51234debug test message")
}

func TestLoggerViaPackageInactiveDebugCtx(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestLoggerViaPackageDebugf(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	Debugf("debug test %s", "message")

	assertMessageViaPackage(t, "Debugf", "5debug test message")
}

func TestLoggerViaPackageInactiveDebugf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestLoggerViaPackageDebugWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaPackage(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestLoggerViaPackageInactiveDebugWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestLoggerViaPackageDebugCustomf(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaPackage(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestLoggerViaPackageInactiveDebugCustomf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestLoggerViaPackageDebugCtxf(t *testing.T) {
	initLoggerViaPackageTest("DEBUG", "DEBUG")

	DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaPackage(t, "DebugCtxf", "51234debug test message")
}

func TestLoggerViaPackageInactiveDebugCtxf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Debug Common Block
//
// -------------------

func TestLoggerViaCommonDebug(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	Debug("debug test message")

	assertMessageViaCommon(t, "Debug", "5debug test message")
}

func TestLoggerViaCommonInactiveDebug(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestLoggerViaCommonDebugWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	DebugWithCorrelation("1234", "debug test message")

	assertMessageViaCommon(t, "DebugWithCorrelation", "51234debug test message")
}

func TestLoggerViaCommonInactiveDebugWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestLoggerViaCommonDebugCustom(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaCommon(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestLoggerViaCommonInactiveDebugCustom(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestLoggerViaCommonDebugCtx(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	DebugCtx(testDummyContext, "debug test message")

	assertMessageViaCommon(t, "DebugCtx", "51234debug test message")
}

func TestLoggerViaCommonInactiveDebugCtx(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}

func TestLoggerViaCommonDebugf(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	Debugf("debug test %s", "message")

	assertMessageViaCommon(t, "Debugf", "5debug test message")
}

func TestLoggerViaCommonInactiveDebugf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestLoggerViaCommonDebugWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaCommon(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestLoggerViaCommonInactiveDebugWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestLoggerViaCommonDebugCustomf(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestLoggerViaCommonInactiveDebugCustomf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestLoggerViaCommonDebugCtxf(t *testing.T) {
	initLoggerViaCommonTest("DEBUG", "DEBUG")

	DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCtxf", "51234debug test message")
}

func TestLoggerViaCommonInactiveDebugCtxf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Debug Only Common Block
//
// -------------------

func TestLoggerOnlyCommonDebug(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	Debug("debug test message")

	assertMessageViaCommon(t, "Debug", "5debug test message")
}

func TestLoggerOnlyCommonInactiveDebug(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	Debug("debug test message")

	assertNoMessage(t, "Debug")
}

func TestLoggerOnlyCommonDebugWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	DebugWithCorrelation("1234", "debug test message")

	assertMessageViaCommon(t, "DebugWithCorrelation", "51234debug test message")
}

func TestLoggerOnlyCommonInactiveDebugWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	DebugWithCorrelation("1234", "debug test message")

	assertNoMessage(t, "DebugWithCorrelation")
}

func TestLoggerOnlyCommonDebugCustom(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertMessageViaCommon(t, "DebugCustom", "5 map[test:123]debug test message")
}

func TestLoggerOnlyCommonInactiveDebugCustom(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	DebugCustom(map[string]any{"test": 123}, "debug test message")

	assertNoMessage(t, "DebugCustom")
}

func TestLoggerOnlyCommonDebugCtx(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	DebugCtx(testDummyContext, "debug test message")

	assertMessageViaCommon(t, "DebugCtx", "51234debug test message")
}

func TestLoggerOnlyCommonInactiveDebugCtx(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	DebugCtx(testDummyContext, "debug test message")

	assertNoMessage(t, "DebugCtx")
}

func TestLoggerOnlyCommonDebugf(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	Debugf("debug test %s", "message")

	assertMessageViaCommon(t, "Debugf", "5debug test message")
}

func TestLoggerOnlyCommonInactiveDebugf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	Debugf("debug test %s", "message")

	assertNoMessage(t, "Debugf")
}

func TestLoggerOnlyCommonDebugWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	DebugWithCorrelationf("1234", "debug test %s", "message")

	assertMessageViaCommon(t, "DebugWithCorrelationf", "51234debug test message")
}

func TestLoggerOnlyCommonInactiveDebugWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	DebugWithCorrelationf("1234", "debug test %s", "message")

	assertNoMessage(t, "DebugWithCorrelationf")
}

func TestLoggerOnlyCommonDebugCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCustomf", "5 map[test:123]debug test message")
}

func TestLoggerOnlyCommonInactiveDebugCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	DebugCustomf(map[string]any{"test": 123}, "debug test %s", "message")

	assertNoMessage(t, "DebugCustomf")
}

func TestLoggerOnlyCommonDebugCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("DEBUG")

	DebugCtxf(testDummyContext, "debug test %s", "message")

	assertMessageViaCommon(t, "DebugCtxf", "51234debug test message")
}

func TestLoggerOnlyCommonInactiveDebugCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	DebugCtxf(testDummyContext, "debug test %s", "message")

	assertNoMessage(t, "DebugCtxf")
}

// -------------------
//
// Information Package Block
//
// -------------------

func TestLoggerViaPackageInformation(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	Information("info test message")

	assertMessageViaPackage(t, "Information", "4info test message")
}

func TestLoggerViaPackageInactiveInformation(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	Information("info test message")

	assertNoMessage(t, "Information")
}

func TestLoggerViaPackageInformationWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	InformationWithCorrelation("1234", "info test message")

	assertMessageViaPackage(t, "InformationWithCorrelation", "41234info test message")
}

func TestLoggerViaPackageInactiveInformationWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	InformationWithCorrelation("1234", "info test message")

	assertNoMessage(t, "InformationWithCorrelation")
}

func TestLoggerViaPackageInformationCustom(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	InformationCustom(map[string]any{"test": 123}, "info test message")

	assertMessageViaPackage(t, "InformationCustom", "4 map[test:123]info test message")
}

func TestLoggerViaPackageInactiveInformationCustom(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	InformationCustom(map[string]any{"test": 123}, "info test message")

	assertNoMessage(t, "InformationCustom")
}

func TestLoggerViaPackageInformationCtx(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	InformationCtx(testDummyContext, "info test message")

	assertMessageViaPackage(t, "InformationCtx", "41234info test message")
}

func TestLoggerViaPackageInactiveInformationCtx(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	InformationCtx(testDummyContext, "info test message")

	assertNoMessage(t, "InformationCtx")
}

func TestLoggerViaPackageInformationf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	Informationf("info test %s", "message")

	assertMessageViaPackage(t, "Informationf", "4info test message")
}

func TestLoggerViaPackageInactiveInformationf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	Informationf("info test %s", "message")

	assertNoMessage(t, "Informationf")
}

func TestLoggerViaPackageInformationWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	InformationWithCorrelationf("1234", "info test %s", "message")

	assertMessageViaPackage(t, "InformationWithCorrelationf", "41234info test message")
}

func TestLoggerViaPackageInactiveInformationWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	InformationWithCorrelationf("1234", "info test %s", "message")

	assertNoMessage(t, "InformationWithCorrelationf")
}

func TestLoggerViaPackageInformationCustomf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertMessageViaPackage(t, "InformationCustomf", "4 map[test:123]info test message")
}

func TestLoggerViaPackageInactiveInformationCustomf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertNoMessage(t, "InformationCustomf")
}

func TestLoggerViaPackageInformationCtxf(t *testing.T) {
	initLoggerViaPackageTest("INFO", "INFO")

	InformationCtxf(testDummyContext, "info test %s", "message")

	assertMessageViaPackage(t, "InformationCtxf", "41234info test message")
}

func TestLoggerViaPackageInactiveInformationCtxf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	InformationCtxf(testDummyContext, "info test %s", "message")

	assertNoMessage(t, "InformationCtxf")
}

// -------------------
//
// Information Common Block
//
// -------------------

func TestLoggerViaCommonInformation(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	Information("info test message")

	assertMessageViaCommon(t, "Information", "4info test message")
}

func TestLoggerViaCommonInactiveInformation(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	Information("info test message")

	assertNoMessage(t, "Information")
}

func TestLoggerViaCommonInformationWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	InformationWithCorrelation("1234", "info test message")

	assertMessageViaCommon(t, "InformationWithCorrelation", "41234info test message")
}

func TestLoggerViaCommonInactiveInformationWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	InformationWithCorrelation("1234", "info test message")

	assertNoMessage(t, "InformationWithCorrelation")
}

func TestLoggerViaCommonInformationCustom(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	InformationCustom(map[string]any{"test": 123}, "info test message")

	assertMessageViaCommon(t, "InformationCustom", "4 map[test:123]info test message")
}

func TestLoggerViaCommonInactiveInformationCustom(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	InformationCustom(map[string]any{"test": 123}, "info test message")

	assertNoMessage(t, "InformationCustom")
}

func TestLoggerViaCommonInformationCtx(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	InformationCtx(testDummyContext, "info test message")

	assertMessageViaCommon(t, "InformationCtx", "41234info test message")
}

func TestLoggerViaCommonInactiveInformationCtx(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	InformationCtx(testDummyContext, "info test message")

	assertNoMessage(t, "InformationCtx")
}

func TestLoggerViaCommonInformationf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	Informationf("info test %s", "message")

	assertMessageViaCommon(t, "Informationf", "4info test message")
}

func TestLoggerViaCommonInactiveInformationf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	Informationf("info test %s", "message")

	assertNoMessage(t, "Informationf")
}

func TestLoggerViaCommonInformationWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	InformationWithCorrelationf("1234", "info test %s", "message")

	assertMessageViaCommon(t, "InformationWithCorrelationf", "41234info test message")
}

func TestLoggerViaCommonInactiveInformationWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	InformationWithCorrelationf("1234", "info test %s", "message")

	assertNoMessage(t, "InformationWithCorrelationf")
}

func TestLoggerViaCommonInformationCustomf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCustomf", "4 map[test:123]info test message")
}

func TestLoggerViaCommonInactiveInformationCustomf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertNoMessage(t, "InformationCustomf")
}

func TestLoggerViaCommonInformationCtxf(t *testing.T) {
	initLoggerViaCommonTest("INFO", "INFO")

	InformationCtxf(testDummyContext, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCtxf", "41234info test message")
}

func TestLoggerViaCommonInactiveInformationCtxf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	InformationCtxf(testDummyContext, "info test %s", "message")

	assertNoMessage(t, "InformationCtxf")
}

// -------------------
//
// Information Only Common Block
//
// -------------------

func TestLoggerOnlyCommonInformation(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	Information("info test message")

	assertMessageViaCommon(t, "Information", "4info test message")
}

func TestLoggerOnlyCommonInactiveInformation(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	Information("info test message")

	assertNoMessage(t, "Information")
}

func TestLoggerOnlyCommonInformationWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	InformationWithCorrelation("1234", "info test message")

	assertMessageViaCommon(t, "InformationWithCorrelation", "41234info test message")
}

func TestLoggerOnlyCommonInactiveInformationWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	InformationWithCorrelation("1234", "info test message")

	assertNoMessage(t, "InformationWithCorrelation")
}

func TestLoggerOnlyCommonInformationCustom(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	InformationCustom(map[string]any{"test": 123}, "info test message")

	assertMessageViaCommon(t, "InformationCustom", "4 map[test:123]info test message")
}

func TestLoggerOnlyCommonInactiveInformationCustom(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	InformationCustom(map[string]any{"test": 123}, "info test message")

	assertNoMessage(t, "InformationCustom")
}

func TestLoggerOnlyCommonInformationCtx(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	InformationCtx(testDummyContext, "info test message")

	assertMessageViaCommon(t, "InformationCtx", "41234info test message")
}

func TestLoggerOnlyCommonInactiveInformationCtx(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	InformationCtx(testDummyContext, "info test message")

	assertNoMessage(t, "InformationCtx")
}

func TestLoggerOnlyCommonInformationf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	Informationf("info test %s", "message")

	assertMessageViaCommon(t, "Informationf", "4info test message")
}

func TestLoggerOnlyCommonInactiveInformationf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	Informationf("info test %s", "message")

	assertNoMessage(t, "Informationf")
}

func TestLoggerOnlyCommonInformationWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	InformationWithCorrelationf("1234", "info test %s", "message")

	assertMessageViaCommon(t, "InformationWithCorrelationf", "41234info test message")
}

func TestLoggerOnlyCommonInactiveInformationWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	InformationWithCorrelationf("1234", "info test %s", "message")

	assertNoMessage(t, "InformationWithCorrelationf")
}

func TestLoggerOnlyCommonInformationCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCustomf", "4 map[test:123]info test message")
}

func TestLoggerOnlyCommonInactiveInformationCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	InformationCustomf(map[string]any{"test": 123}, "info test %s", "message")

	assertNoMessage(t, "InformationCustomf")
}

func TestLoggerOnlyCommonInformationCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("INFO")

	InformationCtxf(testDummyContext, "info test %s", "message")

	assertMessageViaCommon(t, "InformationCtxf", "41234info test message")
}

func TestLoggerOnlyCommonInactiveInformationCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	InformationCtxf(testDummyContext, "info test %s", "message")

	assertNoMessage(t, "InformationCtxf")
}

// -------------------
//
// Warning Package Block
//
// -------------------

func TestLoggerViaPackageWarning(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	Warning("warn test message")

	assertMessageViaPackage(t, "Warning", "3warn test message")
}

func TestLoggerViaPackageInactiveWarning(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	Warning("warn test message")

	assertNoMessage(t, "Warning")
}

func TestLoggerViaPackageWarningWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	WarningWithCorrelation("1234", "warn test message")

	assertMessageViaPackage(t, "WarningWithCorrelation", "31234warn test message")
}

func TestLoggerViaPackageInactiveWarningWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	WarningWithCorrelation("1234", "warn test message")

	assertNoMessage(t, "WarningWithCorrelation")
}

func TestLoggerViaPackageWarningCustom(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertMessageViaPackage(t, "WarningCustom", "3 map[test:123]warn test message")
}

func TestLoggerViaPackageInactiveWarningCustom(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertNoMessage(t, "WarningCustom")
}

func TestLoggerViaPackageWarningCtx(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	WarningCtx(testDummyContext, "warn test message")

	assertMessageViaPackage(t, "WarningCtx", "31234warn test message")
}

func TestLoggerViaPackageInactivWarningCtx(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	WarningCtx(testDummyContext, "warn test message")

	assertNoMessage(t, "WarningCtx")
}

func TestLoggerViaPackageWarningf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	Warningf("warn test %s", "message")

	assertMessageViaPackage(t, "Warningf", "3warn test message")
}

func TestLoggerViaPackageInactiveWarningf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	Warningf("warn test %s", "message")

	assertNoMessage(t, "Warningf")
}

func TestLoggerViaPackageWarningWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	WarningWithCorrelationf("1234", "warn test %s", "message")

	assertMessageViaPackage(t, "WarningWithCorrelationf", "31234warn test message")
}

func TestLoggerViaPackageInactiveWarningWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	WarningWithCorrelationf("1234", "warn test %s", "message")

	assertNoMessage(t, "WarningWithCorrelationf")
}

func TestLoggerViaPackageWarningCustomf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaPackage(t, "WarningCustomf", "3 map[test:123]warn test message")
}

func TestLoggerViaPackageInactiveWarningCustomf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessage(t, "WarningCustomf")
}

func TestLoggerViaPackageWarningWarningCtxf(t *testing.T) {
	initLoggerViaPackageTest("WARN", "WARN")

	WarningCtxf(testDummyContext, "warn test %s", "message")

	assertMessageViaPackage(t, "WarningCtxf", "31234warn test message")
}

func TestLoggerViaPackageInactiveWarningCtxf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	WarningCtxf(testDummyContext, "warn test %s", "message")

	assertNoMessage(t, "WarningCtxf")
}

// -------------------
//
// Warning Common Block
//
// -------------------

func TestLoggerViaCommonWarning(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	Warning("warn test message")

	assertMessageViaCommon(t, "Warning", "3warn test message")
}

func TestLoggerViaCommonInactiveWarning(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	Warning("warn test message")

	assertNoMessage(t, "Warning")
}

func TestLoggerViaCommonWarningWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningWithCorrelation("1234", "warn test message")

	assertMessageViaCommon(t, "WarningWithCorrelation", "31234warn test message")
}

func TestLoggerViaCommonInactiveWarningWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningWithCorrelation("1234", "warn test message")

	assertNoMessage(t, "WarningWithCorrelation")
}

func TestLoggerViaCommonWarningCustom(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertMessageViaCommon(t, "WarningCustom", "3 map[test:123]warn test message")
}

func TestLoggerViaCommonInactiveWarningCustom(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertNoMessage(t, "WarningCustom")
}

func TestLoggerViaCommoWarningCtx(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCtx(testDummyContext, "warn test message")

	assertMessageViaCommon(t, "WarningCtx", "31234warn test message")
}

func TestLoggerViaCommonInactivWarningCtx(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCtx(testDummyContext, "warn test message")

	assertNoMessage(t, "WarningCtx")
}

func TestLoggerViaCommonWarningf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	Warningf("warn test %s", "message")

	assertMessageViaCommon(t, "Warningf", "3warn test message")
}

func TestLoggerViaCommonInactiveWarningf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	Warningf("warn test %s", "message")

	assertNoMessage(t, "Warningf")
}

func TestLoggerViaCommonWarningWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningWithCorrelationf("1234", "warn test %s", "message")

	assertMessageViaCommon(t, "WarningWithCorrelationf", "31234warn test message")
}

func TestLoggerViaCommonInactiveWarningWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningWithCorrelationf("1234", "warn test %s", "message")

	assertNoMessage(t, "WarningWithCorrelationf")
}

func TestLoggerViaCommonWarningCustomf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCustomf", "3 map[test:123]warn test message")
}

func TestLoggerViaCommonInactiveWarningCustomf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessage(t, "WarningCustomf")
}

func TestLoggerViaCommonWarningCtxf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCtxf(testDummyContext, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCtxf", "31234warn test message")
}

func TestLoggerViaCommonInactiveWarningCtxf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCtxf(testDummyContext, "warn test %s", "message")

	assertNoMessage(t, "WarningCtxf")
}

// -------------------
//
// Warning Only Common Block
//
// -------------------

func TestLoggerOnlyCommonWarning(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	Warning("warn test message")

	assertMessageViaCommon(t, "Warning", "3warn test message")
}

func TestLoggerOnlyCommonInactiveWarning(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	Warning("warn test message")

	assertNoMessage(t, "Warning")
}

func TestLoggerOnlyCommonWarningWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	WarningWithCorrelation("1234", "warn test message")

	assertMessageViaCommon(t, "WarningWithCorrelation", "31234warn test message")
}

func TestLoggerOnlyCommonInactiveWarningWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	WarningWithCorrelation("1234", "warn test message")

	assertNoMessage(t, "WarningWithCorrelation")
}

func TestLoggerOnlyCommonWarningCustom(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertMessageViaCommon(t, "WarningCustom", "3 map[test:123]warn test message")
}

func TestLoggerOnlyCommonInactiveWarningCustom(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	WarningCustom(map[string]any{"test": 123}, "warn test message")

	assertNoMessage(t, "WarningCustom")
}

func TestLoggerOnlyCommonWarningCtx(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	WarningCtx(testDummyContext, "warn test message")

	assertMessageViaCommon(t, "WarningCtx", "31234warn test message")
}

func TestLoggerOnlyCommonInactivWarningCtx(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	WarningCtx(testDummyContext, "warn test message")

	assertNoMessage(t, "WarningCtx")
}

func TestLoggerOnlyCommonWarningf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	Warningf("warn test %s", "message")

	assertMessageViaCommon(t, "Warningf", "3warn test message")
}

func TestLoggerOnlyCommonInactiveWarningf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	Warningf("warn test %s", "message")

	assertNoMessage(t, "Warningf")
}

func TestLoggerOnlyCommonWarningWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	WarningWithCorrelationf("1234", "warn test %s", "message")

	assertMessageViaCommon(t, "WarningWithCorrelationf", "31234warn test message")
}

func TestLoggerOnlyCommonInactiveWarningWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	WarningWithCorrelationf("1234", "warn test %s", "message")

	assertNoMessage(t, "WarningWithCorrelationf")
}

func TestLoggerOnlyCommonWarningCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	WarningCtxf(testDummyContext, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCtxf", "31234warn test message")
}

func TestLoggerOnlyCommonInactiveWarningCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	WarningCtxf(testDummyContext, "warn test %s", "message")

	assertNoMessage(t, "WarningCtxf")
}

// -------------------
//
// Warning With Panic Block
//
// -------------------

func TestLoggerWarningWithPanic(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningWithPanic("warn test message")

	assertMessageWithPanic(t, "WarningWithPanic", "3warn test message")
}

func TestLoggerInactiveWarningWithPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningWithPanic("warn test message")

	assertNoMessageWithPanic(t, "WarningWithPanic")
}

func TestLoggerWarningWithCorrelationAndPanic(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningWithCorrelationAndPanic("1234", "warn test message")

	assertMessageWithPanic(t, "WarningWithCorrelationAndPanic", "31234warn test message")
}

func TestLoggerInactiveWarningWithCorrelationAndPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningWithCorrelationAndPanic("1234", "warn test message")

	assertNoMessageWithPanic(t, "WarningWithCorrelationAndPanic")
}

func TestLoggerWarningCustomWithPanic(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCustomWithPanic(map[string]any{"test": 123}, "warn test message")

	assertMessageWithPanic(t, "WarningCustomWithPanic", "3 map[test:123]warn test message")
}

func TestLoggerInactiveWarningCustomWithPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCustomWithPanic(map[string]any{"test": 123}, "warn test message")

	assertNoMessageWithPanic(t, "WarningCustomWithPanic")
}

func TestLoggerWarningCtxWithPanic(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCtxWithPanic(testDummyContext, "warn test message")

	assertMessageWithPanic(t, "WarningCtxWithPanic", "31234warn test message")
}

func TestLoggerInactiveWarningCtxWithPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCtxWithPanic(testDummyContext, "warn test message")

	assertNoMessageWithPanic(t, "WarningCtxWithPanic")
}

func TestLoggerWarningWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningWithPanicf("warn test %s", "message")

	assertMessageWithPanic(t, "WarningWithPanicf", "3warn test message")
}

func TestLoggerInactiveWarningWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningWithPanicf("warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningWithPanicf")
}

func TestLoggerWarningWithCorrelationAndPanicf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningWithCorrelationAndPanicf("1234", "warn test %s", "message")

	assertMessageWithPanic(t, "WarningWithCorrelationAndPanicf", "31234warn test message")
}

func TestLoggerInactiveWarningWithCorrelationAndPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningWithCorrelationAndPanicf("1234", "warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningWithCorrelationAndPanicf")
}

func TestLoggerWarningCustomWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCustomWithPanicf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageWithPanic(t, "WarningCustomWithPanicf", "3 map[test:123]warn test message")
}

func TestLoggerInactiveWarningCustomWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCustomWithPanicf(map[string]any{"test": 123}, "warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningCustomWithPanicf")
}

func TestLoggerWarningCtxWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("WARN", "WARN")

	WarningCtxWithPanicf(testDummyContext, "warn test %s", "message")

	assertMessageWithPanic(t, "WarningCtxWithPanicf", "31234warn test message")
}

func TestLoggerInactiveWarningCtxWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	WarningCtxWithPanicf(testDummyContext, "warn test %s", "message")

	assertNoMessageWithPanic(t, "WarningCtxWithPanicf")
}

// -------------------
//
// Error Package Block
//
// -------------------

func TestLoggerViaPackageError(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	Error("error test message")

	assertMessageViaPackage(t, "Error", "2error test message")
}

func TestLoggerViaPackageInactiveError(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	Error("error test message")

	assertNoMessage(t, "Error")
}

func TestLoggerViaPackageErrorWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	ErrorWithCorrelation("1234", "error test message")

	assertMessageViaPackage(t, "ErrorWithCorrelation", "21234error test message")
}

func TestLoggerViaPackageInactiveErrorWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	ErrorWithCorrelation("1234", "error test message")

	assertNoMessage(t, "ErrorWithCorrelation")
}

func TestLoggerViaPackageErrorCustom(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertMessageViaPackage(t, "ErrorCustom", "2 map[test:123]error test message")
}

func TestLoggerViaPackageInactiveErrorCustom(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertNoMessage(t, "ErrorCustom")
}

func TestLoggerViaPackageErrorCtx(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	ErrorCtx(testDummyContext, "error test message")

	assertMessageViaPackage(t, "ErrorCtx", "21234error test message")
}

func TestLoggerViaPackageInactivErrorCtx(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	ErrorCtx(testDummyContext, "error test message")

	assertNoMessage(t, "ErrorCtx")
}

func TestLoggerViaPackageErrorf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	Errorf("error test %s", "message")

	assertMessageViaPackage(t, "Errorf", "2error test message")
}

func TestLoggerViaPackageInactiveErrorf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	Errorf("error test %s", "message")

	assertNoMessage(t, "Errorf")
}

func TestLoggerViaPackageErrorWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	ErrorWithCorrelationf("1234", "error test %s", "message")

	assertMessageViaPackage(t, "ErrorWithCorrelationf", "21234error test message")
}

func TestLoggerViaPackageInactiveErrorWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	ErrorWithCorrelationf("1234", "error test %s", "message")

	assertNoMessage(t, "ErrorWithCorrelationf")
}

func TestLoggerViaPackageErrorCustomf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageViaPackage(t, "ErrorCustomf", "2 map[test:123]error test message")
}

func TestLoggerViaPackageInactiveErrorCustomf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessage(t, "ErrorCustomf")
}

func TestLoggerViaPackageErrorCtxf(t *testing.T) {
	initLoggerViaPackageTest("ERROR", "ERROR")

	ErrorCtxf(testDummyContext, "error test %s", "message")

	assertMessageViaPackage(t, "ErrorCtxf", "21234error test message")
}

func TestLoggerViaPackageInactiveErrorCtxf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	ErrorCtxf(testDummyContext, "error test %s", "message")

	assertNoMessage(t, "ErrorCtxf")
}

// -------------------
//
// Error Common Block
//
// -------------------

func TestLoggerViaCommonError(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	Error("error test message")

	assertMessageViaCommon(t, "Error", "2error test message")
}

func TestLoggerViaCommonInactiveError(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	Error("error test message")

	assertNoMessage(t, "Error")
}

func TestLoggerViaCommonErrorWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorWithCorrelation("1234", "error test message")

	assertMessageViaCommon(t, "ErrorWithCorrelation", "21234error test message")
}

func TestLoggerViaCommonInactiveErrorWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorWithCorrelation("1234", "error test message")

	assertNoMessage(t, "ErrorWithCorrelation")
}

func TestLoggerViaCommonErrorCustom(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertMessageViaCommon(t, "ErrorCustom", "2 map[test:123]error test message")
}

func TestLoggerViaCommonInactiveErrorCustom(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertNoMessage(t, "ErrorCustom")
}

func TestLoggerViaCommonErrorCtx(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCtx(testDummyContext, "error test message")

	assertMessageViaCommon(t, "ErrorCtx", "21234error test message")
}

func TestLoggerViaCommonInactiveErrorCtx(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCtx(testDummyContext, "error test message")

	assertNoMessage(t, "ErrorCtx")
}

func TestLoggerViaCommonErrorf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	Errorf("error test %s", "message")

	assertMessageViaCommon(t, "Errorf", "2error test message")
}

func TestLoggerViaCommonInactiveErrorf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	Errorf("error test %s", "message")

	assertNoMessage(t, "Errorf")
}

func TestLoggerViaCommonErrorWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorWithCorrelationf("1234", "error test %s", "message")

	assertMessageViaCommon(t, "ErrorWithCorrelationf", "21234error test message")
}

func TestLoggerViaCommonInactiveErrorWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorWithCorrelationf("1234", "error test %s", "message")

	assertNoMessage(t, "ErrorWithCorrelationf")
}

func TestLoggerViaCommonErrorCustomf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCustomf", "2 map[test:123]error test message")
}

func TestLoggerViaCommonInactiveErrorCustomf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessage(t, "ErrorCustomf")
}

func TestLoggerViaCommonErrorCtxf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCtxf(testDummyContext, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCtxf", "21234error test message")
}

func TestLoggerViaCommonInactiveErrorCtxf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCtxf(testDummyContext, "error test %s", "message")

	assertNoMessage(t, "ErrorCtxf")
}

// -------------------
//
// Error Only Common Block
//
// -------------------

func TestLoggerOnlyCommonError(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	Error("error test message")

	assertMessageViaCommon(t, "Error", "2error test message")
}

func TestLoggerOnlyCommonInactiveError(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	Error("error test message")

	assertNoMessage(t, "Error")
}

func TestLoggerOnlyCommonErrorWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	ErrorWithCorrelation("1234", "error test message")

	assertMessageViaCommon(t, "ErrorWithCorrelation", "21234error test message")
}

func TestLoggerOnlyCommonInactiveErrorWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	ErrorWithCorrelation("1234", "error test message")

	assertNoMessage(t, "ErrorWithCorrelation")
}

func TestLoggerOnlyCommonErrorCustom(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertMessageViaCommon(t, "ErrorCustom", "2 map[test:123]error test message")
}

func TestLoggerOnlyCommonInactiveErrorCustom(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	ErrorCustom(map[string]any{"test": 123}, "error test message")

	assertNoMessage(t, "ErrorCustom")
}

func TestLoggerOnlyCommonErrorCtx(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	ErrorCtx(testDummyContext, "error test message")

	assertMessageViaCommon(t, "ErrorCtx", "21234error test message")
}

func TestLoggerOnlyCommonInactiveErrorCtx(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	ErrorCtx(testDummyContext, "error test message")

	assertNoMessage(t, "ErrorCtx")
}
func TestLoggerOnlyCommonErrorf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	Errorf("error test %s", "message")

	assertMessageViaCommon(t, "Errorf", "2error test message")
}

func TestLoggerOnlyCommonInactiveErrorf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	Errorf("error test %s", "message")

	assertNoMessage(t, "Errorf")
}

func TestLoggerOnlyCommonErrorWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	ErrorWithCorrelationf("1234", "error test %s", "message")

	assertMessageViaCommon(t, "ErrorWithCorrelationf", "21234error test message")
}

func TestLoggerOnlyCommonInactiveErrorWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	ErrorWithCorrelationf("1234", "error test %s", "message")

	assertNoMessage(t, "ErrorWithCorrelationf")
}

func TestLoggerOnlyCommonErrorCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCustomf", "2 map[test:123]error test message")
}

func TestLoggerOnlyCommonInactiveErrorCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	ErrorCustomf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessage(t, "ErrorCustomf")
}

func TestLoggerOnlyCommonErrorCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("ERROR")

	ErrorCtxf(testDummyContext, "error test %s", "message")

	assertMessageViaCommon(t, "ErrorCtxf", "21234error test message")
}

func TestLoggerOnlyCommonInactiveErrorCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	ErrorCtxf(testDummyContext, "error test %s", "message")

	assertNoMessage(t, "ErrorCtxf")
}

// -------------------
//
// Error With Panic Block
//
// -------------------

func TestLoggerErrorWithPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorWithPanic("error test message")

	assertMessageWithPanic(t, "ErrorWithPanic", "2error test message")
}

func TestLoggerInactiveErrorWithPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorWithPanic("error test message")

	assertNoMessageWithPanic(t, "ErrorWithPanic")
}

func TestLoggerErrorWithCorrelationAndPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorWithCorrelationAndPanic("1234", "error test message")

	assertMessageWithPanic(t, "ErrorWithCorrelationAndPanic", "21234error test message")
}

func TestLoggerInactiveErrorWithCorrelationAndPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorWithCorrelationAndPanic("1234", "error test message")

	assertNoMessageWithPanic(t, "ErrorWithCorrelationAndPanic")
}

func TestLoggerErrorCustomWithPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCustomWithPanic(map[string]any{"test": 123}, "error test message")

	assertMessageWithPanic(t, "ErrorCustomWithPanic", "2 map[test:123]error test message")
}

func TestLoggerInactiveErrorCustomWithPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCustomWithPanic(map[string]any{"test": 123}, "error test message")

	assertNoMessageWithPanic(t, "ErrorCustomWithPanic")
}

func TestLoggerErrorCtxWithPanic(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCtxWithPanic(testDummyContext, "error test message")

	assertMessageWithPanic(t, "ErrorCtxWithPanic", "21234error test message")
}

func TestLoggerInactivErrorCtxWithPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCtxWithPanic(testDummyContext, "error test message")

	assertNoMessageWithPanic(t, "ErrorCtxWithPanic")
}

func TestLoggerErrorWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorWithPanicf("error test %s", "message")

	assertMessageWithPanic(t, "ErrorWithPanicf", "2error test message")
}

func TestLoggerInactiveErrorWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorWithPanicf("error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorWithPanicf")
}

func TestLoggerErrorWithCorrelationAndPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorWithCorrelationAndPanicf("1234", "error test %s", "message")

	assertMessageWithPanic(t, "ErrorWithCorrelationAndPanicf", "21234error test message")
}

func TestLoggerInactiveErrorWithCorrelationAndPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorWithCorrelationAndPanicf("1234", "error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorWithCorrelationAndPanicf")
}

func TestLoggerErrorCustomWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCustomWithPanicf(map[string]any{"test": 123}, "error test %s", "message")

	assertMessageWithPanic(t, "ErrorCustomWithPanicf", "2 map[test:123]error test message")
}

func TestLoggerInactiveErrorCustomWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCustomWithPanicf(map[string]any{"test": 123}, "error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorCustomWithPanicf")
}

func TestLoggerErrorCtxWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("ERROR", "ERROR")

	ErrorCtxWithPanicf(testDummyContext, "error test %s", "message")

	assertMessageWithPanic(t, "ErrorCtxWithPanicf", "21234error test message")
}

func TestLoggerInactiveErrorCtxWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	ErrorCtxWithPanicf(testDummyContext, "error test %s", "message")

	assertNoMessageWithPanic(t, "ErrorCtxWithPanicf")
}

// -------------------
//
// Fatal Package Block
//
// -------------------

func TestLoggerViaPackageFatal(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	Fatal("fatal test message")

	assertMessageViaPackage(t, "Fatal", "1fatal test message")
}

func TestLoggerViaPackageInactiveFatal(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	Fatal("fatal test message")

	assertNoMessage(t, "Fatal")
}

func TestLoggerViaPackageFatalWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalWithCorrelation("1234", "fatal test message")

	assertMessageViaPackage(t, "FatalWithCorrelation", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalWithCorrelation(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalWithCorrelation("1234", "fatal test message")

	assertNoMessage(t, "FatalWithCorrelation")
}

func TestLoggerViaPackageFatalCustom(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaPackage(t, "FatalCustom", "1 map[test:123]fatal test message")
}

func TestLoggerViaPackageInactiveFatalCustom(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertNoMessage(t, "FatalCustom")
}

func TestLoggerViaPackageFatalCtx(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCtx(testDummyContext, "fatal test message")

	assertMessageViaPackage(t, "FatalCtx", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalCtx(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCtx(testDummyContext, "fatal test message")

	assertNoMessage(t, "FatalCtx")
}

func TestLoggerViaPackageFatalf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	Fatalf("fatal test %s", "message")

	assertMessageViaPackage(t, "Fatalf", "1fatal test message")
}

func TestLoggerViaPackageInactiveFatalf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	Fatalf("fatal test %s", "message")

	assertNoMessage(t, "Fatalf")
}

func TestLoggerViaPackageFatalWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertMessageViaPackage(t, "FatalWithCorrelationf", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalWithCorrelationf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertNoMessage(t, "FatalWithCorrelationf")
}

func TestLoggerViaPackageFatalCustomf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaPackage(t, "FatalCustomf", "1 map[test:123]fatal test message")
}

func TestLoggerViaPackageInactiveFatalCustomf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessage(t, "FatalCustomf")
}

func TestLoggerViaPackageFatalCtxf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertMessageViaPackage(t, "FatalCtxf", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalCtxf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertNoMessage(t, "FatalCtxf")
}

// -------------------
//
// Fatal Common Block
//
// -------------------

func TestLoggerViaCommonFatal(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	Fatal("fatal test message")

	assertMessageViaCommon(t, "Fatal", "1fatal test message")
}

func TestLoggerViaCommonInactiveFatal(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	Fatal("fatal test message")

	assertNoMessage(t, "Fatal")
}

func TestLoggerViaCommonFatalWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithCorrelation("1234", "fatal test message")

	assertMessageViaCommon(t, "FatalWithCorrelation", "11234fatal test message")
}

func TestLoggerViaCommonInactiveFatalWithCorrelation(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithCorrelation("1234", "fatal test message")

	assertNoMessage(t, "FatalWithCorrelation")
}

func TestLoggerViaCommonFatalCustom(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaCommon(t, "FatalCustom", "1 map[test:123]fatal test message")
}

func TestLoggerViaCommonInactiveFatalCustom(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertNoMessage(t, "FatalCustom")
}

func TestLoggerViaCommonFatalCtx(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCtx(testDummyContext, "fatal test message")

	assertMessageViaCommon(t, "FatalCtx", "11234fatal test message")
}

func TestLoggerViaCommonInactiveFatalCtx(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCtx(testDummyContext, "fatal test message")

	assertNoMessage(t, "FatalCtx")
}

func TestLoggerViaCommonFatalf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	Fatalf("fatal test %s", "message")

	assertMessageViaCommon(t, "Fatalf", "1fatal test message")
}

func TestLoggerViaCommonInactiveFatalf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	Fatalf("fatal test %s", "message")

	assertNoMessage(t, "Fatalf")
}

func TestLoggerViaCommonFatalWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalWithCorrelationf", "11234fatal test message")
}

func TestLoggerViaCommonInactiveFatalWithCorrelationf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertNoMessage(t, "FatalWithCorrelationf")
}

func TestLoggerViaCommonFatalCustomf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCustomf", "1 map[test:123]fatal test message")
}

func TestLoggerViaCommonInactiveFatalCustomf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessage(t, "FatalCustomf")
}

func TestLoggerViaCommonFatalCtxf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCtxf", "11234fatal test message")
}

func TestLoggerViaCommonInactiveFatalCtxf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertNoMessage(t, "FatalCtxf")
}

// -------------------
//
// Fatal Only Common Block
//
// -------------------

func TestLoggerOnlyCommonFatal(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	Fatal("fatal test message")

	assertMessageViaCommon(t, "Fatal", "1fatal test message")
}

func TestLoggerOnlyCommonInactiveFatal(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	Fatal("fatal test message")

	assertNoMessage(t, "Fatal")
}

func TestLoggerOnlyCommonFatalWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalWithCorrelation("1234", "fatal test message")

	assertMessageViaCommon(t, "FatalWithCorrelation", "11234fatal test message")
}

func TestLoggerOnlyCommonInactiveFatalWithCorrelation(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	FatalWithCorrelation("1234", "fatal test message")

	assertNoMessage(t, "FatalWithCorrelation")
}

func TestLoggerOnlyCommonFatalCustom(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaCommon(t, "FatalCustom", "1 map[test:123]fatal test message")
}

func TestLoggerOnlyCommonInactiveFatalCustom(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	FatalCustom(map[string]any{"test": 123}, "fatal test message")

	assertNoMessage(t, "FatalCustom")
}

func TestLoggerOnlyCommonFatalCtx(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalCtx(testDummyContext, "fatal test message")

	assertMessageViaCommon(t, "FatalCtx", "11234fatal test message")
}

func TestLoggerOnlyCommonInactiveFatalCtx(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	FatalCtx(testDummyContext, "fatal test message")

	assertNoMessage(t, "FatalCtx")
}

func TestLoggerOnlyCommonFatalf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	Fatalf("fatal test %s", "message")

	assertMessageViaCommon(t, "Fatalf", "1fatal test message")
}

func TestLoggerOnlyCommonInactiveFatalf(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	Fatalf("fatal test %s", "message")

	assertNoMessage(t, "Fatalf")
}

func TestLoggerOnlyCommonFatalWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalWithCorrelationf", "11234fatal test message")
}

func TestLoggerOnlyCommonInactiveFatalWithCorrelationf(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	FatalWithCorrelationf("1234", "fatal test %s", "message")

	assertNoMessage(t, "FatalWithCorrelationf")
}

func TestLoggerOnlyCommonFatalCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalCustomf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCustomf", "1 map[test:123]fatal test message")
}

func TestLoggerOnlyCommonWarningCustomf(t *testing.T) {
	initLoggerOnlyCommonTest("WARN")

	WarningCustomf(map[string]any{"test": 123}, "warn test %s", "message")

	assertMessageViaCommon(t, "WarningCustomf", "3 map[test:123]warn test message")
}

func TestLoggerOnlyCommonFatalCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertMessageViaCommon(t, "FatalCtxf", "11234fatal test message")
}

func TestLoggerOnlyCommonInactiveFatalCtxf(t *testing.T) {
	initLoggerOnlyCommonTest("OFF")

	FatalCtxf(testDummyContext, "fatal test %s", "message")

	assertNoMessage(t, "FatalCtxf")
}

// -------------------
//
// Fatal With Panic Block
//
// -------------------

func TestLoggerFatalWithPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithPanic("fatal test message")

	assertMessageWithPanic(t, "FatalWithPanic", "1fatal test message")
}

func TestLoggerInactiveFatalWithPanic(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithPanic("fatal test message")

	assertNoMessageWithPanic(t, "FatalWithPanic")
}

func TestLoggerFatalWithCorrelationAndPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithCorrelationAndPanic("1234", "fatal test message")

	assertMessageWithPanic(t, "FatalWithCorrelationAndPanic", "11234fatal test message")
}

func TestLoggerInactiveFatalWithCorrelationAndPanic(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithCorrelationAndPanic("1234", "fatal test message")

	assertNoMessageWithPanic(t, "FatalWithCorrelationAndPanic")
}

func TestLoggerFatalCustomWithPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCustomWithPanic(map[string]any{"test": 123}, "fatal test message")

	assertMessageWithPanic(t, "FatalCustomWithPanic", "1 map[test:123]fatal test message")
}

func TestLoggerInactiveFatalCustomWithPanic(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCustomWithPanic(map[string]any{"test": 123}, "fatal test message")

	assertNoMessageWithPanic(t, "FatalCustomWithPanic")
}

func TestLoggerFatalCtxWithPanic(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCtxWithPanic(testDummyContext, "fatal test message")

	assertMessageWithPanic(t, "FatalCtxWithPanic", "11234fatal test message")
}

func TestLoggerInactiveFatalCtxWithPanic(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCtxWithPanic(testDummyContext, "fatal test message")

	assertNoMessageWithPanic(t, "FatalCtxWithPanic")
}

func TestLoggerFatalWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithPanicf("fatal test %s", "message")

	assertMessageWithPanic(t, "FatalWithPanicf", "1fatal test message")
}

func TestLoggerInactiveFatalWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithPanicf("fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalWithPanicf")
}

func TestLoggerFatalWithCorrelationAndPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithCorrelationAndPanicf("1234", "fatal test %s", "message")

	assertMessageWithPanic(t, "FatalWithCorrelationAndPanicf", "11234fatal test message")
}

func TestLoggerInactiveFatalWithCorrelationAndPanicf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithCorrelationAndPanicf("1234", "fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalWithCorrelationAndPanicf")
}

func TestLoggerFatalCustomWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCustomWithPanicf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageWithPanic(t, "FatalCustomWithPanicf", "1 map[test:123]fatal test message")
}

func TestLoggerInactiveFatalCustomWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCustomWithPanicf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalCustomWithPanicf")
}

func TestLoggerFatalCtxWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCtxWithPanicf(testDummyContext, "fatal test %s", "message")

	assertMessageWithPanic(t, "FatalCtxWithPanicf", "11234fatal test message")
}

func TestLoggerInactiveFatalCtxWithPanicf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCtxWithPanicf(testDummyContext, "fatal test %s", "message")

	assertNoMessageWithPanic(t, "FatalCtxWithPanicf")
}

// -------------------
//
// Fatal With Exit Block
//
// -------------------

func TestLoggerFatalWithExit(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithExit("fatal test message")

	assertMessageWithExit(t, "FatalWithExit", "1fatal test message")
}

func TestLoggerInactiveFatalWithExit(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithExit("fatal test message")

	assertNoMessageWithExit(t, "FatalWithExit")
}

func TestLoggerFatalWithCorrelationAndExit(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithCorrelationAndExit("1234", "fatal test message")

	assertMessageWithExit(t, "FatalWithCorrelationAndExit", "11234fatal test message")
}

func TestLoggerInactiveFatalWithCorrelationAndExit(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithCorrelationAndExit("1234", "fatal test message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExit")
}

func TestLoggerFatalCustomWithExit(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertMessageWithExit(t, "FatalCustomWithExit", "1 map[test:123]fatal test message")
}

func TestLoggerInactiveFatalCustomWithExit(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertNoMessageWithExit(t, "FatalCustomWithExit")
}

func TestLoggerFatalCtxWithExit(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCtxWithExit(testDummyContext, "fatal test message")

	assertMessageWithExit(t, "FatalCtxWithExit", "11234fatal test message")
}

func TestLoggerInactiveFatalCtxWithExit(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCtxWithExit(testDummyContext, "fatal test message")

	assertNoMessageWithExit(t, "FatalCtxWithExit")
}

func TestLoggerFatalWithExitf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithExitf("fatal test %s", "message")

	assertMessageWithExit(t, "FatalWithExitf", "1fatal test message")
}

func TestLoggerInactiveFatalWithExitf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithExitf("fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithExitf")
}

func TestLoggerFatalWithCorrelationAndExitf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertMessageWithExit(t, "FatalWithCorrelationAndExitf", "11234fatal test message")
}

func TestLoggerInactiveFatalWithCorrelationAndExitf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExitf")
}

func TestLoggerFatalCustomWithExitf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageWithExit(t, "FatalCustomWithExitf", "1 map[test:123]fatal test message")
}

func TestLoggerInactiveFatalCustomWithExitf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCustomWithExitf")
}

func TestLoggerFatalCtxWithExitf(t *testing.T) {
	initLoggerViaCommonTest("FATAL", "FATAL")

	FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertMessageWithExit(t, "FatalCtxWithExitf", "11234fatal test message")
}

func TestLoggerInactiveFatalCtxWithExitf(t *testing.T) {
	initLoggerViaCommonTest("OFF", "OFF")

	FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCtxWithExitf")
}

// -------------------
//
// Fatal Package With Exit Block
//
// -------------------

func TestLoggerViaPackageFatalWithExit(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalWithExit("fatal test message")

	assertMessageViaPackageWithExit(t, "FatalWithExit", "1fatal test message")
}

func TestLoggerViaPackageInactiveFatalWithExit(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalWithExit("fatal test message")

	assertNoMessageWithExit(t, "FatalWithExit")
}

func TestLoggerViaPackageFatalWithCorrelationAndExit(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalWithCorrelationAndExit("1234", "fatal test message")

	assertMessageViaPackageWithExit(t, "FatalWithCorrelationAndExit", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalWithCorrelationAndExit(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalWithCorrelationAndExit("1234", "fatal test message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExit")
}

func TestLoggerViaPackageFatalCustomWithExit(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertMessageViaPackageWithExit(t, "FatalCustomWithExit", "1 map[test:123]fatal test message")
}

func TestLoggerViaPackageInactiveFatalCustomWithExit(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCustomWithExit(map[string]any{"test": 123}, "fatal test message")

	assertNoMessageWithExit(t, "FatalCustomWithExit")
}

func TestLoggerViaPackageFatalCtxWithExit(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCtxWithExit(testDummyContext, "fatal test message")

	assertMessageViaPackageWithExit(t, "FatalCtxWithExit", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalCtxWithExit(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCtxWithExit(testDummyContext, "fatal test message")

	assertNoMessageWithExit(t, "FatalCtxWithExit")
}

func TestLoggerViaPackageFatalWithExitf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalWithExitf("fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalWithExitf", "1fatal test message")
}

func TestLoggerViaPackageInactiveFatalWithExitf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalWithExitf("fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithExitf")
}

func TestLoggerViaPackageFatalWithCorrelationAndExitf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalWithCorrelationAndExitf", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalWithCorrelationAndExitf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalWithCorrelationAndExitf("1234", "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalWithCorrelationAndExitf")
}

func TestLoggerViaPackageFatalCustomWithExitf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalCustomWithExitf", "1 map[test:123]fatal test message")
}

func TestLoggerViaPackageInactiveFatalCustomWithExitf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCustomWithExitf(map[string]any{"test": 123}, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCustomWithExitf")
}

func TestLoggerViaPackageFatalCtxWithExitf(t *testing.T) {
	initLoggerViaPackageTest("FATAL", "FATAL")

	FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertMessageViaPackageWithExit(t, "FatalCtxWithExitf", "11234fatal test message")
}

func TestLoggerViaPackageInactiveFatalCtxWithExitf(t *testing.T) {
	initLoggerViaPackageTest("OFF", "OFF")

	FatalCtxWithExitf(testDummyContext, "fatal test %s", "message")

	assertNoMessageWithExit(t, "FatalCtxWithExitf")
}

// -------------------
//
// Fatal Only Common Block
//
// -------------------

func TestLoggerOnlyCommonFatalWithExit(t *testing.T) {
	initLoggerOnlyCommonTest("FATAL")

	FatalWithExit("fatal test message")

	assertMessageWithExit(t, "Fatal", "1fatal test message")
}

// -------------------
//
// De-/Register appender and formatter
//
// -------------------

type DummyAppenderConfig struct {
	common *config.CommonAppenderConfig
}

func (c DummyAppenderConfig) Id() string {
	return c.common.Id
}
func (c DummyAppenderConfig) AppenderType() string {
	return c.common.AppenderType
}
func (c DummyAppenderConfig) IsDefault() bool {
	return c.common.IsDefault
}
func (c DummyAppenderConfig) PackageParameter() string {
	return c.common.PackageParameter
}
func (c DummyAppenderConfig) GetCommon() *config.CommonAppenderConfig {
	return c.common
}
func (c DummyAppenderConfig) Equals(other *config.AppenderConfig) bool {
	return c.common.Equals((*other).GetCommon())
}
func (c DummyAppenderConfig) CreateFullCopy() config.AppenderConfig {
	return c
}

func createDummyAppenderConfig(relevantKeyValues *map[string]string, commonConfig *config.CommonAppenderConfig) (*config.AppenderConfig, error) {
	var result config.AppenderConfig = DummyAppenderConfig{common: commonConfig}
	return &result, nil
}

type DummyAppender struct {
	formatter *format.Formatter
}

func (s DummyAppender) Write(logValues *common.LogValues) {}
func (s DummyAppender) Close()                            {}

func createDummyAppender(appenderConfig *config.AppenderConfig, formatter *format.Formatter) (*appender.Appender, error) {
	var result appender.Appender = DummyAppender{formatter: formatter}
	return &result, nil
}

func TestAppenderRegistration(t *testing.T) {
	customAppenderName := "CUSTOM_APPENDER"

	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, customAppenderName)

	Reset()

	// Check fallback to default one if CUSTOM_FORMATTER is not registered yet
	result := getLoggers()
	testutil.AssertNotNil(result, t, "before register - result")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "before register - generalLogger.appender.Name")

	// Register custom one
	err := RegisterAppenderWithConfig(customAppenderName, []string{}, createDummyAppender, createDummyAppenderConfig)
	testutil.AssertNil(err, t, "err of RegisterAppenderWithConfig")

	result = getLoggers()
	testutil.AssertNotNil(result, t, "after register - result")
	testutil.AssertEquals("DummyAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "after register - generalLogger.appender.Name")

	// Deregister custom one
	err = DeregisterAppenderTogetherWithConfig(customAppenderName)
	testutil.AssertNil(err, t, "err of DeregisterAppenderTogetherWithConfig")

	// Load config without registered appender: fallback to default one
	result = getLoggers()
	testutil.AssertNotNil(result, t, "deregister - result")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "deregister - generalLogger.appender.Name")
}

func TestAppenderRegistrationKnownAppender(t *testing.T) {
	Reset()
	err := RegisterAppenderWithConfig(config.APPENDER_STDOUT, []string{}, createDummyAppender, createDummyAppenderConfig)
	testutil.AssertNotNil(err, t, "registered known standard output appender")

	err = RegisterAppenderWithConfig(config.APPENDER_FILE, []string{}, createDummyAppender, createDummyAppenderConfig)
	testutil.AssertNotNil(err, t, "registered known file appender")
}

func TestAppenderRegistrationWithError(t *testing.T) {
	customAppenderName := "CUSTOM_APPENDER"

	Reset()
	// Appender already exist
	err := RegisterAppender(customAppenderName, createDummyAppender)
	testutil.AssertNil(err, t, "RegisterFormatter")
	err = RegisterAppenderWithConfig(customAppenderName, []string{}, createDummyAppender, createDummyAppenderConfig)
	testutil.AssertNotNil(err, t, "RegisterAppenderWithConfig")
}

func TestAppenderDeregistrationBuildInAppender(t *testing.T) {
	Reset()
	err := DeregisterAppenderTogetherWithConfig(config.APPENDER_STDOUT)
	testutil.AssertNotNil(err, t, "deregistered standard output appender")

	err = DeregisterAppenderTogetherWithConfig(config.APPENDER_FILE)
	testutil.AssertNotNil(err, t, "deregistered file appender")
}

type DummyFormatterConfig struct {
	common *config.CommonFormatterConfig
}

func (c DummyFormatterConfig) Id() string {
	return c.common.Id
}
func (c DummyFormatterConfig) FormatterType() string {
	return c.common.FormatterType
}
func (c DummyFormatterConfig) IsDefault() bool {
	return c.common.IsDefault
}
func (c DummyFormatterConfig) PackageParameter() string {
	return c.common.PackageParameter
}
func (c DummyFormatterConfig) TimeLayout() string {
	return c.common.TimeLayout
}
func (c DummyFormatterConfig) GetCommon() *config.CommonFormatterConfig {
	return c.common
}
func (c DummyFormatterConfig) Equals(other *config.FormatterConfig) bool {
	return c.common.Equals((*other).GetCommon())
}
func (c DummyFormatterConfig) CreateFullCopy() config.FormatterConfig {
	return c
}

func createDummyFormatterConfig(relevantKeyValues *map[string]string, commonConfig *config.CommonFormatterConfig) (*config.FormatterConfig, error) {
	var result config.FormatterConfig = DummyFormatterConfig{common: commonConfig}
	return &result, nil
}

type DummyFormatter struct {
}

func (f DummyFormatter) Format(logValues *common.LogValues) string {
	return logValues.Message
}

func createDummyFormatter(formatterConfig *config.FormatterConfig) (*format.Formatter, error) {
	var result format.Formatter = DummyFormatter{}
	return &result, nil
}

func TestFormatterRegistration(t *testing.T) {

	customAppenderName := "CUSTOM_APPENDER"
	customFormatterName := "CUSTOM_FORMATTER"

	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, customAppenderName)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, customFormatterName)

	Reset()

	err := RegisterAppenderWithConfig(customAppenderName, []string{}, createDummyAppender, createDummyAppenderConfig)
	testutil.AssertNil(err, t, "err of RegisterAppenderWithConfig")

	// Check fallback to default one if CUSTOM_FORMATTER is not registered yet
	result := getLoggers()
	testutil.AssertNotNil(result, t, "before register - result")
	testutil.AssertEquals("DelimiterFormatter", reflect.TypeOf(*(*result.generalLogger.appender).(DummyAppender).formatter).Name(), t, "before register - generalLogger.appender.formatter.Name")

	// Register custom one
	err = RegisterFormatterWithConfig(customFormatterName, []string{}, createDummyFormatter, createDummyFormatterConfig)
	testutil.AssertNil(err, t, "err of RegisterFormatterWithConfig")

	result = getLoggers()
	testutil.AssertNotNil(result, t, "after register - result")
	testutil.AssertEquals("DummyFormatter", reflect.TypeOf(*(*result.generalLogger.appender).(DummyAppender).formatter).Name(), t, "after register - generalLogger.appender.formatter.Name")

	// Deregister custom one
	err = DeregisterFormatterTogetherWithConfig(customFormatterName)
	testutil.AssertNil(err, t, "err of DeregisterFormatter")

	// Load config without registered formatter: fallback to default one
	result = getLoggers()
	testutil.AssertNotNil(result, t, "deregister - result")
	testutil.AssertEquals("DelimiterFormatter", reflect.TypeOf(*(*result.generalLogger.appender).(DummyAppender).formatter).Name(), t, "deregister - generalLogger.appender.formatter.Name")
}

func TestFormatterRegistrationKnownFormatter(t *testing.T) {
	Reset()
	err := RegisterFormatterWithConfig(config.FORMATTER_DELIMITER, []string{}, createDummyFormatter, createDummyFormatterConfig)
	testutil.AssertNotNil(err, t, "registered known delimiter formatter")

	err = RegisterFormatterWithConfig(config.FORMATTER_TEMPLATE, []string{}, createDummyFormatter, createDummyFormatterConfig)
	testutil.AssertNotNil(err, t, "registered known template formatter")

	err = RegisterFormatterWithConfig(config.FORMATTER_JSON, []string{}, createDummyFormatter, createDummyFormatterConfig)
	testutil.AssertNotNil(err, t, "registered known json formatter")
}

func TestFormatterRegistrationWithError(t *testing.T) {
	customFormatterName := "CUSTOM_FORMATTER"

	Reset()
	// Formatter already exist
	err := RegisterFormatter(customFormatterName, createDummyFormatter)
	testutil.AssertNil(err, t, "RegisterFormatter")
	err = RegisterFormatterWithConfig(customFormatterName, []string{}, createDummyFormatter, createDummyFormatterConfig)
	testutil.AssertNotNil(err, t, "RegisterFormatterWithConfig")
}

func TestFormatterDeregistrationBuildInFormatter(t *testing.T) {
	Reset()
	err := DeregisterFormatterTogetherWithConfig(config.FORMATTER_DELIMITER)
	testutil.AssertNotNil(err, t, "deregistered delimiter formatter")

	err = DeregisterFormatterTogetherWithConfig(config.FORMATTER_TEMPLATE)
	testutil.AssertNotNil(err, t, "deregistered template formatter")

	err = DeregisterFormatterTogetherWithConfig(config.FORMATTER_JSON)
	testutil.AssertNotNil(err, t, "deregistered json formatter")
}

func TestClose(t *testing.T) {
	initLoggerViaPackageTest("OFF", "DEBUG")

	Close()

	testutil.AssertEquals(2, testGeneralLoggerCounterAppenderClosed, t, "number of closed appenders")
	testutil.AssertFalse(config.IsConfigInitialized(), t, "config initialized")
}
