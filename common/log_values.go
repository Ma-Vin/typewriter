package common

import (
	"time"
)

// Value container which holds the values to log. External given values CorrelationId and CustomValues are provided as pointers
// (Message might be concat by [fmt.Sprint] / [fmt.Sprintf] from several values, so Message is not given external)
type LogValues struct {
	Time           time.Time
	Severity       int
	CorrelationId  *string
	Message        string
	CustomValues   *map[string]any
}

var lockValuesMockTime *time.Time = nil

// Creates default log values
func CreateLogValues(severity int, message string) LogValues {
	return LogValues{
		Time:           getNow(),
		Severity:       severity,
		CorrelationId:  nil,
		Message:        message,
		CustomValues:   nil}
}

// Creates log values with a correlation id
func CreateLogValuesWithCorrelation(severity int, correlationId *string, message string) LogValues {
	rec := CreateLogValues(severity, message)
	rec.CorrelationId = correlationId
	return rec
}

// Creates log values with a custom value map
func CreateLogValuesCustom(severity int, message string, customValues *map[string]any) LogValues {
	rec := CreateLogValues(severity, message)
	rec.CustomValues = customValues
	return rec
}

func getNow() time.Time {
	if lockValuesMockTime != nil {
		return *lockValuesMockTime
	}
	return time.Now()
}

// For test usage only! Sets constant mock time. If this parameter is nil [time.Now] will be calculated
func SetLogValuesMockTime(mockTime *time.Time) {
	lockValuesMockTime = mockTime
}
