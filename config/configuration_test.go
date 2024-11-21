package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/testutil"
)

//
// Get Config
//

func TestGetConfigCreateNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigExistingFromNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	GetConfig()

	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_DEBUG)

	result := GetConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromFile(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_CONFIG_FILE_ENV_NAME, "any/path/to/config.yaml")
	configInitialized = false

	result := GetConfig()

	testutil.AssertNil(result, t, "result")
}

func TestGetConfigCreateFromEnvDefaultDelimiter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_DELIMITER)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME, ":")
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(":", result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromEnvDefaultTemplate(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_1", "time: %s severity: %s message: %s")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_2", "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_3", "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_4", time.RFC1123Z)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[0].Template, t, "result.formatter[0].template")
	testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[0].CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
	testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[0].CustomTemplate, t, "result.formatter[0].customTemplate")
	testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultTemplateWithoutParameter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_TEMPLATE)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("[%s] %s: %s", result.Formatter[0].Template, t, "result.formatter[0].template")
	testutil.AssertEquals("[%s] %s %s: %s", result.Formatter[0].CorrelationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
	testutil.AssertEquals("[%s] %s: %s", result.Formatter[0].CustomTemplate, t, "result.formatter[0].customTemplate")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultJson(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_1", "timing")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_2", "level")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_3", "cor")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_4", "msg")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_5", "customValues")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_6", "true")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_7", time.RFC1123Z)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("timing", result.Formatter[0].TimeKey, t, "result.formatter[0].timeKey")
	testutil.AssertEquals("level", result.Formatter[0].SeverityKey, t, "result.formatter[0].severityKey")
	testutil.AssertEquals("cor", result.Formatter[0].CorrelationKey, t, "result.formatter[0].correlationKey")
	testutil.AssertEquals("msg", result.Formatter[0].MessageKey, t, "result.formatter[0].messageKey")
	testutil.AssertEquals("customValues", result.Formatter[0].CustomValuesKey, t, "result.formatter[0].customValuesKey")
	testutil.AssertTrue(result.Formatter[0].CustomValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
	testutil.AssertEquals(time.RFC1123Z, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultJsonWithoutParameter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, FORMATTER_JSON)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_JSON, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("time", result.Formatter[0].TimeKey, t, "result.formatter[0].timeKey")
	testutil.AssertEquals("severity", result.Formatter[0].SeverityKey, t, "result.formatter[0].severityKey")
	testutil.AssertEquals("correlation", result.Formatter[0].CorrelationKey, t, "result.formatter[0].correlationKey")
	testutil.AssertEquals("message", result.Formatter[0].MessageKey, t, "result.formatter[0].messageKey")
	testutil.AssertEquals("custom", result.Formatter[0].CustomValuesKey, t, "result.formatter[0].customValuesKey")
	testutil.AssertFalse(result.Formatter[0].CustomValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
	testutil.AssertEquals(time.RFC3339, result.Formatter[0].TimeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultFileAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
	os.Setenv(DEFAULT_LOG_APPENDER_PARAMETER_PROPERTY_NAME, logFilePath)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_FILE, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertEquals(logFilePath, result.Appender[0].PathToLogFile, t, "result.appender[0].pathToLogFile")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromEnvDefaultFileAppenderMissingPath(t *testing.T) {
	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, APPENDER_FILE)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromEnvDefaultUnkown(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_PROPERTY_NAME, "abc")
	os.Setenv(DEFAULT_LOG_FORMATTER_PROPERTY_NAME, "123")
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(1, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals("", result.Appender[0].AppenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals("", result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("", result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
}

//
// Get Config with packages
//

func TestGetConfigCreateFromEnvPackageDelimiter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName, "_")
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("_", result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackageTemplate(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_TEMPLATE)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_1", "time: %s severity: %s message: %s")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_2", "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_3", "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_4", time.RFC1123Z)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_TEMPLATE, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("time: %s severity: %s message: %s", result.Formatter[1].Template, t, "result.formatter[1].template")
	testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.Formatter[1].CorrelationIdTemplate, t, "result.formatter[1].correlationIdTemplate")
	testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.Formatter[1].CustomTemplate, t, "result.formatter[1].customTemplate")
	testutil.AssertEquals(time.RFC1123Z, result.Formatter[1].TimeLayout, t, "result.formatter[1].timeLayout")
}

func TestGetConfigCreateFromEnvPackageJson(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_JSON)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_1", "timing")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_2", "level")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_3", "cor")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_4", "msg")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_5", "customValues")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_6", "true")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+"_7", time.RFC1123Z)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_JSON, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("timing", result.Formatter[1].TimeKey, t, "result.formatter[1].timeKey")
	testutil.AssertEquals("level", result.Formatter[1].SeverityKey, t, "result.formatter[1].severityKey")
	testutil.AssertEquals("cor", result.Formatter[1].CorrelationKey, t, "result.formatter[1].correlationKey")
	testutil.AssertEquals("msg", result.Formatter[1].MessageKey, t, "result.formatter[1].messageKey")
	testutil.AssertEquals("customValues", result.Formatter[1].CustomValuesKey, t, "result.formatter[1].customValuesKey")
	testutil.AssertTrue(result.Formatter[1].CustomValuesAsSubElement, t, "result.formatter[1].customValuesAsSubElement")
	testutil.AssertEquals(time.RFC1123Z, result.Formatter[1].TimeLayout, t, "result.formatter[1].timeLayout")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyLevel(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, LOG_LEVEL_DEBUG)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyAppender(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_STDOUT)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyFromatter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyFromatterWithParamterDelimiter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, FORMATTER_DELIMITER)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName, "_")
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("_", result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackageFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, APPENDER_FILE)
	os.Setenv(PACKAGE_LOG_APPENDER_PARAMETER_PROPERTY_NAME+packageName, logFilePath)
	configInitialized = false

	result := GetConfig()

	testutil.AssertEquals(2, len(result.Logger), t, "len(result.logger)")
	testutil.AssertTrue(result.Logger[0].IsDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.Logger[0].PackageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[0].Severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.Logger[1].IsDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Logger[1].PackageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.Logger[1].Severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.Appender), t, "len(result.appender)")
	testutil.AssertTrue(result.Appender[0].IsDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.Appender[0].PackageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.Appender[0].AppenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.Appender[1].IsDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Appender[1].PackageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_FILE, result.Appender[1].AppenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.Formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.Formatter[0].IsDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.Formatter[0].PackageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[0].FormatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[0].Delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.Formatter[1].IsDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.Formatter[1].PackageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.Formatter[1].FormatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.Formatter[1].Delimiter, t, "result.formatter[1].delimiter")
}
