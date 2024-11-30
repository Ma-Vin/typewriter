package format

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ma-vin/typewriter/common"
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
func (j JsonFormatter) Format(logValues *common.LogValues) string {
	jsonEntries := make(map[string]any, j.getJsonEntriesMapCapacity(logValues))

	jsonEntries[j.timeKey] = logValues.Time.Format(j.timeLayout)
	jsonEntries[j.severityKey] = severityTrimTextMap[logValues.Severity]
	jsonEntries[j.messageKey] = logValues.Message

	if logValues.CorrelationId != nil {
		jsonEntries[j.correlationKey] = *logValues.CorrelationId
	}

	if logValues.CustomValues != nil {
		if j.customValuesAsSubElement {
			jsonEntries[j.customValuesKey] = *logValues.CustomValues
		} else {
			for k, v := range *logValues.CustomValues {
				jsonEntries[k] = v
			}
		}
	}

	return j.formatJsonEntriesMap(&jsonEntries)
}

func (j *JsonFormatter) getJsonEntriesMapCapacity(logValues *common.LogValues) int {
	result := 3
	
	if logValues.CorrelationId != nil {
		result++
	}

	if logValues.CustomValues != nil {
		if j.customValuesAsSubElement {
			result++
		} else {
			result += len(*logValues.CustomValues)
		}
	}

	return result
}


func (j *JsonFormatter) formatJsonEntriesMap(jsonEntries *map[string]any) string {
	jsonByteArray, err := json.Marshal(jsonEntries)

	if err != nil {
		return j.formatWithError(jsonEntries, err)
	}

	return string(jsonByteArray)
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
