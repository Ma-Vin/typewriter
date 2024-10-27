package appender

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/format"
	testutil "github.com/ma-vin/typewriter/util"
)

var testTime = time.Date(2024, time.October, 1, 13, 20, 0, 0, time.UTC)
var testTimeText = testTime.Local().Format(time.RFC3339)

var testDelimiterFormatter = format.CreateDelimiterFormatter(" - ")

func TestDefaultIsStdOut(t *testing.T) {
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)

	testutil.AssertEquals(os.Stdout, appender.writer, t, "default output")
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	appender.Write(testTime, constants.INFORMATION_SEVERITY, "Testmessage")

	testutil.AssertEquals(testTimeText+" - INFO  - Testmessage", strings.TrimSpace(buf.String()), t, "Write")
}

func TestWriteWithCorrelation(t *testing.T) {
	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	appender.WriteWithCorrelation(testTime, constants.INFORMATION_SEVERITY, "someCorrelationId", "Testmessage")

	testutil.AssertEquals(testTimeText+" - INFO  - someCorrelationId - Testmessage", strings.TrimSpace(buf.String()), t, "WriteWithCorrelation")
}

func TestWriteCustom(t *testing.T) {
	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	customProperties := map[string]any{
		"first":  "abc",
		"third":  true,
		"second": 1,
	}

	appender.WriteCustom(testTime, constants.INFORMATION_SEVERITY, "Testmessage", customProperties)

	testutil.AssertEquals(testTimeText+" - INFO  - Testmessage - abc - 1 - true", strings.TrimSpace(buf.String()), t, "WriteCustom")
}
