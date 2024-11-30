package format

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ma-vin/typewriter/common"
)

// Formatter which append given parameter with a delimter. Since name the of the parameter will not be contained, the keys of customValues at FormatCustom neither.
type DelimiterFormatter struct {
	delimiter string
}

// Creates a new formater with a given delimiter
func CreateDelimiterFormatter(delimiter string) Formatter {
	return DelimiterFormatter{delimiter}
}

// Formats the given parameter to a string to log
func (d DelimiterFormatter) Format(logValues *common.LogValues) string {
	var sb strings.Builder
	sb.WriteString(logValues.Time.Format(time.RFC3339))
	sb.WriteString(d.delimiter)
	sb.WriteString(severityTextMap[logValues.Severity])
	if logValues.CorrelationId != nil {
		sb.WriteString(d.delimiter)
		sb.WriteString(*logValues.CorrelationId)
	}
	sb.WriteString(d.delimiter)
	sb.WriteString(logValues.Message)
	if logValues.CustomValues != nil && len(*logValues.CustomValues) > 0 {
		sb.WriteString(d.delimiter)
		sb.WriteString(formatMapToString(logValues.CustomValues, &d.delimiter))
	}
	return sb.String()
}

func formatMapToString(customValues *map[string]any, delimiter *string) string {
	var keys []string
	for key := range *customValues {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, key := range keys {
		if sb.Len() > 0 {
			sb.WriteString(*delimiter)
		}
		sb.WriteString(fmt.Sprint((*customValues)[key]))
	}
	return sb.String()
}
