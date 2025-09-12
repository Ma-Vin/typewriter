package logger

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
)

var loggersInitialized = false
var loggerCreationMu = sync.Mutex{}
var mLogger MainLogger
var cLoggers []GeneralLogger
var appenders []appender.Appender
var formatters []format.Formatter

// map containing creator functions for appender
var registeredAppenders map[string]func(config *config.AppenderConfig, formatter *format.Formatter) (*appender.Appender, error)

// map containing creator functions for formatter
var registeredFormatters map[string]func(config *config.FormatterConfig) (*format.Formatter, error)

// returns the main logger. If not initialized, a new one will be created from config
func getLoggers() *MainLogger {
	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	if loggersInitialized && config.IsConfigInitialized() {
		return &mLogger
	}

	if len(registeredAppenders) == 0 {
		initializeRegisteredAppenders()
	}
	if len(registeredFormatters) == 0 {
		initializeRegisteredFormatters()
	}

	conf := config.GetConfig()
	if conf == nil {
		return nil
	}

	formatterConfigMapping := make(map[string]*format.Formatter, len(conf.Formatter))
	appenderConfigMapping := make(map[string]*appender.Appender, len(conf.Appender))
	loggerConfigMapping := make(map[string]*GeneralLogger, len(conf.Logger))

	createFormatters(&conf.Formatter, &formatterConfigMapping)
	createAppenders(conf, &appenderConfigMapping, &formatterConfigMapping)
	createGeneralLoggers(conf, &loggerConfigMapping, &appenderConfigMapping)
	createMainLogger(conf, &loggerConfigMapping)

	return &mLogger
}

// initializes the registered appenders and formatters. Marks the whole logger and configuration as not initialized also
func ResetRegisteredAppenderAndFormatters() {
	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	initializeRegisteredAppenders()
	initializeRegisteredFormatters()

	loggersInitialized = false
}

// Creates the all relevant formatters from config elements
func createFormatters(formatterConfigs *[]config.FormatterConfig, formatterConfigMapping *map[string]*format.Formatter) {
	for _, fc1 := range *formatterConfigs {
		alreadyCreated := false

		for _, fc2 := range *formatterConfigs {
			if fc1.Id() == fc2.Id() {
				_, alreadyCreated = (*formatterConfigMapping)[fc1.Id()]
				break
			}
		}

		if alreadyCreated {
			continue
		}

		result := createFormatter(fc1)
		appendFormatter(*result)
		setLastFormatter(fc1.Id(), formatterConfigMapping)
	}
}

func createFormatter(formatterConfig config.FormatterConfig) *format.Formatter {
	var result *format.Formatter
	var err error
	if creator, exist := registeredFormatters[formatterConfig.FormatterType()]; exist {
		result, err = creator(&formatterConfig)
		if err != nil {
			fmt.Printf("Fail to create formatter %s, because of error: %s", formatterConfig.FormatterType(), err)
			fmt.Println()
		}
	}
	if result == nil {
		var defaultFormatConfig config.FormatterConfig = config.DelimiterFormatterConfig{
			Delimiter: config.DEFAULT_DELIMITER,
			Common:    &config.CommonFormatterConfig{FormatterType: config.FORMATTER_DELIMITER, TimeLayout: config.DEFAULT_TIME_LAYOUT},
		}
		result, _ = format.CreateDelimiterFormatterFromConfig(&defaultFormatConfig)
	}
	return result

}

func appendFormatter(formatter format.Formatter) {
	formatters = append(formatters, formatter)
}

func setLastFormatter(formatterId string, formatterConfigMapping *map[string]*format.Formatter) {
	(*formatterConfigMapping)[formatterId] = &formatters[len(formatters)-1]
}

// Registers the out of the box delimiter-, template- and json-formatter
func initializeRegisteredFormatters() {
	registeredFormatters = make(map[string]func(formatterConfig *config.FormatterConfig) (*format.Formatter, error))

	registerFormatterInternal(config.FORMATTER_DELIMITER, format.CreateDelimiterFormatterFromConfig)
	registerFormatterInternal(config.FORMATTER_TEMPLATE, format.CreateTemplateFormatterFromConfig)
	registerFormatterInternal(config.FORMATTER_JSON, format.CreateJsonFormatterFromConfig)
}

