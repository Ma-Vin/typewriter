package appender

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/common"
)

func TestMultiAppenderWrite(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	appender, buf1, buf2 := createTestMultiAppender()

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage", strings.TrimSpace(buf1.String()), t, "Write")
	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage", strings.TrimSpace(buf2.String()), t, "Write")
}

func TestMultiAppenderWriteWithCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)
	correlation := "someCorrelationId"

	appender, buf1, buf2 := createTestMultiAppender()

	logValuesToFormat := common.CreateLogValuesWithCorrelation(common.INFORMATION_SEVERITY, &correlation, "Testmessage")
	appender.Write(&logValuesToFormat)
	appender.Close()

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - someCorrelationId - Testmessage", strings.TrimSpace(buf1.String()), t, "WriteWithCorrelation")
	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - someCorrelationId - Testmessage", strings.TrimSpace(buf2.String()), t, "WriteWithCorrelation")
}

func TestMultiAppenderWriteCustom(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	appender, buf1, buf2 := createTestMultiAppender()

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)
	appender.Write(&logValuesToFormat)
	appender.Close()

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage - abc - 1 - true", strings.TrimSpace(buf1.String()), t, "WriteCustom")
	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage - abc - 1 - true", strings.TrimSpace(buf2.String()), t, "WriteCustom")
}

func createTestMultiAppender() (Appender, *bytes.Buffer, *bytes.Buffer) {
	subAppender1, buf1 := createTestStandardOutputAppender()
	subAppender2, buf2 := createTestStandardOutputAppender()

	appender := MultiAppender{[]*Appender{&subAppender1, &subAppender2}}

	return appender, buf1, buf2
}
