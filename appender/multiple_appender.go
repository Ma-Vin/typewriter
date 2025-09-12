package appender

import (
	"reflect"
	"slices"

	"github.com/ma-vin/typewriter/common"
)

// Delegates the messages to its sub-appenders
type MultiAppender struct {
	appenders []*Appender
}

// Writes the given logValues at its sub-appenders
func (m MultiAppender) Write(logValues *common.LogValues) {
	for _, a := range m.appenders {
		(*a).Write(logValues)
	}
}

// Closes all sub-appenders
func (m MultiAppender) Close() {
	for _, a := range m.appenders {
		(*a).Close()
	}
}

// creates a new multiple appender with an initial capacity
func CreateMultiAppenderWithCapacity(capacity int) *MultiAppender {
	return &MultiAppender{appenders: make([]*Appender, 0, capacity)}
}

// adds a pointer to a sub-appender
func (m *MultiAppender) AddSubAppender(appenderToAdd *Appender) {
	m.appenders = append(m.appenders, appenderToAdd)
}

// Function to test existence of sub appender types. Return true if all subAppenderTypes exits a m.appenders
func (m MultiAppender) CheckSubAppenderTypesForTest(subAppenderTypes []string) bool {
	if len(subAppenderTypes) != len(m.appenders) {
		return false
	}
	for _, s := range subAppenderTypes {
		if !slices.ContainsFunc(m.appenders, func(a *Appender) bool { return reflect.TypeOf(*a).Name() == s }) {
			return false
		}
	}
	return true
}
