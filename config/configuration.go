// This package provides functions to derive configuration from environment variables or configuration files
package config

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
)

// root config element
type Config struct {
	Logger                      []LoggerConfig
	Appender                    []AppenderConfig
	Formatter                   []FormatterConfig
	UseFullQualifiedPackageName bool
}

// config of a single logger
type LoggerConfig struct {
	Id               string
	IsDefault        bool
	PackageParameter string
	PackageName      string
	Severity         int
	IsCallerToSet    bool
}

// config of an appender
type AppenderConfig struct {
	Id               string
	AppenderType     string
	IsDefault        bool
	PackageParameter string
	PathToLogFile    string
	CronExpression   string
	LimitByteSize    string
}

// config of a formatter
type FormatterConfig struct {
	Id                          string
	FormatterType               string
	IsDefault                   bool
	PackageParameter            string
	Delimiter                   string
	Template                    string
	CorrelationIdTemplate       string
	CustomTemplate              string
	CallerTemplate              string
	CallerCorrelationIdTemplate string
	CallerCustomTemplate        string
	TrimSeverityText            bool
	TimeKey                     string
	SeverityKey                 string
	MessageKey                  string
	CorrelationKey              string
	CustomValuesKey             string
	CustomValuesAsSubElement    bool
	CallerFunctionKey           string
	CallerFileKey               string
	CallerFileLineKey           string
	TimeLayout                  string
}

