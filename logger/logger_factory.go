package logger

import (
	"fmt"
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

		switch fc1.FormatterType() {
		case config.FORMATTER_DELIMITER:
			appendFormatter(format.CreateDelimiterFormatterFromConfig(fc1.(config.DelimiterFormatterConfig)))
			setLastFormatter(fc1.Id(), formatterConfigMapping)
		case config.FORMATTER_TEMPLATE:
			appendFormatter(format.CreateTemplateFormatterFromConfig(fc1.(config.TemplateFormatterConfig)))
			setLastFormatter(fc1.Id(), formatterConfigMapping)
		case config.FORMATTER_JSON:
			appendFormatter(format.CreateJsonFormatterFromConfig(fc1.(config.JsonFormatterConfig)))
			setLastFormatter(fc1.Id(), formatterConfigMapping)
		default:
			// not relevant: handled at config load
		}
	}
}

func appendFormatter(formatter format.Formatter) {
	formatters = append(formatters, formatter)
}

func setLastFormatter(formatterId string, formatterConfigMapping *map[string]*format.Formatter) {
	(*formatterConfigMapping)[formatterId] = &formatters[len(formatters)-1]
}

// Creates the all relevant appenders from config elements
func createAppenders(conf *config.Config, appenderConfigMapping *map[string]*appender.Appender, formatterConfigMapping *map[string]*format.Formatter) {
	appender.CleanFileDeductions()
	for _, ac1 := range conf.Appender {
		alreadyCreated := false
		ac1FormatterId := (*getFormatterConfigForPackage(&ac1.GetCommon().PackageParameter, &conf.Formatter)).Id()

		for _, ac2 := range conf.Appender {
			if ac1.Id() == ac2.Id() && ac1FormatterId == (*getFormatterConfigForPackage(&ac2.GetCommon().PackageParameter, &conf.Formatter)).Id() {

				_, alreadyCreated = (*appenderConfigMapping)[ac1.Id()+ac1FormatterId]
				break
			}
		}

		if alreadyCreated {
			continue
		}

		result := createAppender(ac1, (*formatterConfigMapping)[ac1FormatterId])
		appendAppender(*result)
		setLastAppender(ac1.Id()+ac1FormatterId, appenderConfigMapping)
	}
}

func createAppender(appenderConfig config.AppenderConfig, formatter *format.Formatter) *appender.Appender {
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
