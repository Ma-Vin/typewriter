package format

import (
	"time"

	"github.com/ma-vin/typewriter/constants"
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

type Formatter interface {
	// Formats the given parameter to a string to log
	Format(severity int, message string) string
	// Formats the given default parameter and a correleation id to a string to log
	FormatWithCorrelation(severity int, correlationId string, message string) string
	// Formats the given parameter to a string to log and he customValues will be added at the end
	FormatCustom(severity int, message string, customValues map[string]any) string
}

var severityTextMap = map[int]string{
	constants.DEBUG_SEVERITY:       DEBUG_PREFIX,
	constants.INFORMATION_SEVERITY: INFORMATION_PREFIX,
	constants.WARNING_SEVERITY:     WARNING_PREFIX,
	constants.ERROR_SEVERITY:       ERROR_PREFIX,
	constants.FATAL_SEVERITY:       FATAL_PREFIX,
}

var severityTrimTextMap = map[int]string{
	constants.DEBUG_SEVERITY:       DEBUG_PREFIX,
	constants.INFORMATION_SEVERITY: INFORMATION_TRIM_PREFIX,
	constants.WARNING_SEVERITY:     WARNING_TRIM_PREFIX,
	constants.ERROR_SEVERITY:       ERROR_PREFIX,
	constants.FATAL_SEVERITY:       FATAL_PREFIX,
}

var formatterMockTime *time.Time = nil

func getNowAsStringFromLayout(template string) string {
	timeToFormat := time.Now()
	if formatterMockTime != nil {
		timeToFormat = *formatterMockTime
	}
	return timeToFormat.Local().Format(template)
}

func getNowAsStringDefaultLayout() string {
	return getNowAsStringFromLayout(time.RFC3339)
}
