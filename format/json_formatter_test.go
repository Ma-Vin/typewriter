package format

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/testutil"
)

var jsonFormatTestTime = time.Date(2024, time.November, 15, 20, 00, 0, 0, time.UTC)
var jsonFormatTestTimeText = jsonFormatTestTime.Format(time.RFC3339Nano)

var jsonFormatter Formatter = CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339Nano, "caller", "file", "line", false)
var jsonFormatterSub Formatter = CreateJsonFormatter("time", "severity", "message", "correlation", "custom", time.RFC3339Nano, "caller", "file", "line", true)

func TestJsonFormat(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "{\"message\":\"Testmessage\",\"severity\":\"DEBUG\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.INFORMATION_SEVERITY: "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.WARNING_SEVERITY:     "{\"message\":\"Testmessage\",\"severity\":\"WARN\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.ERROR_SEVERITY:       "{\"message\":\"Testmessage\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.FATAL_SEVERITY:       "{\"message\":\"Testmessage\",\"severity\":\"FATAL\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, jsonFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestJsonFormatCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)
	correleation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"DEBUG\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.INFORMATION_SEVERITY: "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.WARNING_SEVERITY:     "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"WARN\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.ERROR_SEVERITY:       "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.FATAL_SEVERITY:       "{\"correlation\":\"someCorrelationId\",\"message\":\"Testmessage\",\"severity\":\"FATAL\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correleation, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, jsonFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestJsonFormatCustom(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	customProperties := map[string]any{
		"first": "abc",
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"DEBUG\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.INFORMATION_SEVERITY: "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.WARNING_SEVERITY:     "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"WARN\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.ERROR_SEVERITY:       "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.FATAL_SEVERITY:       "{\"first\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"FATAL\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		result := jsonFormatter.Format(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, result, t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestJsonFormatCustomAllTypes(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

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
	var stringValue string = "abc"

	customProperties := map[string]any{
		"boolValue":    boolValue,
		"byteValue":    byteValue,
		"intValue":     intValue,
		"int8Value":    int8Value,
		"int16Value":   int16Value,
		"int32Value":   int32Value,
		"int64Value":   int64Value,
		"uintValue":    uintValue,
		"uint8Value":   uint8Value,
		"uint16Value":  uint16Value,
		"uint32Value":  uint32Value,
		"uint64Value":  uint64Value,
		"float32Value": float32Value,
		"float64Value": float64Value,
		"stringValue":  stringValue,
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)

	result := jsonFormatter.Format(&logValuesToFormat)

	testutil.AssertTrue(strings.Contains(result, "\"boolValue\":true"), t, "Format conatians bool")
	testutil.AssertTrue(strings.Contains(result, "\"byteValue\":1"), t, "Format conatians byteValue")
	testutil.AssertTrue(strings.Contains(result, "\"intValue\":2"), t, "Format conatians intValue")
	testutil.AssertTrue(strings.Contains(result, "\"int8Value\":3"), t, "Format conatians int8Value")
	testutil.AssertTrue(strings.Contains(result, "\"int16Value\":4"), t, "Format conatians int16Value")
	testutil.AssertTrue(strings.Contains(result, "\"int32Value\":5"), t, "Format conatians int32Value")
	testutil.AssertTrue(strings.Contains(result, "\"int64Value\":6"), t, "Format conatians int64Value")
	testutil.AssertTrue(strings.Contains(result, "\"uintValue\":7"), t, "Format conatians uintValue")
	testutil.AssertTrue(strings.Contains(result, "\"uint8Value\":8"), t, "Format conatians uint8Value")
	testutil.AssertTrue(strings.Contains(result, "\"uint16Value\":9"), t, "Format conatians uint16Value")
	testutil.AssertTrue(strings.Contains(result, "\"uint32Value\":10"), t, "Format conatians uint32Value")
	testutil.AssertTrue(strings.Contains(result, "\"uint64Value\":11"), t, "Format conatians uint64Value")
	testutil.AssertTrue(strings.Contains(result, "\"float32Value\":1.1"), t, "Format conatians float32Value")
	testutil.AssertTrue(strings.Contains(result, "\"float64Value\":1.2"), t, "Format conatians float64Value")
	testutil.AssertTrue(strings.Contains(result, "\"stringValue\":\"abc\""), t, "Format conatians stringValue")
}

func TestJsonFormatUnsupported(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	customProperties := map[string]any{
		"complex64Value":  complex(1, 1),
		"complex128Value": complex(1, 2),
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)

	result := jsonFormatter.Format(&logValuesToFormat)

	expectedResult := "{\"message\":\"Fail to marshal to json, use custom formatter: json: unsupported type: complex128\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}" +
		fmt.Sprintln() + "{\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}"

	testutil.AssertEquals(expectedResult, result, t, "Format unsupported")
}

func TestJsonFormatUnsupportedAndCorrelationId(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	customProperties := map[string]any{
		"correlation":     "abc",
		"complex64Value":  complex(1, 1),
		"complex128Value": complex(1, 2),
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)

	result := jsonFormatter.Format(&logValuesToFormat)

	expectedResult := "{\"correlation\":\"abc\",\"message\":\"Fail to marshal to json, use custom formatter: json: unsupported type: complex128\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}" +
		fmt.Sprintln() + "{\"correlation\":\"abc\",\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}"

	testutil.AssertEquals(expectedResult, result, t, "Format unsupported")
}

func TestJsonFormatCustomSub(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	customProperties := map[string]any{
		"first": "abc",
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "{\"custom\":{\"first\":\"abc\"},\"message\":\"Testmessage\",\"severity\":\"DEBUG\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.INFORMATION_SEVERITY: "{\"custom\":{\"first\":\"abc\"},\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.WARNING_SEVERITY:     "{\"custom\":{\"first\":\"abc\"},\"message\":\"Testmessage\",\"severity\":\"WARN\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.ERROR_SEVERITY:       "{\"custom\":{\"first\":\"abc\"},\"message\":\"Testmessage\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.FATAL_SEVERITY:       "{\"custom\":{\"first\":\"abc\"},\"message\":\"Testmessage\",\"severity\":\"FATAL\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)

		result := jsonFormatterSub.Format(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, result, t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestJsonFormatCaller(t *testing.T) {
	common.SetLogValuesMockTime(&jsonFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "{\"caller\":\"f1\",\"file\":\"abc\",\"line\":3,\"message\":\"Testmessage\",\"severity\":\"DEBUG\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.INFORMATION_SEVERITY: "{\"caller\":\"f1\",\"file\":\"abc\",\"line\":3,\"message\":\"Testmessage\",\"severity\":\"INFO\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.WARNING_SEVERITY:     "{\"caller\":\"f1\",\"file\":\"abc\",\"line\":3,\"message\":\"Testmessage\",\"severity\":\"WARN\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.ERROR_SEVERITY:       "{\"caller\":\"f1\",\"file\":\"abc\",\"line\":3,\"message\":\"Testmessage\",\"severity\":\"ERROR\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
		common.FATAL_SEVERITY:       "{\"caller\":\"f1\",\"file\":\"abc\",\"line\":3,\"message\":\"Testmessage\",\"severity\":\"FATAL\",\"time\":\"" + jsonFormatTestTimeText + "\"}",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		logValuesToFormat.CallerFunction = "f1"
		logValuesToFormat.CallerFile = "abc"
		logValuesToFormat.CallerFileLine = 3
		logValuesToFormat.IsCallerSet = true

		result := jsonFormatterSub.Format(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, result, t, fmt.Sprintf("Format severity %d", severity))
	}
}
