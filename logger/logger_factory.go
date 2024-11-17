package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/format"
)

const (
	DEFAULT_LOG_LEVEL_ENV_NAME    = "TYPEWRITER_LOG_LEVEL"
	DEFAULT_LOG_APPENDER_ENV_NAME = "TYPEWRITER_LOG_APPENDER_TYPE"
	//DEFAULT_LOG_APPENDER_PARAMETER_ENV_NAME  = "TYPEWRITER_LOG_APPENDER_PARAMETER"
	DEFAULT_LOG_FORMATTER_ENV_NAME           = "TYPEWRITER_LOG_FORMATTER_TYPE"
	DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME = "TYPEWRITER_LOG_FORMATTER_PARAMETER"

	PACKAGE_LOG_LEVEL_ENV_NAME    = "TYPEWRITER_PACKAGE_LOG_LEVEL_"
	PACKAGE_LOG_APPENDER_ENV_NAME = "TYPEWRITER_PACKAGE_LOG_APPENDER_TYPE_"
	//PACKAGE_LOG_APPENDER_PARAMETER_ENV_NAME  = "TYPEWRITER_PACKAGE_LOG_APPENDER_PARAMETER_"
	PACKAGE_LOG_FORMATTER_ENV_NAME           = "TYPEWRITER_PACKAGE_LOG_FORMATTER_TYPE_"
	PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME = "TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_"

	DEFAULT_LOG_CONFIG_FILE_ENV_NAME = "TYPEWRITER_CONFIG_FILE"

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

	DEFAULT_DELIMITER             = " - "
	DEFAULT_TEMPLATE              = "[%s] %s: %s"
	DEFAULT_CORRELATION_TEMPLATE  = "[%s] %s %s: %s"
	DEFAULT_CUSTOM_TEMPLATE       = DEFAULT_TEMPLATE
	DEFAULT_TIME_KEY              = "time"
	DEFAULT_SEVERITY_KEY          = "severity"
	DEFAULT_MESSAGE_KEY           = "message"
	DEFAULT_CORRELATION_KEY       = "correlation"
	DEFAULT_CUSTOM_VALUES_KEY     = "custom"
	DEFAULT_CUSTOM_AS_SUB_ELEMENT = "false"
	DEFAULT_TIME_LAYOUT           = time.RFC3339
)

// root config element
type overallConfig struct {
	logger    []loggerConfig
	appender  []appenderConfig
	formatter []formatterConfig
}

// config of a single logger
type loggerConfig struct {
	isDefault   bool
	packageName string
	severity    int
	logger      *CommonLogger
}

// config of an appender
type appenderConfig struct {
	appenderType string
	isDefault    bool
	packageName  string
	appender     *appender.Appender
}

// config of a formatter
type formatterConfig struct {
	formatterType            string
	isDefault                bool
	packageName              string
	delimiter                string
	template                 string
	correlationIdTemplate    string
	customTemplate           string
	timeKey                  string
	severityKey              string
	messageKey               string
	correlationKey           string
	customValuesKey          string
	customValuesAsSubElement bool
	timeLayout               string
	formatter                *format.Formatter
}

var configInitialized = false
var config overallConfig

var severityLevelMap = map[string]int{
	LOG_LEVEL_DEBUG:       constants.DEBUG_SEVERITY,
	LOG_LEVEL_INFO:        constants.INFORMATION_SEVERITY,
	LOG_LEVEL_INFORMATION: constants.INFORMATION_SEVERITY,
	LOG_LEVEL_WARN:        constants.WARNING_SEVERITY,
	LOG_LEVEL_WARNING:     constants.WARNING_SEVERITY,
	LOG_LEVEL_ERROR:       constants.ERROR_SEVERITY,
	LOG_LEVEL_FATAL:       constants.FATAL_SEVERITY,
}

var loggersInitialized = false
var mLogger MainLogger
var cLoggers []CommonLogger
var appenders []appender.Appender
var formatters []format.Formatter

