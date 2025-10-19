// This package provides formatter to transform log parameter to an entry
package format

import (
	"os"
	"strconv"
	"strings"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

const (
	DEBUG_PREFIX            string = "DEBUG"
	INFORMATION_TRIM_PREFIX string = "INFO"
	INFORMATION_PREFIX      string = INFORMATION_TRIM_PREFIX + " "
	WARNING_TRIM_PREFIX     string = "WARN"
	WARNING_PREFIX          string = WARNING_TRIM_PREFIX + " "
	ERROR_PREFIX            string = "ERROR"
	FATAL_PREFIX            string = "FATAL"
)

// Interface to format record values
type Formatter interface {
	// Formats the given parameter to a string to log
	Format(logValues *common.LogValues) string
}

// Shared common properties of a formatter
type CommonFormatterProperties struct {
	envNamesToLog  []string
	envValuesToLog []any
}

// Creates the CommonFormatterProperties from common config
func CreateCommonFormatterProperties(config *config.CommonFormatterConfig) *CommonFormatterProperties {
	result := CommonFormatterProperties{envNamesToLog: make([]string, 0, len(config.EnvNamesToLog)), envValuesToLog: make([]any, 0, len(config.EnvNamesToLog))}
	for _, s := range config.EnvNamesToLog {
		envValue := os.Getenv(s)
		if envValue == "" {
			continue
		}
		result.envNamesToLog = append(result.envNamesToLog, s)
		intValue, err := strconv.ParseInt(envValue, 10, 0)
		if err == nil {
			result.envValuesToLog = append(result.envValuesToLog, intValue)
			continue
		}
		floatValue, err := strconv.ParseFloat(envValue, 64)
		if err == nil {
			result.envValuesToLog = append(result.envValuesToLog, floatValue)
			continue
		}
		upperValue := strings.ToUpper(envValue)
		if upperValue == "TRUE" || upperValue == "FALSE" {
			result.envValuesToLog = append(result.envValuesToLog, upperValue == "TRUE")
			continue
		}
		result.envValuesToLog = append(result.envValuesToLog, envValue)
	}
	return &result
}

var severityTextMap = map[int]string{
	common.DEBUG_SEVERITY:       DEBUG_PREFIX,
	common.INFORMATION_SEVERITY: INFORMATION_PREFIX,
	common.WARNING_SEVERITY:     WARNING_PREFIX,
	common.ERROR_SEVERITY:       ERROR_PREFIX,
	common.FATAL_SEVERITY:       FATAL_PREFIX,
}

var severityTrimTextMap = map[int]string{
	common.DEBUG_SEVERITY:       DEBUG_PREFIX,
	common.INFORMATION_SEVERITY: INFORMATION_TRIM_PREFIX,
	common.WARNING_SEVERITY:     WARNING_TRIM_PREFIX,
	common.ERROR_SEVERITY:       ERROR_PREFIX,
	common.FATAL_SEVERITY:       FATAL_PREFIX,
}