const (
	DEFAULT_LOG_LEVEL_PROPERTY_NAME                  = "TYPEWRITER_LOG_LEVEL"
	DEFAULT_LOG_APPENDER_PROPERTY_NAME               = "TYPEWRITER_LOG_APPENDER_TYPE"
	DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME          = "TYPEWRITER_LOG_APPENDER_FILE"
	DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME = "TYPEWRITER_LOG_APPENDER_CRON_RENAMING"
	DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME = "TYPEWRITER_LOG_APPENDER_SIZE_RENAMING"
	DEFAULT_LOG_FORMATTER_PROPERTY_NAME              = "TYPEWRITER_LOG_FORMATTER_TYPE"
	DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME    = "TYPEWRITER_LOG_FORMATTER_PARAMETER"

	PACKAGE_LOG_PACKAGE_PROPERTY_NAME                = "TYPEWRITER_PACKAGE_LOG_PACKAGE_"
	PACKAGE_LOG_LEVEL_PROPERTY_NAME                  = "TYPEWRITER_PACKAGE_LOG_LEVEL_"
	PACKAGE_LOG_APPENDER_PROPERTY_NAME               = "TYPEWRITER_PACKAGE_LOG_APPENDER_TYPE_"
	PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME          = "TYPEWRITER_PACKAGE_LOG_APPENDER_FILE_"
	PACKAGE_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME = "TYPEWRITER_PACKAGE_LOG_APPENDER_CRON_RENAMING_"
	PACKAGE_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME = "TYPEWRITER_PACKAGE_LOG_APPENDER_SIZE_RENAMING_"
	PACKAGE_LOG_FORMATTER_PROPERTY_NAME              = "TYPEWRITER_PACKAGE_LOG_FORMATTER_TYPE_"
	PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME    = "TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_"

	TIME_LAYOUT_PARAMETER                 = "_TIME_LAYOUT"
	DELIMITER_PARAMETER                   = "_DELIMITER"
	JSON_CALLER_FUNCTION_KEY_PARAMETER    = "_JSON_CALLER_FUNCTION_KEY"
	JSON_CALLER_FILE_KEY_PARAMETER        = "_JSON_CALLER_FILE_KEY"
	JSON_CALLER_LINE_KEY_PARAMETER        = "_JSON_CALLER_LINE_KEY"
	JSON_CORRELATION_KEY_PARAMETER        = "_JSON_CORRELATION_KEY"
	JSON_CUSTOM_VALUES_KEY_PARAMETER      = "_JSON_CUSTOM_VALUES_KEY"
	JSON_CUSTOM_VALUES_SUB_PARAMETER      = "_JSON_CUSTOM_VALUES_SUB"
	JSON_MESSAGE_KEY_PARAMETER            = "_JSON_MESSAGE_KEY"
	JSON_SEVERITY_KEY_PARAMETER           = "_JSON_SEVERITY_KEY"
	JSON_TIME_KEY_PARAMETER               = "_JSON_TIME_KEY"
	TEMPLATE_PARAMETER                    = "_TEMPLATE"
	TEMPLATE_CORRELATION_PARAMETER        = "_TEMPLATE_CORRELATION"
	TEMPLATE_CUSTOM_PARAMETER             = "_TEMPLATE_CUSTOM"
	TEMPLATE_TRIM_SEVERITY_PARAMETER      = "_TEMPLATE_TRIM_SEVERITY"
	TEMPLATE_CALLER_PARAMETER             = "_TEMPLATE_CALLER"
	TEMPLATE_CALLER_CORRELATION_PARAMETER = "_TEMPLATE_CALLER_CORRELATION"
	TEMPLATE_CALLER_CUSTOM_PARAMETER      = "_TEMPLATE_CALLER_CUSTOM"

	LOG_CONFIG_FILE_ENV_NAME                   = "TYPEWRITER_CONFIG_FILE"
	LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME       = "TYPEWRITER_LOG_CALLER"
	LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME = "TYPEWRITER_PACKAGE_FULL_QUALIFIED"

	LOG_LEVEL_DEBUG       = "DEBUG"
	LOG_LEVEL_INFO        = "INFO"
	LOG_LEVEL_INFORMATION = "INFORMATION"
	LOG_LEVEL_WARN        = "WARN"
	LOG_LEVEL_WARNING     = "WARNING"
	LOG_LEVEL_ERROR       = "ERROR"
	LOG_LEVEL_FATAL       = "FATAL"

	FORMATTER_DELIMITER = "DELIMITER"
	FORMATTER_TEMPLATE  = "TEMPLATE"
	FORMATTER_JSON      = "JSON"

	APPENDER_STDOUT = "STDOUT"
	APPENDER_FILE   = "FILE"

	DEFAULT_DELIMITER                   = " - "
	DEFAULT_TEMPLATE                    = format.DEFAULT_TEMPLATE
	DEFAULT_CORRELATION_TEMPLATE        = "[%s] %s %s: %s"
	DEFAULT_CUSTOM_TEMPLATE             = DEFAULT_TEMPLATE
	DEFAULT_CALLER_TEMPLATE             = format.DEFAULT_CALLER_TEMPLATE
	DEFAULT_CALLER_CORRELATION_TEMPLATE = "[%s] %s %s %s(%s.%d): %s"
	DEFAULT_CALLER_CUSTOM_TEMPLATE      = DEFAULT_CALLER_TEMPLATE
	DEFAULT_TRIM_SEVERITY_TEXT          = "false"
	DEFAULT_TIME_KEY                    = "time"
	DEFAULT_SEVERITY_KEY                = "severity"
	DEFAULT_MESSAGE_KEY                 = "message"
	DEFAULT_CORRELATION_KEY             = "correlation"
	DEFAULT_CUSTOM_VALUES_KEY           = "custom"
	DEFAULT_CUSTOM_AS_SUB_ELEMENT       = "false"
	DEFAULT_CALLER_FUNCTION_KEY         = "caller"
	DEFAULT_CALLER_FILE_KEY             = "file"
	DEFAULT_CALLER_FILE_LINE_KEY        = "line"
	DEFAULT_TIME_LAYOUT                 = time.RFC3339
)

var configInitialized = false
var config Config

