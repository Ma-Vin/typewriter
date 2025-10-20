package logger

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "generalLogger.appender.Name")

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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "generalLogger.appender.Name")

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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "generalLogger.appender.Name")

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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")
	testutil.AssertEquals("FileAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "generalLogger.appender.Name")

	testutil.AssertFalse(result.existPackageLogger, t, "existPackageLogger")
	testutil.AssertFalse(result.useFullQualifiedPackage, t, "useFullQualifiedPackage")
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaultMultiAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	logFilePath := "pathToLogFile"
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_FILE+","+config.APPENDER_STDOUT)
	os.Setenv(config.DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME, logFilePath)
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")
	testutil.AssertEquals("MultiAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "generalLogger.appender.Name")
	testutil.AssertTrue((*result.generalLogger.appender).(appender.MultiAppender).CheckSubAppenderTypesForTest([]string{"FileAppender", "StandardOutputAppender"}), t, "CheckSubAppenderTypesForTest")

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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")

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

	testutil.AssertNotEquals(result.generalLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")

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

	testutil.AssertNotEquals(result.generalLogger.appender, result.packageLoggers[fullQualifiedPackageName].appender, t, "packageLoggers[packageName].appender.")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyLevel(t *testing.T) {
	packageName := "testPackage"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, config.LOG_LEVEL_DEBUG)
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")

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

	testutil.AssertEquals(result.generalLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
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

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")

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

	testutil.AssertNotEquals(result.generalLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
}
func TestGetLoggersCreateFromEnvPackagePartialOnlyMultiAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	packageName := "testPackage"
	logFilePath := "pathToLogFile"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, config.APPENDER_FILE+","+config.APPENDER_STDOUT)
	os.Setenv(config.PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME+packageName, logFilePath)
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "generalLogger.appender.Name")

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
	testutil.AssertEquals("MultiAppender", reflect.TypeOf(*result.packageLoggers[packageName].appender).Name(), t, "generalLogger.appender.Name")
	testutil.AssertTrue((*result.packageLoggers[packageName].appender).(appender.MultiAppender).CheckSubAppenderTypesForTest([]string{"FileAppender", "StandardOutputAppender"}), t, "CheckSubAppenderTypesForTest")
}

func TestGetLoggersCreateFromEnvPackagePartialOnlyFormatterWithParameter(t *testing.T) {
	packageName := "testPackage"
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageName, packageName)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, config.FORMATTER_DELIMITER)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName+config.DELIMITER_PARAMETER, "_")
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
	testutil.AssertNotNil(result.generalLogger.appender, t, "generalLogger.appender")

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

	testutil.AssertNotEquals(result.generalLogger.appender, result.packageLoggers[packageName].appender, t, "packageLoggers[packageName].appender.")
}

