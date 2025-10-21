package config

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/common"
)

type initConfigTest func(t *testing.T) *os.File
type addValueConfigTest func(targetFile *os.File, key string, value string)
type postInitConfigTest func(*os.File)

var envInitConfigTest initConfigTest = func(t *testing.T) *os.File {
	os.Clearenv()
	return nil
}

var propertiesFileInitConfigTest initConfigTest = func(t *testing.T) *os.File {
	os.Clearenv()

	pathToFile := testutil.DetermineTestCaseFilePathAt("LogConfig", "properties", true, true, "genTestResources")

	propertiesFile, err := os.Create(pathToFile)
	if err != nil {
		t.Errorf("Fail to create properties file: %v", err)
	}

	os.Setenv(LOG_CONFIG_FILE_ENV_NAME, pathToFile)

	return propertiesFile
}

var envAddValueConfigTest addValueConfigTest = func(targetFile *os.File, key string, value string) {
	os.Setenv(key, value)
}

var propertiesFileAddValueConfigTest addValueConfigTest = func(targetFile *os.File, key string, value string) {
	fmt.Fprintln(targetFile, key, "=", value)
}

var envPostInitConfigTest postInitConfigTest = func(*os.File) {
	// Nothing to Do
}

var propertiesFilePostInitConfigTest postInitConfigTest = func(targetFile *os.File) {
	targetFile.Close()
}

var allInitConfigTest = []initConfigTest{envInitConfigTest, propertiesFileInitConfigTest}
var allAddValueConfigTest = []addValueConfigTest{envAddValueConfigTest, propertiesFileAddValueConfigTest}
var allPostInitConfigTest = []postInitConfigTest{envPostInitConfigTest, propertiesFilePostInitConfigTest}

const countOfConfigTests = 2

//
// Get Config
//

func TestGetConfigNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[0].(GeneralLoggerConfig).IsCallerToSet, t, "result.logger[0].IsCallerToSet")
	testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
}

func TestGetConfigAlreadyExistingFromNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	GetConfig()

	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_DEBUG)

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
	testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
}

func TestGetConfigNonExistingFile(t *testing.T) {
	os.Clearenv()
	os.Setenv(LOG_CONFIG_FILE_ENV_NAME, "any/path/to/config.yaml")
	configInitialized = false

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
	testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
}

func TestGetConfigCaller(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	os.Setenv(LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME, "true")

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
	testutil.AssertTrue(result.Logger[0].(GeneralLoggerConfig).IsCallerToSet, t, "result.logger[0].IsCallerToSet")
	testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
}