// Mapping external severity levels to internal ids
var SeverityLevelMap = map[string]int{
	LOG_LEVEL_DEBUG:       common.DEBUG_SEVERITY,
	LOG_LEVEL_INFO:        common.INFORMATION_SEVERITY,
	LOG_LEVEL_INFORMATION: common.INFORMATION_SEVERITY,
	LOG_LEVEL_WARN:        common.WARNING_SEVERITY,
	LOG_LEVEL_WARNING:     common.WARNING_SEVERITY,
	LOG_LEVEL_ERROR:       common.ERROR_SEVERITY,
	LOG_LEVEL_FATAL:       common.FATAL_SEVERITY,
}

// prefixes to filter all given environment variables or properties from file
var relevantKeyPrefixes = []string{
	DEFAULT_LOG_LEVEL_PROPERTY_NAME,
	DEFAULT_LOG_APPENDER_PROPERTY_NAME,
	DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME,
	DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME,
	DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME,
	DEFAULT_LOG_FORMATTER_PROPERTY_NAME,
	DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME,
	PACKAGE_LOG_PACKAGE_PROPERTY_NAME,
	PACKAGE_LOG_LEVEL_PROPERTY_NAME,
	PACKAGE_LOG_APPENDER_PROPERTY_NAME,
	PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME,
	PACKAGE_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME,
	PACKAGE_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME,
	PACKAGE_LOG_FORMATTER_PROPERTY_NAME,
	PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME,
	LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME,
	LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME,
}

// returns the config or creates it if it was not initialized yet
func GetConfig() *Config {
	if configInitialized {
		return &config
	}

	config = Config{}
	var relevantKeyValues map[string]string

	if deriveConfigFromFile() {
		relevantKeyValues = determineRelevantPropertyFileKeyValues()
	} else if deriveConfigFromEnv() {
		relevantKeyValues = determineRelevantEnvKeyValues()
	} else {
		relevantKeyValues = map[string]string{}
	}
	if len(relevantKeyValues) > 0 {
		createFormatterConfig(&relevantKeyValues)
		createAppenderConfig(&relevantKeyValues)
		createLoggerConfig(&relevantKeyValues)
	}
	completeConfig(&relevantKeyValues)

	configInitialized = true
	return &config
}

// Resets initialization status
func ClearConfig() {
	configInitialized = false
}

// Checks whether two formatter config equals without regarding pointers to formatter or package
func FormatterConfigEquals(fc1 *FormatterConfig, fc2 *FormatterConfig) bool {
	return fc1.FormatterType == fc2.FormatterType &&
		fc1.Delimiter == fc2.Delimiter &&
		fc1.Template == fc2.Template && fc1.CorrelationIdTemplate == fc2.CorrelationIdTemplate && fc1.CustomTemplate == fc2.CustomTemplate &&
		fc1.TimeLayout == fc2.TimeLayout
}

// Checks whether two appender config equals without regarding pointers to appender or package
func AppenderConfigEquals(ac1 *AppenderConfig, ac2 *AppenderConfig) bool {
	return ac1.AppenderType == ac2.AppenderType && (ac1.AppenderType != APPENDER_FILE || ac1.PathToLogFile == ac2.PathToLogFile)
}

// Checks whether two logger config equals without regarding pointers to logger or package
func LoggerConfigEquals(lc1 *LoggerConfig, lc2 *LoggerConfig) bool {
	return lc1.Severity == lc2.Severity
}

// Checks whether the environment variable of the config is empty or not
func deriveConfigFromFile() bool {
	return existsAnyAtEnv(LOG_CONFIG_FILE_ENV_NAME)
}

// Checks whether any environment variable of a severity log level is set
func deriveConfigFromEnv() bool {
	return existsAnyAtEnv(DEFAULT_LOG_LEVEL_PROPERTY_NAME, DEFAULT_LOG_APPENDER_PROPERTY_NAME, DEFAULT_LOG_FORMATTER_PROPERTY_NAME, LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME) ||
		existsAnyPrefixAtEnv(PACKAGE_LOG_LEVEL_PROPERTY_NAME, PACKAGE_LOG_APPENDER_PROPERTY_NAME, PACKAGE_LOG_FORMATTER_PROPERTY_NAME)
}

