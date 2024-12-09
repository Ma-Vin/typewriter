package appender

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/common"
	"github.com/ma-vin/typewriter/format"
	"github.com/ma-vin/typewriter/testutil"
)

var testDelimiterFormatter = format.CreateDelimiterFormatter(" - ", time.RFC3339)
var delimiterFormatterTestTime = time.Date(2024, time.November, 30, 19, 0, 0, 0, time.UTC)
var delimiterFormatterTestTimeText = delimiterFormatterTestTime.Format(time.RFC3339)

func TestDefaultIsStdOut(t *testing.T) {
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)

	testutil.AssertEquals(os.Stdout, appender.writer, t, "default output")
}

func TestWrite(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	logValuesToFormat := common.CreateLogValues(common.INFORMATION_SEVERITY, "Testmessage")
	appender.Write(&logValuesToFormat)

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage", strings.TrimSpace(buf.String()), t, "Write")
}

func TestWriteWithCorrelation(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)
	correleation := "someCorrelationId"

	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	logValuesToFormat := common.CreateLogValuesWithCorrelation(common.INFORMATION_SEVERITY, &correleation, "Testmessage")
	appender.Write(&logValuesToFormat)

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - someCorrelationId - Testmessage", strings.TrimSpace(buf.String()), t, "WriteWithCorrelation")
}

func TestWriteCustom(t *testing.T) {
	common.SetLogValuesMockTime(&delimiterFormatterTestTime)

	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	logValuesToFormat := common.CreateLogValuesCustom(common.INFORMATION_SEVERITY, "Testmessage", &customProperties)
	appender.Write(&logValuesToFormat)

	testutil.AssertEquals(delimiterFormatterTestTimeText+" - INFO  - Testmessage - abc - 1 - true", strings.TrimSpace(buf.String()), t, "WriteCustom")
}
