package format

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

// template formatter placeholder
const (
	PLACEHOLDER_CALLER_FUNCTION  = "$func"
	PLACEHOLDER_CALLER_FILE      = "$file"
	PLACEHOLDER_CALLER_FILE_LINE = "$line"
	PLACEHOLDER_CORRELATION      = "$corr"
	PLACEHOLDER_CUSTOM_VALUES    = "$cust"
	PLACEHOLDER_ENV_VALUES       = "$env"
	PLACEHOLDER_MSG              = "$msg"
	PLACEHOLDER_SEQUENCE         = "$seq"
	PLACEHOLDER_SEVERITY         = "$sev"
	PLACEHOLDER_TIME             = "$time"
)

// map with replacement verbs for template formatter placeholders
var replacements map[string]string = map[string]string{
	PLACEHOLDER_CALLER_FUNCTION:  "s",
	PLACEHOLDER_CALLER_FILE:      "s",
	PLACEHOLDER_CALLER_FILE_LINE: "d",
	PLACEHOLDER_CORRELATION:      "s",
	PLACEHOLDER_MSG:              "s",
	PLACEHOLDER_SEQUENCE:         "d",
	PLACEHOLDER_SEVERITY:         "s",
	PLACEHOLDER_TIME:             "s",
}

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
	commonProperties                     *CommonFormatterProperties
	template                             string
	isDefaultTemplate                    bool
	callerTemplate                       string
	isDefaultCallerTemplate              bool
	correlationIdTemplate                string
	isDefaultCorrelationIdTemplate       bool
	callerCorrelationIdTemplate          string
	isDefaultCallerCorrelationIdTemplate bool
	customTemplate                       string
	isDefaultCustomTemplate              bool
	callerCustomTemplate                 string
	isDefaultCallerCustomTemplate        bool
	timeLayout                           string
	isSequenceActive                     bool
	trimSeverityText                     bool
	envNameAndValuesToLog                []any
}

// Creates a new formatter from a given config
func CreateTemplateFormatterFromConfig(formatterConfig *config.FormatterConfig) (*Formatter, error) {
	templateFormatterConfig, ok := (*formatterConfig).(config.TemplateFormatterConfig)
	if !ok {
		return nil, fmt.Errorf("failed to convert interface to TemplateFormatterConfig for formatter %s", (*formatterConfig).FormatterType())
	}

	result := TemplateFormatter{
		commonProperties:                     CreateCommonFormatterProperties(templateFormatterConfig.Common),
		template:                             templateFormatterConfig.Template,
		isDefaultTemplate:                    templateFormatterConfig.IsDefaultTemplate,
		callerTemplate:                       templateFormatterConfig.CallerTemplate,
		isDefaultCallerTemplate:              templateFormatterConfig.IsDefaultCallerTemplate,
		correlationIdTemplate:                templateFormatterConfig.CorrelationIdTemplate,
		isDefaultCorrelationIdTemplate:       templateFormatterConfig.IsDefaultCorrelationIdTemplate,
		callerCorrelationIdTemplate:          templateFormatterConfig.CallerCorrelationIdTemplate,
		isDefaultCallerCorrelationIdTemplate: templateFormatterConfig.IsDefaultCallerCorrelationIdTemplate,
		customTemplate:                       templateFormatterConfig.CustomTemplate,
		isDefaultCustomTemplate:              templateFormatterConfig.IsDefaultCustomTemplate,
		callerCustomTemplate:                 templateFormatterConfig.CallerCustomTemplate,
		isDefaultCallerCustomTemplate:        templateFormatterConfig.IsDefaultCallerCustomTemplate,
		timeLayout:                           templateFormatterConfig.TimeLayout(),
		trimSeverityText:                     templateFormatterConfig.TrimSeverityText,
		isSequenceActive:                     templateFormatterConfig.Common.IsSequenceActive,
		envNameAndValuesToLog:                make([]any, 0, 2*len(templateFormatterConfig.Common.EnvNamesToLog)),
	}

	for i, s := range result.commonProperties.envNamesToLog {
		result.envNameAndValuesToLog = append(result.envNameAndValuesToLog, s, result.commonProperties.envValuesToLog[i])
	}

	appendEnvNameValue(&result.template, &result.commonProperties.envNamesToLog, result.isDefaultTemplate)
	appendEnvNameValue(&result.callerTemplate, &result.commonProperties.envNamesToLog, result.isDefaultCallerTemplate)
	appendEnvNameValue(&result.correlationIdTemplate, &result.commonProperties.envNamesToLog, result.isDefaultCorrelationIdTemplate)
	appendEnvNameValue(&result.callerCorrelationIdTemplate, &result.commonProperties.envNamesToLog, result.isDefaultCallerCorrelationIdTemplate)
	appendEnvNameValue(&result.customTemplate, &result.commonProperties.envNamesToLog, result.isDefaultCustomTemplate)
	appendEnvNameValue(&result.callerCustomTemplate, &result.commonProperties.envNamesToLog, result.isDefaultCallerCustomTemplate)

	adjustTemplates(&result)

	var resultFormatter Formatter = result
	return &resultFormatter, nil
}

