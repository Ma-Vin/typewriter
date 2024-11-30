// This package provides formatter to transform log parameter to an entry
package format

import (
	"github.com/ma-vin/typewriter/common"
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

// Interface to format recdord values
type Formatter interface {
	// Formats the given parameter to a string to log
	Format(logValues *common.LogValues) string
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