// Checks if least one entry environment variables matches at least one entry of a given list of environment variables names
func existsAnyAtEnv(envNames ...string) bool {
	for _, envName := range envNames {
		if strings.TrimSpace(os.Getenv(envName)) != "" {
			return true
		}
	}
	return false
}

// Checks if least one entry environment variables matches as prefix at least one entry of a given list of environment variables names
func existsAnyPrefixAtEnv(envNames ...string) bool {
	for _, envEntry := range os.Environ() {
		keyValue := strings.SplitN(strings.ToUpper(envEntry), "=", 2)
		if len(keyValue) == 2 {
			for _, envName := range envNames {
				if strings.HasPrefix(keyValue[0], envName) {
					return true
				}
			}
		}
	}
	return false
}

// Determines a reduced map of key-value pairs from environment, which only contains relevant keys and non empty values
func determineRelevantEnvKeyValues() map[string]string {
	return createMapFromSliceWithKeyValues(os.Environ())
}

// Determines a reduced map of key-value pairs from property file, which only contains relevant keys and non empty values
func determineRelevantPropertyFileKeyValues() map[string]string {
	pathToFile := os.Getenv(LOG_CONFIG_FILE_ENV_NAME)
	propertiesFile, err := os.Open(pathToFile)
	if err != nil {
		fmt.Printf("Failed to load configuration file defined by environment variable %s with value \"%s\". Use default config instead", LOG_CONFIG_FILE_ENV_NAME, pathToFile)
		fmt.Println()
	}
	defer propertiesFile.Close()

	fileContent := []string{}
	scanner := bufio.NewScanner(propertiesFile)
	multilineCommentOpen := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if multilineCommentOpen && strings.Contains(line, "*/") {
			multilineCommentOpen = false
			continue
		}
		if multilineCommentOpen || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "--") {
			continue
		}
		if strings.HasPrefix(line, "/*") {
			multilineCommentOpen = !strings.HasSuffix(line, "*/")
			continue
		}
		fileContent = append(fileContent, line)
	}
	return createMapFromSliceWithKeyValues(fileContent)
}

// creates a map with the relevant key and values from a given slice
// The key will be compared after transforming in upper way and the keys at result map will be upper case also.
func createMapFromSliceWithKeyValues(sliceToConvert []string) map[string]string {
	result := make(map[string]string, len(sliceToConvert))

	for _, entry := range sliceToConvert {
		keyValue := strings.SplitN(entry, "=", 2)
		if len(keyValue) == 2 && strings.TrimSpace(keyValue[1]) != "" {
			key := strings.ToUpper(strings.TrimSpace(keyValue[0]))
			for _, keyPrefix := range relevantKeyPrefixes {
				if strings.HasPrefix(key, keyPrefix) {
					result[key] = strings.TrimSpace(keyValue[1])
					break
				}
			}
		}
	}

	return result
}

// creates all relevant appender config elements derived from relevant properties
func createAppenderConfig(relevantKeyValues *map[string]string) {
	config.Appender = append(config.Appender, AppenderConfig{IsDefault: true, PackageParameter: ""})
	appenderIndex := len(config.Appender) - 1
	configureAppender(relevantKeyValues, &config.Appender[appenderIndex], nil)

	for key := range *relevantKeyValues {
		packageOfAppender, found := strings.CutPrefix(key, PACKAGE_LOG_APPENDER_PROPERTY_NAME)
		if found {
			config.Appender = append(config.Appender, AppenderConfig{IsDefault: false, PackageParameter: packageOfAppender})
			appenderIndex++
			configureAppender(relevantKeyValues, &config.Appender[appenderIndex], &packageOfAppender)
		}
	}
}

