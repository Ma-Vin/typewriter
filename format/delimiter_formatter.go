package format

import (
	"fmt"
	"sort"
	"strings"
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
func (d DelimiterFormatter) Format(severity int, message string) string {
	return concatWithDelimiter(&d.delimiter, getNowAsStringDefaultLayout(), severityTextMap[severity], message)
}

// Formats the given default parameter and a correlation id to a string to log
func (d DelimiterFormatter) FormatWithCorrelation(severity int, correlationId string, message string) string {
	return concatWithDelimiter(&d.delimiter, getNowAsStringDefaultLayout(), severityTextMap[severity], correlationId, message)
}

// Formats the given parameter to a string to log and the customValues will be added at the end
func (d DelimiterFormatter) FormatCustom(severity int, message string, customValues map[string]any) string {
	return concatWithDelimiter(&d.delimiter, getNowAsStringDefaultLayout(), severityTextMap[severity], message, formatMapToString(&customValues, &d.delimiter))
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
