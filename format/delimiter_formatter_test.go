package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/testutil"
)

var delimiterFormatter Formatter = CreateDelimiterFormatter(" - ")

var delimiterFormatTestTime = time.Date(2024, time.October, 1, 13, 20, 0, 0, time.UTC)
var delimiterFormatTestTimeText = delimiterFormatTestTime.Format(time.RFC3339)

func TestDelimiterFormat(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expectedMessage, delimiterFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestDelimiterFormatCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)
	correleation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - someCorrelationId - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - someCorrelationId - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - someCorrelationId - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - someCorrelationId - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - someCorrelationId - Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correleation, "Testmessage")
		testutil.AssertEquals(expectedMessage, delimiterFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestDelimiterFormatCustom(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage - abc - 1 - true",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage - abc - 1 - true",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage - abc - 1 - true",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage - abc - 1 - true",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage - abc - 1 - true",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expectedMessage, delimiterFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestDelimiterFormatEmptyCustom(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)
	customProperties := map[string]any{}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expectedMessage, delimiterFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}