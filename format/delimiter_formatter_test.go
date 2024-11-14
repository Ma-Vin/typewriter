package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/testutil"
)

var delimiterFormatter Formatter = CreateDelimiterFormatter(" - ")

var delimiterFormatTestTime = time.Date(2024, time.October, 1, 13, 20, 0, 0, time.UTC)
var delimiterFormatTestTimeText = delimiterFormatTestTime.Local().Format(time.RFC3339)

func TestDelimiterFormat(t *testing.T) {
	formatterMockTime = &delimiterFormatTestTime
	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage",
		constants.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage",
		constants.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage",
		constants.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage",
		constants.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, delimiterFormatter.Format(severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestDelimiterFormatCorrelation(t *testing.T) {
	formatterMockTime = &delimiterFormatTestTime
	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - someCorrelationId - Testmessage",
		constants.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - someCorrelationId - Testmessage",
		constants.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - someCorrelationId - Testmessage",
		constants.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - someCorrelationId - Testmessage",
		constants.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - someCorrelationId - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, delimiterFormatter.FormatWithCorrelation(severity, "someCorrelationId", "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestDelimiterFormatCustom(t *testing.T) {
	formatterMockTime = &delimiterFormatTestTime
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage - abc - 1 - true",
		constants.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage - abc - 1 - true",
		constants.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage - abc - 1 - true",
		constants.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage - abc - 1 - true",
		constants.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage - abc - 1 - true",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, delimiterFormatter.FormatCustom(severity, "Testmessage", customProperties), t, fmt.Sprintf("Format severity %d", severity))
	}
}