// returns the main logger. If not initialized, a new one will be created from config
func getLoggers() *MainLogger {
	if loggersInitialized {
		return &mLogger
	}

	getConfig()
	if !configInitialized {
		return nil
	}

	createFormatters()
	createAppenders()
	createCommonLoggers()
	createMainLogger()

	return &mLogger
}

// Creates the all relevant fromatters from config elements
func createFormatters() {
	for i, fc1 := range config.formatter {
		alreadyCreated := false

		for _, fc2 := range config.formatter {
			if fc2.formatter != nil && formatterConfigEquals(&fc1, &fc2) {
				config.formatter[i].formatter = fc2.formatter
				alreadyCreated = true
				break
			}
		}

		if alreadyCreated {
			continue
		}

		switch fc1.formatterType {
		case FORMATTER_DELIMITER:
			appendFormatter(format.CreateDelimiterFormatter(fc1.delimiter))
			setLastFormatter(i)
		case FORMATTER_TEMPLATE:
			appendFormatter(format.CreateTemplateFormatter(fc1.template, fc1.correlationIdTemplate, fc1.customTemplate, fc1.timeLayout))
			setLastFormatter(i)
		case FORMATTER_JSON:
			appendFormatter(format.CreateJsonFormatter(fc1.timeKey, fc1.severityKey, fc1.messageKey, fc1.correlationKey, fc1.customValuesKey, fc1.timeLayout, fc1.customValuesAsSubElement))
			setLastFormatter(i)
		default:
			// not relevant: handled at config load
		}
	}
}

func appendFormatter(formatter format.Formatter) {
	formatters = append(formatters, formatter)
}

func setLastFormatter(index int) {
	config.formatter[index].formatter = &formatters[len(formatters)-1]
}

// Creates the all relevant appenders from config elements
func createAppenders() {
	for i, ac1 := range config.appender {
		alreadyCreated := false

		for _, ac2 := range config.appender {
			if ac2.appender != nil && appenderConfigEquals(&ac2, &ac1) &&
				formatterConfigEquals(getFormatterConfigForPackage(&ac2.packageName), getFormatterConfigForPackage(&ac1.packageName)) {

				config.appender[i].appender = ac2.appender
				alreadyCreated = true
				break
			}
		}

		if alreadyCreated {
			continue
		}

		switch ac1.appenderType {
		case APPENDER_STDOUT:
			appenders = append(appenders, appender.CreateStandardOutputAppender(getFormatterConfigForPackage(&ac1.packageName).formatter))
			config.appender[i].appender = &appenders[len(appenders)-1]
		case APPENDER_FILE:
			// not supported yet
		default:
			// not relevant: handled at config load
		}
	}
}

// Creates the all relevant common logger from config elements
func createCommonLoggers() {
	for i, lc1 := range config.logger {
		alreadyCreated := false

		for _, lc2 := range config.logger {
			if lc2.logger != nil && loggerConfigEquals(&lc2, &lc1) &&
				appenderConfigEquals(getAppenderConfigForPackage(&lc2.packageName), getAppenderConfigForPackage(&lc1.packageName)) &&
				formatterConfigEquals(getFormatterConfigForPackage(&lc2.packageName), getFormatterConfigForPackage(&lc1.packageName)) {

				config.logger[i].logger = lc2.logger
				alreadyCreated = true
				break

			}
		}

		if alreadyCreated {
			continue
		}

		cLoggers = append(cLoggers, CreateCommonLogger(getAppenderConfigForPackage(&lc1.packageName).appender, lc1.severity))
		config.logger[i].logger = &cLoggers[len(cLoggers)-1]
	}
}

