// This package provides some utility functions to assert correct testing results
package testutil

import (
	"reflect"
	"strings"
	"testing"
)

// Checks that 'actual' is not nil. If it is nil [testing.T.Errorf] will be called
func AssertNotNil(actual any, t *testing.T, objectName string) {
	if isNil(actual) {
		t.Errorf("Element %s should not be nil", objectName)
	}
}

// Checks that 'actual' is nil. If it is not nil [testing.T.Errorf] will be called
func AssertNil(actual any, t *testing.T, objectName string) {
	if !isNil(actual) {
		t.Errorf("Element %s should be nil, but has value '%v'", objectName, actual)
	}
}

// Checks that 'actual' is true. If it is false [testing.T.Errorf] will be called
func AssertTrue(actual bool, t *testing.T, objectName string) {
	AssertEquals(true, actual, t, objectName)
}

// Checks that 'actual' is false. If it is true [testing.T.Errorf] will be called
func AssertFalse(actual bool, t *testing.T, objectName string) {
	AssertEquals(false, actual, t, objectName)
}

// Checks that 'actual' is not nil and equal to 'expected'. If it is nil or not equal [testing.T.Errorf] will be called
func AssertEquals(expected any, actual any, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if expected != actual {
		t.Errorf("Element %s is not as expected. expected: '%v' actual: '%v'", objectName, expected, actual)
	}
}

// Checks that 'actual' is not nil and not equal to 'expected'. If it is nil or equal [testing.T.Errorf] will be called
func AssertNotEquals(notExpected any, actual any, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if notExpected == actual {
		t.Errorf("Element %s equals the unexpected. unexpected: '%v' actual: '%v'", objectName, notExpected, actual)
	}
}

// Checks that 'actual' is not nil and has 'expected' as suffix. If it is nil or does not have the correct suffix [testing.T.Errorf] will be called
func AssertHasSuffix(expected string, actual string, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if !strings.HasSuffix(actual, expected) {
		t.Errorf("Element %s is not as expected. expected suffix: '%v' actual: '%v'", objectName, expected, actual)
	}
}

func isNil(toCheck any) bool {
	v := reflect.ValueOf(toCheck)
	k := v.Kind()

	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return toCheck == nil
	}
}
