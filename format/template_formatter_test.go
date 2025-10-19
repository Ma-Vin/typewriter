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

// Creates a new formatter with given templates and time layout
func createTemplateFormatterForTest(template string, correlationIdTemplate string, customTemplate string,
	callerTemplate string, callerCorrelationIdTemplate string, callerCustomTemplate string,
	timeLayout string, trimSeverityText bool, isSequenceActive bool, envNamesToLog []string) Formatter {

	commonConfig := config.CommonFormatterConfig{TimeLayout: timeLayout, IsSequenceActive: isSequenceActive, EnvNamesToLog: envNamesToLog}
	var config config.FormatterConfig = config.TemplateFormatterConfig{
		Common:                               &commonConfig,
		Template:                             template,
		IsDefaultTemplate:                    (!isSequenceActive && template == config.DEFAULT_TEMPLATE) || (isSequenceActive && template == config.DEFAULT_SEQUENCE_TEMPLATE),
		CallerTemplate:                       callerTemplate,
		IsDefaultCallerTemplate:              (!isSequenceActive && callerTemplate == config.DEFAULT_CALLER_TEMPLATE) || (isSequenceActive && callerTemplate == config.DEFAULT_SEQUENCE_CALLER_TEMPLATE),
		CorrelationIdTemplate:                correlationIdTemplate,
		IsDefaultCorrelationIdTemplate:       (!isSequenceActive && correlationIdTemplate == config.DEFAULT_CORRELATION_TEMPLATE) || (isSequenceActive && correlationIdTemplate == config.DEFAULT_SEQUENCE_CORRELATION_TEMPLATE),
		CallerCorrelationIdTemplate:          callerCorrelationIdTemplate,
		IsDefaultCallerCorrelationIdTemplate: (!isSequenceActive && callerCorrelationIdTemplate == config.DEFAULT_CALLER_CORRELATION_TEMPLATE) || (isSequenceActive && callerCorrelationIdTemplate == config.DEFAULT_SEQUENCE_CALLER_CORRELATION_TEMPLATE),
		CustomTemplate:                       customTemplate,
		IsDefaultCustomTemplate:              (!isSequenceActive && customTemplate == config.DEFAULT_CUSTOM_TEMPLATE) || (isSequenceActive && customTemplate == config.DEFAULT_SEQUENCE_CUSTOM_TEMPLATE),
		CallerCustomTemplate:                 callerCustomTemplate,
		IsDefaultCallerCustomTemplate:        (!isSequenceActive && callerCustomTemplate == config.DEFAULT_CALLER_CUSTOM_TEMPLATE) || (isSequenceActive && callerCustomTemplate == config.DEFAULT_SEQUENCE_CALLER_CUSTOM_TEMPLATE),
		TrimSeverityText:                     trimSeverityText,
	}
	result, _ := CreateTemplateFormatterFromConfig(&config)
	return *result
}

func createTemplateFormatterWithEnvValuesForTest(template string, correlationIdTemplate string, customTemplate string,
	callerTemplate string, callerCorrelationIdTemplate string, callerCustomTemplate string,
	timeLayout string, trimSeverityText bool, isSequenceActive bool) Formatter {

	os.Clearenv()
	os.Setenv("test1", "abc")
	os.Setenv("TEST2", "1")
	os.Setenv("Test3", "2.1")
	os.Setenv("test4", "true")

	return createTemplateFormatterForTest(template, correlationIdTemplate, customTemplate, callerTemplate, callerCorrelationIdTemplate,
		callerCustomTemplate, timeLayout, trimSeverityText, isSequenceActive, []string{"test1", "TEST2", "Test3", "test4"})
}

var templateFormatter Formatter = createTemplateFormatterForTest(
	"time: %s severity: %s message: %s",
	"time: %s severity: %s correlation: %s message: %s",
	"time: %s severity: %s message: %s %s: %s %s: %d %s: %t",
	"time: %s severity: %s caller: %s file: %s line: %d message: %s",
	"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
	"time: %s severity: %s caller: %s file: %s line: %d message: %s %s: %s %s: %d %s: %t",
	time.RFC1123Z,
	false,
	false,
	[]string{})

