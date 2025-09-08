// This package provides functions to derive configuration from environment variables or configuration files
package config

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ma-vin/typewriter/common"
)

// root config element
type Config struct {
	Logger                      []LoggerConfig
	Appender                    []AppenderConfig
	Formatter                   []FormatterConfig
	UseFullQualifiedPackageName bool
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

	APPENDER_STDOUT   = "STDOUT"
	APPENDER_FILE     = "FILE"
	APPENDER_MULTIPLE = "MULTIPLE"

	LOGGER_GENERAL = "GENERAL"

	DEFAULT_DELIMITER                   = " - "
	DEFAULT_TEMPLATE                    = "[%s] %s: %s"
	DEFAULT_CORRELATION_TEMPLATE        = "[%s] %s %s: %s"
	DEFAULT_CUSTOM_TEMPLATE             = DEFAULT_TEMPLATE
	DEFAULT_CALLER_TEMPLATE             = "[%s] %s %s(%s.%d): %s"
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
var configCreationMu = sync.Mutex{}

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
	DEFAULT_LOG_FORMATTER_PROPERTY_NAME,
	PACKAGE_LOG_PACKAGE_PROPERTY_NAME,
	PACKAGE_LOG_LEVEL_PROPERTY_NAME,
	PACKAGE_LOG_APPENDER_PROPERTY_NAME,
	PACKAGE_LOG_FORMATTER_PROPERTY_NAME,
	LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME,
	LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME,
}

// map containing creator functions of configurations of formatter
var registeredFormatterConfigs map[string]func(relevantKeyValues *map[string]string, commonConfig *CommonFormatterConfig) (*FormatterConfig, error)

// map containing creator functions of configurations of appender
var registeredAppenderConfigs map[string]func(relevantKeyValues *map[string]string, commonConfig *CommonAppenderConfig) (*AppenderConfig, error)

