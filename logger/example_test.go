package logger

import (
	"os"
	"time"

	"github.com/ma-vin/typewriter/common"
)

func ExampleLog_withDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will not be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be ", "printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - ERROR - Error will be printed
	// 2024-11-25T08:30:00Z - 2 - FATAL - Fatal will be printed
}

func ExampleLog_withDefaultConfigurationWithoutInterface() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	Debug("Debug will not be printed")
	Information("Information will not be printed")
	Warning("Warning will not be printed")
	Error("Error will be printed")
	Fatal("Fatal will be ", "printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - ERROR - Error will be printed
	// 2024-11-25T08:30:00Z - 2 - FATAL - Fatal will be printed
}

func ExampleLog_formatWithDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	Log().Debugf("Debug %s %s %s %s", "will", "not", "be", "printed")
	Log().Informationf("Information %s %s %s %s", "will", "not", "be", "printed")
	Log().Warningf("Warning %s %s %s %s", "will", "not", "be ", "printed")
	Log().Errorf("Error %s %s %s", "will", "be", "printed")
	Log().Fatalf("Fatal %s %s %s", "will", "be", "printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - ERROR - Error will be printed
	// 2024-11-25T08:30:00Z - 2 - FATAL - Fatal will be printed
}

func ExampleLog_correlationIdWithDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	Log().DebugWithCorrelation("CorrelationId123", "Debug will not be printed")
	Log().InformationWithCorrelation("CorrelationId123", "Information will not be printed")
	Log().WarningWithCorrelation("CorrelationId123", "Warning will not be printed")
	Log().ErrorWithCorrelation("CorrelationId123", "Error will be ", "printed")
	Log().FatalWithCorrelation("CorrelationId123", "Fatal will be printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - ERROR - CorrelationId123 - Error will be printed
	// 2024-11-25T08:30:00Z - 2 - FATAL - CorrelationId123 - Fatal will be printed
}

func ExampleLog_customValuesWithDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	customValueMap := make(map[string]any, 3)
	customValueMap["b"] = true
	customValueMap["c"] = 1.2
	customValueMap["a"] = "firstEntry"

	Log().DebugCustom(customValueMap, "Debug will not be printed")
	Log().InformationCustom(customValueMap, "Information will not be printed")
	Log().WarningCustom(customValueMap, "Warning will not be printed")
	Log().ErrorCustom(customValueMap, "Error will be ", "printed")
	Log().FatalCustom(customValueMap, "Fatal will be printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - ERROR - Error will be printed - firstEntry - true - 1.2
	// 2024-11-25T08:30:00Z - 2 - FATAL - Fatal will be printed - firstEntry - true - 1.2
}

func ExampleLog_enableAllLevels() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "DEBUG")

	Log().Debug("Debug will be printed")
	Log().Information("Information will be printed")
	Log().Warning("Warning will be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - DEBUG - Debug will be printed
	// 2024-11-25T08:30:00Z - 2 - INFO  - Information will be printed
	// 2024-11-25T08:30:00Z - 3 - WARN  - Warning will be printed
	// 2024-11-25T08:30:00Z - 4 - ERROR - Error will be printed
	// 2024-11-25T08:30:00Z - 5 - FATAL - Fatal will be printed
}

func ExampleLog_levelRestrictedByPackage() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "DEBUG")
	os.Setenv("TYPEWRITER_PACKAGE_LOG_PACKAGE_LOGGER", "logger")
	os.Setenv("TYPEWRITER_PACKAGE_LOG_LEVEL_LOGGER", "WARN")

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - WARN  - Warning will be printed
	// 2024-11-25T08:30:00Z - 2 - ERROR - Error will be printed
	// 2024-11-25T08:30:00Z - 3 - FATAL - Fatal will be printed
}

func ExampleLog_levelRestrictedByPackageFullQualified() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "DEBUG")
	os.Setenv("TYPEWRITER_PACKAGE_FULL_QUALIFIED", "TRUE")
	os.Setenv("TYPEWRITER_PACKAGE_LOG_PACKAGE_LOGGER", "github.com/ma-vin/typewriter/logger")
	os.Setenv("TYPEWRITER_PACKAGE_LOG_LEVEL_LOGGER", "WARN")

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T08:30:00Z - 1 - WARN  - Warning will be printed
	// 2024-11-25T08:30:00Z - 2 - ERROR - Error will be printed
	// 2024-11-25T08:30:00Z - 3 - FATAL - Fatal will be printed
}

func ExampleLog_otherDelimiter() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "DELIMITER")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_DELIMITER", ",")

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will not be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T08:30:00Z,1,ERROR,Error will be printed
	// 2024-11-25T08:30:00Z,2,FATAL,Fatal will be printed
}

