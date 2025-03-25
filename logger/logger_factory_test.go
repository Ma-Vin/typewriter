package logger

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/testutil"
)

//
// Get Loggers
//

func TestGetLoggersCreateFromEnvDefaultDelimiter(t *testing.T) {
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_STDOUT)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, config.FORMATTER_DELIMITER)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME, ":")
	Reset()

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
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaultTemplate(t *testing.T) {
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_STDOUT)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, config.FORMATTER_TEMPLATE)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.TEMPLATE_PARAMETER, "time: %s severity: %s message: %s")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.TEMPLATE_CORRELATION_PARAMETER, "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.TEMPLATE_CUSTOM_PARAMETER, "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.TIME_LAYOUT_PARAMETER, time.RFC1123Z)
	Reset()

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
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaultJson(t *testing.T) {
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_STDOUT)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, config.FORMATTER_JSON)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_TIME_KEY_PARAMETER, "timing")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_SEVERITY_KEY_PARAMETER, "level")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_CORRELATION_KEY_PARAMETER, "cor")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_MESSAGE_KEY_PARAMETER, "msg")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_CUSTOM_VALUES_KEY_PARAMETER, "customValues")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_CUSTOM_VALUES_SUB_PARAMETER, "true")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.TIME_LAYOUT_PARAMETER, time.RFC1123Z)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_CALLER_FUNCTION_KEY_PARAMETER, "callerFunction")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_CALLER_FILE_KEY_PARAMETER, "callerFileName")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+config.JSON_CALLER_LINE_KEY_PARAMETER, "callerFileLine")
	Reset()

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
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaultFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	logFilePath := "pathToLogFile"
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_FILE)
	os.Setenv(config.DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
	Reset()

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
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

//
// Get Loggers with packages
//

func TestGetLoggersCreateFromEnvPackage(t *testing.T) {
	packageName := "testPackage"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, config.LOG_LEVEL_DEBUG)
	os.Setenv(config.PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, config.APPENDER_STDOUT)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, config.FORMATTER_DELIMITER)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+config.DELIMITER_PARAMETER, "_")
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageName], t, "packageLoggers[packageName]")
	testutil.AssertTrue(result.packageLoggers[packageName].debugEnabled, t, "packageLoggers[packageName].debugEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].informationEnabled, t, "packageLoggers[packageName].informationEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].warningEnabled, t, "packageLoggers[packageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].errorEnabled, t, "packageLoggers[packageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].fatalEnabled, t, "packageLoggers[packageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.packageLoggers[packageName].appender).Name(), t, "packageLoggers[packageName].appender.Name")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
}

func TestGetLoggersCreateFromEnvPackageFullQualified(t *testing.T) {
	packageParameterName := "testPackage"
	fullQualifiedPackageName := "github.com/ma-vin/typewriter/testPackage"
	os.Clearenv()
	os.Setenv(config.LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME, "true")
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageParameterName, fullQualifiedPackageName)
	os.Setenv(config.PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageParameterName, config.LOG_LEVEL_DEBUG)
	os.Setenv(config.PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageParameterName, config.APPENDER_STDOUT)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageParameterName, config.FORMATTER_DELIMITER)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageParameterName+config.DELIMITER_PARAMETER, "_")
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertTrue(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[fullQualifiedPackageName], t, "packageLoggers[packageName]")
	testutil.AssertTrue(result.packageLoggers[fullQualifiedPackageName].debugEnabled, t, "packageLoggers[packageName].debugEnabled")
	testutil.AssertTrue(result.packageLoggers[fullQualifiedPackageName].informationEnabled, t, "packageLoggers[packageName].informationEnabled")
	testutil.AssertTrue(result.packageLoggers[fullQualifiedPackageName].warningEnabled, t, "packageLoggers[packageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[fullQualifiedPackageName].errorEnabled, t, "packageLoggers[packageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[fullQualifiedPackageName].fatalEnabled, t, "packageLoggers[packageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[fullQualifiedPackageName].appender, t, "packageLoggers[packageName].appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.packageLoggers[fullQualifiedPackageName].appender).Name(), t, "packageLoggers[packageName].appender.Name")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[fullQualifiedPackageName].appender, t, "packageLoggers[packageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyLevel(t *testing.T) {
	packageName := "testPackage"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, config.LOG_LEVEL_DEBUG)
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageName], t, "packageLoggers[packageName]")
	testutil.AssertTrue(result.packageLoggers[packageName].debugEnabled, t, "packageLoggers[packageName].debugEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].informationEnabled, t, "packageLoggers[packageName].informationEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].warningEnabled, t, "packageLoggers[packageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].errorEnabled, t, "packageLoggers[packageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].fatalEnabled, t, "packageLoggers[packageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender")

	testutil.AssertEquals(result.commonLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, config.APPENDER_FILE)
	os.Setenv(config.PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageName], t, "packageLoggers[packageName]")
	testutil.AssertFalse(result.packageLoggers[packageName].debugEnabled, t, "packageLoggers[packageName].debugEnabled")
	testutil.AssertFalse(result.packageLoggers[packageName].informationEnabled, t, "packageLoggers[packageName].informationEnabled")
	testutil.AssertFalse(result.packageLoggers[packageName].warningEnabled, t, "packageLoggers[packageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].errorEnabled, t, "packageLoggers[packageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].fatalEnabled, t, "packageLoggers[packageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyFormatterWithParameter(t *testing.T) {
	packageName := "testPackage"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, config.FORMATTER_DELIMITER)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+config.DELIMITER_PARAMETER, "_")
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.commonLogger, t, "commonLogger")
	testutil.AssertFalse(result.commonLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.commonLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.commonLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.commonLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.commonLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.commonLogger.appender, t, "commonLogger.appender")

	testutil.AssertTrue(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(1, len(result.packageLoggers), t, "len(result.packageLoggers)")
	testutil.AssertNotNil(result.packageLoggers[packageName], t, "packageLoggers[packageName]")
	testutil.AssertFalse(result.packageLoggers[packageName].debugEnabled, t, "packageLoggers[packageName].debugEnabled")
	testutil.AssertFalse(result.packageLoggers[packageName].informationEnabled, t, "packageLoggers[packageName].informationEnabled")
	testutil.AssertFalse(result.packageLoggers[packageName].warningEnabled, t, "packageLoggers[packageName].warningEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].errorEnabled, t, "packageLoggers[packageName].errorEnabled")
	testutil.AssertTrue(result.packageLoggers[packageName].fatalEnabled, t, "packageLoggers[packageName].fatalEnabled")
	testutil.AssertNotNil(result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender")

	testutil.AssertNotEquals(result.commonLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
}