var templateSequenceFormatter Formatter = createTemplateFormatterForTest(
	"time: %s sequence: %d severity: %s message: %s",
	"time: %s sequence: %d severity: %s correlation: %s message: %s",
	"time: %s sequence: %d severity: %s message: %s %s: %s %s: %d %s: %t",
	"time: %s sequence: %d severity: %s caller: %s file: %s line: %d message: %s",
	"time: %s sequence: %d severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
	"time: %s sequence: %d severity: %s caller: %s file: %s line: %d message: %s %s: %s %s: %d %s: %t",
	time.RFC1123Z,
	false,
	true,
	[]string{})

var templateFormatterOrder Formatter = createTemplateFormatterForTest(
	"severity: %[2]s message: %[3]s time: %[1]s",
	"severity: %[2]s correlation: %[3]s message: %[4]s time: %[1]s",
	"severity: %[2]s message: %[3]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %[9]t time: %[1]s",
	"caller: %[3]s file: %[4]s line: %[5]d severity: %[2]s message: %[6]s time: %[1]s",
	"caller: %[4]s file: %[5]s line: %[6]d severity: %[2]s correlation: %[3]s message: %[7]s time: %[1]s",
	"caller: %[3]s file: %[4]s line: %[5]d severity: %[2]s message: %[6]s %[7]s: %[8]s %[9]s: %[10]d %[11]s: %[12]t time: %[1]s",
	time.RFC1123Z,
	false,
	false,
	[]string{})

var templateSequenceFormatterOrder Formatter = createTemplateFormatterForTest(
	"severity: %[3]s message: %[4]s time: %[1]s sequence: %[2]d",
	"severity: %[3]s correlation: %[4]s message: %[5]s time: %[1]s sequence: %[2]d",
	"severity: %[3]s message: %[4]s %[5]s: %[6]s %[7]s: %[8]d %[9]s: %[10]t time: %[1]s sequence: %[2]d",
	"caller: %[4]s file: %[5]s line: %[6]d severity: %[3]s message: %[7]s time: %[1]s sequence: %[2]d",
	"caller: %[5]s file: %[6]s line: %[7]d severity: %[3]s correlation: %[4]s message: %[8]s time: %[1]s sequence: %[2]d",
	"caller: %[4]s file: %[5]s line: %[6]d severity: %[3]s message: %[7]s %[8]s: %[9]s %[10]s: %[11]d %[12]s: %[13]t time: %[1]s sequence: %[2]d",
	time.RFC1123Z,
	false,
	true,
	[]string{})

var templateFormatterDefaultEnvValues Formatter = createTemplateFormatterWithEnvValuesForTest(
	config.DEFAULT_TEMPLATE,
	config.DEFAULT_CORRELATION_TEMPLATE,
	config.DEFAULT_CUSTOM_TEMPLATE,
	config.DEFAULT_CALLER_TEMPLATE,
	config.DEFAULT_CALLER_CORRELATION_TEMPLATE,
	config.DEFAULT_CALLER_CUSTOM_TEMPLATE,
	time.RFC1123Z,
	false,
	false)