// returns the config or creates it if it was not initialized yet
func GetConfig() *Config {
	configCreationMu.Lock()
	defer configCreationMu.Unlock()

	if configInitialized {
		return &config
	}

	if len(registeredAppenderConfigs) == 0 {
		initializeRegisteredAppenderConfigs()
	}
	if len(registeredFormatterConfigs) == 0 {
		initializeRegisteredFormatterConfigs()
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
	configCreationMu.Lock()
	defer configCreationMu.Unlock()
	configInitialized = false
}

// initializes the registered appender and formatter configurations. Marks the whole config as not initialized also
func ResetRegisteredAppenderAndFormatterConfigs() {
	configCreationMu.Lock()
	defer configCreationMu.Unlock()

	relevantKeyPrefixes = []string{
		DEFAULT_LOG_LEVEL_PROPERTY_NAME,
		DEFAULT_LOG_APPENDER_PROPERTY_NAME,
		DEFAULT_LOG_FORMATTER_PROPERTY_NAME,
		PACKAGE_LOG_PACKAGE_PROPERTY_NAME,
		PACKAGE_LOG_LEVEL_PROPERTY_NAME,
		PACKAGE_LOG_APPENDER_PROPERTY_NAME,
		PACKAGE_LOG_FORMATTER_PROPERTY_NAME,
		LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME,
		LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME,
	}

	initializeRegisteredAppenderConfigs()
	initializeRegisteredFormatterConfigs()

	configInitialized = false
}

func IsConfigInitialized() bool {
	configCreationMu.Lock()
	defer configCreationMu.Unlock()
	return configInitialized
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
	createAndAppendAppenderConfig(relevantKeyValues, "")
	for key := range *relevantKeyValues {
		packageOfAppender, found := strings.CutPrefix(key, PACKAGE_LOG_APPENDER_PROPERTY_NAME)
		if found {
			createAndAppendAppenderConfig(relevantKeyValues, packageOfAppender)
		}
	}
}

func createAndAppendAppenderConfig(relevantKeyValues *map[string]string, packageOfFormatter string) {
	appenderConfig := createAppenderConfigEntry(relevantKeyValues, packageOfFormatter)
	if appenderConfig != nil {
		config.Appender = append(config.Appender, *appenderConfig)
	}
}

// configures a given appender config element from properties concerning a given package name
func createAppenderConfigEntry(relevantKeyValues *map[string]string, packageParameter string) *AppenderConfig {
	var appenderKey string
	if len(packageParameter) > 0 {
		appenderKey = PACKAGE_LOG_APPENDER_PROPERTY_NAME + packageParameter
	} else {
		appenderKey = DEFAULT_LOG_APPENDER_PROPERTY_NAME
	}

	commonAppenderConfig := CommonAppenderConfig{
		AppenderType:     getValueFromMapOrDefault(relevantKeyValues, appenderKey, APPENDER_STDOUT),
		IsDefault:        len(packageParameter) == 0,
		PackageParameter: packageParameter,
	}

	if creator, exist := getAppenderConfigCreator(&commonAppenderConfig); exist {
		result, err := creator(relevantKeyValues, &commonAppenderConfig)
		if err != nil {
			fmt.Printf("Fail to create appender config %s, because of error: %s", commonAppenderConfig.AppenderType, err)
			fmt.Println()
		}
		return result
	}

	printHint(commonAppenderConfig.AppenderType, appenderKey)
	return nil
}

func getAppenderConfigCreator(commonAppenderConfig *CommonAppenderConfig) (func(relevantKeyValues *map[string]string, commonConfig *CommonAppenderConfig) (*AppenderConfig, error), bool) {
	if strings.Contains(commonAppenderConfig.AppenderType, ",") {
		for s := range strings.SplitSeq(commonAppenderConfig.AppenderType, ",") {
			_, exist := registeredAppenderConfigs[strings.TrimSpace(s)]
			if !exist {
				return nil, false
			}
		}
		creator, exist := registeredAppenderConfigs[APPENDER_MULTIPLE]
		return creator, exist
	}
	creator, exist := registeredAppenderConfigs[commonAppenderConfig.AppenderType]
	return creator, exist
}

// Creates a standard output appender configuration
func createStdOutAppenderConfig(relevantKeyValues *map[string]string, commonAppenderConfig *CommonAppenderConfig) (*AppenderConfig, error) {
	var result AppenderConfig = StdOutAppenderConfig{Common: commonAppenderConfig}
	return &result, nil
}

// Creates a file appender configuration
func createFileAppenderConfig(relevantKeyValues *map[string]string, commonAppenderConfig *CommonAppenderConfig) (*AppenderConfig, error) {
	var fileParameterKey string
	var cronParameterKey string
	var sizeParameterKey string

	if len(commonAppenderConfig.PackageParameter) > 0 {
		fileParameterKey = PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME + commonAppenderConfig.PackageParameter
		cronParameterKey = PACKAGE_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME + commonAppenderConfig.PackageParameter
		sizeParameterKey = PACKAGE_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME + commonAppenderConfig.PackageParameter
	} else {
		fileParameterKey = DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME
		cronParameterKey = DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME
		sizeParameterKey = DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME
	}

	if fileValue, fileFound := (*relevantKeyValues)[fileParameterKey]; fileFound {
		var result AppenderConfig = FileAppenderConfig{
			Common:         commonAppenderConfig,
			PathToLogFile:  fileValue,
			CronExpression: getValueFromMapOrDefault(relevantKeyValues, cronParameterKey, ""),
			LimitByteSize:  getValueFromMapOrDefault(relevantKeyValues, sizeParameterKey, ""),
		}
		return &result, nil
	}
	return nil, fmt.Errorf("cannot use file appender, because there is no value at %s. Use %s appender instead", fileParameterKey, APPENDER_STDOUT)
}

func createMultipleAppenderConfig(relevantKeyValues *map[string]string, commonAppenderConfig *CommonAppenderConfig) (*AppenderConfig, error) {
	appenderTypes := strings.Split(commonAppenderConfig.AppenderType, ",")
	appenderConfigs := make([]AppenderConfig, len(appenderTypes))

	for i, s := range appenderTypes {
		s = strings.TrimSpace(s)
		subCommonAppenderConfig := CommonAppenderConfig{
			AppenderType:     s,
			IsDefault:        commonAppenderConfig.IsDefault,
			PackageParameter: commonAppenderConfig.PackageParameter,
		}
		appenderConfig, err := registeredAppenderConfigs[s](relevantKeyValues, &subCommonAppenderConfig)
		if err != nil {
			return nil, err
		}
		appenderConfigs[i] = *appenderConfig
	}

	commonAppenderConfig.AppenderType = APPENDER_MULTIPLE
	var result AppenderConfig = MultiAppenderConfig{
		Common:          commonAppenderConfig,
		AppenderConfigs: &appenderConfigs,
	}

	return &result, nil
}

// Registers the out of the box stdOut- and file-appender
func initializeRegisteredAppenderConfigs() {
	registeredAppenderConfigs = make(map[string]func(relevantKeyValues *map[string]string, commonAppenderConfig *CommonAppenderConfig) (*AppenderConfig, error))

	fileKeyPrefixes := []string{
		DEFAULT_LOG_APPENDER_FILE_PROPERTY_NAME,
		DEFAULT_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME,
		DEFAULT_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME,
		PACKAGE_LOG_APPENDER_FILE_PROPERTY_NAME,
		PACKAGE_LOG_APPENDER_CRON_RENAMING_PROPERTY_NAME,
		PACKAGE_LOG_APPENDER_SIZE_RENAMING_PROPERTY_NAME,
	}

	registerAppenderConfigInternal(APPENDER_STDOUT, []string{}, createStdOutAppenderConfig)
	registerAppenderConfigInternal(APPENDER_FILE, fileKeyPrefixes, createFileAppenderConfig)
	registerAppenderConfigInternal(APPENDER_MULTIPLE, []string{}, createMultipleAppenderConfig)
}

// Registers a creator of a configuration for an appender with the given identifier 'appenderType'. But without initialization check and mutex lock.
// The relevant key values pairs, from environment or file, will be filtered by given key prefixes and 'keyPrefixes' will be added to these relevant ones, so that they will be available at 'relevantKeyValues'.
func registerAppenderConfigInternal(appenderType string, keyPrefixes []string, configCreator func(relevantKeyValues *map[string]string, commonConfig *CommonAppenderConfig) (*AppenderConfig, error)) error {
	registeredAppenderConfigs[appenderType] = configCreator

	for _, keyPrefix := range keyPrefixes {
		keyPrefix = strings.ToUpper(keyPrefix)
		if !slices.Contains(relevantKeyPrefixes, keyPrefix) {
			relevantKeyPrefixes = append(relevantKeyPrefixes, keyPrefix)
		}
	}

	return nil
}

// Registers a creator of a configuration for an appender with the given identifier 'appenderType' and resets any existing configuration.
// The relevant key values pairs, from environment or file, will be filtered by given key prefixes and 'keyPrefixes' will be added to these relevant ones, so that they will be available at 'relevantKeyValues'.
func RegisterAppenderConfig(appenderType string, keyPrefixes []string, configCreator func(relevantKeyValues *map[string]string, commonConfig *CommonAppenderConfig) (*AppenderConfig, error)) error {
	configCreationMu.Lock()
	defer configCreationMu.Unlock()

	configInitialized = false

	if len(registeredAppenderConfigs) == 0 {
		initializeRegisteredAppenderConfigs()
	}

	appenderType = strings.ToUpper(appenderType)

	if _, exists := registeredAppenderConfigs[appenderType]; exists {
		return fmt.Errorf("there exists a creator of a configuration for appender %s already. The existing one will not be replaced", appenderType)
	}

	return registerAppenderConfigInternal(appenderType, keyPrefixes, configCreator)
}

// Removes a registered creator of a configuration for an appender. Build in ones will not be removed
func DeregisterAppenderConfig(appenderType string) error {
	appenderType = strings.ToUpper(appenderType)

	if appenderType == APPENDER_STDOUT || appenderType == APPENDER_FILE {
		return fmt.Errorf("the build in config creator for appender %s is not deletable", appenderType)
	}
	if _, exists := registeredAppenderConfigs[appenderType]; !exists {
		return fmt.Errorf("there does not exists any creator of configuration for appender %s", appenderType)
	}

	configCreationMu.Lock()
	defer configCreationMu.Unlock()

	configInitialized = false

	delete(registeredAppenderConfigs, appenderType)
	return nil
}

// creates all relevant formatter config elements derived from relevant properties
func createFormatterConfig(relevantKeyValues *map[string]string) {
	createAndAppendFormatterConfig(relevantKeyValues, "")

	for key := range *relevantKeyValues {
		packageOfFormatter, found := strings.CutPrefix(key, PACKAGE_LOG_FORMATTER_PROPERTY_NAME)
		if found {
			createAndAppendFormatterConfig(relevantKeyValues, packageOfFormatter)
		}
	}
}

func createAndAppendFormatterConfig(relevantKeyValues *map[string]string, packageOfFormatter string) {
	formatterConfig := createFormatterConfigEntry(relevantKeyValues, packageOfFormatter)
	if formatterConfig != nil {
		config.Formatter = append(config.Formatter, *formatterConfig)
	}
}

// creates and configures a formatter config element from properties concerning a given package name
func createFormatterConfigEntry(relevantKeyValues *map[string]string, packageParameter string) *FormatterConfig {
	var formatterKey string
	var formatterParameterKey string
	if len(packageParameter) > 0 {
		formatterKey = PACKAGE_LOG_FORMATTER_PROPERTY_NAME + packageParameter
		formatterParameterKey = PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME + packageParameter
	} else {
		formatterKey = DEFAULT_LOG_FORMATTER_PROPERTY_NAME
		formatterParameterKey = DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME
	}

	commonFormatterConfig := CommonFormatterConfig{
		FormatterType:    getValueFromMapOrDefault(relevantKeyValues, formatterKey, FORMATTER_DELIMITER),
		IsDefault:        len(packageParameter) == 0,
		PackageParameter: packageParameter,
		TimeLayout:       getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TIME_LAYOUT_PARAMETER, DEFAULT_TIME_LAYOUT),
	}

	if creator, exist := registeredFormatterConfigs[commonFormatterConfig.FormatterType]; exist {
		result, err := creator(relevantKeyValues, &commonFormatterConfig)
		if err != nil {
			fmt.Printf("Fail to create formatter config %s, because of error: %s", commonFormatterConfig.FormatterType, err)
			fmt.Println()
		}
		return result
	}

	printHint(commonFormatterConfig.FormatterType, formatterKey)
	return nil
}

