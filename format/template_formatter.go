package format

import (
	"fmt"
	"sort"
)

// Formatter which uses given different templates for Format, FormatWithCorrelation and FormatCustom
// The templates will be used in combination with [fmt.Sprintf]
//
// The arguments of the templates will be passed in the following order (paramter skipped if not relevant)
//
// 1. time (derived by [time.Format] with timeLayout parameter)
//
// 2. severity
//
// 3. correlationId
//
// 4. message
//
// 5. custom values as key-value pairs sorted by key
//
// Because of explicit argument indices can be used at templates
type TemplateFormatter struct {
	template              string
	correlationIdTemplate string
	customTemplate        string
	timeLayout            string
}

// Creates a new formater with a given delimiter
func CreateTemplateFormatter(template string, correlationIdTemplate string, customTemplate string, timeLayout string) Formatter {
	return TemplateFormatter{
		template:              template,
		correlationIdTemplate: correlationIdTemplate,
		customTemplate:        customTemplate,
		timeLayout:            timeLayout,
	}
}

func (t TemplateFormatter) Format(severity int, message string) string {
	return formatValues(t.template, getNowAsStringFromLayout(t.timeLayout), severityTextMap[severity], message)
}

func (t TemplateFormatter) FormatWithCorrelation(severity int, correlationId string, message string) string {
	return formatValues(t.correlationIdTemplate, getNowAsStringFromLayout(t.timeLayout), severityTextMap[severity], correlationId, message)
}

func (t TemplateFormatter) FormatCustom(severity int, message string, customValues map[string]any) string {
	args := make([]any, 0, 2*len(customValues)+3)

	args = append(args, getNowAsStringFromLayout(t.timeLayout))
	args = append(args, severityTextMap[severity])
	args = append(args, message)
	args = appendCustomValues(args, &customValues)

	return formatValues(t.customTemplate, args...)
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
