package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/ma-vin/typewriter"
	testutil "github.com/ma-vin/typewriter/util"
)

var formatter Formatter = DelimiterFormatter{
	Delimiter: " - ",
}

var testTime = time.Date(2024, time.October, 1, 13, 20, 0, 0, time.UTC)
var testTimeText = testTime.Local().Format(time.RFC3339)

func TestFormat(t *testing.T) {
	expectedResults := map[int]string{
		typewriter.DEBUG_SEVERITY:       testTimeText + " - DEBUG - Testmessage",
		typewriter.INFORMATION_SEVERITY: testTimeText + " - INFO  - Testmessage",
		typewriter.WARNING_SEVERITY:     testTimeText + " - WARN  - Testmessage",
		typewriter.ERROR_SEVERITY:       testTimeText + " - ERROR - Testmessage",
		typewriter.FATAL_SEVERITY:       testTimeText + " - FATAL - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, formatter.Format(testTime, severity, "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestFormatCorrelation(t *testing.T) {
	expectedResults := map[int]string{
		typewriter.DEBUG_SEVERITY:       testTimeText + " - DEBUG - someCorrelationId - Testmessage",
		typewriter.INFORMATION_SEVERITY: testTimeText + " - INFO  - someCorrelationId - Testmessage",
		typewriter.WARNING_SEVERITY:     testTimeText + " - WARN  - someCorrelationId - Testmessage",
		typewriter.ERROR_SEVERITY:       testTimeText + " - ERROR - someCorrelationId - Testmessage",
		typewriter.FATAL_SEVERITY:       testTimeText + " - FATAL - someCorrelationId - Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, formatter.FormatWithCorrelation(testTime, severity, "someCorrelationId", "Testmessage"), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestFormatCustom(t *testing.T) {
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		typewriter.DEBUG_SEVERITY:       testTimeText + " - DEBUG - Testmessage - abc - 1 - true",
		typewriter.INFORMATION_SEVERITY: testTimeText + " - INFO  - Testmessage - abc - 1 - true",
		typewriter.WARNING_SEVERITY:     testTimeText + " - WARN  - Testmessage - abc - 1 - true",
		typewriter.ERROR_SEVERITY:       testTimeText + " - ERROR - Testmessage - abc - 1 - true",
		typewriter.FATAL_SEVERITY:       testTimeText + " - FATAL - Testmessage - abc - 1 - true",
	}

	for severity, expexpectedMessage := range expectedResults {
		testutil.AssertEquals(expexpectedMessage, formatter.FormatCustom(testTime, severity, "Testmessage", customProperties), t, fmt.Sprintf("Format severity %d", severity))
	}
}
