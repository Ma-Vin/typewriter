package format

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
)

func createDelimiterFormatterForTest(withSequence bool, withAnsiColor bool, envNamesToLog []string) Formatter {
	common.InitSequenceCounter()
	var config config.FormatterConfig = config.DelimiterFormatterConfig{
		Common: createCommonFormatterConfigForTest(
			time.RFC3339,
			withSequence,
			withAnsiColor,
			envNamesToLog),
		Delimiter: " - "}
	result, _ := CreateDelimiterFormatterFromConfig(&config)
	return *result
}

var delimiterFormatter Formatter = createDelimiterFormatterForTest(false, false, []string{})

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
	correlation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - DEBUG - someCorrelationId - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - INFO  - someCorrelationId - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - WARN  - someCorrelationId - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - ERROR - someCorrelationId - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - FATAL - someCorrelationId - Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correlation, "Testmessage")
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

func TestDelimiterFormatCaller(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	logValuesToFormat.CallerFile = "someFile"
	logValuesToFormat.CallerFileLine = 42
	logValuesToFormat.CallerFunction = "someFunction"
	logValuesToFormat.IsCallerSet = true

	testutil.AssertEquals(delimiterFormatTestTimeText+" - INFO  - someFunction at someFile (Line 42) - Testmessage", delimiterFormatter.Format(&logValuesToFormat), t, "Format caller")
}

func TestDelimiterFormatWithSequence(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)

	delimiterFormatterWithSequence := createDelimiterFormatterForTest(true, false, []string{})

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		expectedMessage := fmt.Sprintf(delimiterFormatTestTimeText+" - %d - %s - Testmessage", i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, delimiterFormatterWithSequence.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestDelimiterFormatWithEnvNames(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)

	os.Clearenv()
	os.Setenv("test1", "abc")
	os.Setenv("TEST2", "1")
	os.Setenv("Test3", "2.1")
	os.Setenv("test4", "true")

	delimiterFormatterWithSequence := createDelimiterFormatterForTest(true, false, []string{"test1", "TEST2", "Test3", "test4"})

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		expectedMessage := fmt.Sprintf(delimiterFormatTestTimeText+" - %d - %s - abc - 1 - 2.1 - true - Testmessage", i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, delimiterFormatterWithSequence.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestDElimiterFormatWithSeverityAnsiColored(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatTestTime)

	delimiterFormatterWithAnsiColor := createDelimiterFormatterForTest(false, true, []string{})

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       delimiterFormatTestTimeText + " - \033[34mDEBUG\033[0m - Testmessage",
		common.INFORMATION_SEVERITY: delimiterFormatTestTimeText + " - \033[32mINFO \033[0m - Testmessage",
		common.WARNING_SEVERITY:     delimiterFormatTestTimeText + " - \033[33mWARN \033[0m - Testmessage",
		common.ERROR_SEVERITY:       delimiterFormatTestTimeText + " - \033[31mERROR\033[0m - Testmessage",
		common.FATAL_SEVERITY:       delimiterFormatTestTimeText + " - \033[35mFATAL\033[0m - Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expectedMessage, delimiterFormatterWithAnsiColor.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}