// Registers a creator of a formatter with the given identifier 'formatterType' and resets any existing configuration. But without initiation check and mutex lock
func registerFormatterInternal(formatterType string, formatterCreator func(formatterConfig *config.FormatterConfig) (*format.Formatter, error)) error {
	registeredFormatters[formatterType] = formatterCreator
	return nil
}

// Registers a creator of a formatter with the given identifier 'formatterType' and resets any existing configuration.
func RegisterFormatter(formatterType string, formatterCreator func(formatterConfig *config.FormatterConfig) (*format.Formatter, error)) error {
	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	loggersInitialized = false

	if len(registeredFormatters) == 0 {
		initializeRegisteredFormatters()
	}

	formatterType = strings.ToUpper(formatterType)

	if _, exists := registeredFormatters[formatterType]; exists {
		return fmt.Errorf("there exists a creator of a formatter %s already. The existing one will not be replaced", formatterType)
	}

	return registerFormatterInternal(formatterType, formatterCreator)
}

// Removes a registered creator of a configuration for a formatter. Build in ones will not be removed
func DeregisterFormatter(formatterType string) error {
	formatterType = strings.ToUpper(formatterType)

	if formatterType == config.FORMATTER_DELIMITER || formatterType == config.FORMATTER_TEMPLATE || formatterType == config.FORMATTER_JSON {
		return fmt.Errorf("the build in formatter creator for  %s is not deletable", formatterType)
	}
	if _, exists := registeredFormatters[formatterType]; !exists {
		return fmt.Errorf("there does not exists any creator of formatter %s", formatterType)
	}

	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	loggersInitialized = false

	delete(registeredFormatters, formatterType)
	return nil
}

// Creates the all relevant appenders from config elements
func createAppenders(conf *config.Config, appenderConfigMapping *map[string]*appender.Appender, formatterConfigMapping *map[string]*format.Formatter) {
	appender.CleanFileDeductions()
	flattenAppenderConfigs := flattenAppenderConfigs(&conf.Appender)
	for i, ac1 := range *flattenAppenderConfigs {
		if getExistingAppender(i, ac1, flattenAppenderConfigs, &conf.Formatter, appenderConfigMapping) != nil {
			continue
		}

		ac1FormatterId := (*getFormatterConfigForPackage(&ac1.GetCommon().PackageParameter, &conf.Formatter)).Id()
		appender := createAppender(ac1, (*formatterConfigMapping)[ac1FormatterId], ac1FormatterId, appenderConfigMapping)
		appendAppender(*appender)
		setLastAppender(ac1.Id()+ac1FormatterId, appenderConfigMapping)
	}
}

// Flattens a list of appender configs by adding the sub-appenders to the result list also
func flattenAppenderConfigs(appenderConfigs *[]config.AppenderConfig) *[]config.AppenderConfig {
	result := make([]config.AppenderConfig, 0, len(*appenderConfigs))
	for _, c := range *appenderConfigs {
		if containsAppenderConfig(&result, &c) {
			continue
		}
		if c.GetCommon().AppenderType == config.APPENDER_MULTIPLE {
			for _, c2 := range *c.(config.MultiAppenderConfig).AppenderConfigs {
				if !containsAppenderConfig(&result, &c2) {
					result = append(result, c2)
				}
			}
		}
		result = append(result, c)
	}
	return &result
}

// Checks whether an appender config is contained by id and package name at a slice of configs or not
func containsAppenderConfig(appenderConfigs *[]config.AppenderConfig, appenderConfig *config.AppenderConfig) bool {
	return slices.ContainsFunc(*appenderConfigs, func(c config.AppenderConfig) bool {
		return (*appenderConfig).Id() == c.Id() && (*appenderConfig).PackageParameter() == c.PackageParameter()
	})
}

