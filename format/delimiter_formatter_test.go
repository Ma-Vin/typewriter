package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/constants"
	testutil "github.com/ma-vin/typewriter/util"
)

var formatter Formatter = CreateDelimiterFormatter(" - ")

var testTime = time.Date(2024, time.October, 1, 13, 20, 0, 0, time.UTC)
var testTimeText = testTime.Local().Format(time.RFC3339)

func TestFormat(t *testing.T) {
	delimiterFormatterMockTime = &testTime
	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       testTimeText + " - DEBUG - Testmessage",
		constants.INFORMATION_SEVERITY: testTimeText + " - INFO  - Testmessage",
		constants.WARNING_SEVERITY:     testTimeText + " - WARN  - Testmessage",
		constants.ERROR_SEVERITY:       testTimeText + " - ERROR - Testmessage",
		constants.FATAL_SEVERITY:       testTimeText + " - FATAL - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, formatter.Format(severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestFormatCorrelation(t *testing.T) {
	delimiterFormatterMockTime = &testTime
	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       testTimeText + " - DEBUG - someCorrelationId - Testmessage",
		constants.INFORMATION_SEVERITY: testTimeText + " - INFO  - someCorrelationId - Testmessage",
		constants.WARNING_SEVERITY:     testTimeText + " - WARN  - someCorrelationId - Testmessage",
		constants.ERROR_SEVERITY:       testTimeText + " - ERROR - someCorrelationId - Testmessage",
		constants.FATAL_SEVERITY:       testTimeText + " - FATAL - someCorrelationId - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, formatter.FormatWithCorrelation(severity, "someCorrelationId", "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestFormatCustom(t *testing.T) {
	delimiterFormatterMockTime = &testTime
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		constants.DEBUG_SEVERITY:       testTimeText + " - DEBUG - Testmessage - abc - 1 - true",
		constants.INFORMATION_SEVERITY: testTimeText + " - INFO  - Testmessage - abc - 1 - true",
		constants.WARNING_SEVERITY:     testTimeText + " - WARN  - Testmessage - abc - 1 - true",
		constants.ERROR_SEVERITY:       testTimeText + " - ERROR - Testmessage - abc - 1 - true",
		constants.FATAL_SEVERITY:       testTimeText + " - FATAL - Testmessage - abc - 1 - true",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, formatter.FormatCustom(severity, "Testmessage", customProperties), t, fmt.Sprintf("Format severity %d", severity))
	}
}