// configures a given appender config element from properties concerning a given package name
func configureAppender(relevantKeyValues *map[string]string, appenderConfig *AppenderConfig, packageParameter *string) {
	var appenderKey string
	if packageParameter != nil && len(*packageParameter) > 0 {
		appenderKey = PACKAGE_LOG_APPENDER_PROPERTY_NAME + *packageParameter
	} else {
		appenderKey = DEFAULT_LOG_APPENDER_PROPERTY_NAME
	}

	appenderName := getValueFromMapOrDefault(relevantKeyValues, appenderKey, APPENDER_STDOUT)
	appenderConfig.AppenderType = appenderName
	switch appenderName {
	case APPENDER_STDOUT:
		// Nothing to do
	case APPENDER_FILE:
		setFileAppenderConfig(relevantKeyValues, appenderConfig, packageParameter)
	default:
		appenderConfig.AppenderType = ""
		printHint(appenderName, appenderKey)
	}
}

func setFileAppenderConfig(relevantKeyValues *map[string]string, appenderConfig *AppenderConfig, packageParameter *string) {
	var fileParameterKey string
	var cronParameterKey string
	var sizeParameterKey string

	if packageParameter != nil && len(*packageParameter) > 0 {
		fileParameterKey = PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME + *packageParameter
		cronParameterKey = PACKAGE_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME + *packageParameter
		sizeParameterKey = PACKAGE_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME + *packageParameter
	} else {
		fileParameterKey = DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME
		cronParameterKey = DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME
		sizeParameterKey = DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME
	}

	if fileValue, fileFound := (*relevantKeyValues)[fileParameterKey]; fileFound {
		appenderConfig.PathToLogFile = fileValue
		if cronValue, cronFound := (*relevantKeyValues)[cronParameterKey]; cronFound {
			appenderConfig.CronExpression = cronValue
		}
		if sizeValue, sizeFound := (*relevantKeyValues)[sizeParameterKey]; sizeFound {
			appenderConfig.LimitByteSize = sizeValue
		}
	} else {
		fmt.Printf("Cannot use file appender, because there is no value at %s. Use %s appender instead", fileParameterKey, APPENDER_STDOUT)
		fmt.Println()
		appenderConfig.AppenderType = APPENDER_STDOUT
	}

}

// creates all relevant formatter config elements derived from relevant properties
func createFormatterConfig(relevantKeyValues *map[string]string) {
	config.Formatter = append(config.Formatter, FormatterConfig{IsDefault: true, PackageParameter: ""})
	formatterIndex := len(config.Formatter) - 1
	configureFormatter(relevantKeyValues, &config.Formatter[formatterIndex], "")

	for key := range *relevantKeyValues {
		packageOfFormatter, found := strings.CutPrefix(key, PACKAGE_LOG_FORMATTER_PROPERTY_NAME)
		if found {
			config.Formatter = append(config.Formatter, FormatterConfig{IsDefault: false, PackageParameter: packageOfFormatter})
			formatterIndex++
			configureFormatter(relevantKeyValues, &config.Formatter[formatterIndex], packageOfFormatter)
		}
	}
}