// Determines if there exists an appender for a given config with respect to its id and the package relevant formatter id. If there does not exist any, nil will be returned.
func getExistingAppender(appenderConfigIndex int, appenderConfig config.AppenderConfig, flattenAppenderConfigs *[]config.AppenderConfig,
	formatterConfigs *[]config.FormatterConfig, appenderConfigMapping *map[string]*appender.Appender) *appender.Appender {

	formatterId := (*getFormatterConfigForPackage(&appenderConfig.GetCommon().PackageParameter, formatterConfigs)).Id()

	for i, a := range *flattenAppenderConfigs {
		if i >= appenderConfigIndex {
			return nil
		}
		if appenderConfig.Id() == a.Id() && formatterId == (*getFormatterConfigForPackage(&a.GetCommon().PackageParameter, formatterConfigs)).Id() {
			return (*appenderConfigMapping)[appenderConfig.Id()+formatterId]
		}
	}
	return nil
}

// Creates the appender for a given config and formatter. To be able to find existing sub-appenders for a multi appender, the id of the formatter and appender id map is relevant also
func createAppender(appenderConfig config.AppenderConfig, formatter *format.Formatter, formatterId string, appenderConfigMapping *map[string]*appender.Appender) *appender.Appender {
	if appenderConfig.GetCommon().AppenderType == config.APPENDER_MULTIPLE {
		return createMultiAppender(appenderConfig.(config.MultiAppenderConfig), formatterId, appenderConfigMapping)
	}

	var result *appender.Appender
	var err error

	if creator, exist := registeredAppenders[appenderConfig.AppenderType()]; exist {
		result, err = creator(&appenderConfig, formatter)
		if err != nil {
			fmt.Printf("Fail to create appender %s, because of error: %s", appenderConfig.AppenderType(), err)
			fmt.Println()
		}
	}
	if result == nil {
		result, _ = appender.CreateStandardOutputAppenderFromConfig(nil, formatter)
	}
	return result
}

// Creates a new multi appender. All sub-appenders have to created before usage
func createMultiAppender(appenderConfig config.MultiAppenderConfig, formatterId string, appenderConfigMapping *map[string]*appender.Appender) *appender.Appender {
	multiAppender := appender.CreateMultiAppenderWithCapacity(len(*appenderConfig.AppenderConfigs))
	for _, ac := range *appenderConfig.AppenderConfigs {
		subAppender := (*appenderConfigMapping)[ac.Id()+formatterId]
		multiAppender.AddSubAppender(subAppender)
	}
	var appender appender.Appender = *multiAppender
	return &appender
}

func appendAppender(appender appender.Appender) {
	appenders = append(appenders, appender)
}

func setLastAppender(appenderId string, appenderConfigMapping *map[string]*appender.Appender) {
	(*appenderConfigMapping)[appenderId] = &appenders[len(appenders)-1]
}

// Registers the out of the box stdOut- and file-appender
func initializeRegisteredAppenders() {
	registeredAppenders = make(map[string]func(appenderConfig *config.AppenderConfig, formatter *format.Formatter) (*appender.Appender, error))

	registerAppenderInternal(config.APPENDER_STDOUT, appender.CreateStandardOutputAppenderFromConfig)
	registerAppenderInternal(config.APPENDER_FILE, appender.CreateFileAppenderFromConfig)
}

// Registers a creator of an appender with the given identifier 'appenderType' and resets any existing configuration. But without initiation check and mutex lock
func registerAppenderInternal(appenderType string, appenderCreator func(config *config.AppenderConfig, formatter *format.Formatter) (*appender.Appender, error)) error {
	registeredAppenders[appenderType] = appenderCreator
	return nil
}

// Registers a creator of an appender with the given identifier 'appenderType' and resets any existing configuration.
func RegisterAppender(appenderType string, appenderCreator func(appenderConfig *config.AppenderConfig, formatter *format.Formatter) (*appender.Appender, error)) error {
	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	loggersInitialized = false

	if len(registeredAppenders) == 0 {
		initializeRegisteredAppenders()
	}

	appenderType = strings.ToUpper(appenderType)

	if _, exists := registeredAppenders[appenderType]; exists {
		return fmt.Errorf("there exists a creator of an appender %s already. The existing one will not be replaced", appenderType)
	}

	return registerAppenderInternal(appenderType, appenderCreator)
}