func TestGetLoggersClearConfig(t *testing.T) {
	os.Clearenv()
	Reset()

	result := getLoggers()

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertFalse(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertFalse(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")

	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	config.ClearConfig()

	result = getLoggers()

	testutil.AssertNotNil(result.generalLogger, t, "generalLogger")
	testutil.AssertFalse(result.generalLogger.debugEnabled, t, "debugEnabled")
	testutil.AssertTrue(result.generalLogger.informationEnabled, t, "informationEnabled")
	testutil.AssertTrue(result.generalLogger.warningEnabled, t, "warningEnabled")
	testutil.AssertTrue(result.generalLogger.errorEnabled, t, "errorEnabled")
	testutil.AssertTrue(result.generalLogger.fatalEnabled, t, "fatalEnabled")
}

func createFileAppenderConfigForTest(relevantKeyValues *map[string]string, commonAppenderConfig *config.CommonAppenderConfig) (*config.AppenderConfig, error) {
	var result config.AppenderConfig = config.FileAppenderConfig{
		Common:         commonAppenderConfig,
		PathToLogFile:  "pathToLogFile",
		CronExpression: "",
		LimitByteSize:  "",
	}
	return &result, nil
}

func TestRegisterAndDeregisterAppender(t *testing.T) {
	customAppenderName := "CUSTOM_APPENDER"

	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, customAppenderName)

	Reset()

	err := config.RegisterAppenderConfig(customAppenderName, []string{}, createFileAppenderConfigForTest)
	testutil.AssertNil(err, t, "err of RegisterAppenderConfig")

	// Check fallback to default one if CUSTOM_APPENDER is not registered yet
	result := getLoggers()
	testutil.AssertNotNil(result, t, "before register - result")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "before register - generalLogger.appender.Name")

	// Register custom one
	err = RegisterAppender(customAppenderName, appender.CreateFileAppenderFromConfig)
	testutil.AssertNil(err, t, "err of RegisterAppender")
	result = getLoggers()
	testutil.AssertNotNil(result, t, "after register - result")
	testutil.AssertEquals("FileAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "after register - generalLogger.appender.Name")

	// Deregister custom one
	err = DeregisterAppender(customAppenderName)
	testutil.AssertNil(err, t, "err of DeregisterAppender")

	// Load config without registered appender: fallback to default one
	result = getLoggers()
	testutil.AssertNotNil(result, t, "deregister - result")
	testutil.AssertEquals("StandardOutputAppender", reflect.TypeOf(*result.generalLogger.appender).Name(), t, "deregister - generalLogger.appender.Name")
}

func TestRegisterKnownAppender(t *testing.T) {
	Reset()
	err := RegisterAppender(config.APPENDER_FILE, appender.CreateFileAppenderFromConfig)
	testutil.AssertNotNil(err, t, "registered known appender")
}

func TestDeregisterBuildInAppender(t *testing.T) {
	Reset()
	err := DeregisterAppender(config.APPENDER_STDOUT)
	testutil.AssertNotNil(err, t, "deregistered standard output appender")

	err = DeregisterAppender(config.APPENDER_FILE)
	testutil.AssertNotNil(err, t, "deregistered file appender")
}

func TestDeregisterUnknownAppender(t *testing.T) {
	Reset()
	err := DeregisterAppender("Anything")
	testutil.AssertNotNil(err, t, "deregistered unknown appender")
}

func createJsonConfigForTest(relevantKeyValues *map[string]string, commonFormatterConfig *config.CommonFormatterConfig) (*config.FormatterConfig, error) {
	var result config.FormatterConfig = config.JsonFormatterConfig{
		Common:                   commonFormatterConfig,
		TimeKey:                  config.DEFAULT_TIME_KEY,
		SeverityKey:              config.DEFAULT_SEVERITY_KEY,
		CorrelationKey:           config.DEFAULT_CORRELATION_KEY,
		MessageKey:               config.DEFAULT_MESSAGE_KEY,
		CustomValuesKey:          config.DEFAULT_CUSTOM_VALUES_KEY,
		CustomValuesAsSubElement: config.DEFAULT_CUSTOM_AS_SUB_ELEMENT,
		CallerFunctionKey:        config.DEFAULT_CALLER_FUNCTION_KEY,
		CallerFileKey:            config.DEFAULT_CALLER_FILE_KEY,
		CallerFileLineKey:        config.DEFAULT_CALLER_FILE_LINE_KEY,
	}
	return &result, nil
}

func TestRegisterAndDeregisterFormatter(t *testing.T) {

	customAppenderName := "CUSTOM_APPENDER"
	customFormatterName := "CUSTOM_FORMATTER"

	appender.SkipFileCreationForTest = true
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, customAppenderName)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, customFormatterName)

	Reset()

	err := config.RegisterAppenderConfig(customAppenderName, []string{}, createFileAppenderConfigForTest)
	testutil.AssertNil(err, t, "err of RegisterAppenderConfig")
	err = RegisterAppender(customAppenderName, createDummyAppender)
	testutil.AssertNil(err, t, "err of RegisterAppender")
	err = config.RegisterFormatterConfig(customFormatterName, []string{}, createJsonConfigForTest)
	testutil.AssertNil(err, t, "err of RegisterFormatterConfig")

	// Check fallback to default one if CUSTOM_FORMATTER is not registered yet
	result := getLoggers()
	testutil.AssertNotNil(result, t, "before register - result")

	testutil.AssertEquals("DelimiterFormatter", reflect.TypeOf(*(*result.generalLogger.appender).(DummyAppender).formatter).Name(), t, "before register - generalLogger.appender.formatter.Name")

	// Register custom one
	RegisterFormatter(customFormatterName, format.CreateJsonFormatterFromConfig)
	result = getLoggers()
	testutil.AssertNotNil(result, t, "after register - result")
	testutil.AssertEquals("JsonFormatter", reflect.TypeOf(*(*result.generalLogger.appender).(DummyAppender).formatter).Name(), t, "after register - generalLogger.appender.formatter.Name")

	// Deregister custom one
	err = DeregisterFormatter(customFormatterName)
	testutil.AssertNil(err, t, "err of DeregisterFormatter")

	// Load config without registered formatter: fallback to default one
	result = getLoggers()
	testutil.AssertNotNil(result, t, "deregister - result")
	testutil.AssertEquals("DelimiterFormatter", reflect.TypeOf(*(*result.generalLogger.appender).(DummyAppender).formatter).Name(), t, "deregister - generalLogger.appender.formatter.Name")
}

func TestRegisterKnownFormatter(t *testing.T) {
	Reset()
	err := RegisterFormatter(config.FORMATTER_JSON, format.CreateJsonFormatterFromConfig)
	testutil.AssertNotNil(err, t, "registered known formatter")
}

func TestDeregisterBuildInFormatter(t *testing.T) {
	Reset()
	err := DeregisterFormatter(config.FORMATTER_DELIMITER)
	testutil.AssertNotNil(err, t, "deregistered delimiter formatter")

	err = DeregisterFormatter(config.FORMATTER_TEMPLATE)
	testutil.AssertNotNil(err, t, "deregistered template formatter")

	err = DeregisterFormatter(config.FORMATTER_JSON)
	testutil.AssertNotNil(err, t, "deregistered json formatter")
}

func TestDeregisterUnknownFormatter(t *testing.T) {
	Reset()
	err := DeregisterFormatter("Anything")
	testutil.AssertNotNil(err, t, "deregistered unknown formatter")
}
