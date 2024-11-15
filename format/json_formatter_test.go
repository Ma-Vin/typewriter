package format

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/testutil"
)

var jsonFormatTestTime = time.Date(2024, time.November, 15, 20, 00, 0, 0, time.UTC)
var jsonFormatTestTimeText = jsonFormatTestTime.Local().Format(time.RFC3339Nano)

var jsonFormatter Formatter = CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339Nano, false)
var jsonFormatterSub Formatter = CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339Nano, true)

func TestJsonFormat(t *testing.T) {
	formatterMockTime = &jsonFormatTestTime

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"DEBUG\", \"message\": \"Testmessage\" }",
		constants.INFORMATION_SEVERITY: "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"INFO\", \"message\": \"Testmessage\" }",
		constants.WARNING_SEVERITY:     "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"WARN\", \"message\": \"Testmessage\" }",
		constants.ERROR_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"ERROR\", \"message\": \"Testmessage\" }",
		constants.FATAL_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"FATAL\", \"message\": \"Testmessage\" }",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, jsonFormatter.Format(severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestJsonFormatCorrelation(t *testing.T) {
	formatterMockTime = &jsonFormatTestTime

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"DEBUG\", \"correlation\": \"someCorrelationId\", \"message\": \"Testmessage\" }",
		constants.INFORMATION_SEVERITY: "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"INFO\", \"correlation\": \"someCorrelationId\", \"message\": \"Testmessage\" }",
		constants.WARNING_SEVERITY:     "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"WARN\", \"correlation\": \"someCorrelationId\", \"message\": \"Testmessage\" }",
		constants.ERROR_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"ERROR\", \"correlation\": \"someCorrelationId\", \"message\": \"Testmessage\" }",
		constants.FATAL_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"FATAL\", \"correlation\": \"someCorrelationId\", \"message\": \"Testmessage\" }",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, jsonFormatter.FormatWithCorrelation(severity, "someCorrelationId", "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestJsonFormatCustom(t *testing.T) {
	formatterMockTime = &jsonFormatTestTime

	customProperties := map[string]any{
		"first": "abc",
	}

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"DEBUG\", \"message\": \"Testmessage\", \"first\": \"abc\" }",
		constants.INFORMATION_SEVERITY: "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"INFO\", \"message\": \"Testmessage\", \"first\": \"abc\" }",
		constants.WARNING_SEVERITY:     "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"WARN\", \"message\": \"Testmessage\", \"first\": \"abc\" }",
		constants.ERROR_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"ERROR\", \"message\": \"Testmessage\", \"first\": \"abc\" }",
		constants.FATAL_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"FATAL\", \"message\": \"Testmessage\", \"first\": \"abc\" }",
	}

	for severity, expectedMessage := range expectedResults {
		result := jsonFormatter.FormatCustom(severity, "Testmessage", customProperties)
		testutil.AssertEquals(expectedMessage, result, t, fmt.Sprintf("Format severity %d", severity))
	}
}

// There is no guarantee that customProperties will be iterated in order of insert
// So we check contains instead of equals
func TestJsonFormatCustomAllTypes(t *testing.T) {
	formatterMockTime = &jsonFormatTestTime

	var boolValue bool = true
	var byteValue byte = 1
	var intValue int = 2
	var int8Value int8 = 3
	var int16Value int16 = 4
	var int32Value int32 = 5
	var int64Value int64 = 6
	var uintValue uint = 7
	var uint8Value uint8 = 8
	var uint16Value uint16 = 9
	var uint32Value uint32 = 10
	var uint64Value uint64 = 11
	var float32Value float32 = 1.1
	var float64Value float64 = 1.2
	var complex64Value complex64 = complex(1, 1)
	var complex128Value complex128 = complex(1, 2)
	var stringValue string = "abc"

	customProperties := map[string]any{
		"boolValue":       boolValue,
		"byteValue":       byteValue,
		"intValue":        intValue,
		"int8Value":       int8Value,
		"int16Value":      int16Value,
		"int32Value":      int32Value,
		"int64Value":      int64Value,
		"uintValue":       uintValue,
		"uint8Value":      uint8Value,
		"uint16Value":     uint16Value,
		"uint32Value":     uint32Value,
		"uint64Value":     uint64Value,
		"float32Value":    float32Value,
		"float64Value":    float64Value,
		"complex64Value":  complex64Value,
		"complex128Value": complex128Value,
		"stringValue":     stringValue,
	}

	result := jsonFormatter.FormatCustom(constants.INFORMATION_SEVERITY, "Testmessage", customProperties)

	testutil.AssertTrue(strings.Contains(result, "\"boolValue\": true"), t, "Format conatians bool")
	testutil.AssertTrue(strings.Contains(result, "\"byteValue\": 1"), t, "Format conatians byteValue")
	testutil.AssertTrue(strings.Contains(result, "\"intValue\": 2"), t, "Format conatians intValue")
	testutil.AssertTrue(strings.Contains(result, "\"int8Value\": 3"), t, "Format conatians int8Value")
	testutil.AssertTrue(strings.Contains(result, "\"int16Value\": 4"), t, "Format conatians int16Value")
	testutil.AssertTrue(strings.Contains(result, "\"int32Value\": 5"), t, "Format conatians int32Value")
	testutil.AssertTrue(strings.Contains(result, "\"int64Value\": 6"), t, "Format conatians int64Value")
	testutil.AssertTrue(strings.Contains(result, "\"uintValue\": 7"), t, "Format conatians uintValue")
	testutil.AssertTrue(strings.Contains(result, "\"uint8Value\": 8"), t, "Format conatians uint8Value")
	testutil.AssertTrue(strings.Contains(result, "\"uint16Value\": 9"), t, "Format conatians uint16Value")
	testutil.AssertTrue(strings.Contains(result, "\"uint32Value\": 10"), t, "Format conatians uint32Value")
	testutil.AssertTrue(strings.Contains(result, "\"uint64Value\": 11"), t, "Format conatians uint64Value")
	testutil.AssertTrue(strings.Contains(result, "\"float32Value\": 1.1"), t, "Format conatians float32Value")
	testutil.AssertTrue(strings.Contains(result, "\"float64Value\": 1.2"), t, "Format conatians float64Value")
	testutil.AssertTrue(strings.Contains(result, "\"complex64Value\": (1+1i)"), t, "Format conatians complex64Value")
	testutil.AssertTrue(strings.Contains(result, "\"complex128Value\": (1+2i)"), t, "Format conatians complex128Value")
	testutil.AssertTrue(strings.Contains(result, "\"stringValue\": \"abc\""), t, "Format conatians stringValue")

}

func TestJsonFormatCustomSub(t *testing.T) {
	formatterMockTime = &jsonFormatTestTime

	customProperties := map[string]any{
		"first": "abc",
	}

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"DEBUG\", \"message\": \"Testmessage\", \"custom\": { \"first\": \"abc\" } }",
		constants.INFORMATION_SEVERITY: "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"INFO\", \"message\": \"Testmessage\", \"custom\": { \"first\": \"abc\" } }",
		constants.WARNING_SEVERITY:     "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"WARN\", \"message\": \"Testmessage\", \"custom\": { \"first\": \"abc\" } }",
		constants.ERROR_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"ERROR\", \"message\": \"Testmessage\", \"custom\": { \"first\": \"abc\" } }",
		constants.FATAL_SEVERITY:       "{ \"time\": \"" + jsonFormatTestTimeText + "\", \"severity\": \"FATAL\", \"message\": \"Testmessage\", \"custom\": { \"first\": \"abc\" } }",
	}

	for severity, expectedMessage := range expectedResults {
		result := jsonFormatterSub.FormatCustom(severity, "Testmessage", customProperties)
		testutil.AssertEquals(expectedMessage, result, t, fmt.Sprintf("Format severity %d", severity))
	}
}
