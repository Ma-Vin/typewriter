package logger

import (
	"github.com/ma-vin/typewriter/appender"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
)

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

	conf := config.GetConfig()
	if conf == nil {
		return nil
	}

	formatterConfigMapping := make(map[string]*format.Formatter, len(conf.Formatter))
	appenderConfigMapping := make(map[string]*appender.Appender, len(conf.Appender))
	loggerConfigMapping := make(map[string]*CommonLogger, len(conf.Logger))

	createFormatters(&conf.Formatter, &formatterConfigMapping)
	createAppenders(conf, &appenderConfigMapping, &formatterConfigMapping)
	createCommonLoggers(conf, &loggerConfigMapping, &appenderConfigMapping)
	createMainLogger(conf, &loggerConfigMapping)

	return &mLogger
}

// Creates the all relevant fromatters from config elements
func createFormatters(formatterConfigs *[]config.FormatterConfig, formatterConfigMapping *map[string]*format.Formatter) {
	for _, fc1 := range *formatterConfigs {
		alreadyCreated := false

		for _, fc2 := range *formatterConfigs {
			if fc1.Id == fc2.Id {
				_, alreadyCreated = (*formatterConfigMapping)[fc1.Id]
				break
			}
		}

		if alreadyCreated {
			continue
		}

		switch fc1.FormatterType {
		case config.FORMATTER_DELIMITER:
			appendFormatter(format.CreateDelimiterFormatter(fc1.Delimiter))
			setLastFormatter(fc1.Id, formatterConfigMapping)
		case config.FORMATTER_TEMPLATE:
			appendFormatter(format.CreateTemplateFormatter(fc1.Template, fc1.CorrelationIdTemplate, fc1.CustomTemplate, 
				fc1.CallerTemplate, fc1.CallerCorrelationIdTemplate, fc1.CallerCustomTemplate, 
				fc1.TimeLayout, fc1.TrimSeverityText))
			setLastFormatter(fc1.Id, formatterConfigMapping)
		case config.FORMATTER_JSON:
			appendFormatter(format.CreateJsonFormatter(fc1.TimeKey, fc1.SeverityKey, fc1.MessageKey, fc1.CorrelationKey, fc1.CustomValuesKey, fc1.TimeLayout,
				fc1.CallerFunctionKey, fc1.CallerFileKey, fc1.CallerFileLineKey, fc1.CustomValuesAsSubElement))
			setLastFormatter(fc1.Id, formatterConfigMapping)
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
	appender.CleanFileToMutex()
	for _, ac1 := range conf.Appender {
		alreadyCreated := false
		ac1FormatterId := getFormatterConfigForPackage(&ac1.PackageName, &conf.Formatter).Id

		for _, ac2 := range conf.Appender {
			if ac1.Id == ac2.Id && ac1FormatterId == getFormatterConfigForPackage(&ac2.PackageName, &conf.Formatter).Id {

				_, alreadyCreated = (*appenderConfigMapping)[ac1.Id+ac1FormatterId]
				break
			}
		}

		if alreadyCreated {
			continue
		}

		formatter := (*formatterConfigMapping)[ac1FormatterId]
		switch ac1.AppenderType {
		case config.APPENDER_STDOUT:
			appendAppender(appender.CreateStandardOutputAppender(formatter))
			setLastAppender(ac1.Id+ac1FormatterId, appenderConfigMapping)
		case config.APPENDER_FILE:
			appendAppender(appender.CreateFileAppender(ac1.PathToLogFile, formatter))
			setLastAppender(ac1.Id+ac1FormatterId, appenderConfigMapping)
		default:
			// not relevant: handled at config load
		}
	}
}

func appendAppender(appender appender.Appender) {
	appenders = append(appenders, appender)
}

func setLastAppender(appenderId string, appenderConfigMapping *map[string]*appender.Appender) {
	(*appenderConfigMapping)[appenderId] = &appenders[len(appenders)-1]
}

// Creates the all relevant common logger from config elements
func createCommonLoggers(conf *config.Config, loggerConfigMapping *map[string]*CommonLogger, appenderConfigMapping *map[string]*appender.Appender) {
	for _, lc1 := range (*conf).Logger {
		alreadyCreated := false
		lc1FormatterId := getFormatterConfigForPackage(&lc1.PackageName, &conf.Formatter).Id
		lc1AppenderId := getAppenderConfigForPackage(&lc1.PackageName, &conf.Appender).Id

		for _, lc2 := range (*conf).Logger {
			if lc1.Id == lc2.Id &&
				lc1FormatterId == getAppenderConfigForPackage(&lc2.PackageName, &conf.Appender).Id &&
				lc1AppenderId == getFormatterConfigForPackage(&lc2.PackageName, &conf.Formatter).Id {

				_, alreadyCreated = (*loggerConfigMapping)[lc1.Id+lc1AppenderId+lc1FormatterId]
				break

			}
		}

		if alreadyCreated {
			continue
		}

		appender := (*appenderConfigMapping)[lc1AppenderId+lc1FormatterId]
		cLoggers = append(cLoggers, CreateCommonLogger(appender, lc1.Severity, lc1.IsCallerToSet))
		(*loggerConfigMapping)[lc1.Id+lc1AppenderId+lc1FormatterId] = &cLoggers[len(cLoggers)-1]
	}
}

// Creates the main logger from config elements
func createMainLogger(conf *config.Config, loggerConfigMapping *map[string]*CommonLogger) {
	mLogger = MainLogger{}
	mLogger.existPackageLogger = len(conf.Logger) > 1
	mLogger.packageLoggers = make(map[string]*CommonLogger, len(conf.Logger)-1)

	for _, lc := range conf.Logger {
		lc1FormatterId := getFormatterConfigForPackage(&lc.PackageName, &conf.Formatter).Id
		lc1AppenderId := getAppenderConfigForPackage(&lc.PackageName, &conf.Appender).Id
		if lc.IsDefault {
			mLogger.commonLogger = (*loggerConfigMapping)[lc.Id+lc1AppenderId+lc1FormatterId]
		} else {
			mLogger.packageLoggers[lc.PackageName] = (*loggerConfigMapping)[lc.Id+lc1AppenderId+lc1FormatterId]
		}
	}

	loggersInitialized = true
}

// returns a pointer to the formatter config for a given package
func getFormatterConfigForPackage(packageName *string, formatterConfig *[]config.FormatterConfig) *config.FormatterConfig {
	for i, fc := range *formatterConfig {
		if fc.PackageName == *packageName {
			return &(*formatterConfig)[i]
		}
	}
	return nil
}

// returns a pointer to the appender config for a given package
func getAppenderConfigForPackage(packageName *string, appenderConfig *[]config.AppenderConfig) *config.AppenderConfig {
	for i, ac := range *appenderConfig {
		if ac.PackageName == *packageName {
			return &(*appenderConfig)[i]
		}
	}
	return nil
}
