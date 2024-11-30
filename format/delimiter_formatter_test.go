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
	formatterMockTime = &delimiterFormatTestTime
	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, delimiterFormatter.Format(severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestDelimiterFormatCorrelation(t *testing.T) {
	formatterMockTime = &delimiterFormatTestTime
	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - someCorrelationId - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - someCorrelationId - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - someCorrelationId - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - someCorrelationId - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - someCorrelationId - Testmessage",
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
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - Testmessage - abc - 1 - true",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - Testmessage - abc - 1 - true",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - Testmessage - abc - 1 - true",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - Testmessage - abc - 1 - true",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - Testmessage - abc - 1 - true",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, delimiterFormatter.FormatCustom(severity, "Testmessage", customProperties), t, fmt.Sprintf("Format severity %d", severity))
	}
}
