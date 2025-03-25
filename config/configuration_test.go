package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/testutil"
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

	_, filename, _, _ := runtime.Caller(0)
	pathToFile := strings.Replace(filename, ".go", "_LogConfig_scratch.properties", 1)

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
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[0].IsCallerToSet, t, "result.logger[0].IsCallerToSet")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].TimeLayout")
}

func TestGetConfigAlreadyExistingFromNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	GetConfig()

	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_DEBUG)

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].TimeLayout")
}

func TestGetConfigNonExistingFile(t *testing.T) {
	os.Clearenv()
	os.Setenv(LOG_CONFIG_FILE_ENV_NAME, "any/path/to/config.yaml")
	configInitialized = false

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].TimeLayout")
}

func TestGetConfigCaller(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	os.Setenv(LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME, "true")

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertTrue(result.Logger[0].IsCallerToSet, t, "result.logger[0].IsCallerToSet")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].TimeLayout")
}

func TestGetConfigDefaultDelimiter(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_DELIMITER)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+DELIMITER_PARAMETER, ":")
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+TIME_LAYOUT_PARAMETER, time.RFC1123Z)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(":", result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout, t, "result.formatter[0].TimeLayout")
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
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[0].Template, t, "result.formatter[0].template")
		testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[0].CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
		testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[0].CustomTemplate, t, "result.formatter[0].customTemplate")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
		testutil.AssertEquals("time: %s severity: %s caller:%s file:%s line:%d message: %s", result.Formatter[0].CallerTemplate, t, "result.formatter[0].CallerTemplate")
		testutil.AssertEquals("time: %s severity: %s correlation: %s caller:%s file:%s line:%d message: %s", result.Formatter[0].CallerCorrelationIdTemplate, t, "result.formatter[0].CallerCorrelationIdTemplate")
		testutil.AssertEquals("time: %s severity: %s caller:%s file:%s line:%d message: %s %s: %s %s: %d %s: %t", result.Formatter[0].CallerCustomTemplate, t, "result.formatter[0].CallerCustomTemplate")
	}
}

func TestGetConfigDefaultTemplateWithoutParameter(t *testing.T) {
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
		allPostInitConfigTest[i](optionalFile)
		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals("[%s] %s: %s", result.Formatter[0].Template, t, "result.formatter[0].template")
		testutil.AssertEquals("[%s] %s %s: %s", result.Formatter[0].CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
		testutil.AssertEquals("[%s] %s: %s", result.Formatter[0].CustomTemplate, t, "result.formatter[0].customTemplate")
		testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
		testutil.AssertEquals("[%s] %s %s(%s.%d): %s", result.Formatter[0].CallerTemplate, t, "result.formatter[0].CallerTemplate")
		testutil.AssertEquals("[%s] %s %s %s(%s.%d): %s", result.Formatter[0].CallerCorrelationIdTemplate, t, "result.formatter[0].CallerCorrelationIdTemplate")
		testutil.AssertEquals("[%s] %s %s(%s.%d): %s", result.Formatter[0].CallerCustomTemplate, t, "result.formatter[0].CallerCustomTemplate")
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
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals("timing", result.Formatter[0].TimeKey, t, "result.formatter[0].timeKey")
		testutil.AssertEquals("level", result.Formatter[0].SeverityKey, t, "result.formatter[0].severityKey")
		testutil.AssertEquals("cor", result.Formatter[0].CorrelationKey, t, "result.formatter[0].correlationKey")
		testutil.AssertEquals("msg", result.Formatter[0].MessageKey, t, "result.formatter[0].messageKey")
		testutil.AssertEquals("customValues", result.Formatter[0].CustomValuesKey, t, "result.formatter[0].customValuesKey")
		testutil.AssertTrue(result.Formatter[0].CustomValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
		testutil.AssertEquals("callerFunction", result.Formatter[0].CallerFunctionKey, t, "result.formatter[0].CallerFunctionKey")
		testutil.AssertEquals("callerFile", result.Formatter[0].CallerFileKey, t, "result.formatter[0].CallerFileKey")
		testutil.AssertEquals("callerFileLine", result.Formatter[0].CallerFileLineKey, t, "result.formatter[0].CallerFileLineKey")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals("time", result.Formatter[0].TimeKey, t, "result.formatter[0].timeKey")
		testutil.AssertEquals("severity", result.Formatter[0].SeverityKey, t, "result.formatter[0].severityKey")
		testutil.AssertEquals("correlation", result.Formatter[0].CorrelationKey, t, "result.formatter[0].correlationKey")
		testutil.AssertEquals("message", result.Formatter[0].MessageKey, t, "result.formatter[0].messageKey")
		testutil.AssertEquals("custom", result.Formatter[0].CustomValuesKey, t, "result.formatter[0].customValuesKey")
		testutil.AssertFalse(result.Formatter[0].CustomValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
		testutil.AssertEquals("caller", result.Formatter[0].CallerFunctionKey, t, "result.formatter[0].CallerFunctionKey")
		testutil.AssertEquals("file", result.Formatter[0].CallerFileKey, t, "result.formatter[0].CallerFileKey")
		testutil.AssertEquals("line", result.Formatter[0].CallerFileLineKey, t, "result.formatter[0].CallerFileLineKey")
		testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
	}
}

func TestGetConfigDefaultFileAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	appender.SkipFileCreationForTest = true
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[0].PathToLogFile, t, "result.appender[0].pathToLogFile")
		testutil.AssertEquals("", result.Appender[0].CronExpression, t, "result.appender[0].CronExpression")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	}
}

