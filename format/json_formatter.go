package format

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Formats the log entries as JSON
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

// Formats the given parameter to a string to log
func (j JsonFormatter) Format(severity int, message string) string {
	return j.formatJsonEntriesMap(j.createJsonEntriesMap(3, severity, "", message))
}

// Formats the given default parameter and a correlation id to a string to log
func (j JsonFormatter) FormatWithCorrelation(severity int, correlationId string, message string) string {
	return j.formatJsonEntriesMap(j.createJsonEntriesMap(4, severity, correlationId, message))
}

// Formats the given parameter to a string to log and the customValues will be added
func (j JsonFormatter) FormatCustom(severity int, message string, customValues map[string]any) string {
	var jsonEntries *map[string]any
	if j.customValuesAsSubElement {
		jsonEntries = j.createJsonEntriesMap(4, severity, "", message)
		(*jsonEntries)[j.customValuesKey] = customValues
	} else {
		jsonEntries = j.createJsonEntriesMap(3+len(customValues), severity, "", message)
		for k, v := range customValues {
			(*jsonEntries)[k] = v
		}
	}

	return j.formatJsonEntriesMap(jsonEntries)
}

func (j *JsonFormatter) formatJsonEntriesMap(jsonEntries *map[string]any) string {
	jsonByteArray, err := json.Marshal(jsonEntries)

	if err != nil {
		return j.formatWithError(jsonEntries, err)
	}

	return string(jsonByteArray)
}

func (j *JsonFormatter) createJsonEntriesMap(capacity int, severity int, correlationId string, message string) *map[string]any {
	result := make(map[string]any, capacity)

	if correlationId != "" {
		result[j.correlationKey] = correlationId
	}

	result[j.timeKey] = getNowAsStringFromLayout(j.timeLayout)
	result[j.severityKey] = severityTrimTextMap[severity]
	result[j.messageKey] = message

	return &result
}

func (j *JsonFormatter) formatWithError(jsonEntries *map[string]any, err error) string {
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "{")

	correlationId, correlationFound := (*jsonEntries)[j.correlationKey]
	if correlationFound {
		fmt.Fprintf(buf, "\"%s\":\"%s\",", j.correlationKey, correlationId)

	}
	fmt.Fprintf(buf, "\"%s\":\"%s: %v\",\"%s\":\"%s\",\"%s\":\"%s\"}",
		j.messageKey, "Fail to marshal to json, use custom formatter", err,
		j.severityKey, ERROR_PREFIX,
		j.timeKey, (*jsonEntries)[j.timeKey])

	fmt.Fprintln(buf)

	fmt.Fprint(buf, "{")
	if correlationFound {
		fmt.Fprintf(buf, "\"%s\":\"%s\",", j.correlationKey, correlationId)

	}
	fmt.Fprintf(buf, "\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}",
		j.messageKey, (*jsonEntries)[j.messageKey],
		j.severityKey, (*jsonEntries)[j.severityKey],
		j.timeKey, (*jsonEntries)[j.timeKey])

	return buf.String()
}
