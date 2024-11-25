package logger

import (
	"os"
	"reflect"
	"strings"
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
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaulTemplate(t *testing.T) {
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_STDOUT)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, config.FORMATTER_TEMPLATE)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_1", "time: %s severity: %s message: %s")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_2", "time: %s severity: %s correlation: %s message: %s")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_3", "time: %s severity: %s message: %s %s: %s %s: %d %s: %t")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_4", time.RFC1123Z)
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
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaulJson(t *testing.T) {
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_STDOUT)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PROPERTY_NAME, config.FORMATTER_JSON)
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_1", "timing")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_2", "level")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_3", "cor")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_4", "msg")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_5", "customValues")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_6", "true")
	os.Setenv(config.DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+"_7", time.RFC1123Z)
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
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

func TestGetLoggersCreateFromEnvDefaultFileAppender(t *testing.T) {
	appender.SkipFileCreationForTest = true
	logFilePath := "pathToLogFile"
	os.Clearenv()
	os.Setenv(config.DEFAULT_LOG_LEVEL_PROPERTY_NAME, config.LOG_LEVEL_INFO)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PROPERTY_NAME, config.APPENDER_FILE)
	os.Setenv(config.DEFAULT_LOG_APPENDER_PARAMETER_PROPERTY_NAME, logFilePath)
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
	testutil.AssertEquals(0, len(result.packageLoggers), t, "len(result.packageLoggers)")
}

//
// Get Loggers with packages
//

func TestGetLoggersCreateFromEnvPackage(t *testing.T) {
	packageName := "testPackage"
	packageNameUpper := strings.ToUpper(packageName)
	os.Clearenv()
	os.Setenv(config.PACKAGE_LOG_LEVEL_PROPERTY_NAME+packageName, config.LOG_LEVEL_DEBUG)
	os.Setenv(config.PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, config.APPENDER_STDOUT)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, config.FORMATTER_DELIMITER)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName, "_")
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
	os.Setenv(config.PACKAGE_LOG_APPENDER_PROPERTY_NAME+packageName, config.APPENDER_FILE)
	os.Setenv(config.PACKAGE_LOG_APPENDER_PARAMETER_PROPERTY_NAME+packageName, logFilePath)
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
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PROPERTY_NAME+packageName, config.FORMATTER_DELIMITER)
	os.Setenv(config.PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME+packageName, "_")
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
