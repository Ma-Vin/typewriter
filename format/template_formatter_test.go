package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/testutil"
)

var templateFormatter Formatter = CreateTemplateFormatter(
	"time: %s severity: %s message: %s",
	"time: %s severity: %s correlation: %s message: %s",
	"time: %s severity: %s message: %s %s: %s %s: %d %s: %t",
	"time: %s severity: %s caller: %s file: %s line: %d message: %s",
	"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
	"time: %s severity: %s caller: %s file: %s line: %d message: %s %s: %s %s: %d %s: %t",
	time.RFC1123Z,
	false)

var templateFormatterOrder Formatter = CreateTemplateFormatter(
	"severity: %[2]s message: %[3]s time: %[1]s",
	"severity: %[2]s correlation: %[3]s message: %[4]s time: %[1]s",
	"severity: %[2]s message: %[3]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %[9]t time: %[1]s",
	"caller: %[3]s file: %[4]s line: %[5]d severity: %[2]s message: %[6]s time: %[1]s",
	"caller: %[4]s file: %[5]s line: %[6]d severity: %[2]s correlation: %[3]s message: %[7]s time: %[1]s",
	"caller: %[3]s file: %[4]s line: %[5]d severity: %[2]s message: %[6]s %[7]s: %[8]s %[9]s: %[10]d %[11]s: %[12]t time: %[1]s",
	time.RFC1123Z,
	false)

var templateFormatTestTime = time.Date(2024, time.November, 1, 20, 15, 0, 0, time.UTC)
var templateFormatTestTimeText = templateFormatTestTime.Format(time.RFC1123Z)

func TestTemplateFormat(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "severity: DEBUG message: Testmessage time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "severity: INFO  message: Testmessage time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "severity: WARN  message: Testmessage time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "severity: ERROR message: Testmessage time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "severity: FATAL message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correleation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG correlation: someCorrelationId message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  correlation: someCorrelationId message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  correlation: someCorrelationId message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR correlation: someCorrelationId message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL correlation: someCorrelationId message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correleation, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelationOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correleation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "severity: DEBUG correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "severity: INFO  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "severity: WARN  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "severity: ERROR correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "severity: FATAL correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correleation, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustom(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG message: Testmessage first: abc second: 1 third: true",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  message: Testmessage first: abc second: 1 third: true",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  message: Testmessage first: abc second: 1 third: true",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR message: Testmessage first: abc second: 1 third: true",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL message: Testmessage first: abc second: 1 third: true",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomDefaultFormat(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	var templateFormatterDefaultCustom Formatter = CreateTemplateFormatter(
		"time: %s severity: %s message: %s",
		"time: %s severity: %s correlation: %s message: %s",
		DEFAULT_TEMPLATE,
		"time: %s severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		DEFAULT_CALLER_TEMPLATE,
		time.RFC1123Z,
		false)

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "[" + templateFormatTestTimeText + "] DEBUG: Testmessage [first]: abc [second]: 1 [third]: true",
		common.INFORMATION_SEVERITY: "[" + templateFormatTestTimeText + "] INFO : Testmessage [first]: abc [second]: 1 [third]: true",
		common.WARNING_SEVERITY:     "[" + templateFormatTestTimeText + "] WARN : Testmessage [first]: abc [second]: 1 [third]: true",
		common.ERROR_SEVERITY:       "[" + templateFormatTestTimeText + "] ERROR: Testmessage [first]: abc [second]: 1 [third]: true",
		common.FATAL_SEVERITY:       "[" + templateFormatTestTimeText + "] FATAL: Testmessage [first]: abc [second]: 1 [third]: true",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expexpectedMessage, templateFormatterDefaultCustom.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "severity: DEBUG message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "severity: INFO  message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "severity: WARN  message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "severity: ERROR message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "severity: FATAL message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatTrimSeverity(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	var templateFormatterTrim Formatter = CreateTemplateFormatter(
		"time: %s severity: %s message: %s",
		"time: %s severity: %s correlation: %s message: %s",
		"time: %s severity: %s message: %s %s: %s %s: %d %s: %t",
		"time: %s severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s caller: %s file: %s line: %d message: %s %s: %s %s: %d %s: %t",
		time.RFC1123Z,
		true)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expexpectedMessage, templateFormatterTrim.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG caller: someFunction file: someFile line: 42 message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  caller: someFunction file: someFile line: 42 message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  caller: someFunction file: someFile line: 42 message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR caller: someFunction file: someFile line: 42 message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL caller: someFunction file: someFile line: 42 message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: DEBUG message: Testmessage time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "caller: someFunction file: someFile line: 42 severity: INFO  message: Testmessage time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "caller: someFunction file: someFile line: 42 severity: WARN  message: Testmessage time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: ERROR message: Testmessage time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: FATAL message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelationCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correleation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correleation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelationOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correleation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: DEBUG correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "caller: someFunction file: someFile line: 42 severity: INFO  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "caller: someFunction file: someFile line: 42 severity: WARN  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: ERROR correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: FATAL correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correleation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG caller: someFunction file: someFile line: 42 message: Testmessage first: abc second: 1 third: true",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  caller: someFunction file: someFile line: 42 message: Testmessage first: abc second: 1 third: true",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  caller: someFunction file: someFile line: 42 message: Testmessage first: abc second: 1 third: true",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR caller: someFunction file: someFile line: 42 message: Testmessage first: abc second: 1 third: true",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL caller: someFunction file: someFile line: 42 message: Testmessage first: abc second: 1 third: true",
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expexpectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: DEBUG message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "caller: someFunction file: someFile line: 42 severity: INFO  message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "caller: someFunction file: someFile line: 42 severity: WARN  message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: ERROR message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: FATAL message: Testmessage first: abc second: 1 third: true time: " + templateFormatTestTimeText,
	}

	for severity, expexpectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expexpectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func setCallerValues(logValuesToFormat *common.LogValues) {
	logValuesToFormat.CallerFile = "someFile"
	logValuesToFormat.CallerFileLine = 42
	logValuesToFormat.CallerFunction = "someFunction"
	logValuesToFormat.IsCallerSet = true
}