// Creates the main logger from config elements
func createMainLogger() {
	mLogger = MainLogger{}
	mLogger.existPackageLogger = len(config.logger) > 1
	mLogger.packageLoggers = make(map[string]*CommonLogger, len(config.logger)-1)

	for _, lc := range config.logger {
		if lc.isDefault {
			mLogger.commonLogger = lc.logger
		} else {
			mLogger.packageLoggers[lc.packageName] = lc.logger
		}
	}

	loggersInitialized = true
}

// Checks whether two formatter config equals without regarding pointers to formatter or package
func formatterConfigEquals(fc1 *formatterConfig, fc2 *formatterConfig) bool {
	return fc1.formatterType == fc2.formatterType &&
		fc1.delimiter == fc2.delimiter &&
		fc1.template == fc2.template && fc1.correlationIdTemplate == fc2.correlationIdTemplate && fc1.customTemplate == fc2.customTemplate &&
		fc1.timeLayout == fc2.timeLayout
}

// Checks whether two appender config equals without regarding pointers to appender or package
func appenderConfigEquals(ac1 *appenderConfig, ac2 *appenderConfig) bool {
	return ac1.appenderType == ac2.appenderType
}

// Checks whether two logger config equals without regarding pointers to logger or package
func loggerConfigEquals(lc1 *loggerConfig, lc2 *loggerConfig) bool {
	return lc1.severity == lc2.severity
}

// returns a pointer to the formatter config for a given package
func getFormatterConfigForPackage(packageName *string) *formatterConfig {
	for i, fc := range config.formatter {
		if fc.packageName == *packageName {
			return &config.formatter[i]
		}
	}
	return nil
}

// returns a pointer to the appender config for a given package
func getAppenderConfigForPackage(packageName *string) *appenderConfig {
	for i, ac := range config.appender {
		if ac.packageName == *packageName {
			return &config.appender[i]
		}
	}
	return nil
}

// Checks whether the environment variable of the config is empty or not
func deriveConfigFromFile() bool {
	return existsAnyAtEnv(DEFAULT_LOG_CONFIG_FILE_ENV_NAME)
}

// Checks whether any environment variable of a severity log level is set
func deriveConfigFromEnv() bool {
	return existsAnyAtEnv(DEFAULT_LOG_LEVEL_ENV_NAME, DEFAULT_LOG_APPENDER_ENV_NAME, DEFAULT_LOG_FORMATTER_ENV_NAME) ||
		existsAnyPrefixAtEnv(PACKAGE_LOG_LEVEL_ENV_NAME, PACKAGE_LOG_APPENDER_ENV_NAME, PACKAGE_LOG_FORMATTER_ENV_NAME)
}

// Checks if least one entry environment variables matchs at least one entry of a given list of environment variables names
func existsAnyAtEnv(envNames ...string) bool {
	for _, envName := range envNames {
		if strings.TrimSpace(os.Getenv(envName)) != "" {
			return true
		}
	}
	return false
}

// Checks if least one entry environment variables matchs as prefix at least one entry of a given list of environment variables names
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

// returns the config or creates it if it was not initialized yet
func getConfig() *overallConfig {
	if configInitialized {
		return &config
	}

	config = overallConfig{}
	if deriveConfigFromFile() {
		// not supported yet
		fmt.Println("Logger configuration is not yet supported:", DEFAULT_LOG_CONFIG_FILE_ENV_NAME, os.Getenv(DEFAULT_LOG_CONFIG_FILE_ENV_NAME))
		return nil
	} else if deriveConfigFromEnv() {
		relevantEnvKeyValues := determineRelevantEnvKeyValues()
		createFormatterConfigFromEnv(&relevantEnvKeyValues)
		createAppenderConfigFromEnv(&relevantEnvKeyValues)
		createLoggerConfigFromEnv(&relevantEnvKeyValues)
	}

	completeConfig()

	configInitialized = true
	return &config
}

