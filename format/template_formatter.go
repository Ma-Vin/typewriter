package format

import (
	"fmt"
	"sort"

	"github.com/ma-vin/typewriter/common"
)

const DEFAULT_TEMPLATE = "[%s] %s: %s"

// Formatter which uses given different templates for Format, FormatWithCorrelation and FormatCustom
// The templates will be used in combination with [fmt.Sprintf]
//
// The arguments of the templates will be passed in the following order (paramter skipped if not relevant)
//
//  1. time (derived by [time.Format] with timeLayout parameter)
//  2. severity
//  3. correlationId
//  4. message
//  5. custom values as key-value pairs sorted by key
//
// Because of explicit argument indices can be used at templates
type TemplateFormatter struct {
	template              string
	correlationIdTemplate string
	customTemplate        string
	timeLayout            string
	trimSeverityText      bool
}

// Creates a new formater with given templates and time layout
func CreateTemplateFormatter(template string, correlationIdTemplate string, customTemplate string, timeLayout string, trimSeverityText bool) Formatter {
	return TemplateFormatter{
		template:              template,
		correlationIdTemplate: correlationIdTemplate,
		customTemplate:        customTemplate,
		timeLayout:            timeLayout,
		trimSeverityText:      trimSeverityText,
	}
}

// Formats the given parameter to a string to log
func (t TemplateFormatter) Format(logValues *common.LogValues) string {
	if logValues.CorrelationId != nil {
		return formatValues(t.correlationIdTemplate, t.formatTime(logValues), t.getSeverityText(logValues.Severity), *logValues.CorrelationId, logValues.Message)
	}
	if logValues.CustomValues != nil {
		return t.formatCustom(logValues)
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
	args = append(args, logValues.Message)
	args = appendCustomValues(args, logValues.CustomValues)

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