func TestGetConfigDefaultDelimiter(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+DELIMITER_PARAMETER, ":")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(":", result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
		testutil.AssertEquals(2, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[0].GetCommon().EnvNamesToLog[0], t, "result.Formatter[0].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[0].GetCommon().EnvNamesToLog[1], t, "result.Formatter[0].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigDefaultTemplate(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_PARAMETER, "time: %s severity: %s message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CUSTOM_PARAMETER, "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_PARAMETER, "time: %s severity: %s caller:%s file:%s line:%d message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s caller:%s file:%s line:%d message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_CUSTOM_PARAMETER, "time: %s severity: %s caller:%s file:%s line:%d message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[0].(TemplateFormatterConfig).Template, t, "result.formatter[0].template")
		testutil.AssertFalse(result.Formatter[0].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[0].IsDefaultTemplate")
		testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[0].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
		testutil.AssertFalse(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCorrelationIdTemplate")
		testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[0].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[0].customTemplate")
		testutil.AssertFalse(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[0].IsDefaultCustomTemplate")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
		testutil.AssertEquals("time: %s severity: %s caller:%s file:%s line:%d message: %s", result.Formatter[0].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[0].CallerTemplate")
		testutil.AssertFalse(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[0].IsDefaultCallerTemplate")
		testutil.AssertEquals("time: %s severity: %s correlation: %s caller:%s file:%s line:%d message: %s", result.Formatter[0].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[0].CallerCorrelationIdTemplate")
		testutil.AssertFalse(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCallerCorrelationIdTemplate")
		testutil.AssertEquals("time: %s severity: %s caller:%s file:%s line:%d message: %s %s: %s %s: %d %s: %t", result.Formatter[0].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[0].CallerCustomTemplate")
		testutil.AssertFalse(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[0].IsDefaultCallerCustomTemplate")
		testutil.AssertEquals(2, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[0].GetCommon().EnvNamesToLog[0], t, "result.Formatter[0].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[0].GetCommon().EnvNamesToLog[1], t, "result.Formatter[0].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigDefaultTemplateInactiveSequenceWithoutParameter(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+SEQUENCE_ACTIVE_PARAMETER, "false")
		allPostInitConfigTest[i](optionalFile)
		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals("[%s] %s: %s", result.Formatter[0].(TemplateFormatterConfig).Template, t, "result.formatter[0].template")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[0].IsDefaultTemplate")
		testutil.AssertEquals("[%s] %s %s: %s", result.Formatter[0].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCorrelationIdTemplate")
		testutil.AssertEquals("[%s] %s: %s", result.Formatter[0].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[0].customTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[0].IsDefaultCustomTemplate")
		testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
		testutil.AssertEquals("[%s] %s %s(%s.%d): %s", result.Formatter[0].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[0].CallerTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[0].IsDefaultCallerTemplate")
		testutil.AssertEquals("[%s] %s %s %s(%s.%d): %s", result.Formatter[0].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[0].CallerCorrelationIdTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCallerCorrelationIdTemplate")
		testutil.AssertEquals("[%s] %s %s(%s.%d): %s", result.Formatter[0].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[0].CallerCustomTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[0].IsDefaultCallerCustomTemplate")
		testutil.AssertEquals(0, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
	}
}

func TestGetConfigDefaultTemplateActiveSequenceWithoutParameter(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
		allPostInitConfigTest[i](optionalFile)
		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals("[%s] %d %s: %s", result.Formatter[0].(TemplateFormatterConfig).Template, t, "result.formatter[0].template")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[0].IsDefaultTemplate")
		testutil.AssertEquals("[%s] %d %s %s: %s", result.Formatter[0].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCorrelationIdTemplate")
		testutil.AssertEquals("[%s] %d %s: %s", result.Formatter[0].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[0].customTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[0].IsDefaultCustomTemplate")
		testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
		testutil.AssertEquals("[%s] %d %s %s(%s.%d): %s", result.Formatter[0].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[0].CallerTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[0].IsDefaultCallerTemplate")
		testutil.AssertEquals("[%s] %d %s %s %s(%s.%d): %s", result.Formatter[0].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[0].CallerCorrelationIdTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCallerCorrelationIdTemplate")
		testutil.AssertEquals("[%s] %d %s %s(%s.%d): %s", result.Formatter[0].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[0].CallerCustomTemplate")
		testutil.AssertTrue(result.Formatter[0].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[0].IsDefaultCallerCustomTemplate")
		testutil.AssertEquals(0, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
	}
}

func TestGetConfigDefaultJson(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_TIME_KEY_PARAMETER, "timing")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_SEVERITY_KEY_PARAMETER, "level")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CORRELATION_KEY_PARAMETER, "cor")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_MESSAGE_KEY_PARAMETER, "msg")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CUSTOM_VALUES_KEY_PARAMETER, "customValues")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CUSTOM_VALUES_SUB_PARAMETER, "true")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_FUNCTION_KEY_PARAMETER, "callerFunction")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_FILE_KEY_PARAMETER, "callerFile")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_LINE_KEY_PARAMETER, "callerFileLine")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals("timing", result.Formatter[0].(JsonFormatterConfig).TimeKey, t, "result.formatter[0].timeKey")
		testutil.AssertEquals("level", result.Formatter[0].(JsonFormatterConfig).SeverityKey, t, "result.formatter[0].severityKey")
		testutil.AssertEquals("cor", result.Formatter[0].(JsonFormatterConfig).CorrelationKey, t, "result.formatter[0].correlationKey")
		testutil.AssertEquals("msg", result.Formatter[0].(JsonFormatterConfig).MessageKey, t, "result.formatter[0].messageKey")
		testutil.AssertEquals("customValues", result.Formatter[0].(JsonFormatterConfig).CustomValuesKey, t, "result.formatter[0].customValuesKey")
		testutil.AssertTrue(result.Formatter[0].(JsonFormatterConfig).CustomValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
		testutil.AssertEquals("callerFunction", result.Formatter[0].(JsonFormatterConfig).CallerFunctionKey, t, "result.formatter[0].CallerFunctionKey")
		testutil.AssertEquals("callerFile", result.Formatter[0].(JsonFormatterConfig).CallerFileKey, t, "result.formatter[0].CallerFileKey")
		testutil.AssertEquals("callerFileLine", result.Formatter[0].(JsonFormatterConfig).CallerFileLineKey, t, "result.formatter[0].CallerFileLineKey")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
		testutil.AssertEquals(2, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[0].GetCommon().EnvNamesToLog[0], t, "result.Formatter[0].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[0].GetCommon().EnvNamesToLog[1], t, "result.Formatter[0].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigDefaultJsonWithoutParameter(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals("time", result.Formatter[0].(JsonFormatterConfig).TimeKey, t, "result.formatter[0].timeKey")
		testutil.AssertEquals("severity", result.Formatter[0].(JsonFormatterConfig).SeverityKey, t, "result.formatter[0].severityKey")
		testutil.AssertEquals("correlation", result.Formatter[0].(JsonFormatterConfig).CorrelationKey, t, "result.formatter[0].correlationKey")
		testutil.AssertEquals("message", result.Formatter[0].(JsonFormatterConfig).MessageKey, t, "result.formatter[0].messageKey")
		testutil.AssertEquals("custom", result.Formatter[0].(JsonFormatterConfig).CustomValuesKey, t, "result.formatter[0].customValuesKey")
		testutil.AssertFalse(result.Formatter[0].(JsonFormatterConfig).CustomValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
		testutil.AssertEquals("caller", result.Formatter[0].(JsonFormatterConfig).CallerFunctionKey, t, "result.formatter[0].CallerFunctionKey")
		testutil.AssertEquals("file", result.Formatter[0].(JsonFormatterConfig).CallerFileKey, t, "result.formatter[0].CallerFileKey")
		testutil.AssertEquals("line", result.Formatter[0].(JsonFormatterConfig).CallerFileLineKey, t, "result.formatter[0].CallerFileLineKey")
		testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout(), t, "result.formatter[0].TimeLayout()")
		testutil.AssertEquals(0, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
	}
}

func TestGetConfigDefaultFileAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[0].(FileAppenderConfig).PathToLogFile, t, "result.appender[0].pathToLogFile")
		testutil.AssertEquals("", result.Appender[0].(FileAppenderConfig).CronExpression, t, "result.appender[0].CronExpression")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	}
}

func TestGetConfigDefaultFileAppenderMissingPath(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	}
}

func TestGetConfigDefaultUnknown(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, "abc")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, "123")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	}
}

func TestGetConfigCronFileAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	cronExpression := "* * * * *"
	limitByteSize := "64"
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME, cronExpression)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME, limitByteSize)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[0].(FileAppenderConfig).PathToLogFile, t, "result.appender[0].pathToLogFile")
		testutil.AssertEquals(cronExpression, result.Appender[0].(FileAppenderConfig).CronExpression, t, "result.appender[0].CronExpression")
		testutil.AssertEquals(limitByteSize, result.Appender[0].(FileAppenderConfig).LimitByteSize, t, "result.appender[0].LimitByteSize")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
	}
}

//
// Get Config with packages
//

func TestGetConfigPackageDelimiter(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+DELIMITER_PARAMETER, "_")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+STATIC_ENV_NAMES, "param1,param2")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertEquals(0, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals("_", result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
		testutil.AssertEquals(2, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[1].GetCommon().EnvNamesToLog[0], t, "result.Formatter[1].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[1].GetCommon().EnvNamesToLog[1], t, "result.Formatter[1].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigPackageDelimiterInherit(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+DELIMITER_PARAMETER, ":")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")

		for j := range 2 {
			testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[j].FormatterType(), t, "result.formatter[j].formatterType")
			testutil.AssertEquals(":", result.Formatter[j].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[j].delimiter")
			testutil.AssertEquals(time.RFC1123Z, result.Formatter[j].TimeLayout(), t, "result.formatter[j].delimiter")
			testutil.AssertEquals(2, len(result.Formatter[j].GetCommon().EnvNamesToLog), t, "len(result.Formatter[j].GetCommon().EnvNamesToLog)")
			testutil.AssertEquals("param1", result.Formatter[j].GetCommon().EnvNamesToLog[0], t, "result.Formatter[j].GetCommon().EnvNamesToLog[0]")
			testutil.AssertEquals("param2", result.Formatter[j].GetCommon().EnvNamesToLog[1], t, "result.Formatter[j].GetCommon().EnvNamesToLog[1]")
		}
	}
}

func TestGetConfigPackageDelimiterInheritOtherDefaultType(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType(), t, "result.formatter[0].formatterType")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
		testutil.AssertEquals(DEFAULT_TIME_LAYOUT, result.Formatter[1].TimeLayout(), t, "result.formatter[1].delimiter")
		testutil.AssertEquals(2, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[1].GetCommon().EnvNamesToLog[0], t, "result.Formatter[1].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[1].GetCommon().EnvNamesToLog[1], t, "result.Formatter[1].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigPackageDelimiterNoInheriting(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_INHERIT_CONFIG_ENV_NAME, "false")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+DELIMITER_PARAMETER, ":")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].formatterType")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
		testutil.AssertEquals(DEFAULT_TIME_LAYOUT, result.Formatter[1].TimeLayout(), t, "result.formatter[1].delimiter")
		testutil.AssertEquals(0, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
	}
}

func TestGetConfigPackageDelimiterFullQualifiedPackageName(t *testing.T) {
	packageParameterName := "testPackage"
	fullQualifiedPackageName := "github.com/ma-vin/typewriter/testPackage"
	packageParameterUpperName := strings.ToUpper(packageParameterName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME, "true")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageParameterName, fullQualifiedPackageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageParameterName, LOG_LEVEL_DEBUG)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageParameterName, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageParameterName, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageParameterName+DELIMITER_PARAMETER, "_")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertTrue(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameterUpperName, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(fullQualifiedPackageName, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameterUpperName, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameterUpperName, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals("_", result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackageTemplate(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_TEMPLATE)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+TEMPLATE_PARAMETER, "time: %s severity: %s message: %s")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+TEMPLATE_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s message: %s")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+TEMPLATE_CUSTOM_PARAMETER, "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+STATIC_ENV_NAMES, "param1,param2")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertEquals(0, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[1].(TemplateFormatterConfig).Template, t, "result.formatter[1].template")
		testutil.AssertFalse(result.Formatter[1].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[1].IsDefaultTemplate")
		testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[1].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[1].correlationIdTemplate")
		testutil.AssertFalse(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[0].IsDefaultCorrelationIdTemplate")
		testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[1].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[1].customTemplate")
		testutil.AssertFalse(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[1].IsDefaultCustomTemplate")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[1].TimeLayout(), t, "result.formatter[1].timeLayout")
		testutil.AssertEquals("[%s] %d %s %s(%s.%d): %s", result.Formatter[1].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[1].CallerTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[1].IsDefaultCallerTemplate")
		testutil.AssertEquals("[%s] %d %s %s %s(%s.%d): %s", result.Formatter[1].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[1].CallerCorrelationIdTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[1].IsDefaultCallerCorrelationIdTemplate")
		testutil.AssertEquals("[%s] %d %s %s(%s.%d): %s", result.Formatter[1].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[1].CallerCustomTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[1].IsDefaultCallerCustomTemplate")
		testutil.AssertEquals(2, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[1].GetCommon().EnvNamesToLog[0], t, "result.Formatter[1].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[1].GetCommon().EnvNamesToLog[1], t, "result.Formatter[1].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigPackageTemplateInherit(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_PARAMETER, "time: %s severity: %s message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CUSTOM_PARAMETER, "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_PARAMETER, "time: %s severity: %s caller:%s file:%s line:%d message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s caller:%s file:%s line:%d message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_CUSTOM_PARAMETER, "time: %s severity: %s caller:%s file:%s line:%d message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_TEMPLATE)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].IsDefault()")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter()")

		for j := range 2 {
			testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[j].FormatterType(), t, "result.formatter[j].FormatterType()")
			testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[j].(TemplateFormatterConfig).Template, t, "result.formatter[j].template")
			testutil.AssertFalse(result.Formatter[j].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[j].IsDefaultTemplate")
			testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[j].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[j].correlationIdTemplate")
			testutil.AssertFalse(result.Formatter[j].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[j].IsDefaultCorrelationIdTemplate")
			testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[j].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[j].customTemplate")
			testutil.AssertFalse(result.Formatter[j].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[j].IsDefaultCustomTemplate")
			testutil.AssertEquals(time.RFC1123Z, result.Formatter[j].TimeLayout(), t, "result.formatter[j].TimeLayout()")
			testutil.AssertEquals("time: %s severity: %s caller:%s file:%s line:%d message: %s", result.Formatter[j].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[j].CallerTemplate")
			testutil.AssertFalse(result.Formatter[j].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[j].IsDefaultCallerTemplate")
			testutil.AssertEquals("time: %s severity: %s correlation: %s caller:%s file:%s line:%d message: %s", result.Formatter[j].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[j].CallerCorrelationIdTemplate")
			testutil.AssertFalse(result.Formatter[j].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[j].IsDefaultCallerCorrelationIdTemplate")
			testutil.AssertEquals("time: %s severity: %s caller:%s file:%s line:%d message: %s %s: %s %s: %d %s: %t", result.Formatter[j].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[j].CallerCustomTemplate")
			testutil.AssertFalse(result.Formatter[j].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[j].IsDefaultCallerCustomTemplate")
			testutil.AssertEquals(2, len(result.Formatter[j].GetCommon().EnvNamesToLog), t, "len(result.Formatter[j].GetCommon().EnvNamesToLog)")
			testutil.AssertEquals("param1", result.Formatter[j].GetCommon().EnvNamesToLog[0], t, "result.Formatter[j].GetCommon().EnvNamesToLog[0]")
			testutil.AssertEquals("param2", result.Formatter[j].GetCommon().EnvNamesToLog[1], t, "result.Formatter[j].GetCommon().EnvNamesToLog[1]")
		}
	}
}

func TestGetConfigPackageTemplateInheritOtherDefaultType(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_TEMPLATE)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].IsDefault()")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter()")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[1].FormatterType(), t, "result.formatter[1].FormatterType()")
		testutil.AssertEquals(DEFAULT_SEQUENCE_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).Template, t, "result.formatter[1].template")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[1].IsDefaultTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CORRELATION_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[1].correlationIdTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[1].IsDefaultCorrelationIdTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CUSTOM_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[1].customTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[1].IsDefaultCustomTemplate")
		testutil.AssertEquals(DEFAULT_TIME_LAYOUT, result.Formatter[1].TimeLayout(), t, "result.formatter[1].TimeLayout()")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CALLER_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[1].CallerTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[1].IsDefaultCallerTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CALLER_CORRELATION_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[1].CallerCorrelationIdTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[1].IsDefaultCallerCorrelationIdTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CALLER_CUSTOM_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[1].CallerCustomTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[1].IsDefaultCallerCustomTemplate")
		testutil.AssertEquals(2, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[1].GetCommon().EnvNamesToLog[0], t, "result.Formatter[1].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[1].GetCommon().EnvNamesToLog[1], t, "result.Formatter[1].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigPackageTemplateNoInheriting(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_INHERIT_CONFIG_ENV_NAME, "false")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_PARAMETER, "time: %s severity: %s message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CUSTOM_PARAMETER, "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_PARAMETER, "time: %s severity: %s caller:%s file:%s line:%d message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s caller:%s file:%s line:%d message: %s")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CALLER_CUSTOM_PARAMETER, "time: %s severity: %s caller:%s file:%s line:%d message: %s %s: %s %s: %d %s: %t")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_TEMPLATE)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].IsDefault()")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter()")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[1].FormatterType(), t, "result.formatter[1].FormatterType()")
		testutil.AssertEquals(DEFAULT_SEQUENCE_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).Template, t, "result.formatter[1].template")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultTemplate, t, "result.formatter[1].IsDefaultTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CORRELATION_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CorrelationIdTemplate, t, "result.formatter[1].correlationIdTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate, t, "result.formatter[1].IsDefaultCorrelationIdTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CUSTOM_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CustomTemplate, t, "result.formatter[1].customTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCustomTemplate, t, "result.formatter[1].IsDefaultCustomTemplate")
		testutil.AssertEquals(DEFAULT_TIME_LAYOUT, result.Formatter[1].TimeLayout(), t, "result.formatter[1].TimeLayout()")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CALLER_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CallerTemplate, t, "result.formatter[1].CallerTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerTemplate, t, "result.formatter[1].IsDefaultCallerTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CALLER_CORRELATION_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CallerCorrelationIdTemplate, t, "result.formatter[1].CallerCorrelationIdTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate, t, "result.formatter[1].IsDefaultCallerCorrelationIdTemplate")
		testutil.AssertEquals(DEFAULT_SEQUENCE_CALLER_CUSTOM_TEMPLATE, result.Formatter[1].(TemplateFormatterConfig).CallerCustomTemplate, t, "result.formatter[1].CallerCustomTemplate")
		testutil.AssertTrue(result.Formatter[1].(TemplateFormatterConfig).IsDefaultCallerCustomTemplate, t, "result.formatter[1].IsDefaultCallerCustomTemplate")
		testutil.AssertEquals(0, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
	}
}