// Creates a delimiter formatter configuration
func createDelimiterFormatterConfig(relevantKeyValues *map[string]string, commonFormatterConfig *CommonFormatterConfig) (*FormatterConfig, error) {
	var formatterParameterKey string
	if len(commonFormatterConfig.PackageParameter) > 0 {
		formatterParameterKey = PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME + commonFormatterConfig.PackageParameter
	} else {
		formatterParameterKey = DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME
	}

	var result FormatterConfig = DelimiterFormatterConfig{
		Common:    commonFormatterConfig,
		Delimiter: getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+DELIMITER_PARAMETER, DEFAULT_DELIMITER),
	}

	return &result, nil
}

// Creates a template formatter configuration
func createTemplateFormatterConfig(relevantKeyValues *map[string]string, commonFormatterConfig *CommonFormatterConfig) (*FormatterConfig, error) {
	var formatterParameterKey string
	if len(commonFormatterConfig.PackageParameter) > 0 {
		formatterParameterKey = PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME + commonFormatterConfig.PackageParameter
	} else {
		formatterParameterKey = DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME
	}

	var result FormatterConfig = TemplateFormatterConfig{
		Common:                      commonFormatterConfig,
		Template:                    getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_PARAMETER, DEFAULT_TEMPLATE),
		CorrelationIdTemplate:       getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CORRELATION_PARAMETER, DEFAULT_CORRELATION_TEMPLATE),
		CustomTemplate:              getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CUSTOM_PARAMETER, DEFAULT_CUSTOM_TEMPLATE),
		TrimSeverityText:            strings.ToLower(getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_TRIM_SEVERITY_PARAMETER, DEFAULT_TRIM_SEVERITY_TEXT)) == "true",
		CallerTemplate:              getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CALLER_PARAMETER, DEFAULT_CALLER_TEMPLATE),
		CallerCorrelationIdTemplate: getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CALLER_CORRELATION_PARAMETER, DEFAULT_CALLER_CORRELATION_TEMPLATE),
		CallerCustomTemplate:        getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+TEMPLATE_CALLER_CUSTOM_PARAMETER, DEFAULT_CALLER_CUSTOM_TEMPLATE),
	}

	return &result, nil
}