// Determines a reduced slice of key-value pairs, which only contains relevant keys and non empty values
func determineRelevantEnvKeyValues() map[string]string {
	result := make(map[string]string, len(os.Environ()))

	relevantKeyPrefixes := []string{DEFAULT_LOG_LEVEL_ENV_NAME, DEFAULT_LOG_APPENDER_ENV_NAME, DEFAULT_LOG_FORMATTER_ENV_NAME,
		DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME, PACKAGE_LOG_LEVEL_ENV_NAME, PACKAGE_LOG_APPENDER_ENV_NAME,
		PACKAGE_LOG_FORMATTER_ENV_NAME, PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME}

	for _, envEntry := range os.Environ() {
		keyValue := strings.SplitN(envEntry, "=", 2)
		if len(keyValue) == 2 && strings.TrimSpace(keyValue[1]) != "" {
			key := strings.ToUpper(keyValue[0])
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

// creates all relevant appender config elements derived from relevant environment variables
func createAppenderConfigFromEnv(relevantEnvKeyValues *map[string]string) {
	config.appender = append(config.appender, appenderConfig{isDefault: true, packageName: ""})
	appenderIndex := len(config.appender) - 1
	configureAppenderFromEnv(relevantEnvKeyValues, &config.appender[appenderIndex], "")

	for key, _ := range *relevantEnvKeyValues {
		packageOfAppender, found := strings.CutPrefix(key, PACKAGE_LOG_APPENDER_ENV_NAME)
		if found {
			config.appender = append(config.appender, appenderConfig{isDefault: false, packageName: packageOfAppender})
			appenderIndex++
			configureAppenderFromEnv(relevantEnvKeyValues, &config.appender[appenderIndex], packageOfAppender)
		}
	}
}

// configures a given appender config element from environment variables concerning a given package name
func configureAppenderFromEnv(relevantEnvKeyValues *map[string]string, appenderConfig *appenderConfig, packageName string) {
	var appenderEnvKey string
	if len(packageName) > 0 {
		appenderEnvKey = PACKAGE_LOG_APPENDER_ENV_NAME + packageName
	} else {
		appenderEnvKey = DEFAULT_LOG_APPENDER_ENV_NAME
	}

	appenderName := getValueFromMapOrDefault(relevantEnvKeyValues, appenderEnvKey, APPENDER_STDOUT)

	switch appenderName {
	case APPENDER_STDOUT:
		appenderConfig.appenderType = appenderName
	case APPENDER_FILE:
		// not supported yet
		printHint(true, false, appenderName, appenderEnvKey, "appender")
	default:
		printHint(false, true, appenderName, appenderEnvKey, "appender")
	}
}

// creates all relevant formatter config elements derived from relevant environment variables
func createFormatterConfigFromEnv(relevantEnvKeyValues *map[string]string) {
	config.formatter = append(config.formatter, formatterConfig{isDefault: true, packageName: ""})
	formatterIndex := len(config.formatter) - 1
	configureFormatterFromEnv(relevantEnvKeyValues, &config.formatter[formatterIndex], "")

	for key, _ := range *relevantEnvKeyValues {
		packageOfFormatter, found := strings.CutPrefix(key, PACKAGE_LOG_FORMATTER_ENV_NAME)
		if found {
			config.formatter = append(config.formatter, formatterConfig{isDefault: false, packageName: packageOfFormatter})
			formatterIndex++
			configureFormatterFromEnv(relevantEnvKeyValues, &config.formatter[formatterIndex], packageOfFormatter)
		}
	}
}

// configures a given formatter config element from environment variables concerning a given package name
func configureFormatterFromEnv(relevantEnvKeyValues *map[string]string, formatterConfig *formatterConfig, packageName string) {
	var formatterEnvKey string
	var formatterPackageEnvKey string
	if len(packageName) > 0 {
		formatterEnvKey = PACKAGE_LOG_FORMATTER_ENV_NAME + packageName
		formatterPackageEnvKey = PACKAGE_LOG_FORMATTER_PARAMETER_ENV_NAME + packageName
	} else {
		formatterEnvKey = DEFAULT_LOG_FORMATTER_ENV_NAME
		formatterPackageEnvKey = DEFAULT_LOG_FORMATTER_PARAMETER_ENV_NAME
	}

	formatterName := getValueFromMapOrDefault(relevantEnvKeyValues, formatterEnvKey, FORMATTER_DELIMITER)
	switch formatterName {
	case FORMATTER_DELIMITER:
		formatterConfig.formatterType = formatterName
		formatterConfig.delimiter = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey, DEFAULT_DELIMITER)
	case FORMATTER_TEMPLATE:
		formatterConfig.formatterType = formatterName
		formatterConfig.template = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_1", DEFAULT_TEMPLATE)
		formatterConfig.correlationIdTemplate = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_2", DEFAULT_CORRELATION_TEMPLATE)
		formatterConfig.customTemplate = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_3", DEFAULT_CUSTOM_TEMPLATE)
		formatterConfig.timeLayout = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_4", DEFAULT_TIME_LAYOUT)
	case FORMATTER_JSON:
		formatterConfig.formatterType = formatterName
		formatterConfig.timeKey = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_1", DEFAULT_TIME_KEY)
		formatterConfig.severityKey = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_2", DEFAULT_SEVERITY_KEY)
		formatterConfig.correlationKey = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_3", DEFAULT_CORRELATION_KEY)
		formatterConfig.messageKey = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_4", DEFAULT_MESSAGE_KEY)
		formatterConfig.customValuesKey = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_5", DEFAULT_CUSTOM_VALUES_KEY)
		formatterConfig.customValuesAsSubElement = strings.ToLower(getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_6", DEFAULT_CUSTOM_AS_SUB_ELEMENT)) == "true"
		formatterConfig.timeLayout = getValueFromMapOrDefault(relevantEnvKeyValues, formatterPackageEnvKey+"_7", DEFAULT_TIME_LAYOUT)
	default:
		printHint(false, true, formatterName, formatterEnvKey, "formatter")
	}
}

// creates all relevant logger config elements derived from relevant environment variables
func createLoggerConfigFromEnv(relevantEnvKeyValues *map[string]string) {
	config.logger = append(config.logger, loggerConfig{isDefault: true, packageName: ""})
	loggerIndex := len(config.logger) - 1
	configureLoggerFromEnv(relevantEnvKeyValues, &config.logger[loggerIndex], "")

	for key, _ := range *relevantEnvKeyValues {
		packageOfLogger, found := strings.CutPrefix(key, PACKAGE_LOG_LEVEL_ENV_NAME)
		if found {
			config.logger = append(config.logger, loggerConfig{isDefault: false, packageName: packageOfLogger})
			loggerIndex++
			configureLoggerFromEnv(relevantEnvKeyValues, &config.logger[loggerIndex], packageOfLogger)
		}
	}
}

// configures a given logger config element from environment variables concerning a given package name
func configureLoggerFromEnv(relevantEnvKeyValues *map[string]string, loggerConfig *loggerConfig, packageName string) {
	var formatterEnvKey string
	if len(packageName) > 0 {
		formatterEnvKey = PACKAGE_LOG_LEVEL_ENV_NAME + packageName
	} else {
		formatterEnvKey = DEFAULT_LOG_LEVEL_ENV_NAME
	}

	loglevel := getValueFromMapOrDefault(relevantEnvKeyValues, formatterEnvKey, LOG_LEVEL_ERROR)
	severity, found := severityLevelMap[loglevel]
	if !found {
		severity = constants.ERROR_SEVERITY
	}
	loggerConfig.severity = severity
}

// Returns the value from a map for a given key. If there is none, the default will be returned
func getValueFromMapOrDefault(source *map[string]string, key string, defaultValue string) string {
	value, found := (*source)[key]
	if found {
		return value
	}
	return defaultValue
}

func printHint(isNotSupported bool, isUnkown bool, propertyName string, propertyEnvName string, objectType string) {
	if isNotSupported {
		fmt.Println(objectType, propertyName, "for logger at env variable", propertyEnvName, "is not supported yet")
	}
	if isUnkown {
		fmt.Println("Unkown", objectType, propertyName, "for logger at env variable", propertyEnvName)
	}
}

// creates default configs if missing and adds package specfic copies of defaults if at least one of the other config types exists as package variant
func completeConfig() {
	completeDefaults()

	completeAppenderConfigPackageForward()
	completeFormatterConfigPackageForward()

	completeAppenderConfigPackageBackward()
	completeLoggerConfigPackageBackward()
}

// creates default configs if missing
func completeDefaults() {
	found := false

	for _, fc := range config.formatter {
		if fc.isDefault {
			found = true
			break
		}
	}
	if !found {
		config.formatter = append(config.formatter, formatterConfig{formatterType: FORMATTER_DELIMITER, isDefault: true, packageName: "", delimiter: DEFAULT_DELIMITER})
	}

	for _, ac := range config.appender {
		if ac.isDefault {
			found = true
			break
		}
	}
	if !found {
		config.appender = append(config.appender, appenderConfig{appenderType: APPENDER_STDOUT, isDefault: true, packageName: ""})
	}

	for _, lc := range config.logger {
		if lc.isDefault {
			found = true
			break
		}
	}
	if !found {
		config.logger = append(config.logger, loggerConfig{isDefault: true, packageName: "", severity: constants.ERROR_SEVERITY})
	}
}

// creates appender configs if there exists a logger package variant
func completeAppenderConfigPackageForward() {
	for _, lc := range config.logger {
		if lc.isDefault {
			continue
		}
		createAppenderConfigIfNecessary(&lc.packageName)
	}
}

// creates appender configs if there exists a formatter package variant
func completeAppenderConfigPackageBackward() {
	for _, fc := range config.formatter {
		if fc.isDefault {
			continue
		}
		createAppenderConfigIfNecessary(&fc.packageName)
	}
}

// creates an appender config if it does not exists for a given package name
func createAppenderConfigIfNecessary(packageName *string) {
	found := false
	for _, ac := range config.appender {
		if ac.packageName == *packageName {
			found = true
			break
		}
	}
	if !found {
		for _, ac := range config.appender {
			if ac.isDefault {
				acp := ac
				acp.isDefault = false
				acp.packageName = *packageName
				config.appender = append(config.appender, acp)
				break
			}
		}
	}
}

// creates formatter configs if there exists a appender package variant
func completeFormatterConfigPackageForward() {
	for _, ac := range config.appender {
		if ac.isDefault {
			continue
		}
		createFormatterConfigIfNecessary(&ac.packageName)
	}
}

// creates an formatter config if it does not exists for a given package name
func createFormatterConfigIfNecessary(packageName *string) {
	found := false
	for _, fc := range config.formatter {
		if fc.packageName == *packageName {
			found = true
			break
		}
	}
	if !found {
		for _, fc := range config.formatter {
			if fc.isDefault {
				fcp := fc
				fcp.isDefault = false
				fcp.packageName = *packageName
				config.formatter = append(config.formatter, fcp)
				break
			}
		}
	}
}

// creates logger configs if there exists a appender package variant
func completeLoggerConfigPackageBackward() {
	for _, ac := range config.appender {
		if ac.isDefault {
			continue
		}
		createLoggerConfigIfNecessary(&ac.packageName)
	}
}

// creates an logger config if it does not exists for a given package name
func createLoggerConfigIfNecessary(packageName *string) {
	found := false
	for _, lc := range config.logger {
		if lc.packageName == *packageName {
			found = true
			break
		}
	}
	if !found {
		for _, lc := range config.logger {
			if lc.isDefault {
				lcp := lc
				lcp.isDefault = false
				lcp.packageName = *packageName
				config.logger = append(config.logger, lcp)
				break
			}
		}
	}
}
