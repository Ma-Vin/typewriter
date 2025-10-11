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
	template := *t.determineTemplate(logValues)
	args := *t.createSliceFromLogValues(logValues)
	return formatValues(template, args...)
}

// Creates a pointer to a slice with all relevant log entries. The values will be added as key-value pairs
func (t TemplateFormatter) createSliceFromLogValues(logValues *common.LogValues) *[]any {
	capacity := 3
	if logValues.IsCallerSet {
		capacity += 3
	}
	if logValues.CorrelationId != nil {
		capacity += 1
	}
	if logValues.CustomValues != nil {
		capacity += 2 * len(*logValues.CustomValues)
	}

	args := make([]any, 0, capacity)

	args = append(args, t.formatTime(logValues))
	args = append(args, t.getSeverityText(logValues.Severity))
	if logValues.CorrelationId != nil {
		args = append(args, *logValues.CorrelationId)
	}
	if logValues.IsCallerSet {
		args = append(args, logValues.CallerFunction)
		args = append(args, logValues.CallerFile)
		args = append(args, logValues.CallerFileLine)
	}
	args = append(args, logValues.Message)
	if logValues.CustomValues != nil {
		args = appendCustomValues(args, logValues.CustomValues)
	}

	return &args
}

// Determines the template to use custom, with caller, correlation or default one
func (t TemplateFormatter) determineTemplate(logValues *common.LogValues) *string {
	if logValues.CustomValues != nil {
		return t.determineCustomTemplate(logValues)
	} else {
		if logValues.IsCallerSet {
			if logValues.CorrelationId != nil {
				return &t.callerCorrelationIdTemplate
			}
			return &t.callerTemplate
		}
		if logValues.CorrelationId != nil {
			return &t.correlationIdTemplate
		}
		return &t.template
	}
}

// Determines the custom template
func (t TemplateFormatter) determineCustomTemplate(logValues *common.LogValues) *string {
	var result string
	if logValues.IsCallerSet {
		result = t.callerCustomTemplate
		if result == DEFAULT_CALLER_TEMPLATE {
			for i := 0; i < len(*logValues.CustomValues); i++ {
				result += " [%s]: %v"
			}
		}
	} else {
		result = t.customTemplate
		if result == DEFAULT_TEMPLATE {
			for i := 0; i < len(*logValues.CustomValues); i++ {
				result += " [%s]: %v"
			}
		}
	}
	return &result
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