// Creates a json formatter configuration
func createJsonFormatterConfig(relevantKeyValues *map[string]string, commonFormatterConfig *CommonFormatterConfig) (*FormatterConfig, error) {
	var formatterParameterKey string
	if len(commonFormatterConfig.PackageParameter) > 0 {
		formatterParameterKey = PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME + commonFormatterConfig.PackageParameter
	} else {
		formatterParameterKey = DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME
	}

	var result FormatterConfig = JsonFormatterConfig{
		Common:                   commonFormatterConfig,
		TimeKey:                  getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_TIME_KEY_PARAMETER, DEFAULT_TIME_KEY),
		SeverityKey:              getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_SEVERITY_KEY_PARAMETER, DEFAULT_SEVERITY_KEY),
		CorrelationKey:           getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CORRELATION_KEY_PARAMETER, DEFAULT_CORRELATION_KEY),
		MessageKey:               getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_MESSAGE_KEY_PARAMETER, DEFAULT_MESSAGE_KEY),
		CustomValuesKey:          getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CUSTOM_VALUES_KEY_PARAMETER, DEFAULT_CUSTOM_VALUES_KEY),
		CustomValuesAsSubElement: strings.ToLower(getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CUSTOM_VALUES_SUB_PARAMETER, DEFAULT_CUSTOM_AS_SUB_ELEMENT)) == "true",
		CallerFunctionKey:        getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CALLER_FUNCTION_KEY_PARAMETER, DEFAULT_CALLER_FUNCTION_KEY),
		CallerFileKey:            getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CALLER_FILE_KEY_PARAMETER, DEFAULT_CALLER_FILE_KEY),
		CallerFileLineKey:        getValueFromMapOrDefault(relevantKeyValues, formatterParameterKey+JSON_CALLER_LINE_KEY_PARAMETER, DEFAULT_CALLER_FILE_LINE_KEY),
	}

	return &result, nil
}

