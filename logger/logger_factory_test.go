package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/testutil"
)

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

func TestGetConfigExistingNoEnv(t *testing.T) {
	os.Clearenv()
	configInitialized = false

	result := getConfig()

	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_DEBUG)

	result = getConfig()

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

func TestGetConfigCreateFromEnvDefault(t *testing.T) {
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

func TestGetConfigCreateFromEnvDefaultNotSupportedYet(t *testing.T) {
	os.Clearenv()
	os.Setenv(DEFAULT_LOG_LEVEL_ENV_NAME, LOG_LEVEL_INFO)
	os.Setenv(DEFAULT_LOG_APPENDER_ENV_NAME, APPENDER_FILE)
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
	testutil.AssertEquals("", result.appender[0].appenderType, t, "result.appender[0].appenderType")

	testutil.AssertEquals(1, len(result.formatter), t, "len(result.formatter)")
	testutil.AssertTrue(result.formatter[0].isDefault, t, "result.formatter[0].isDefault")
	testutil.AssertEquals("", result.formatter[0].packageName, t, "result.formatter[0].packageName")
	testutil.AssertEquals("", result.formatter[0].formatterType, t, "result.formatter[0].formatterType")
	testutil.AssertEquals("", result.formatter[0].delimiter, t, "result.formatter[0].delimiter")
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

func TestGetConfigCreateFromEnvPackage(t *testing.T) {
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

func TestGetConfigCreateFromEnvPackagePartialOnlyFromatterWithParamter(t *testing.T) {
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
