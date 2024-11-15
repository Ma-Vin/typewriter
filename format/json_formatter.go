package format

import (
	"bytes"
	"fmt"
)

type JsonFormatter struct {
	timeKey                  string
	severityKey              string
	messageKey               string
	correlationKey           string
	customValuesKey          string
	timeLayout               string
	customValuesAsSubElement bool
}

// Creates a new formater with given key names and time layout
func CreateJsonFormatter(timeKey string, severityKey string, messageKey string, correlationKey string,
	customValuesKey string, timeLayout string, customValuesAsSubElement bool) Formatter {

	return JsonFormatter{
		timeKey:                  timeKey,
		severityKey:              severityKey,
		messageKey:               messageKey,
		correlationKey:           correlationKey,
		customValuesKey:          customValuesKey,
		timeLayout:               timeLayout,
		customValuesAsSubElement: customValuesAsSubElement,
	}
}

func (j JsonFormatter) Format(severity int, message string) string {
	return fmt.Sprintf("{ \"%s\": \"%s\", \"%s\": \"%s\", \"%s\": \"%s\" }",
		j.timeKey, getNowAsStringFromLayout(j.timeLayout),
		j.severityKey, severityTrimTextMap[severity],
		j.messageKey, message)
}

func (j JsonFormatter) FormatWithCorrelation(severity int, correlationId string, message string) string {
	return fmt.Sprintf("{ \"%s\": \"%s\", \"%s\": \"%s\", \"%s\": \"%s\", \"%s\": \"%s\" }",
		j.timeKey, getNowAsStringFromLayout(j.timeLayout),
		j.severityKey, severityTrimTextMap[severity],
		j.correlationKey, correlationId,
		j.messageKey, message)
}

func (j JsonFormatter) FormatCustom(severity int, message string, customValues map[string]any) string {

	if j.customValuesAsSubElement {
		return fmt.Sprintf("{ \"%s\": \"%s\", \"%s\": \"%s\", \"%s\": \"%s\", \"%s\": { %s } }",
			j.timeKey, getNowAsStringFromLayout(j.timeLayout),
			j.severityKey, severityTrimTextMap[severity],
			j.messageKey, message,
			j.customValuesKey, formatCustomValuesToJson(&customValues))
	}

	return fmt.Sprintf("{ \"%s\": \"%s\", \"%s\": \"%s\", \"%s\": \"%s\", %s }",
		j.timeKey, getNowAsStringFromLayout(j.timeLayout),
		j.severityKey, severityTrimTextMap[severity],
		j.messageKey, message,
		formatCustomValuesToJson(&customValues))
}

func formatCustomValuesToJson(customValues *map[string]any) string {
	var customValuesBuffer bytes.Buffer
	addComma := false
	for key, value := range *customValues {
		if addComma {
			customValuesBuffer.WriteString(", ")
		}
		customValuesBuffer.WriteString(fmt.Sprintf("\"%s\": ", key))
		switch value.(type) {
		case bool:
			customValuesBuffer.WriteString(fmt.Sprintf("%t", value))
		case byte:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case int:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case int8:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case int16:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case int32:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case int64:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case uint:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		// case uint8: equal to byte
		case uint16:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case uint32:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case uint64:
			customValuesBuffer.WriteString(fmt.Sprintf("%d", value))
		case float32:
			customValuesBuffer.WriteString(fmt.Sprintf("%g", value))
		case float64:
			customValuesBuffer.WriteString(fmt.Sprintf("%g", value))
		case complex64:
			customValuesBuffer.WriteString(fmt.Sprintf("%g", value))
		case complex128:
			customValuesBuffer.WriteString(fmt.Sprintf("%g", value))
		case string:
			customValuesBuffer.WriteString(fmt.Sprintf("\"%s\"", value))
		default:
			customValuesBuffer.WriteString(fmt.Sprintf("\"%v\"", value))
		}
		addComma = true
	}

	return customValuesBuffer.String()
}
