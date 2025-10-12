package format

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

// Formatter which append given parameter with a delimiter. Since name the of the parameter will not be contained, the keys of customValues at FormatCustom neither.
type DelimiterFormatter struct {
	delimiter        string
	timeLayout       string
	isSequenceActive bool
}

// Creates a new formatter from a given config
func CreateDelimiterFormatterFromConfig(formatterConfig *config.FormatterConfig) (*Formatter, error) {
	delimiterFormatterConfig, ok := (*formatterConfig).(config.DelimiterFormatterConfig)
	if !ok {
		return nil, fmt.Errorf("failed to convert interface to DelimiterFormatterConfig for formatter %s", (*formatterConfig).FormatterType())
	}

	var result Formatter = DelimiterFormatter{delimiter: delimiterFormatterConfig.Delimiter, timeLayout: delimiterFormatterConfig.TimeLayout(), isSequenceActive: delimiterFormatterConfig.Common.IsSequenceActive}
	return &result, nil
}

// Formats the given parameter to a string to log
func (d DelimiterFormatter) Format(logValues *common.LogValues) string {
	var sb strings.Builder

	sb.WriteString(logValues.Time.Format(d.timeLayout))
	if d.isSequenceActive {
		sb.WriteString(d.delimiter)
		sb.WriteString(strconv.FormatUint(logValues.Sequence, 10))
	}
	sb.WriteString(d.delimiter)
	sb.WriteString(severityTextMap[logValues.Severity])

	if logValues.CorrelationId != nil {
		sb.WriteString(d.delimiter)
		sb.WriteString(*logValues.CorrelationId)
	}

	if logValues.IsCallerSet {
		sb.WriteString(d.delimiter)
		sb.WriteString(fmt.Sprintf("%s at %s (Line %d)", logValues.CallerFunction, logValues.CallerFile, logValues.CallerFileLine))
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
