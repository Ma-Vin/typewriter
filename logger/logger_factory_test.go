package logger

import (
	"os"
	"reflect"
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

	result := getConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigExistingFromNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	getConfig()

	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_DEBUG)

	result := getConfig()

	testutil.AssertNotNil(result, t, "result")

	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromFile(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_CONFIG_FILE_ENV_NAME, "any/path/to/config.yaml")
	configInitialized = false

	result := getConfig()

	testutil.AssertNil(result, t, "result")
}

func TestGetConfigCreateFromEnvDefaultDelimiter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_DELIMITER)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME, ":")
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(":", result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromEnvDefaultTemplate(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_TEMPLATE)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_1", "time: %s severity: %s message: %s")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_2", "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_3", "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_4", time.RFC1123Z)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_TEMPLATE, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("time: %s severity: %s message: %s", result.formatter[0].template, t, "result.formatter[0].template")
	testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.formatter[0].correlationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
	testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.formatter[0].customTemplate, t, "result.formatter[0].customTemplate")
	testutil.AssertEquals(time.RFC1123Z, result.formatter[0].timeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultTemplateWithoutParameter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_TEMPLATE)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_TEMPLATE, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("[%s] %s: %s", result.formatter[0].template, t, "result.formatter[0].template")
	testutil.AssertEquals("[%s] %s %s: %s", result.formatter[0].correlationIdTemplate, t, "result.formatter[0].correlationIdTemplate")
	testutil.AssertEquals("[%s] %s: %s", result.formatter[0].customTemplate, t, "result.formatter[0].customTemplate")
	testutil.AssertEquals(time.RFC3339, result.formatter[0].timeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultJson(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_JSON)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_1", "timing")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_2", "level")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_3", "cor")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_4", "msg")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_5", "customValues")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_6", "true")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_7", time.RFC1123Z)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_JSON, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("timing", result.formatter[0].timeKey, t, "result.formatter[0].timeKey")
	testutil.AssertEquals("level", result.formatter[0].severityKey, t, "result.formatter[0].severityKey")
	testutil.AssertEquals("cor", result.formatter[0].correlationKey, t, "result.formatter[0].correlationKey")
	testutil.AssertEquals("msg", result.formatter[0].messageKey, t, "result.formatter[0].messageKey")
	testutil.AssertEquals("customValues", result.formatter[0].customValuesKey, t, "result.formatter[0].customValuesKey")
	testutil.AssertTrue(result.formatter[0].customValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
	testutil.AssertEquals(time.RFC1123Z, result.formatter[0].timeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultJsonWithoutParameter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_JSON)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_JSON, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("time", result.formatter[0].timeKey, t, "result.formatter[0].timeKey")
	testutil.AssertEquals("severity", result.formatter[0].severityKey, t, "result.formatter[0].severityKey")
	testutil.AssertEquals("correlation", result.formatter[0].correlationKey, t, "result.formatter[0].correlationKey")
	testutil.AssertEquals("message", result.formatter[0].messageKey, t, "result.formatter[0].messageKey")
	testutil.AssertEquals("custom", result.formatter[0].customValuesKey, t, "result.formatter[0].customValuesKey")
	testutil.AssertFalse(result.formatter[0].customValuesAsSubElement, t, "result.formatter[0].customValuesAsSubElement")
	testutil.AssertEquals(time.RFC3339, result.formatter[0].timeLayout, t, "result.formatter[0].timeLayout")
}

func TestGetConfigCreateFromEnvDefaultFileAppender(t *testing.T) {
	logFilePath := "pathToLogFile"
	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_FILE)
	os.Setenv(DEFAULT_LOG_APPENDER_PARAMETER_ENV_NAME, logFilePath)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_FILE, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertEquals(logFilePath, result.appender[0].pathToLogFile, t, "result.appender[0].pathToLogFile")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromEnvDefaultFileAppenderMissingPath(t *testing.T) {
	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_FILE)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
}

