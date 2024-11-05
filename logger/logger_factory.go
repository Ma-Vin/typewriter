package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/ma-vin/typewriter/constants"
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
	FORMATTER_JSON      = "JSON"

	APPENDER_STDOUT = "STDOUT"
	APPENDER_FILE   = "FILE"

	DEFAULT_DELIMITER = " - "
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
}

// config of an appender
type appenderConfig struct {
	appenderType string
	isDefault    bool
	packageName  string
}

// config of a formatter
type formatterConfig struct {
	formatterType string
	isDefault     bool
	packageName   string
	delimiter     string
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
		keyValue := strings.SplitN(strings.ToUpper(envEntry), "=", 2)
		if len(keyValue) == 2 && strings.TrimSpace(keyValue[1]) != "" {
			for _, keyPrefix := range relevantKeyPrefixes {
				if strings.HasPrefix(keyValue[0], keyPrefix) {
					result[keyValue[0]] = strings.ToUpper(strings.TrimSpace(keyValue[1]))
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
	case FORMATTER_JSON:
		// not supported yet
		printHint(true, false, formatterName, formatterEnvKey, "formatter")
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