func TestGetConfigPackageJson(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_JSON)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+JSON_TIME_KEY_PARAMETER, "timing")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+JSON_SEVERITY_KEY_PARAMETER, "level")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+JSON_CORRELATION_KEY_PARAMETER, "cor")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+JSON_MESSAGE_KEY_PARAMETER, "msg")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+JSON_CUSTOM_VALUES_KEY_PARAMETER, "customValues")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+JSON_CUSTOM_VALUES_SUB_PARAMETER, "true")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+STATIC_ENV_NAMES, "param1,param2")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertEquals(0, len(result.Formatter[0].GetCommon().EnvNamesToLog), t, "len(result.Formatter[0].GetCommon().EnvNamesToLog)")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals("timing", result.Formatter[1].(JsonFormatterConfig).TimeKey, t, "result.formatter[1].timeKey")
		testutil.AssertEquals("level", result.Formatter[1].(JsonFormatterConfig).SeverityKey, t, "result.formatter[1].severityKey")
		testutil.AssertEquals("cor", result.Formatter[1].(JsonFormatterConfig).CorrelationKey, t, "result.formatter[1].correlationKey")
		testutil.AssertEquals("msg", result.Formatter[1].(JsonFormatterConfig).MessageKey, t, "result.formatter[1].messageKey")
		testutil.AssertEquals("customValues", result.Formatter[1].(JsonFormatterConfig).CustomValuesKey, t, "result.formatter[1].customValuesKey")
		testutil.AssertTrue(result.Formatter[1].(JsonFormatterConfig).CustomValuesAsSubElement, t, "result.formatter[1].customValuesAsSubElement")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[1].TimeLayout(), t, "result.formatter[1].timeLayout")
		testutil.AssertEquals(2, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[1].GetCommon().EnvNamesToLog[0], t, "result.Formatter[1].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[1].GetCommon().EnvNamesToLog[1], t, "result.Formatter[1].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigPackageJsonInherit(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_TIME_KEY_PARAMETER, "timing")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_SEVERITY_KEY_PARAMETER, "level")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CORRELATION_KEY_PARAMETER, "cor")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_MESSAGE_KEY_PARAMETER, "msg")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CUSTOM_VALUES_KEY_PARAMETER, "customValues")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CUSTOM_VALUES_SUB_PARAMETER, "true")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_FUNCTION_KEY_PARAMETER, "callerFunction")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_FILE_KEY_PARAMETER, "callerFile")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_LINE_KEY_PARAMETER, "callerFileLine")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_JSON)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")

		for j := range 2 {
			testutil.AssertEquals(FORMATTER_JSON, result.Formatter[j].FormatterType(), t, "result.formatter[j].formatterType")
			testutil.AssertEquals("timing", result.Formatter[j].(JsonFormatterConfig).TimeKey, t, "result.formatter[j].timeKey")
			testutil.AssertEquals("level", result.Formatter[j].(JsonFormatterConfig).SeverityKey, t, "result.formatter[j].severityKey")
			testutil.AssertEquals("cor", result.Formatter[j].(JsonFormatterConfig).CorrelationKey, t, "result.formatter[j].correlationKey")
			testutil.AssertEquals("msg", result.Formatter[j].(JsonFormatterConfig).MessageKey, t, "result.formatter[j].messageKey")
			testutil.AssertEquals("customValues", result.Formatter[j].(JsonFormatterConfig).CustomValuesKey, t, "result.formatter[j].customValuesKey")
			testutil.AssertTrue(result.Formatter[j].(JsonFormatterConfig).CustomValuesAsSubElement, t, "result.formatter[j].customValuesAsSubElement")
			testutil.AssertEquals("callerFunction", result.Formatter[j].(JsonFormatterConfig).CallerFunctionKey, t, "result.formatter[j].CallerFunctionKey")
			testutil.AssertEquals("callerFile", result.Formatter[j].(JsonFormatterConfig).CallerFileKey, t, "result.formatter[j].CallerFileKey")
			testutil.AssertEquals("callerFileLine", result.Formatter[j].(JsonFormatterConfig).CallerFileLineKey, t, "result.formatter[j].CallerFileLineKey")
			testutil.AssertEquals(time.RFC1123Z, result.Formatter[j].TimeLayout(), t, "result.formatter[j].TimeLayout()")
			testutil.AssertEquals(2, len(result.Formatter[j].GetCommon().EnvNamesToLog), t, "len(result.Formatter[j].GetCommon().EnvNamesToLog)")
			testutil.AssertEquals("param1", result.Formatter[j].GetCommon().EnvNamesToLog[0], t, "result.Formatter[j].GetCommon().EnvNamesToLog[0]")
			testutil.AssertEquals("param2", result.Formatter[j].GetCommon().EnvNamesToLog[1], t, "result.Formatter[j].GetCommon().EnvNamesToLog[1]")
		}
	}
}