func TestGetConfigCreateFromEnvDefaultUnkown(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, "abc")
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, "123")
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(1, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.INFORMATION_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")

	testutil.AssertEquals(1, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals("", result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals("", result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("", result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
}

//
// Get Config with packages
//

func TestGetConfigCreateFromEnvPackageDelimiter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_ENV_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_DELIMITER)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName, "_")
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("_", result.formatter[1].delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackageTemplate(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_ENV_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_TEMPLATE)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_1", "time: %s severity: %s message: %s")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_2", "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_3", "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_4", time.RFC1123Z)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_TEMPLATE, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("time: %s severity: %s message: %s", result.formatter[1].template, t, "result.formatter[1].template")
	testutil.AssertEquals("time: %s severity: %s correlation: %s message: %s", result.formatter[1].correlationIdTemplate, t, "result.formatter[1].correlationIdTemplate")
	testutil.AssertEquals("time: %s severity: %s message: %s %s: %s %s: %d %s: %t", result.formatter[1].customTemplate, t, "result.formatter[1].customTemplate")
	testutil.AssertEquals(time.RFC1123Z, result.formatter[1].timeLayout, t, "result.formatter[1].timeLayout")
}

func TestGetConfigCreateFromEnvPackageJson(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_ENV_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_JSON)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_1", "timing")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_2", "level")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_3", "cor")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_4", "msg")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_5", "customValues")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_6", "true")
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName+"_7", time.RFC1123Z)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_JSON, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("timing", result.formatter[1].timeKey, t, "result.formatter[1].timeKey")
	testutil.AssertEquals("level", result.formatter[1].severityKey, t, "result.formatter[1].severityKey")
	testutil.AssertEquals("cor", result.formatter[1].correlationKey, t, "result.formatter[1].correlationKey")
	testutil.AssertEquals("msg", result.formatter[1].messageKey, t, "result.formatter[1].messageKey")
	testutil.AssertEquals("customValues", result.formatter[1].customValuesKey, t, "result.formatter[1].customValuesKey")
	testutil.AssertTrue(result.formatter[1].customValuesAsSubElement, t, "result.formatter[1].customValuesAsSubElement")
	testutil.AssertEquals(time.RFC1123Z, result.formatter[1].timeLayout, t, "result.formatter[1].timeLayout")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyLevel(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_ENV_NAME+packageName, LOG_LEVEL_DEBUG)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.DEBUG_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[1].delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyAppender(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_STDOUT)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[1].delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyFromatter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_DELIMITER)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[1].delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackagePartialOnlyFromatterWithParamterDelimiter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_DELIMITER)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName, "_")
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals("_", result.formatter[1].delimiter, t, "result.formatter[1].delimiter")
}

func TestGetConfigCreateFromEnvPackageFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_FILE)
	os.Setenv(PACKAGE_LOG_APPENDER_PARAMETER_ENV_NAME+packageName, logFilePath)
	configInitialized = false

	result := getConfig()

	testutil.AssertEquals(2, len(result.logger), t, "len(result.logger)")
	testutil.AssertTrue(result.logger[0].isDefault, t, "result.logger[0].isDefault")
	testutil.AssertEquals("", result.logger[0].packageName, t, "result.logger[0].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[0].severity, t, "result.logger[0].severity")
	testutil.AssertFalse(result.logger[1].isDefault, t, "result.logger[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.logger[1].packageName, t, "result.logger[1].packageName")
	testutil.AssertEquals(constants.ERROR_SEVERITY, result.logger[1].severity, t, "result.logger[1].severity")

	testutil.AssertEquals(2, len(result.appender), t, "len(result.appender)")
	testutil.AssertTrue(result.appender[0].isDefault, t, "result.appender[0].isDefault")
	testutil.AssertEquals("", result.appender[0].packageName, t, "result.appender[0].packageName")
	testutil.AssertEquals(APPENDER_STDOUT, result.appender[0].appenderType, t, "result.appender[0].appenderType")
	testutil.AssertFalse(result.appender[1].isDefault, t, "result.appender[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.appender[1].packageName, t, "result.appender[1].packageName")
	testutil.AssertEquals(APPENDER_FILE, result.appender[1].appenderType, t, "result.appender[1].appenderType")

	testutil.AssertEquals(2, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
	testutil.AssertFalse(result.formatter[1].isDefault, t, "result.formatter[1].isDefault")
	testutil.AssertEquals(packageNameUpper, result.formatter[1].packageName, t, "result.formatter[1].packageName")
	testutil.AssertEquals(FORMATTER_DELIMITER, result.formatter[1].formatterType, t, "result.formatter[1].formatterType")
	testutil.AssertEquals(DEFAULT_DELIMITER, result.formatter[1].delimiter, t, "result.formatter[1].delimiter")
}

//
// Get Loggers
//

func TestGetLoggersCreateFromEnvDefaultDelimiter(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_DELIMITER)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME, ":")
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.commonLogger.appender).Name(), t, "commonLogger.appender.Name")

	testutil.AssertFalse(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaulTemplate(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_TEMPLATE)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_1", "time: %s severity: %s message: %s")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_2", "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_3", "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_4", time.RFC1123Z)
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.commonLogger.appender).Name(), t, "commonLogger.appender.Name")

	testutil.AssertFalse(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaulJson(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_STDOUT)
	os.Setenv(DEFAULT_LOG_FORMATTER_ENV_NAME, FORMATTER_JSON)
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_1", "timing")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_2", "level")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_3", "cor")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_4", "msg")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_5", "customValues")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_6", "true")
	os.Setenv(DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME+"_7", time.RFC1123Z)
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.commonLogger.appender).Name(), t, "commonLogger.appender.Name")

	testutil.AssertFalse(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaultFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	logFilePath := "pathToLogFile"
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_FILE)
	os.Setenv(DEFAULT_LOG_APPENDER_PARAMETER_ENV_NAME, logFilePath)
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")
	testutil.AssertEquals("FileAppender", reflect.TypeOf(*result.commonLogger.appender).Name(), t, "commonLogger.appender.Name")

	testutil.AssertFalse(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

//
// Get Loggers with packages
//

func TestGetLoggersCreateFromEnvPackage(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_ENV_NAME+packageName, LOG_LEVEL_DEBUG)
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_STDOUT)
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_DELIMITER)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName, "_")
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper], t, "packageLoggers[lowerPackageName]")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].debugEnabled, t, "packageLoggers[lowerPackageName].debugEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].informationEnabled, t, "packageLoggers[lowerPackageName].informationEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].warningEnabled, t, "packageLoggers[lowerPackageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].errorEnabled, t, "packageLoggers[lowerPackageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].fatalEnabled, t, "packageLoggers[lowerPackageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.packageLoggers[packageNameUpper].appender).Name(), t, "packageLoggers[lowerPackageName].appender.Name")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyLevel(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_LEVEL_ENV_NAME+packageName, LOG_LEVEL_DEBUG)
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper], t, "packageLoggers[lowerPackageName]")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].debugEnabled, t, "packageLoggers[lowerPackageName].debugEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].informationEnabled, t, "packageLoggers[lowerPackageName].informationEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].warningEnabled, t, "packageLoggers[lowerPackageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].errorEnabled, t, "packageLoggers[lowerPackageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].fatalEnabled, t, "packageLoggers[lowerPackageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender")

	testutil.AssertEquals(result.commonLogger.appender, result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_APPENDER_ENV_NAME+packageName, APPENDER_FILE)
	os.Setenv(PACKAGE_LOG_APPENDER_PARAMETER_ENV_NAME+packageName, logFilePath)
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper], t, "packageLoggers[lowerPackageName]")
	testutil.AssertFalse(result.packageLoggers[packageNameUpper].debugEnabled, t, "packageLoggers[lowerPackageName].debugEnabled")
	testutil.AssertFalse(result.packageLoggers[packageNameUpper].informationEnabled, t, "packageLoggers[lowerPackageName].informationEnabled")
	testutil.AssertFalse(result.packageLoggers[packageNameUpper].warningEnabled, t, "packageLoggers[lowerPackageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].errorEnabled, t, "packageLoggers[lowerPackageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].fatalEnabled, t, "packageLoggers[lowerPackageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyFromatterWithParamter(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(PACKAGE_LOG_FORMATTER_ENV_NAME+packageName, FORMATTER_DELIMITER)
	os.Setenv(PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME+packageName, "_")
	configInitialized = false
	loggersInitialized = false

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper], t, "packageLoggers[lowerPackageName]")
	testutil.AssertFalse(result.packageLoggers[packageNameUpper].debugEnabled, t, "packageLoggers[lowerPackageName].debugEnabled")
	testutil.AssertFalse(result.packageLoggers[packageNameUpper].informationEnabled, t, "packageLoggers[lowerPackageName].informationEnabled")
	testutil.AssertFalse(result.packageLoggers[packageNameUpper].warningEnabled, t, "packageLoggers[lowerPackageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].errorEnabled, t, "packageLoggers[lowerPackageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageNameUpper].fatalEnabled, t, "packageLoggers[lowerPackageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[packageNameUpper].appender, t, "packageLoggers[lowerPackageName].appender.")
}