// Registers the out of the box delimiter-, template- and json-formatter
func initializeRegisteredFormatterConfigs() {
	registeredFormatterConfigs = make(map[string]func(relevantKeyValues *map[string]string, commonFormatterConfig *CommonFormatterConfig) (*FormatterConfig, error))

	keyPrefixes := []string{DEFAULT_LOG_FORMATTER_PARAMETER_PROPERTY_NAME, PACKAGE_LOG_FORMATTER_PARAMETER_PROPERTY_NAME}

	registerFormatterConfigInternal(FORMATTER_DELIMITER, keyPrefixes, createDelimiterFormatterConfig)
	registerFormatterConfigInternal(FORMATTER_TEMPLATE, keyPrefixes, createTemplateFormatterConfig)
	registerFormatterConfigInternal(FORMATTER_JSON, keyPrefixes, createJsonFormatterConfig)
}

// Registers a creator of a configuration for a formatter with the given identifier 'formatterType'. But without initialization check and mutex lock.
// The relevant key values pairs, from environment or file, will be filtered by given key prefixes and 'keyPrefixes' will be added to these relevant ones, so that they will be available at 'relevantKeyValues'.
func registerFormatterConfigInternal(formatterType string, keyPrefixes []string, configCreator func(relevantKeyValues *map[string]string, commonConfig *CommonFormatterConfig) (*FormatterConfig, error)) error {
	registeredFormatterConfigs[formatterType] = configCreator

	for _, keyPrefix := range keyPrefixes {
		keyPrefix = strings.ToUpper(keyPrefix)
		if !slices.Contains(relevantKeyPrefixes, keyPrefix) {
			relevantKeyPrefixes = append(relevantKeyPrefixes, keyPrefix)
		}
	}

	return nil
}