func TestGetConfigDefaultFileAppenderMissingPath(t *testing.T) {
	appender.SkipFileCreationForTest = true
	for i := range countOfConfigTests {
		optionalFile := allInitConfigTest[i](t)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
		allAddValueConfigTest[i](optionalFile, DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals("", result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals("", result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals("", result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	}
}

func TestGetConfigCronFileAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	cronExpression := "* * * * *"
	limitByteSize := "64"
	appender.SkipFileCreationForTest = true
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

		testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[0].PathToLogFile, t, "result.appender[0].pathToLogFile")
		testutil.AssertEquals(cronExpression, result.Appender[0].CronExpression, t, "result.appender[0].CronExpression")
		testutil.AssertEquals(limitByteSize, result.Appender[0].LimitByteSize, t, "result.appender[0].LimitByteSize")

		testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
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
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals("_", result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameterUpperName, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(fullQualifiedPackageName, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameterUpperName, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameterUpperName, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals("_", result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
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
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[1].Template, t, "result.formatter[1].template")
		testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[1].CorrelationIdTemplate, t, "result.formatter[1].correlationIdTemplate")
		testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[1].CustomTemplate, t, "result.formatter[1].customTemplate")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[1].TimeLayout, t, "result.formatter[1].timeLayout")
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
		allPostInitConfigTest[i](optionalFile)

		configInitialized = false

		result := GetConfig()

		testutil.AssertFalse(result.UseFullQualifiedPackageName, t, "result.UseFullQualifiedPackageName")

		testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_JSON, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals("timing", result.Formatter[1].TimeKey, t, "result.formatter[1].timeKey")
		testutil.AssertEquals("level", result.Formatter[1].SeverityKey, t, "result.formatter[1].severityKey")
		testutil.AssertEquals("cor", result.Formatter[1].CorrelationKey, t, "result.formatter[1].correlationKey")
		testutil.AssertEquals("msg", result.Formatter[1].MessageKey, t, "result.formatter[1].messageKey")
		testutil.AssertEquals("customValues", result.Formatter[1].CustomValuesKey, t, "result.formatter[1].customValuesKey")
		testutil.AssertTrue(result.Formatter[1].CustomValuesAsSubElement, t, "result.formatter[1].customValuesAsSubElement")
		testutil.AssertEquals(time.RFC1123Z, result.Formatter[1].TimeLayout, t, "result.formatter[1].timeLayout")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageNameLower, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals("_", result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigPackageFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[1].PathToLogFile, t, "result.appender[1].pathToLogFile")
		testutil.AssertEquals("", result.Appender[1].CronExpression, t, "result.appender[1].pathToLogFile")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
	}
}

func TestGetConfigFromFileButAllCommentOut(t *testing.T) {
	logFilePath := "pathToLogFile"
	appender.SkipFileCreationForTest = true

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
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
	testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigPackageCronAndSizeRenamerFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
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
		testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
		testutil.AssertEquals("", result.Logger[0].PackageParameter, t, "result.logger[0].PackageParameter")
		testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
		testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Logger[1].PackageParameter, t, "result.logger[1].PackageParameter")
		testutil.AssertEquals(packageName, result.Logger[1].PackageName, t, "result.logger[1].PackageName")
		testutil.AssertEquals(common.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

		testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
		testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
		testutil.AssertEquals("", result.Appender[0].PackageParameter, t, "result.appender[0].PackageParameter")
		testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
		testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Appender[1].PackageParameter, t, "result.appender[1].PackageParameter")
		testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")
		testutil.AssertEquals(logFilePath, result.Appender[1].PathToLogFile, t, "result.appender[1].pathToLogFile")
		testutil.AssertEquals(cronExpression, result.Appender[1].CronExpression, t, "result.appender[1].CronExpression")
		testutil.AssertEquals(limitByteSize, result.Appender[1].LimitByteSize, t, "result.appender[1].LimitByteSize")

		testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
		testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
		testutil.AssertEquals("", result.Formatter[0].PackageParameter, t, "result.formatter[0].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
		testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
		testutil.AssertEquals(packageParameter, result.Formatter[1].PackageParameter, t, "result.formatter[1].PackageParameter")
		testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
		testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
	}
}
