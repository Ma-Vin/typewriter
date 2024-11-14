package format

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Formatter which append given parameter with a delimter. Since name the of the parameter meter will not be contained, the keys of customValues at FormatCustom neither.
type DelimiterFormatter struct {
	delimiter string
}

// Creates a new formater with a given delimiter
func CreateDelimiterFormatter(delimiter string) Formatter {
	return DelimiterFormatter{delimiter}
}

func (d DelimiterFormatter) Format(severity int, message string) string {
	return concatWithDelimiter(&d.delimiter, getNowAsString(), severityTextMap[severity], message)
}

func (d DelimiterFormatter) FormatWithCorrelation(severity int, correlationId string, message string) string {
	return concatWithDelimiter(&d.delimiter, getNowAsString(), severityTextMap[severity], correlationId, message)
}

func (d DelimiterFormatter) FormatCustom(severity int, message string, customValues map[string]any) string {
	return concatWithDelimiter(&d.delimiter, getNowAsString(), severityTextMap[severity], message, formatMapToString(&customValues, &d.delimiter))
}

func concatWithDelimiter(delimiter *string, args ...string) string {
	var sb strings.Builder
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(*delimiter)
		}
		sb.WriteString(arg)
	}

	return sb.String()
}

func getNowAsString() string {
	return getNowAsStringFromLayout(time.RFC3339)
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
