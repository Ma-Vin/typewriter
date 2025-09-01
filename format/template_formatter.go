package format

import (
	"fmt"
	"sort"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

const DEFAULT_TEMPLATE = config.DEFAULT_TEMPLATE
const DEFAULT_CALLER_TEMPLATE = config.DEFAULT_CALLER_TEMPLATE

// Formatter which uses given different templates for Format, FormatWithCorrelation and FormatCustom
// The templates will be used in combination with [fmt.Sprintf]
//
// The arguments of the templates will be passed in the following order (parameter skipped if not relevant)
//
//  1. time (derived by [time.Format] with timeLayout parameter)
//  2. severity
//  3. correlationId
//  4. message
//  5. custom values as key-value pairs sorted by key
//
// Because of explicit argument indices can be used at templates
type TemplateFormatter struct {
	template                    string
	callerTemplate              string
	correlationIdTemplate       string
	callerCorrelationIdTemplate string
	customTemplate              string
	callerCustomTemplate        string
	timeLayout                  string
	trimSeverityText            bool
}

// Creates a new formatter from a given config
func CreateTemplateFormatterFromConfig(formatterConfig *config.FormatterConfig) (*Formatter, error) {
	templateFormatterConfig, ok := (*formatterConfig).(config.TemplateFormatterConfig)
	if !ok {
		return nil, fmt.Errorf("failed to convert interface to TemplateFormatterConfig for formatter %s", (*formatterConfig).FormatterType())
	}

	var result Formatter = TemplateFormatter{
		template:                    templateFormatterConfig.Template,
		callerTemplate:              templateFormatterConfig.CallerTemplate,
		correlationIdTemplate:       templateFormatterConfig.CorrelationIdTemplate,
		callerCorrelationIdTemplate: templateFormatterConfig.CallerCorrelationIdTemplate,
		customTemplate:              templateFormatterConfig.CustomTemplate,
		callerCustomTemplate:        templateFormatterConfig.CallerCustomTemplate,
		timeLayout:                  templateFormatterConfig.TimeLayout(),
		trimSeverityText:            templateFormatterConfig.TrimSeverityText,
	}

	return &result, nil
}

// Formats the given parameter to a string to log
func (t TemplateFormatter) Format(logValues *common.LogValues) string {
	if logValues.CustomValues != nil {
		return t.formatCustom(logValues)
	}
	if logValues.IsCallerSet {
		if logValues.CorrelationId != nil {
			return formatValues(t.callerCorrelationIdTemplate, t.formatTime(logValues), t.getSeverityText(logValues.Severity), *logValues.CorrelationId,
				logValues.CallerFunction, logValues.CallerFile, logValues.CallerFileLine, logValues.Message)
		}
		return formatValues(t.callerTemplate, t.formatTime(logValues), t.getSeverityText(logValues.Severity),
			logValues.CallerFunction, logValues.CallerFile, logValues.CallerFileLine, logValues.Message)
	}
	if logValues.CorrelationId != nil {
		return formatValues(t.correlationIdTemplate, t.formatTime(logValues), t.getSeverityText(logValues.Severity), *logValues.CorrelationId, logValues.Message)
	}
	return formatValues(t.template, t.formatTime(logValues), t.getSeverityText(logValues.Severity), logValues.Message)
}

// Formats the given parameter to a string to log and the customValues will be added at the end
func (t TemplateFormatter) formatCustom(logValues *common.LogValues) string {
	if t.customTemplate == DEFAULT_TEMPLATE {
		for i := 0; i < len(*logValues.CustomValues); i++ {
			t.customTemplate += " [%s]: %v"
		}
	}
	args := make([]any, 0, 2*len(*logValues.CustomValues)+3)

	args = append(args, t.formatTime(logValues))
	args = append(args, t.getSeverityText(logValues.Severity))
	if logValues.IsCallerSet {
		args = append(args, logValues.CallerFunction)
		args = append(args, logValues.CallerFile)
		args = append(args, logValues.CallerFileLine)
	}
	args = append(args, logValues.Message)
	args = appendCustomValues(args, logValues.CustomValues)

	if logValues.IsCallerSet {
		return formatValues(t.callerCustomTemplate, args...)
	}
	return formatValues(t.customTemplate, args...)
}

func (t *TemplateFormatter) formatTime(logValues *common.LogValues) string {
	return logValues.Time.Format(t.timeLayout)
}

func (t *TemplateFormatter) getSeverityText(severity int) string {
	if t.trimSeverityText {
		return severityTrimTextMap[severity]
	}
	return severityTextMap[severity]
}

func formatValues(template string, args ...any) string {
	return fmt.Sprintf(template, args...)
}

func appendCustomValues(args []any, customValues *map[string]any) []any {
	var keys []string
	for key := range *customValues {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		args = append(args, key)
		args = append(args, (*customValues)[key])
	}
	return args
}