// configures a given formatter config element from properties concerning a given package name
func configureFormatter(relevantKeyValues *map[string]string, formatterConfig *FormatterConfig, packageParameter string) {
	var formatterKey string
	var formatterParameterKey string
	if len(packageParameter) > 0 {
		formatterKey = PACKAGE_LOG_FORMATTER_PROPERTY_NAME + packageParameter
		formatterParameterKey = PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME + packageParameter
	} else {
		formatterKey = DEFAULT_LOG_FORMATTER_PROPERTY_NAME
		formatterParameterKey = DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME
	}

	formatterName := getValueFromMapOrDefault(relevantKeyValues, formatterKey, FORMATTER_DELIMITER)
	switch formatterName {
	case FORMATTER_DELIMITER:
		formatterConfig.FormatterType = formatterName
		formatterConfig.Delimiter = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+DELIMITER_PARAMETER, DEFAULT_DELIMITER)
		formatterConfig.TimeLayout = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TIME_LAYOUT_PARAMETER, DEFAULT_TIME_LAYOUT)
	case FORMATTER_TEMPLATE:
		formatterConfig.FormatterType = formatterName
		formatterConfig.Template = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_PARAMETER, DEFAULT_TEMPLATE)
		formatterConfig.CorrelationIdTemplate = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CORRELATION_PARAMETER, DEFAULT_CORRELATION_TEMPLATE)
		formatterConfig.CustomTemplate = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CUSTOM_PARAMETER, DEFAULT_CUSTOM_TEMPLATE)
		formatterConfig.TimeLayout = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TIME_LAYOUT_PARAMETER, DEFAULT_TIME_LAYOUT)
		formatterConfig.TrimSeverityText = strings.ToLower(getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_TRIM_SEVERITY_PARAMETER, DEFAULT_TRIM_SEVERITY_TEXT)) == "true"
		formatterConfig.CallerTemplate = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CALLER_PARAMETER, DEFAULT_CALLER_TEMPLATE)
		formatterConfig.CallerCorrelationIdTemplate = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CALLER_CORRELATION_PARAMETER, DEFAULT_CALLER_CORRELATION_TEMPLATE)
		formatterConfig.CallerCustomTemplate = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CALLER_CUSTOM_PARAMETER, DEFAULT_CALLER_CUSTOM_TEMPLATE)
	case FORMATTER_JSON:
		formatterConfig.FormatterType = formatterName
		formatterConfig.TimeKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_TIME_KEY_PARAMETER, DEFAULT_TIME_KEY)
		formatterConfig.SeverityKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_SEVERITY_KEY_PARAMETER, DEFAULT_SEVERITY_KEY)
		formatterConfig.CorrelationKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CORRELATION_KEY_PARAMETER, DEFAULT_CORRELATION_KEY)
		formatterConfig.MessageKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_MESSAGE_KEY_PARAMETER, DEFAULT_MESSAGE_KEY)
		formatterConfig.CustomValuesKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CUSTOM_VALUES_KEY_PARAMETER, DEFAULT_CUSTOM_VALUES_KEY)
		formatterConfig.CustomValuesAsSubElement = strings.ToLower(getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CUSTOM_VALUES_SUB_PARAMETER, DEFAULT_CUSTOM_AS_SUB_ELEMENT)) == "true"
		formatterConfig.TimeLayout = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TIME_LAYOUT_PARAMETER, DEFAULT_TIME_LAYOUT)
		formatterConfig.CallerFunctionKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CALLER_FUNCTION_KEY_PARAMETER, DEFAULT_CALLER_FUNCTION_KEY)
		formatterConfig.CallerFileKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CALLER_FILE_KEY_PARAMETER, DEFAULT_CALLER_FILE_KEY)
		formatterConfig.CallerFileLineKey = getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CALLER_LINE_KEY_PARAMETER, DEFAULT_CALLER_FILE_LINE_KEY)
	default:
		printHint(formatterName, formatterKey)
	}
}

// creates all relevant logger config elements derived from relevant properties
func createLoggerConfig(relevantKeyValues *map[string]string) {
	config.Logger = append(config.Logger, LoggerConfig{IsDefault: true, PackageParameter: ""})
	loggerIndex := len(config.Logger) - 1
	configureLogger(relevantKeyValues, &config.Logger[loggerIndex], "")

	for key := range *relevantKeyValues {
		packageParameter, found := strings.CutPrefix(key, PACKAGE_LOG_LEVEL_PROPERTY_NAME)
		if found {
			packageName := getValueFromMapOrDefault(relevantKeyValues, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageParameter, strings.ToLower(packageParameter))
			config.Logger = append(config.Logger, LoggerConfig{IsDefault: false, PackageParameter: packageParameter, PackageName: packageName})
			loggerIndex++
			configureLogger(relevantKeyValues, &config.Logger[loggerIndex], packageParameter)
		}
	}

	if fullQualifiedText, found := (*relevantKeyValues)[LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME]; found {
		config.UseFullQualifiedPackageName = strings.ToUpper(fullQualifiedText) == "TRUE"
	}
}