func TestGetConfigPackageJsonInheritOtherDefaultType(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_JSON)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].formatterType")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_TIME_KEY, result.Formatter[1].(JsonFormatterConfig).TimeKey, t, "result.formatter[1].timeKey")
		testutil.AssertEquals(DEFAULT_SEQUENCE_KEY, result.Formatter[1].(JsonFormatterConfig).SequenceKey, t, "result.formatter[1].SequenceKey")
		testutil.AssertEquals(DEFAULT_SEVERITY_KEY, result.Formatter[1].(JsonFormatterConfig).SeverityKey, t, "result.formatter[1].severityKey")
		testutil.AssertEquals(DEFAULT_CORRELATION_KEY, result.Formatter[1].(JsonFormatterConfig).CorrelationKey, t, "result.formatter[1].correlationKey")
		testutil.AssertEquals(DEFAULT_MESSAGE_KEY, result.Formatter[1].(JsonFormatterConfig).MessageKey, t, "result.formatter[1].messageKey")
		testutil.AssertEquals(DEFAULT_CUSTOM_VALUES_KEY, result.Formatter[1].(JsonFormatterConfig).CustomValuesKey, t, "result.formatter[1].customValuesKey")
		testutil.AssertFalse(result.Formatter[1].(JsonFormatterConfig).CustomValuesAsSubElement, t, "result.formatter[1].customValuesAsSubElement")
		testutil.AssertEquals(DEFAULT_CALLER_FUNCTION_KEY, result.Formatter[1].(JsonFormatterConfig).CallerFunctionKey, t, "result.formatter[1].CallerFunctionKey")
		testutil.AssertEquals(DEFAULT_CALLER_FILE_KEY, result.Formatter[1].(JsonFormatterConfig).CallerFileKey, t, "result.formatter[1].CallerFileKey")
		testutil.AssertEquals(DEFAULT_CALLER_FILE_LINE_KEY, result.Formatter[1].(JsonFormatterConfig).CallerFileLineKey, t, "result.formatter[1].CallerFileLineKey")
		testutil.AssertEquals(DEFAULT_TIME_LAYOUT, result.Formatter[1].TimeLayout(), t, "result.formatter[1].TimeLayout()")
		testutil.AssertEquals(2, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
		testutil.AssertEquals("param1", result.Formatter[1].GetCommon().EnvNamesToLog[0], t, "result.Formatter[1].GetCommon().EnvNamesToLog[0]")
		testutil.AssertEquals("param2", result.Formatter[1].GetCommon().EnvNamesToLog[1], t, "result.Formatter[1].GetCommon().EnvNamesToLog[1]")
	}
}

