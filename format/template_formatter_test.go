package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/testutil"
)

var templateFormatter Formatter = CreateTemplateFormatter(
	"time: %s severity: %s message: %s",
	"time: %s severity: %s correlation: %s message: %s",
	"time: %s severity: %s message: %s %s: %s %s: %d %s: %t",
	time.RFC1123Z)

var templateFormatterOrder Formatter = CreateTemplateFormatter(
	"severity: %[2]s message: %[3]s time: %[1]s",
	"severity: %[2]s correlation: %[3]s message: %[4]s time: %[1]s",
	"severity: %[2]s message: %[3]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %[9]t time: %[1]s",
	time.RFC1123Z)

var templateFormatTestTime = time.Date(2024, time.November, 1, 20, 15, 0, 0, time.UTC)
var templateFormatTestTimeText = templateFormatTestTime.Local().Format(time.RFC1123Z)

func TestTemplateFormat(t *testing.T) {
	formatterMockTime = &templateFormatTestTime

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG message: Testmessage",
		constants.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  message: Testmessage",
		constants.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  message: Testmessage",
		constants.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR message: Testmessage",
		constants.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatOrder(t *testing.T) {
	formatterMockTime = &templateFormatTestTime

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "severity: DEBUG message: Testmessage time: " + templateFormatTestTimeText,
		constants.INFORMATION_SEVERITY: "severity: INFO  message: Testmessage time: " + templateFormatTestTimeText,
		constants.WARNING_SEVERITY:     "severity: WARN  message: Testmessage time: " + templateFormatTestTimeText,
		constants.ERROR_SEVERITY:       "severity: ERROR message: Testmessage time: " + templateFormatTestTimeText,
		constants.FATAL_SEVERITY:       "severity: FATAL message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelation(t *testing.T) {
	formatterMockTime = &templateFormatTestTime

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG correlation: someCorrelationId message: Testmessage",
		constants.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  correlation: someCorrelationId message: Testmessage",
		constants.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  correlation: someCorrelationId message: Testmessage",
		constants.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR correlation: someCorrelationId message: Testmessage",
		constants.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL correlation: someCorrelationId message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, templateFormatter.FormatWithCorrelation(severity, "someCorrelationId", "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelationOrder(t *testing.T) {
	formatterMockTime = &templateFormatTestTime

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "severity: DEBUG correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		constants.INFORMATION_SEVERITY: "severity: INFO  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		constants.WARNING_SEVERITY:     "severity: WARN  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		constants.ERROR_SEVERITY:       "severity: ERROR correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		constants.FATAL_SEVERITY:       "severity: FATAL correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.FormatWithCorrelation(severity, "someCorrelationId", "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustom(t *testing.T) {
	formatterMockTime = &templateFormatTestTime

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG message: Testmessage first: abc second: 1 third: true",
		constants.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  message: Testmessage first: abc second: 1 third: true",
		constants.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  message: Testmessage first: abc second: 1 third: true",
		constants.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR message: Testmessage first: abc second: 1 third: true",
		constants.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL message: Testmessage first: abc second: 1 third: true",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, templateFormatter.FormatCustom(severity, "Testmessage", customProperties), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomOrder(t *testing.T) {
	formatterMockTime = &templateFormatTestTime

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       "severity: DEBUG message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		constants.INFORMATION_SEVERITY: "severity: INFO  message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		constants.WARNING_SEVERITY:     "severity: WARN  message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		constants.ERROR_SEVERITY:       "severity: ERROR message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		constants.FATAL_SEVERITY:       "severity: FATAL message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.FormatCustom(severity, "Testmessage", customProperties), t, fmt.Sprintf("Format severity %d", severity))
	}
}