// configures a given logger config element from properties concerning a given package name
func configureLogger(relevantKeyValues *map[string]string, loggerConfig *LoggerConfig, packageParameter string) {
	var formatterKey string
	if len(packageParameter) > 0 {
		formatterKey = PACKAGE_LOG_LEVEL_PROPERTY_NAME + packageParameter
	} else {
		formatterKey = DEFAULT_LOG_LEVEL_PROPERTY_NAME
	}

	logLevel := getValueFromMapOrDefault(relevantKeyValues, formatterKey, LOG_LEVEL_ERROR)
	severity, found := SeverityLevelMap[logLevel]
	if !found {
		severity = common.ERROR_SEVERITY
	}
	loggerConfig.Severity = severity

	loggerConfig.IsCallerToSet = strings.ToLower(getValueFromMapOrDefault(relevantKeyValues, LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME, "false")) == "true"
}

// Returns the value from a map for a given key. If there is none, the default will be returned
func getValueFromMapOrDefault(source *map[string]string, key string, defaultValue string) string {
	value, found := (*source)[key]
	if found {
		return value
	}
	return defaultValue
}

func printHint(propertyName string, objectType string) {
	fmt.Println("Unknown \"", objectType, "\" value at property", propertyName)
}

// creates default configs if missing and adds package specific copies of defaults if at least one of the other config types exists as package variant
func completeConfig(relevantKeyValues *map[string]string) {
	completeDefaults()

	completeAppenderConfigPackageForward()
	completeFormatterConfigPackageForward()

	completeAppenderConfigPackageBackward()
	completeLoggerConfigPackageBackward(relevantKeyValues)

	sortConfig()
	determineIds()
}

// creates default configs if missing
func completeDefaults() {
	found := false

	for _, fc := range config.Formatter {
		if fc.IsDefault {
			found = true
			break
		}
	}
	if !found {
		config.Formatter = append(config.Formatter, FormatterConfig{FormatterType: FORMATTER_DELIMITER, IsDefault: true, PackageParameter: "", Delimiter: DEFAULT_DELIMITER, TimeLayout: DEFAULT_TIME_LAYOUT})
	}

	for _, ac := range config.Appender {
		if ac.IsDefault {
			found = true
			break
		}
	}
	if !found {
		config.Appender = append(config.Appender, AppenderConfig{AppenderType: APPENDER_STDOUT, IsDefault: true, PackageParameter: ""})
	}

	for _, lc := range config.Logger {
		if lc.IsDefault {
			found = true
			break
		}
	}
	if !found {
		config.Logger = append(config.Logger, LoggerConfig{IsDefault: true, PackageParameter: "", Severity: common.ERROR_SEVERITY})
	}
}

// creates appender configs if there exists a logger package variant
func completeAppenderConfigPackageForward() {
	for _, lc := range config.Logger {
		if lc.IsDefault {
			continue
		}
		createAppenderConfigIfNecessary(&lc.PackageParameter)
	}
}

// creates appender configs if there exists a formatter package variant
func completeAppenderConfigPackageBackward() {
	for _, fc := range config.Formatter {
		if fc.IsDefault {
			continue
		}
		createAppenderConfigIfNecessary(&fc.PackageParameter)
	}
}

// creates an appender config if it does not exists for a given package name
func createAppenderConfigIfNecessary(packageParameter *string) {
	found := false
	for _, ac := range config.Appender {
		if ac.PackageParameter == *packageParameter {
			found = true
			break
		}
	}
	if !found {
		for _, ac := range config.Appender {
			if ac.IsDefault {
				acp := ac
				acp.IsDefault = false
				acp.PackageParameter = *packageParameter
				config.Appender = append(config.Appender, acp)
				break
			}
		}
	}
}

// creates formatter configs if there exists a appender package variant
func completeFormatterConfigPackageForward() {
	for _, ac := range config.Appender {
		if ac.IsDefault {
			continue
		}
		createFormatterConfigIfNecessary(&ac.PackageParameter)
	}
}