func TestGetConfigPackageJsonNoInheriting(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_INHERIT_CONFIG_ENV_NAME, "false")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_TIME_KEY_PARAMETER, "timing")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_SEVERITY_KEY_PARAMETER, "level")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CORRELATION_KEY_PARAMETER, "cor")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_MESSAGE_KEY_PARAMETER, "msg")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CUSTOM_VALUES_KEY_PARAMETER, "customValues")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CUSTOM_VALUES_SUB_PARAMETER, "true")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_FUNCTION_KEY_PARAMETER, "callerFunction")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_FILE_KEY_PARAMETER, "callerFile")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+JSON_CALLER_LINE_KEY_PARAMETER, "callerFileLine")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+STATIC_ENV_NAMES, "param1,param2")
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_JSON)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")

		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType(), t, "result.formatter[0].formatterType")

		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_TIME_KEY, result.Formatter[1].(JsonFormatterConfig).TimeKey, t, "result.formatter[1].timeKey")
		testutil.AssertEquals(DEFAULT_SEQUENCE_KEY, result.Formatter[1].(JsonFormatterConfig).SequenceKey, t, "result.formatter[1].SequenceKey")
		testutil.AssertEquals(DEFAULT_SEVERITY_KEY, result.Formatter[1].(JsonFormatterConfig).SeverityKey, t, "result.formatter[1].severityKey")
		testutil.AssertEquals(DEFAULT_CORRELATION_KEY, result.Formatter[1].(JsonFormatterConfig).CorrelationKey, t, "result.formatter[1].correlationKey")
		testutil.AssertEquals(DEFAULT_MESSAGE_KEY, result.Formatter[1].(JsonFormatterConfig).MessageKey, t, "result.formatter[1].messageKey")
		testutil.AssertEquals(DEFAULT_CUSTOM_VALUES_KEY, result.Formatter[1].(JsonFormatterConfig).CustomValuesKey, t, "result.formatter[1].customValuesKey")
		testutil.AssertFalse(result.Formatter[1].(JsonFormatterConfig).CustomValuesAsSubElement, t, "result.formatter[1].customValuesAsSubElement")
		testutil.AssertEquals(DEFAULT_CALLER_FUNCTION_KEY, result.Formatter[1].(JsonFormatterConfig).CallerFunctionKey, t, "result.formatter[1].CallerFunctionKey")
		testutil.AssertEquals(DEFAULT_CALLER_FILE_KEY, result.Formatter[1].(JsonFormatterConfig).CallerFileKey, t, "result.formatter[1].CallerFileKey")
		testutil.AssertEquals(DEFAULT_CALLER_FILE_LINE_KEY, result.Formatter[1].(JsonFormatterConfig).CallerFileLineKey, t, "result.formatter[1].CallerFileLineKey")
		testutil.AssertEquals(DEFAULT_TIME_LAYOUT, result.Formatter[1].TimeLayout(), t, "result.formatter[1].TimeLayout()")
		testutil.AssertEquals(0, len(result.Formatter[1].GetCommon().EnvNamesToLog), t, "len(result.Formatter[1].GetCommon().EnvNamesToLog)")
	}
}