func ExampleLog_jsonFormatDefaultKeys() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "WARN")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "JSON")

	customValueMap := make(map[string]any, 3)
	customValueMap["b"] = true
	customValueMap["c"] = 1.2
	customValueMap["a"] = "firstEntry"

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().ErrorWithCorrelation("CorrelationId123", "Error will be printed")
	Log().FatalCustom(customValueMap, "Fatal will be printed")

	// Output:
	// {"message":"Warning will be printed","sequence":1,"severity":"WARN","time":"2024-11-25T08:30:00Z"}
	// {"correlation":"CorrelationId123","message":"Error will be printed","sequence":2,"severity":"ERROR","time":"2024-11-25T08:30:00Z"}
	// {"a":"firstEntry","b":true,"c":1.2,"message":"Fatal will be printed","sequence":3,"severity":"FATAL","time":"2024-11-25T08:30:00Z"}
}

func ExampleLog_jsonFormatCustomKeys() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "WARN")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "JSON")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_TIME_KEY", "the_time")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_SEQUENCE_KEY", "seq")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_SEVERITY_KEY", "log_level")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_CORRELATION_KEY", "correlation_id")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_MESSAGE_KEY", "msg")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_CUSTOM_VALUES_KEY", "my_values")
	// Parameter 6 creates a sub element for the custom values
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_JSON_CUSTOM_VALUES_SUB", "true")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TIME_LAYOUT", time.RFC822)

	customValueMap := make(map[string]any, 3)
	customValueMap["b"] = true
	customValueMap["c"] = 1.2
	customValueMap["a"] = "firstEntry"

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().ErrorWithCorrelation("CorrelationId123", "Error will be printed")
	Log().FatalCustom(customValueMap, "Fatal will be printed")

	// Output:
	// {"log_level":"WARN","msg":"Warning will be printed","seq":1,"the_time":"25 Nov 24 08:30 DE"}
	// {"correlation_id":"CorrelationId123","log_level":"ERROR","msg":"Error will be printed","seq":2,"the_time":"25 Nov 24 08:30 DE"}
	// {"log_level":"FATAL","msg":"Fatal will be printed","my_values":{"a":"firstEntry","b":true,"c":1.2},"seq":3,"the_time":"25 Nov 24 08:30 DE"}
}

func ExampleLog_templateFormatDefaultTemplates() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "WARN")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "TEMPLATE")

	customValueMap := make(map[string]any, 3)
	customValueMap["b"] = true
	customValueMap["c"] = 1.2
	customValueMap["a"] = "firstEntry"

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().ErrorWithCorrelation("CorrelationId123", "Error will be printed")
	Log().FatalCustom(customValueMap, "Fatal will be printed")

	// Output:
	// [2024-11-25T08:30:00Z] 1 WARN : Warning will be printed
	// [2024-11-25T08:30:00Z] 2 ERROR CorrelationId123: Error will be printed
	// [2024-11-25T08:30:00Z] 3 FATAL: Fatal will be printed [a]: firstEntry [b]: true [c]: 1.2
}

func ExampleLog_templateFormatCustomTemplates() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_LEVEL", "WARN")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "TEMPLATE")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TEMPLATE", "time=[%s], severity=[%s], msg=[%s]")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TEMPLATE_CORRELATION", "time=[%[1]s], severity=[%[2]s], msg=[%[4]s] correlation[%[3]s]")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TEMPLATE_CUSTOM", "time=[%s], severity=[%s], msg=[%s], %s=[%s], %s=[%t], %s=[%g]")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TIME_LAYOUT", time.RFC822)
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_SEQUENCE_ACTIVE", "false")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TEMPLATE_TRIM_SEVERITY", "true")

	customValueMap := make(map[string]any, 3)
	customValueMap["b"] = true
	customValueMap["c"] = 1.2
	customValueMap["a"] = "firstEntry"

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().ErrorWithCorrelation("CorrelationId123", "Error will be printed")
	Log().FatalCustom(customValueMap, "Fatal will be printed")

	// Output:
	// time=[25 Nov 24 08:30 DE], severity=[WARN], msg=[Warning will be printed]
	// time=[25 Nov 24 08:30 DE], severity=[ERROR], msg=[Error will be printed] correlation[CorrelationId123]
	// time=[25 Nov 24 08:30 DE], severity=[FATAL], msg=[Fatal will be printed], a=[firstEntry], b=[true], c=[1.2]
}

func ExampleLog_callerWithIndexedTemplate() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	common.SetLogValuesMockTime(&mockTime)

	// Example begin
	os.Setenv("TYPEWRITER_LOG_CALLER", "true")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "TEMPLATE")
	// Ignore file, because the system depended path can not compared
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_TEMPLATE_CALLER", "time:%[1]s, severity:%[3]s, caller:%[4]s line:%[6]d sequence:%[2]d msg:%[7]s")

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will not be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be ", "printed")

	// Output:
	// time:2024-11-25T08:30:00Z, severity:ERROR, caller:github.com/ma-vin/typewriter/logger.ExampleLog_callerWithIndexedTemplate line:362 sequence:1 msg:Error will be printed
	// time:2024-11-25T08:30:00Z, severity:FATAL, caller:github.com/ma-vin/typewriter/logger.ExampleLog_callerWithIndexedTemplate line:363 sequence:2 msg:Fatal will be printed
}