var templateFormatterEnvValuesOrder Formatter = createTemplateFormatterWithEnvValuesForTest(
	"time: %[1]s severity: %[2]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %.2[9]f %[10]s: %[11]t message: %[3]s",
	"time: %[1]s severity: %[2]s correlation: %[3]s %[5]s: %[6]s %[7]s: %[8]d %[9]s: %.2[10]f %[11]s: %[12]t message: %[4]s",
	"time: %[1]s severity: %[2]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %.2[9]f %[10]s: %[11]t %[12]s: %[13]s %[14]s: %[15]d %[16]s: %[17]t message: %[3]s",
	"time: %[1]s severity: %[2]s caller: %[3]s file: %[4]s line: %[5]d %[7]s: %[8]s %[9]s: %[10]d %[11]s: %.2[12]f %[13]s: %[14]t message: %[6]s",
	"time: %[1]s severity: %[2]s correlation: %[3]s caller: %[4]s file: %[5]s line: %[6]d %[8]s: %[9]s %[10]s: %[11]d %[12]s: %.2[13]f %[14]s: %[15]t message: %[7]s",
	"time: %[1]s severity: %[2]s caller: %[3]s file: %[4]s line: %[5]d %[7]s: %[8]s %[9]s: %[10]d %[11]s: %.2[12]f %[13]s: %[14]t %[15]s: %[16]s %[17]s: %[18]d %[19]s: %[20]t message: %[6]s",
	time.RFC1123Z,
	false,
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequence(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		expectedMessage := fmt.Sprintf("time: %s sequence: %d severity: %s message: Testmessage", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateSequenceFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatSequenceOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		expectedMessage := fmt.Sprintf("severity: %s message: Testmessage time: %s sequence: %d", severityTextMap[i], templateFormatTestTimeText, i)
		testutil.AssertEquals(expectedMessage, templateSequenceFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatEnvValues(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		expectedMessage := fmt.Sprintf("[%s] %s: Testmessage [test1]: abc [TEST2]: 1 [Test3]: 2.1 [test4]: true", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultEnvValues.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatEnvValuesOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		expectedMessage := fmt.Sprintf("time: %s severity: %s test1: abc TEST2: 1 Test3: 2.10 test4: true message: Testmessage", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterEnvValuesOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correlation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG correlation: someCorrelationId message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  correlation: someCorrelationId message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  correlation: someCorrelationId message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR correlation: someCorrelationId message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL correlation: someCorrelationId message: Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correlation, "Testmessage")
		testutil.AssertEquals(expectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelationOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correlation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "severity: DEBUG correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "severity: INFO  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "severity: WARN  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "severity: ERROR correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "severity: FATAL correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correlation, "Testmessage")
		testutil.AssertEquals(expectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequenceCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		expectedMessage := fmt.Sprintf("time: %s sequence: %d severity: %s correlation: someCorrelationId message: Testmessage", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateSequenceFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatSequenceCorrelationOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		expectedMessage := fmt.Sprintf("severity: %s correlation: someCorrelationId message: Testmessage time: %s sequence: %d", severityTextMap[i], templateFormatTestTimeText, i)
		testutil.AssertEquals(expectedMessage, templateSequenceFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCorrelationEnvValues(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		expectedMessage := fmt.Sprintf("[%s] %s %s: Testmessage [test1]: abc [TEST2]: 1 [Test3]: 2.1 [test4]: true", templateFormatTestTimeText, severityTextMap[i], correlation)
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultEnvValues.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCorrelationEnvValuesOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		expectedMessage := fmt.Sprintf("time: %s severity: %s correlation: %s test1: abc TEST2: 1 Test3: 2.10 test4: true message: Testmessage", templateFormatTestTimeText, severityTextMap[i], correlation)
		testutil.AssertEquals(expectedMessage, templateFormatterEnvValuesOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomDefaultFormat(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	var templateFormatterDefaultCustom Formatter = createTemplateFormatterForTest(
		"time: %s severity: %s message: %s",
		"time: %s severity: %s correlation: %s message: %s",
		config.DEFAULT_CUSTOM_TEMPLATE,
		"time: %s severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		config.DEFAULT_CALLER_CUSTOM_TEMPLATE,
		time.RFC1123Z,
		false,
		false,
		[]string{})

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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultCustom.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequenceCustom(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		expectedMessage := fmt.Sprintf("time: %s sequence: %d severity: %s message: Testmessage first: abc second: 1 third: true", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateSequenceFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatSequenceCustomDefaultFormat(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	var templateFormatterDefaultCustom Formatter = createTemplateFormatterForTest(
		"time: %s sequence %d severity: %s message: %s",
		"time: %s sequence %d severity: %s correlation: %s message: %s",
		config.DEFAULT_SEQUENCE_CUSTOM_TEMPLATE,
		"time: %s sequence %d severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s sequence %d severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		config.DEFAULT_SEQUENCE_CALLER_CUSTOM_TEMPLATE,
		time.RFC1123Z,
		false,
		true,
		[]string{})

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		expectedMessage := fmt.Sprintf("[%s] %d %s: Testmessage [first]: abc [second]: 1 [third]: true", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultCustom.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		testutil.AssertEquals(expectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCustomEnvValues(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		expectedMessage := fmt.Sprintf("[%s] %s: Testmessage [test1]: abc [TEST2]: 1 [Test3]: 2.1 [test4]: true [first]: abc [second]: 1 [third]: true", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultEnvValues.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCustomEnvValuesOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		expectedMessage := fmt.Sprintf("time: %s severity: %s test1: abc TEST2: 1 Test3: 2.10 test4: true first: abc second: 1 third: true message: Testmessage", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterEnvValuesOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatTrimSeverity(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	var templateFormatterTrim Formatter = createTemplateFormatterForTest(
		"time: %s severity: %s message: %s",
		"time: %s severity: %s correlation: %s message: %s",
		"time: %s severity: %s message: %s %s: %s %s: %d %s: %t",
		"time: %s severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s caller: %s file: %s line: %d message: %s %s: %s %s: %d %s: %t",
		time.RFC1123Z,
		true,
		false,
		[]string{})

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL message: Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		testutil.AssertEquals(expectedMessage, templateFormatterTrim.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValues(severity, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequenceCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("time: %s sequence: %d severity: %s caller: someFunction file: someFile line: 42 message: Testmessage", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateSequenceFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatSequenceOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("caller: someFunction file: someFile line: 42 severity: %s message: Testmessage time: %s sequence: %d", severityTextMap[i], templateFormatTestTimeText, i)
		testutil.AssertEquals(expectedMessage, templateSequenceFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCallerEnvValues(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("[%s] %s someFunction(someFile.42): Testmessage [test1]: abc [TEST2]: 1 [Test3]: 2.1 [test4]: true", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultEnvValues.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCallerEnvValuesOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValues(i, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("time: %s severity: %s caller: someFunction file: someFile line: 42 test1: abc TEST2: 1 Test3: 2.10 test4: true message: Testmessage", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterEnvValuesOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCorrelationCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correlation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: DEBUG correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.INFORMATION_SEVERITY: "time: " + templateFormatTestTimeText + " severity: INFO  correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.WARNING_SEVERITY:     "time: " + templateFormatTestTimeText + " severity: WARN  correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.ERROR_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: ERROR correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
		common.FATAL_SEVERITY:       "time: " + templateFormatTestTimeText + " severity: FATAL correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correlation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatCorrelationOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	correlation := "someCorrelationId"

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: DEBUG correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.INFORMATION_SEVERITY: "caller: someFunction file: someFile line: 42 severity: INFO  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.WARNING_SEVERITY:     "caller: someFunction file: someFile line: 42 severity: WARN  correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.ERROR_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: ERROR correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
		common.FATAL_SEVERITY:       "caller: someFunction file: someFile line: 42 severity: FATAL correlation: someCorrelationId message: Testmessage time: " + templateFormatTestTimeText,
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(severity, &correlation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequenceCorrelationCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("time: %s sequence: %d severity: %s correlation: someCorrelationId caller: someFunction file: someFile line: 42 message: Testmessage", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateSequenceFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatSequenceCorrelationOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("caller: someFunction file: someFile line: 42 severity: %s correlation: someCorrelationId message: Testmessage time: %s sequence: %d", severityTextMap[i], templateFormatTestTimeText, i)
		testutil.AssertEquals(expectedMessage, templateSequenceFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCallerCorrelationEnvValues(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("[%s] %s %s someFunction(someFile.42): Testmessage [test1]: abc [TEST2]: 1 [Test3]: 2.1 [test4]: true", templateFormatTestTimeText, severityTextMap[i], correlation)
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultEnvValues.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCallerCorrelationEnvValuesOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	correlation := "someCorrelationId"

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesWithCorrelation(i, &correlation, "Testmessage")
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("time: %s severity: %s correlation: %s caller: someFunction file: someFile line: 42 test1: abc TEST2: 1 Test3: 2.10 test4: true message: Testmessage", templateFormatTestTimeText, severityTextMap[i], correlation)
		testutil.AssertEquals(expectedMessage, templateFormatterEnvValuesOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCustomCallerDefaultFormat(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)

	var templateFormatterDefaultCustom Formatter = createTemplateFormatterForTest(
		"time: %s severity: %s message: %s",
		"time: %s severity: %s correlation: %s message: %s",
		config.DEFAULT_CUSTOM_TEMPLATE,
		"time: %s severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		config.DEFAULT_CALLER_CUSTOM_TEMPLATE,
		time.RFC1123Z,
		false,
		false,
		[]string{})

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	expectedResults := map[int]string{
		common.DEBUG_SEVERITY:       "[" + templateFormatTestTimeText + "] DEBUG someFunction(someFile.42): Testmessage [first]: abc [second]: 1 [third]: true",
		common.INFORMATION_SEVERITY: "[" + templateFormatTestTimeText + "] INFO  someFunction(someFile.42): Testmessage [first]: abc [second]: 1 [third]: true",
		common.WARNING_SEVERITY:     "[" + templateFormatTestTimeText + "] WARN  someFunction(someFile.42): Testmessage [first]: abc [second]: 1 [third]: true",
		common.ERROR_SEVERITY:       "[" + templateFormatTestTimeText + "] ERROR someFunction(someFile.42): Testmessage [first]: abc [second]: 1 [third]: true",
		common.FATAL_SEVERITY:       "[" + templateFormatTestTimeText + "] FATAL someFunction(someFile.42): Testmessage [first]: abc [second]: 1 [third]: true",
	}

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultCustom.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequenceCustomCallerDefaultFormat(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	var templateFormatterDefaultCustom Formatter = createTemplateFormatterForTest(
		"time: %s sequence %d severity: %s message: %s",
		"time: %s sequence %d severity: %s correlation: %s message: %s",
		config.DEFAULT_SEQUENCE_CUSTOM_TEMPLATE,
		"time: %s sequence %d severity: %s caller: %s file: %s line: %d message: %s",
		"time: %s sequence %d severity: %s correlation: %s caller: %s file: %s line: %d message: %s",
		config.DEFAULT_SEQUENCE_CALLER_CUSTOM_TEMPLATE,
		time.RFC1123Z,
		false,
		true,
		[]string{})

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("[%s] %d %s someFunction(someFile.42): Testmessage [first]: abc [second]: 1 [third]: true", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultCustom.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
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

	for severity, expectedMessage := range expectedResults {
		logValuesToFormat := common.CreateLogValuesCustom(severity, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		testutil.AssertEquals(expectedMessage, templateFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", severity))
	}
}

func TestTemplateFormatSequenceCustomCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("time: %s sequence: %d severity: %s caller: someFunction file: someFile line: 42 message: Testmessage first: abc second: 1 third: true", templateFormatTestTimeText, i, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateSequenceFormatter.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatSequenceCustomOrderCaller(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("caller: someFunction file: someFile line: 42 severity: %s message: Testmessage first: abc second: 1 third: true time: %s sequence: %d", severityTextMap[i], templateFormatTestTimeText, i)
		testutil.AssertEquals(expectedMessage, templateSequenceFormatterOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCallerCustomEnvValues(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("[%s] %s someFunction(someFile.42): Testmessage [test1]: abc [TEST2]: 1 [Test3]: 2.1 [test4]: true [first]: abc [second]: 1 [third]: true", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterDefaultEnvValues.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func TestTemplateFormatCallerCustomEnvValuesOrder(t *testing.T) {
	common.SetLogValuesMockTime(&templateFormatTestTime)
	common.InitSequenceCounter()
	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	for i := 1; i <= 5; i++ {
		logValuesToFormat := common.CreateLogValuesCustom(i, "Testmessage", &customProperties)
		setCallerValues(&logValuesToFormat)
		expectedMessage := fmt.Sprintf("time: %s severity: %s caller: someFunction file: someFile line: 42 test1: abc TEST2: 1 Test3: 2.10 test4: true first: abc second: 1 third: true message: Testmessage", templateFormatTestTimeText, severityTextMap[i])
		testutil.AssertEquals(expectedMessage, templateFormatterEnvValuesOrder.Format(&logValuesToFormat), t, fmt.Sprintf("Format severity %d", i))
	}
}

func setCallerValues(logValuesToFormat *common.LogValues) {
	logValuesToFormat.CallerFile = "someFile"
	logValuesToFormat.CallerFileLine = 42
	logValuesToFormat.CallerFunction = "someFunction"
	logValuesToFormat.IsCallerSet = true
}
