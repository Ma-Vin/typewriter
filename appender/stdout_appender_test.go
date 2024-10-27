package appender

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/ma-vin/typewriter/constants"
	"github.com/ma-vin/typewriter/format"
	testutil "github.com/ma-vin/typewriter/util"
)

var testDelimiterFormatter = format.CreateDelimiterFormatter(" - ")

func TestDefaultIsStdOut(t *testing.T) {
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)

	testutil.AssertEquals(os.Stdout, appender.writer, t, "default output")
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	appender.Write(constants.INFORMATION_SEVERITY, "Testmessage")

	testutil.AssertHasSuffix(" - INFO  - Testmessage", strings.TrimSpace(buf.String()), t, "Write")
}

func TestWriteWithCorrelation(t *testing.T) {
	buf := new(bytes.Buffer)
	appender := CreateStandardOutputAppender(&testDelimiterFormatter).(StandardOutputAppender)
	appender.writer = buf

	appender.WriteWithCorrelation(constants.INFORMATION_SEVERITY, "someCorrelationId", "Testmessage")

	testutil.AssertHasSuffix(" - INFO  - someCorrelationId - Testmessage", strings.TrimSpace(buf.String()), t, "WriteWithCorrelation")
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

	appender.WriteCustom(constants.INFORMATION_SEVERITY, "Testmessage", customProperties)

	testutil.AssertHasSuffix(" - INFO  - Testmessage - abc - 1 - true", strings.TrimSpace(buf.String()), t, "WriteCustom")
}