func TestGetConfigPackagePartialOnlyLevel(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	packageNameLower := strings.ToLower(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackagePartialOnlyAppender(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	packageNameLower := strings.ToLower(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackagePartialOnlyFormatter(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	packageNameLower := strings.ToLower(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackagePartialOnlyFormatterWithParameterDelimiter(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	packageNameLower := strings.ToLower(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+DELIMITER_PARAMETER, "_")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals("_", result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackageFileAppender(t *testing.T) {
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[1].(FileAppenderConfig).PathToLogFile, t, "result.appender[1].pathToLogFile")
		testutil.AssertEquals("", result.Appender[1].(FileAppenderConfig).CronExpression, t, "result.appender[1].pathToLogFile")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigFromFileButAllCommentOut(t *testing.T) {
	logFilePath := "pathToLogFile"

	propertiesFile := propertiesFileInitConfigTest(t)
	propertiesFileAddValueConfigTest(propertiesFile, "#"+DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	propertiesFileAddValueConfigTest(propertiesFile, "//"+DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
	propertiesFileAddValueConfigTest(propertiesFile, "--"+DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
	propertiesFileAddValueConfigTest(propertiesFile, "/*"+DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE+"*/")
	fmt.Fprintln(propertiesFile, "/*")
	propertiesFileAddValueConfigTest(propertiesFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_PARAMETER, "time: %s severity: %s message: %s")
	propertiesFileAddValueConfigTest(propertiesFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s message: %s")
	propertiesFileAddValueConfigTest(propertiesFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TEMPLATE_CUSTOM_PARAMETER, "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	propertiesFileAddValueConfigTest(propertiesFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
	fmt.Fprintln(propertiesFile, "*/")
	fmt.Fprintln(propertiesFile, "")
	propertiesFilePostInitConfigTest(propertiesFile)
	configInitialized = false

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
	testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigPackageCronAndSizeRenamerFileAppender(t *testing.T) {
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	cronExpression := "* * * * *"
	limitByteSize := "64"
	packageParameter := strings.ToUpper(packageName)

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME+packageName, cronExpression)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME+packageName, limitByteSize)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault(), t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter(), t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName(), t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].(GeneralLoggerConfig).Severity, t, "result.logger[0].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
		testutil.AssertFalse(result.Logger[1].IsDefault(), t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter(), t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName(), t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].(GeneralLoggerConfig).Severity, t, "result.logger[1].severity")
		testutil.AssertEquals(DEFAULT_CONTEXT_CORRELATION_ID_KEY, result.Logger[1].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[1].Common.CorrelationIdKey")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[1].(FileAppenderConfig).PathToLogFile, t, "result.appender[1].pathToLogFile")
		testutil.AssertEquals(cronExpression, result.Appender[1].(FileAppenderConfig).CronExpression, t, "result.appender[1].CronExpression")
		testutil.AssertEquals(limitByteSize, result.Appender[1].(FileAppenderConfig).LimitByteSize, t, "result.appender[1].LimitByteSize")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault(), t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter(), t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType(), t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].(DelimiterFormatterConfig).Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackageFileAppenderInherit(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	logFilePath := "pathToLogFile"
	cronExpression := "* * * * *"
	limitByteSize := "64"

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME, cronExpression)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME, limitByteSize)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")

		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")

		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")

		for j := range 2 {
			testutil.AssertEquals(APPENDER_FILE, result.Appender[j].AppenderType(), t, "result.appender[j].appenderType")
			testutil.AssertEquals(logFilePath, result.Appender[j].(FileAppenderConfig).PathToLogFile, t, "result.appender[j].pathToLogFile")
			testutil.AssertEquals(cronExpression, result.Appender[j].(FileAppenderConfig).CronExpression, t, "result.appender[j].CronExpression")
			testutil.AssertEquals(limitByteSize, result.Appender[j].(FileAppenderConfig).LimitByteSize, t, "result.appender[j].LimitByteSize")
		}
	}
}

func TestGetConfigPackageFileAppenderInheritMissingPath(t *testing.T) {
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	cronExpression := "* * * * *"
	limitByteSize := "64"

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME, cronExpression)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME, limitByteSize)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")

		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")

		testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[0].(FileAppenderConfig).PathToLogFile, t, "result.appender[0].pathToLogFile")
		testutil.AssertEquals(cronExpression, result.Appender[0].(FileAppenderConfig).CronExpression, t, "result.appender[0].CronExpression")
		testutil.AssertEquals(limitByteSize, result.Appender[0].(FileAppenderConfig).LimitByteSize, t, "result.appender[0].LimitByteSize")
	}
}

func TestGetConfigPackageFileAppenderInheritOtherDefaultType(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	logFilePath := "pathToLogFile"

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")

		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[1].(FileAppenderConfig).PathToLogFile, t, "result.appender[1].pathToLogFile")
		testutil.AssertEquals("", result.Appender[1].(FileAppenderConfig).CronExpression, t, "result.appender[1].CronExpression")
		testutil.AssertEquals("", result.Appender[1].(FileAppenderConfig).LimitByteSize, t, "result.appender[1].LimitByteSize")
	}
}

func TestGetConfigPackageFileAppenderNoInheriting(t *testing.T) {
	packageName := "testPackage"
	packageParameter := strings.ToUpper(packageName)
	logFilePath := "pathToLogFile"
	cronExpression := "* * * * *"
	limitByteSize := "64"

	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_INHERIT_CONFIG_ENV_NAME, "false")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME, cronExpression)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME, limitByteSize)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")

		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType(), t, "result.appender[1].appenderType")

		testutil.AssertFalse(result.Appender[1].IsDefault(), t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter(), t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType(), t, "result.appender[1].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[1].(FileAppenderConfig).PathToLogFile, t, "result.appender[1].pathToLogFile")
		testutil.AssertEquals("", result.Appender[1].(FileAppenderConfig).CronExpression, t, "result.appender[1].CronExpression")
		testutil.AssertEquals("", result.Appender[1].(FileAppenderConfig).LimitByteSize, t, "result.appender[1].LimitByteSize")
	}
}

