package format

import (
	"time"

	"github.com/ma-vin/typewriter/constants"
)

const (
	DEBUG_PREFIX       string = "DEBUG"
	INFORMATION_PREFIX string = "INFO "
	WARNING_PREFIX     string = "WARN "
	ERROR_PREFIX       string = "ERROR"
	FATAL_PREFIX       string = "FATAL"
)

type Formatter interface {
	// Formats the given parameter to a string to log
	Format(timestamp time.Time, severity int, message string) string
	// Formats the given default parameter and a correleation id to a string to log
	FormatWithCorrelation(timestamp time.Time, severity int, correlationId string, message string) string
	// Formats the given parameter to a string to log and he customValues will be added at the end
	FormatCustom(timestamp time.Time, severity int, message string, customValues map[string]any) string
}

var severityTextMap = map[int]string{
	constants.DEBUG_SEVERITY:       DEBUG_PREFIX,
	constants.INFORMATION_SEVERITY: INFORMATION_PREFIX,
	constants.WARNING_SEVERITY:     WARNING_PREFIX,
	constants.ERROR_SEVERITY:       ERROR_PREFIX,
	constants.FATAL_SEVERITY:       FATAL_PREFIX,
}