// Removes a registered creator of a configuration for an appender. Build in ones will not be removed
func DeregisterAppender(appenderType string) error {
	appenderType = strings.ToUpper(appenderType)

	if appenderType == config.APPENDER_STDOUT || appenderType == config.APPENDER_FILE {
		return fmt.Errorf("the build in appender creator for  %s is not deletable", appenderType)
	}
	if _, exists := registeredAppenders[appenderType]; !exists {
		return fmt.Errorf("there does not exists any creator of appender %s", appenderType)
	}

	loggerCreationMu.Lock()
	defer loggerCreationMu.Unlock()

	loggersInitialized = false

	delete(registeredAppenders, appenderType)
	return nil
}

// Creates the all relevant general logger from config elements
func createGeneralLoggers(conf *config.Config, loggerConfigMapping *map[string]*GeneralLogger, appenderConfigMapping *map[string]*appender.Appender) {
	for _, lc1 := range (*conf).Logger {
		alreadyCreated := false
		lc1FormatterId := (*getFormatterConfigForPackage(&lc1.GetCommon().PackageParameter, &conf.Formatter)).Id()
		lc1AppenderId := (*getAppenderConfigForPackage(&lc1.GetCommon().PackageParameter, &conf.Appender)).Id()

		for _, lc2 := range (*conf).Logger {
			if lc1.Id() == lc2.Id() &&
				lc1FormatterId == (*getFormatterConfigForPackage(&lc2.GetCommon().PackageParameter, &conf.Formatter)).Id() &&
				lc1AppenderId == (*getAppenderConfigForPackage(&lc2.GetCommon().PackageParameter, &conf.Appender)).Id() {

				_, alreadyCreated = (*loggerConfigMapping)[lc1.Id()+lc1AppenderId+lc1FormatterId]
				break

			}
		}

		if alreadyCreated {
			continue
		}

		appender := (*appenderConfigMapping)[lc1AppenderId+lc1FormatterId]
		// There exists only GeneralLoggerConfig for interface LoggerConfig -> cast without check
		cLoggers = append(cLoggers, CreateGeneralLoggerFromConfig(lc1.(config.GeneralLoggerConfig), appender))
		(*loggerConfigMapping)[lc1.Id()+lc1AppenderId+lc1FormatterId] = &cLoggers[len(cLoggers)-1]
	}
}

// Creates the main logger from config elements
func createMainLogger(conf *config.Config, loggerConfigMapping *map[string]*GeneralLogger) {
	mLogger = MainLogger{}
	mLogger.existPackageLogger = len(conf.Logger) > 1
	mLogger.useFullQualifiedPackage = conf.UseFullQualifiedPackageName
	mLogger.packageLoggers = make(map[string]*GeneralLogger, len(conf.Logger)-1)

	for _, lc := range conf.Logger {
		lc1FormatterId := (*getFormatterConfigForPackage(&lc.GetCommon().PackageParameter, &conf.Formatter)).Id()
		lc1AppenderId := (*getAppenderConfigForPackage(&lc.GetCommon().PackageParameter, &conf.Appender)).Id()
		if lc.IsDefault() {
			mLogger.generalLogger = (*loggerConfigMapping)[lc.Id()+lc1AppenderId+lc1FormatterId]
		} else {
			mLogger.packageLoggers[lc.PackageName()] = (*loggerConfigMapping)[lc.Id()+lc1AppenderId+lc1FormatterId]
		}
	}

	loggersInitialized = true
}

// returns a pointer to the formatter config for a given package
func getFormatterConfigForPackage(PackageParameter *string, formatterConfig *[]config.FormatterConfig) *config.FormatterConfig {
	for i, fc := range *formatterConfig {
		if fc.PackageParameter() == *PackageParameter {
			return &(*formatterConfig)[i]
		}
	}
	return nil
}

// returns a pointer to the appender config for a given package
func getAppenderConfigForPackage(PackageParameter *string, appenderConfig *[]config.AppenderConfig) *config.AppenderConfig {
	for i, ac := range *appenderConfig {
		if ac.PackageParameter() == *PackageParameter {
			return &(*appenderConfig)[i]
		}
	}
	return nil
}