func TestRegisterAndDeregisterAppenderConfig(t *testing.T) {
	os.Clearenv()

	customAppenderName := "CUSTOM_APPENDER"
	customKeyPrefix := "ANY_KEY_PREFIX"

	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, customAppenderName)
	os.Setenv(DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, "pathToLogFile")

	configInitialized = false
	result := GetConfig()
	testutil.AssertNotNil(result, t, "result")

	// Check fallback to default one if CUSTOM_APPENDER is not registered yet
	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].IsDefault()")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter()")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].AppenderType()")

	// Register custom one
	err := RegisterAppenderConfig(customAppenderName, []string{customKeyPrefix}, createFileAppenderConfig)
	testutil.AssertNil(err, t, "err of RegisterAppenderConfig")

	// Load config with registered appender
	result = GetConfig()
	testutil.AssertNotNil(result, t, "registered - result")
	testutil.AssertTrue(slices.Contains(relevantKeyPrefixes, customKeyPrefix), t, "relevantKeyPrefixes contains customKeyPrefix")
	testutil.AssertEquals(1, len(result.Appender), t, "registered - len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "registered - result.appender[0].IsDefault()")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "registered - result.appender[0].PackageParameter()")
	testutil.AssertEquals(customAppenderName, result.Appender[0].AppenderType(), t, "registered - result.appender[0].AppenderType()")

	// Deregister custom one
	err = DeregisterAppenderConfig(customAppenderName)
	testutil.AssertNil(err, t, "err of DeregisterAppenderConfig")

	// Load config without registered appender: fallback to default one
	result = GetConfig()
	testutil.AssertNotNil(result, t, "deregistered - result")
	testutil.AssertEquals(1, len(result.Appender), t, "deregistered - len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault(), t, "deregistered - lresult.appender[0].IsDefault()")
	testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "deregistered - lresult.appender[0].PackageParameter()")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "deregistered - lresult.appender[0].AppenderType()")
}

