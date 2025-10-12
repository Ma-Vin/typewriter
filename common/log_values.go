package common

import (
	"sync"
	"time"
)

const NOT_AVAILABLE_VALUE = "n/a"

var sequenceCounter uint64 = 0
var counterMutex = sync.Mutex{}

// Value container which holds the values to log. External given values CorrelationId and CustomValues are provided as pointers
// (Message might be concat by [fmt.Sprint] / [fmt.Sprintf] from several values, so Message is not given external)
type LogValues struct {
	Time           time.Time
	Sequence       uint64
	Severity       int
	CorrelationId  *string
	Message        string
	CustomValues   *map[string]any
	IsCallerSet    bool
	CallerFile     string
	CallerFileLine int
	CallerFunction string
}

var logValuesMockTime *time.Time = nil

// Creates default log values
func CreateLogValues(severity int, message string) LogValues {
	return LogValues{
		Time:           GetNow(),
		Sequence:       determineNextSequenceElement(),
		Severity:       severity,
		CorrelationId:  nil,
		Message:        message,
		CustomValues:   nil,
		IsCallerSet:    false,
		CallerFile:     NOT_AVAILABLE_VALUE,
		CallerFileLine: -1,
		CallerFunction: NOT_AVAILABLE_VALUE,
	}
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

func determineNextSequenceElement() uint64 {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	sequenceCounter++
	return sequenceCounter
}

func InitSequenceCounter() {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	sequenceCounter = 0
}

// returns the current time. Or a mock time if set.
func GetNow() time.Time {
	if logValuesMockTime != nil {
		return *logValuesMockTime
	}
	return time.Now()
}

// For test usage only! Sets constant mock time. If this parameter is nil [time.Now] will be calculated
func SetLogValuesMockTime(mockTime *time.Time) {
	logValuesMockTime = mockTime
}