func appendEnvNameValue(template *string, envNamesToLog *[]string, isDefault bool) {
	if !isDefault {
		return
	}
	for i := 0; i < len(*envNamesToLog); i++ {
		*template += fmt.Sprintf(" ["+PLACEHOLDER_ENV_VALUES+"_k%[1]d]: "+PLACEHOLDER_ENV_VALUES+"_v%[1]d", i)
	}
}

// Adjust all templates at formatter by replacing template formatter placeholders with go known verbs
func adjustTemplates(formatter *TemplateFormatter) {
	if formatter.isSequenceActive {
		adjustTemplate(formatter, &formatter.template, PLACEHOLDER_TIME, PLACEHOLDER_SEQUENCE, PLACEHOLDER_SEVERITY, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.callerTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEQUENCE, PLACEHOLDER_SEVERITY, PLACEHOLDER_CALLER_FUNCTION, PLACEHOLDER_CALLER_FILE, PLACEHOLDER_CALLER_FILE_LINE, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.correlationIdTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEQUENCE, PLACEHOLDER_SEVERITY, PLACEHOLDER_CORRELATION, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.callerCorrelationIdTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEQUENCE, PLACEHOLDER_SEVERITY, PLACEHOLDER_CORRELATION, PLACEHOLDER_CALLER_FUNCTION, PLACEHOLDER_CALLER_FILE, PLACEHOLDER_CALLER_FILE_LINE, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.customTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEQUENCE, PLACEHOLDER_SEVERITY, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.callerCustomTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEQUENCE, PLACEHOLDER_SEVERITY, PLACEHOLDER_CALLER_FUNCTION, PLACEHOLDER_CALLER_FILE, PLACEHOLDER_CALLER_FILE_LINE, PLACEHOLDER_MSG)
	} else {
		adjustTemplate(formatter, &formatter.template, PLACEHOLDER_TIME, PLACEHOLDER_SEVERITY, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.callerTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEVERITY, PLACEHOLDER_CALLER_FUNCTION, PLACEHOLDER_CALLER_FILE, PLACEHOLDER_CALLER_FILE_LINE, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.correlationIdTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEVERITY, PLACEHOLDER_CORRELATION, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.callerCorrelationIdTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEVERITY, PLACEHOLDER_CORRELATION, PLACEHOLDER_CALLER_FUNCTION, PLACEHOLDER_CALLER_FILE, PLACEHOLDER_CALLER_FILE_LINE, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.customTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEVERITY, PLACEHOLDER_MSG)
		adjustTemplate(formatter, &formatter.callerCustomTemplate, PLACEHOLDER_TIME, PLACEHOLDER_SEVERITY, PLACEHOLDER_CALLER_FUNCTION, PLACEHOLDER_CALLER_FILE, PLACEHOLDER_CALLER_FILE_LINE, PLACEHOLDER_MSG)
	}
}

// Adjust a template by replacing template formatter placeholders with go known verbs
func adjustTemplate(formatter *TemplateFormatter, template *string, placeholders ...string) {
	previousCurrentIndex := -1
	previousLastIndex := -1

	envPairCount, envIsSorted := determineLengthAndSortedPlaceholder(template, PLACEHOLDER_ENV_VALUES)
	customPairCount, customIsSorted := determineLengthAndSortedPlaceholder(template, PLACEHOLDER_CUSTOM_VALUES)

	isSorted := envIsSorted && customIsSorted && envPairCount == len(formatter.commonProperties.envNamesToLog)

	if isSorted {
		extendedPlaceholders := append(placeholders, PLACEHOLDER_ENV_VALUES, PLACEHOLDER_CUSTOM_VALUES)
		for _, s := range extendedPlaceholders {
			currentIndex := strings.Index(*template, s)
			lastIndex := strings.LastIndex(*template, s)
			if (currentIndex != -1 && currentIndex < previousCurrentIndex) || (lastIndex != -1 && lastIndex < previousLastIndex) {
				isSorted = false
				break
			}
			previousCurrentIndex = currentIndex
			previousLastIndex = lastIndex
		}
	}

	for i, s := range placeholders {
		*template = strings.Replace(*template, s, determineReplacement(isSorted, s, i), 1)
	}

	// To avoid replacement of a higher decimal potence by a lower one (e.g. $env_k10 vs. $env_k1): iterate reverse
	envIndex := len(placeholders) + 2*len(formatter.commonProperties.envNamesToLog) - 1
	for i := len(formatter.commonProperties.envNamesToLog) - 1; i >= 0; i-- {
		*template = replacePlaceholder(PLACEHOLDER_ENV_VALUES, template, isSorted, true, i, envIndex-1)
		*template = replacePlaceholder(PLACEHOLDER_ENV_VALUES, template, isSorted, false, i, envIndex)
		envIndex -= 2
	}
	customIndex := len(placeholders) + 2*len(formatter.commonProperties.envNamesToLog) + 2*customPairCount - 1
	for i := customPairCount - 1; i >= 0; i-- {
		*template = replacePlaceholder(PLACEHOLDER_CUSTOM_VALUES, template, isSorted, true, i, customIndex-1)
		*template = replacePlaceholder(PLACEHOLDER_CUSTOM_VALUES, template, isSorted, false, i, customIndex)
		customIndex -= 2
	}
}