func TestRegisterKnownAppenderConfig(t *testing.T) {
	err := RegisterAppenderConfig(APPENDER_FILE, []string{}, createFileAppenderConfig)
	testutil.AssertNotNil(err, t, "registered known appender")
}

func TestDeregisterBuildInAppenderConfig(t *testing.T) {
	err := DeregisterAppenderConfig(APPENDER_STDOUT)
	testutil.AssertNotNil(err, t, "deregistered standard output appender")

	err = DeregisterAppenderConfig(APPENDER_FILE)
	testutil.AssertNotNil(err, t, "deregistered file appender")
}

func TestDeregisterUnknownAppenderConfig(t *testing.T) {
	err := DeregisterAppenderConfig("Anything")
	testutil.AssertNotNil(err, t, "deregistered unknown appender")
}

func TestRegisterAndDeregisterFormatterConfig(t *testing.T) {
	os.Clearenv()

	customFormatterName := "CUSTOM_FORMATTER"
	customKeyPrefix := "ANY_KEY_PREFIX"

	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, customFormatterName)

	configInitialized = false
	result := GetConfig()
	testutil.AssertNotNil(result, t, "result")

	// Check fallback to default one if CUSTOM_FORMATTER is not registered yet
	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "result.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")

	// Register custom one
	err := RegisterFormatterConfig(customFormatterName, []string{customKeyPrefix}, createJsonFormatterConfig)
	testutil.AssertNil(err, t, "err of RegisterFormatterConfig")

	// Load config with registered formatter
	result = GetConfig()
	testutil.AssertNotNil(result, t, "registered - result")
	testutil.AssertTrue(slices.Contains(relevantKeyPrefixes, customKeyPrefix), t, "relevantKeyPrefixes contains customKeyPrefix")
	testutil.AssertEquals(1, len(result.Formatter), t, "registered - len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "registered - result.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "registered - result.formatter[0].PackageParameter()")
	testutil.AssertEquals(customFormatterName, result.Formatter[0].FormatterType(), t, "registered - result.formatter[0].FormatterType()")

	// Deregister custom one
	err = DeregisterFormatterConfig(customFormatterName)
	testutil.AssertNil(err, t, "err of DeregisterFormatterConfig")

	// Load config without registered formatter: fallback to default one
	result = GetConfig()
	testutil.AssertNotNil(result, t, "deregistered - result")
	testutil.AssertEquals(1, len(result.Formatter), t, "deregistered - len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault(), t, "deregistered - lresult.formatter[0].IsDefault()")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter(), t, "deregistered - lresult.formatter[0].PackageParameter()")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "deregistered - lresult.formatter[0].FormatterType()")
}

func TestRegisterKnownFormatterConfig(t *testing.T) {
	err := RegisterFormatterConfig(FORMATTER_DELIMITER, []string{}, createDelimiterFormatterConfig)
	testutil.AssertNotNil(err, t, "registered known formatter")
}

func TestDeregisterBuildInFormatterConfig(t *testing.T) {
	err := DeregisterFormatterConfig(FORMATTER_DELIMITER)
	testutil.AssertNotNil(err, t, "deregistered delimiter formatter")

	err = DeregisterFormatterConfig(FORMATTER_TEMPLATE)
	testutil.AssertNotNil(err, t, "deregistered template formatter")

	err = DeregisterFormatterConfig(FORMATTER_JSON)
	testutil.AssertNotNil(err, t, "deregistered json formatter")
}

func TestDeregisterUnknownFormatterConfig(t *testing.T) {
	err := DeregisterFormatterConfig("Anything")
	testutil.AssertNotNil(err, t, "deregistered unknown formatter")
}

func TestGetConfigMultiAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, fmt.Sprintf("%s, %s", APPENDER_STDOUT, APPENDER_FILE))
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_MULTIPLE, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")
		multiAppender := result.Appender[0].(MultiAppenderConfig)
		testutil.AssertEquals(2, len(*multiAppender.AppenderConfigs), t, "len(*multiAppender.AppenderConfigs)")
		testutil.AssertEquals(APPENDER_STDOUT, (*multiAppender.AppenderConfigs)[0].AppenderType(), t, "(*multiAppender.AppenderConfigs)[0].AppenderType()")
		testutil.AssertEquals(APPENDER_FILE, (*multiAppender.AppenderConfigs)[1].AppenderType(), t, "(*multiAppender.AppenderConfigs)[1].AppenderType()")
		testutil.AssertEquals(logFilePath, (*multiAppender.AppenderConfigs)[1].(FileAppenderConfig).PathToLogFile, t, "(*multiAppender.AppenderConfigs)[1].(FileAppenderConfig).PathToLogFile")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	}
}

func TestGetConfigMultiAppenderUnknown(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, fmt.Sprintf("%s, %s", "abc", APPENDER_FILE))
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault(), t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter(), t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType(), t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType(), t, "result.formatter[0].FormatterType()")
	}
}

func TestGetConfigCorrelationIdKeyWithoutLoggerConfig(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_CONTEXT_CORRELATION_ID_KEY_ENV_NAME, "SomethingNew")
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertEquals("SomethingNew", result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
	}
}

func TestGetConfigCorrelationIdKeyWithLoggerConfig(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, LOG_CONFIG_CONTEXT_CORRELATION_ID_KEY_ENV_NAME, "SomethingNew")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_DEBUG)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertEquals("SomethingNew", result.Logger[0].(GeneralLoggerConfig).Common.CorrelationIdKey, t, "result.logger[0].Common.CorrelationIdKey")
	}
}
