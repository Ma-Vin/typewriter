package appender

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/config"
	"github.com/ma-vin/typewriter/format"
)

func createDelimiterFormatterForTest() format.Formatter {
	commonConfig := config.CommonFormatterConfig{TimeLayout: time.RFC3339}
	var config config.FormatterConfig = config.DelimiterFormatterConfig{Common: &commonConfig, Delimiter: " - "}
	result, _ := format.CreateDelimiterFormatterFromConfig(&config)
	return *result
}

func CreateStandardOutputAppenderForTest(formatter *format.Formatter) Appender {
	commonConfig := config.CommonAppenderConfig{}
	var config config.AppenderConfig = config.StdOutAppenderConfig{Common: &commonConfig}
	appender, _ := CreateStandardOutputAppenderFromConfig(&config, formatter)
	return *appender
}

var testDelimiterFormatter = createDelimiterFormatterForTest()
var delimiterFormatterTestTime = time.Date(2024, time.November, 30, 19, 0, 0, 0, time.UTC)
var delimiterFormatterTestTimeText = delimiterFormatterTestTime.Format(time.RFC3339)

func TestStandardOutputAppenderDefaultIsStdOut(t *testing.T) {
	appender := CreateStandardOutputAppenderForTest(&testDelimiterFormatter).(StandardOutputAppender)

	testutil.AssertEquals(os.Stdout, appender.writer, t, "default output")
}

func TestStandardOutputAppenderWrite(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	appender, buf := createTestStandardOutputAppender()

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()
	testutil.AssertTrue(*appender.(StandardOutputAppender).isClosed, t, "isClosed")
	appender.Close()

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage", strings.TrimSpace(buf.String()), t, "Write")
}

func TestStandardOutputAppenderWriteWithCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)
	correlation := "someCorrelationId"

	appender, buf := createTestStandardOutputAppender()

	logValuesToFormat := common.CreateLogValuesWithCorrelation(common.INFORMATION_SEVERITY, &correlation, "Testmessage")
	appender.Write(&logValuesToFormat)

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - someCorrelationId - Testmessage", strings.TrimSpace(buf.String()), t, "WriteWithCorrelation")
}

func TestStandardOutputAppenderWriteCustom(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	appender, buf := createTestStandardOutputAppender()

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)
	appender.Write(&logValuesToFormat)

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage - abc - 1 - true", strings.TrimSpace(buf.String()), t, "WriteCustom")
}

func TestStandardOutputAppenderClose(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	appender, buf := createTestStandardOutputAppender()

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	testutil.AssertFalse(*appender.(StandardOutputAppender).isClosed, t, "isNotClosed")
	appender.Close()
	appender.Write(&logValuesToFormat)
	testutil.AssertTrue(*appender.(StandardOutputAppender).isClosed, t, "isClosed")

	testutil.AssertEquals("", strings.TrimSpace(buf.String()), t, "Write")
}


func createTestStandardOutputAppender() (Appender, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppenderForTest(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	return appender, buf
}