// creates an formatter config if it does not exists for a given package name
func createFormatterConfigIfNecessary(packageParameter *string) {
	found := false
	for _, fc := range config.Formatter {
		if fc.PackageParameter == *packageParameter {
			found = true
			break
		}
	}
	if !found {
		for _, fc := range config.Formatter {
			if fc.IsDefault {
				fcp := fc
				fcp.IsDefault = false
				fcp.PackageParameter = *packageParameter
				config.Formatter = append(config.Formatter, fcp)
				break
			}
		}
	}
}

// creates logger configs if there exists a appender package variant
func completeLoggerConfigPackageBackward(relevantKeyValues *map[string]string) {
	for _, ac := range config.Appender {
		if ac.IsDefault {
			continue
		}
		createLoggerConfigIfNecessary(relevantKeyValues, &ac.PackageParameter)
	}
}

// creates an logger config if it does not exists for a given package name
func createLoggerConfigIfNecessary(relevantKeyValues *map[string]string, packageParameter *string) {
	found := false
	for _, lc := range config.Logger {
		if lc.PackageParameter == *packageParameter {
			found = true
			break
		}
	}
	if !found {
		for _, lc := range config.Logger {
			if lc.IsDefault {
				lcp := lc
				lcp.IsDefault = false
				lcp.PackageParameter = *packageParameter
				lcp.PackageName = getValueFromMapOrDefault(relevantKeyValues, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+*packageParameter, strings.ToLower(*packageParameter))
				config.Logger = append(config.Logger, lcp)
				break
			}
		}
	}
}

// Sorts the config put the default configs at first index. Because of this a potential cronRenamer of the default appender is used for all equal named files
func sortConfig() {
	sort.Slice(config.Formatter, func(i, j int) bool {
		return (config.Formatter[i].IsDefault && !config.Formatter[j].IsDefault) || (config.Formatter[i].IsDefault == config.Formatter[j].IsDefault && config.Formatter[i].PackageParameter < config.Formatter[j].PackageParameter)
	})
	sort.Slice(config.Appender, func(i, j int) bool {
		return (config.Appender[i].IsDefault && !config.Appender[j].IsDefault) || (config.Appender[i].IsDefault == config.Appender[j].IsDefault && config.Appender[i].PackageParameter < config.Appender[j].PackageParameter)
	})
	sort.Slice(config.Logger, func(i, j int) bool {
		return (config.Logger[i].IsDefault && !config.Logger[j].IsDefault) || (config.Logger[i].IsDefault == config.Logger[j].IsDefault && config.Logger[i].PackageParameter < config.Logger[j].PackageParameter)
	})
}

func determineIds() {
	determineFormatterIds()
	determineAppenderIds()
	determineLoggerIds()
}

func determineFormatterIds() {
	for i := 0; i < len(config.Formatter); i++ {
		if config.Formatter[i].Id != "" {
			continue
		}
		config.Formatter[i].Id = fmt.Sprint("formatter", i)
		for j := i + 1; j < len(config.Formatter); j++ {
			if FormatterConfigEquals(&config.Formatter[i], &config.Formatter[j]) {
				config.Formatter[j].Id = config.Formatter[i].Id
			}
		}
	}
}

func determineAppenderIds() {
	for i := 0; i < len(config.Appender); i++ {
		if config.Appender[i].Id != "" {
			continue
		}
		config.Appender[i].Id = fmt.Sprint("appender", i)
		for j := i + 1; j < len(config.Appender); j++ {
			if AppenderConfigEquals(&config.Appender[i], &config.Appender[j]) {
				config.Appender[j].Id = config.Appender[i].Id
			}
		}
	}
}

func determineLoggerIds() {
	for i := 0; i < len(config.Logger); i++ {
		if config.Logger[i].Id != "" {
			continue
		}
		config.Logger[i].Id = fmt.Sprint("logger", i)
		for j := i + 1; j < len(config.Logger); j++ {
			if LoggerConfigEquals(&config.Logger[i], &config.Logger[j]) {
				config.Logger[j].Id = config.Logger[i].Id
			}
		}
	}
}