// Determines the the replacement token for a given placeholder. If its sorted, an index is provided.
func determineReplacement(isSorted bool, placeholder string, index int) string {
	if isSorted {
		return "%" + replacements[placeholder]
	} else {
		return "%[" + strconv.Itoa(index+1) + "]" + replacements[placeholder]
	}
}

// Replace a token for an environment key name or value. If its sorted, an index is provided.
func replacePlaceholder(placeholder string, template *string, isSorted bool, isKey bool, pairIndex int, index int) string {
	var variableName string
	var typeFormat = "v"
	if isKey {
		variableName = placeholder + "_k" + strconv.Itoa(pairIndex)
		typeFormat = "s"
	} else {
		variableName = placeholder + "_v" + strconv.Itoa(pairIndex)
		variableIndex := strings.Index(*template, variableName)
		formatIndexStart := variableIndex + len(variableName)
		if variableIndex > -1 && len(*template) > formatIndexStart {
			formatIndices := []int{strings.Index((*template)[formatIndexStart:], "["), strings.Index((*template)[formatIndexStart:], "]")}
			if formatIndices[0] == 0 && formatIndices[1] > -1 {
				typeFormat = (*template)[formatIndexStart+1+formatIndices[0] : formatIndexStart+formatIndices[1]]
				variableName = variableName + "[" + typeFormat + "]"
			}
		}
	}

	if isSorted {
		return strings.Replace(*template, variableName, "%"+typeFormat, 1)
	} else {
		formatSpec := ""
		if len(typeFormat) > 1 {
			formatSpec = typeFormat[:len(typeFormat)-1]
			typeFormat = typeFormat[len(typeFormat)-1:]
		}
		return strings.Replace(*template, variableName, "%"+formatSpec+"["+strconv.Itoa(index+1)+"]"+typeFormat, 1)
	}
}

// Determines the maximum length of keyValue placeholder pairs and if the placeholders are sorted
func determineLengthAndSortedPlaceholder(template *string, placeholder string) (int, bool) {
	reachedMaxValue := -1
	isSorted := true
	lengthTemplate := len(*template)
	lengthBasePlaceholder := len(placeholder) + 2 // plus _k or _v

	counter := 0

	index := strings.Index(*template, placeholder)
	for index > -1 {
		counter++
		index += lengthBasePlaceholder
		isKey := (*template)[index-1:index] == "k"
		length := 0
		currentValue := 0
		potence := 1
		checkNext := true
		for checkNext {
			posValue, err := strconv.Atoi((*template)[index+length : index+length+1])
			if err != nil {
				break
			}
			currentValue += potence * posValue
			potence *= 10
			length++
			checkNext = index+length+1 < lengthTemplate
		}
		if reachedMaxValue < currentValue {
			reachedMaxValue = currentValue
		}

		checkPlaceholderSorting(&isSorted, currentValue, reachedMaxValue, isKey, counter)
		determineNextPlaceholderIndex(&index, template, placeholder)
	}

	maxKeyValuePair := reachedMaxValue + 1
	isSorted = isSorted && maxKeyValuePair*2 == counter

	return maxKeyValuePair, isSorted
}

// checks wether placeholders are sorted or not and sets the sorting property
func checkPlaceholderSorting(isSorted *bool, currentValue int, reachedMaxValue int, isKeyPlaceholder bool, currentCount int) {
	if currentValue < reachedMaxValue {
		*isSorted = false
	} else if *isSorted {
		// check key or value position
		*isSorted = (isKeyPlaceholder && currentCount%2 == 1 && currentCount-1 == currentValue*2) || (!isKeyPlaceholder && currentCount%2 == 0 && currentCount-1 == currentValue*2+1)
	}
}

// determines the index of the next placeholder to check
func determineNextPlaceholderIndex(index *int, template *string, placeholder string) {
	next := strings.Index((*template)[*index:], placeholder)
	if next > -1 {
		*index += next
	} else {
		*index = -1
	}
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
	if t.isSequenceActive {
		args = append(args, logValues.Sequence)
	}
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
	if len(t.envNameAndValuesToLog) > 0 {
		args = append(args, t.envNameAndValuesToLog...)
	}
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
	} else {
		result = t.customTemplate
	}
	if t.areCustomKeyValueArgumentsToAppendAtTemplate(logValues) {
		for i := 0; i < len(*logValues.CustomValues); i++ {
			result += " [%s]: %v"
		}
	}
	return &result
}

// Checks whether to add custom key-value pairs at custom templates format or not.
func (t TemplateFormatter) areCustomKeyValueArgumentsToAppendAtTemplate(logValues *common.LogValues) bool {
	withCaller := logValues.IsCallerSet && t.isDefaultCallerCustomTemplate
	withoutCaller := !logValues.IsCallerSet && t.isDefaultCustomTemplate

	return withCaller || withoutCaller
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