// Registers a creator of a configuration for a formatter with the given identifier 'formatterType' and resets any existing configuration.
// The relevant key values pairs, from environment or file, will be filtered by given key prefixes and 'keyPrefixes' will be added to these relevant ones, so that they will be available at 'relevantKeyValues'.
func RegisterFormatterConfig(formatterType string, keyPrefixes []string, configCreator func(relevantKeyValues *map[string]string, commonConfig *CommonFormatterConfig) (*FormatterConfig, error)) error {
	configCreationMu.Lock()
	defer configCreationMu.Unlock()

	configInitialized = false

	if len(registeredFormatterConfigs) == 0 {
		initializeRegisteredFormatterConfigs()
	}

	formatterType = strings.ToUpper(formatterType)

	if _, exists := registeredFormatterConfigs[formatterType]; exists {
		return fmt.Errorf("there exists a creator of a configuration for formatter %s already. The existing one will not be replaced", formatterType)
	}

	return registerFormatterConfigInternal(formatterType, keyPrefixes, configCreator)
}

// Removes a registered creator of a configuration for a formatter. Build in ones will not be removed
func DeregisterFormatterConfig(formatterType string) error {
	formatterType = strings.ToUpper(formatterType)

	if formatterType == FORMATTER_DELIMITER || formatterType == FORMATTER_TEMPLATE || formatterType == FORMATTER_JSON {
		return fmt.Errorf("the build in config creator for formatter %s is not deletable", formatterType)
	}
	if _, exists := registeredFormatterConfigs[formatterType]; !exists {
		return fmt.Errorf("there does not exists any creator of configuration for formatter %s", formatterType)
	}

	configCreationMu.Lock()
	defer configCreationMu.Unlock()

	configInitialized = false

	delete(registeredFormatterConfigs, formatterType)
	return nil
}

// creates all relevant logger config elements derived from relevant properties
func createLoggerConfig(relevantKeyValues *map[string]string) {
	createAndAppendLoggerConfig(relevantKeyValues, "")

	for key := range *relevantKeyValues {
		packageParameter, found := strings.CutPrefix(key, PACKAGE_LOG_LEVEL_PROPERTY_NAME)
		if found {
			createAndAppendLoggerConfig(relevantKeyValues, packageParameter)
		}
	}

	if fullQualifiedText, found := (*relevantKeyValues)[LOG_CONFIG_FULL_QUALIFIED_PACKAGE_ENV_NAME]; found {
		config.UseFullQualifiedPackageName = strings.ToUpper(fullQualifiedText) == "TRUE"
	}
}

func createAndAppendLoggerConfig(relevantKeyValues *map[string]string, packageOfLogger string) {
	loggerConfig := createLoggerConfigEntry(relevantKeyValues, packageOfLogger)
	if loggerConfig != nil {
		config.Logger = append(config.Logger, *loggerConfig)
	}
}

// configures a given logger config element from properties concerning a given package name
func createLoggerConfigEntry(relevantKeyValues *map[string]string, packageParameter string) *LoggerConfig {
	var logLevelKey string
	var packageName string
	if len(packageParameter) > 0 {
		logLevelKey = PACKAGE_LOG_LEVEL_PROPERTY_NAME + packageParameter
		packageName = getValueFromMapOrDefault(relevantKeyValues, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+packageParameter, strings.ToLower(packageParameter))
	} else {
		logLevelKey = DEFAULT_LOG_LEVEL_PROPERTY_NAME
		packageName = ""
	}

	logLevel := getValueFromMapOrDefault(relevantKeyValues, logLevelKey, LOG_LEVEL_ERROR)
	severity, found := SeverityLevelMap[logLevel]
	if !found {
		severity = common.ERROR_SEVERITY
	}
	var result LoggerConfig = GeneralLoggerConfig{
		Common: &CommonLoggerConfig{
			IsDefault:        len(packageParameter) == 0,
			LoggerType:       LOGGER_GENERAL,
			PackageParameter: packageParameter,
			PackageName:      packageName,
		},
		Severity:      severity,
		IsCallerToSet: strings.ToLower(getValueFromMapOrDefault(relevantKeyValues, LOG_CONFIG_IS_CALLER_TO_SET_ENV_NAME, "false")) == "true",
	}

	return &result
}

// Returns the value from a map for a given key. If there is none, the default will be returned
func getValueFromMapOrDefault(source *map[string]string, key string, defaultValue string) string {
	value, found := (*source)[key]
	if found {
		return value
	}
	return defaultValue
}

func printHint(propertyValue string, propertyName string) {
	fmt.Printf("unknown \"%s\" value at property %s", propertyValue, propertyName)
	fmt.Println()
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
		if fc.IsDefault() {
			found = true
			break
		}
	}
	if !found {
		config.Formatter = append(config.Formatter, DelimiterFormatterConfig{Common: &CommonFormatterConfig{FormatterType: FORMATTER_DELIMITER, IsDefault: true, PackageParameter: "", TimeLayout: DEFAULT_TIME_LAYOUT}, Delimiter: DEFAULT_DELIMITER})
	}

	found = false
	for _, ac := range config.Appender {
		if ac.IsDefault() {
			found = true
			break
		}
	}
	if !found {
		config.Appender = append(config.Appender, StdOutAppenderConfig{Common: &CommonAppenderConfig{AppenderType: APPENDER_STDOUT, IsDefault: true, PackageParameter: ""}})
	}

	found = false
	for _, lc := range config.Logger {
		if lc.IsDefault() {
			found = true
			break
		}
	}
	if !found {
		config.Logger = append(config.Logger, GeneralLoggerConfig{Common: &CommonLoggerConfig{LoggerType: LOGGER_GENERAL, IsDefault: true, PackageParameter: ""}, Severity: common.ERROR_SEVERITY})
	}
}

// creates appender configs if there exists a logger package variant
func completeAppenderConfigPackageForward() {
	for _, lc := range config.Logger {
		if lc.IsDefault() {
			continue
		}
		packageParameter := lc.PackageParameter()
		createAppenderConfigIfNecessary(&packageParameter)
	}
}

// creates appender configs if there exists a formatter package variant
func completeAppenderConfigPackageBackward() {
	for _, fc := range config.Formatter {
		if fc.IsDefault() {
			continue
		}
		createAppenderConfigIfNecessary(&fc.GetCommon().PackageParameter)
	}
}

// creates an appender config if it does not exists for a given package name
func createAppenderConfigIfNecessary(packageParameter *string) {
	found := false
	for _, ac := range config.Appender {
		if ac.PackageParameter() == *packageParameter {
			found = true
			break
		}
	}
	if !found {
		for _, ac := range config.Appender {
			if ac.IsDefault() {
				acp := ac.CreateFullCopy()
				acp.GetCommon().IsDefault = false
				acp.GetCommon().PackageParameter = *packageParameter
				config.Appender = append(config.Appender, acp)
				break
			}
		}
	}
}

// creates formatter configs if there exists a appender package variant
func completeFormatterConfigPackageForward() {
	for _, ac := range config.Appender {
		if ac.IsDefault() {
			continue
		}
		createFormatterConfigIfNecessary(&ac.GetCommon().PackageParameter)
	}
}

// creates an formatter config if it does not exists for a given package name
func createFormatterConfigIfNecessary(packageParameter *string) {
	found := false
	for _, fc := range config.Formatter {
		if fc.PackageParameter() == *packageParameter {
			found = true
			break
		}
	}
	if !found {
		for _, fc := range config.Formatter {
			if fc.IsDefault() {
				fcp := fc.CreateFullCopy()
				fcp.GetCommon().IsDefault = false
				fcp.GetCommon().PackageParameter = *packageParameter
				config.Formatter = append(config.Formatter, fcp)
				break
			}
		}
	}
}

// creates logger configs if there exists a appender package variant
func completeLoggerConfigPackageBackward(relevantKeyValues *map[string]string) {
	for _, ac := range config.Appender {
		if ac.IsDefault() {
			continue
		}
		createLoggerConfigIfNecessary(relevantKeyValues, &ac.GetCommon().PackageParameter)
	}
}

// creates an logger config if it does not exists for a given package name
func createLoggerConfigIfNecessary(relevantKeyValues *map[string]string, packageParameter *string) {
	found := false
	for _, lc := range config.Logger {
		if lc.PackageParameter() == *packageParameter {
			found = true
			break
		}
	}
	if !found {
		for _, lc := range config.Logger {
			if lc.IsDefault() {
				lcp := lc.CreateFullCopy()
				lcp.GetCommon().IsDefault = false
				lcp.GetCommon().PackageParameter = *packageParameter
				lcp.GetCommon().PackageName = getValueFromMapOrDefault(relevantKeyValues, PACKAGE_LOG_PACKAGE_PROPERTY_NAME+*packageParameter, strings.ToLower(*packageParameter))
				config.Logger = append(config.Logger, lcp)
				break
			}
		}
	}
}

// Sorts the config put the default configs at first index. Because of this a potential cronRenamer of the default appender is used for all equal named files
func sortConfig() {
	sort.Slice(config.Formatter, func(i, j int) bool {
		return config.Formatter[i].GetCommon().LessCompareForSort(config.Formatter[j].GetCommon())
	})
	sort.Slice(config.Appender, func(i, j int) bool {
		return config.Appender[i].GetCommon().LessCompareForSort(config.Appender[j].GetCommon())
	})
	sort.Slice(config.Logger, func(i, j int) bool {
		return (config.Logger[i].GetCommon().LessCompareForSort(config.Logger[j].GetCommon()))
	})
}

func determineIds() {
	determineFormatterIds()
	determineAppenderIds()
	determineLoggerIds()
}

func determineFormatterIds() {
	for i := 0; i < len(config.Formatter); i++ {
		if config.Formatter[i].Id() != "" {
			continue
		}
		config.Formatter[i].GetCommon().Id = fmt.Sprint("formatter", i)
		for j := i + 1; j < len(config.Formatter); j++ {
			if config.Formatter[i].Equals(&config.Formatter[j]) {
				config.Formatter[j].GetCommon().Id = config.Formatter[i].Id()
			}
		}
	}
}

func determineAppenderIds() {
	allAppenders := make([]AppenderConfig, 0, len(config.Appender))
	for _, a1 := range config.Appender {
		allAppenders = append(allAppenders, a1)
		if a1.GetCommon().AppenderType == APPENDER_MULTIPLE {
			allAppenders = append(allAppenders, *a1.(MultiAppenderConfig).AppenderConfigs...)
		}
	}

	for i := 0; i < len(allAppenders); i++ {
		if allAppenders[i].Id() != "" {
			continue
		}
		allAppenders[i].GetCommon().Id = fmt.Sprint("appender", i)
		for j := i + 1; j < len(allAppenders); j++ {
			if allAppenders[i].Equals(&allAppenders[j]) {
				allAppenders[j].GetCommon().Id = allAppenders[i].Id()
			}
		}
	}
}

func determineLoggerIds() {
	for i := 0; i < len(config.Logger); i++ {
		if config.Logger[i].Id() != "" {
			continue
		}
		config.Logger[i].GetCommon().Id = fmt.Sprint("logger", i)
		for j := i + 1; j < len(config.Logger); j++ {
			if config.Logger[i].Equals(&config.Logger[j]) {
				config.Logger[j].GetCommon().Id = config.Logger[i].Id()
			}
		}
	}
}
