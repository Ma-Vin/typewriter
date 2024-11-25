package logger

import (
	"os"
	"time"

	"github.com/ma-vin/typewriter/format"
)

func ExampleLog_withDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will not be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be ", "printed")

	// Output:
	// 2024-11-25T09:30:00+01:00 - ERROR - Error will be printed
	// 2024-11-25T09:30:00+01:00 - FATAL - Fatal will be printed
}

func ExampleLog_formatWithDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	Log().Debugf("Debug %s %s %s %s", "will", "not", "be", "printed")
	Log().Informationf("Information %s %s %s %s", "will", "not", "be", "printed")
	Log().Warningf("Warning %s %s %s %s", "will", "not", "be ", "printed")
	Log().Errorf("Error %s %s %s", "will", "be", "printed")
	Log().Fatalf("Fatal %s %s %s", "will", "be", "printed")

	// Output:
	// 2024-11-25T09:30:00+01:00 - ERROR - Error will be printed
	// 2024-11-25T09:30:00+01:00 - FATAL - Fatal will be printed
}

func ExampleLog_correlationIdWithDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	Log().DebugWithCorrelation("CorrelationId123", "Debug will not be printed")
	Log().InformationWithCorrelation("CorrelationId123", "Information will not be printed")
	Log().WarningWithCorrelation("CorrelationId123", "Warning will not be printed")
	Log().ErrorWithCorrelation("CorrelationId123", "Error will be ", "printed")
	Log().FatalWithCorrelation("CorrelationId123", "Fatal will be printed")

	// Output:
	// 2024-11-25T09:30:00+01:00 - ERROR - CorrelationId123 - Error will be printed
	// 2024-11-25T09:30:00+01:00 - FATAL - CorrelationId123 - Fatal will be printed
}

func ExampleLog_customValuesWithDefaultConfiguration() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
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
	// 2024-11-25T09:30:00+01:00 - ERROR - Error will be printed - firstEntry - true - 1.2
	// 2024-11-25T09:30:00+01:00 - FATAL - Fatal will be printed - firstEntry - true - 1.2
}

func ExampleLog_enableAllLevels() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	os.Setenv("TYPEWRITER_LOG_LEVEL", "DEBUG")

	Log().Debug("Debug will be printed")
	Log().Information("Information will be printed")
	Log().Warning("Warning will be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T09:30:00+01:00 - DEBUG - Debug will be printed
	// 2024-11-25T09:30:00+01:00 - INFO  - Information will be printed
	// 2024-11-25T09:30:00+01:00 - WARN  - Warning will be printed
	// 2024-11-25T09:30:00+01:00 - ERROR - Error will be printed
	// 2024-11-25T09:30:00+01:00 - FATAL - Fatal will be printed
}

func ExampleLog_levelRestrictedByPackage() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	os.Setenv("TYPEWRITER_LOG_LEVEL", "DEBUG")
	os.Setenv("TYPEWRITER_PACKAGE_LOG_LEVEL_LOGGER", "WARN")

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T09:30:00+01:00 - WARN  - Warning will be printed
	// 2024-11-25T09:30:00+01:00 - ERROR - Error will be printed
	// 2024-11-25T09:30:00+01:00 - FATAL - Fatal will be printed
}

func ExampleLog_otherDdelimiter() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "DELIMITER")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER", ",")

	Log().Debug("Debug will not be printed")
	Log().Information("Information will not be printed")
	Log().Warning("Warning will not be printed")
	Log().Error("Error will be printed")
	Log().Fatal("Fatal will be printed")

	// Output:
	// 2024-11-25T09:30:00+01:00,ERROR,Error will be printed
	// 2024-11-25T09:30:00+01:00,FATAL,Fatal will be printed
}

func ExampleLog_jsonFormatDefaultKeys() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
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
	// {"message":"Warning will be printed","severity":"WARN","time":"2024-11-25T09:30:00+01:00"}
	// {"correlation":"CorrelationId123","message":"Error will be printed","severity":"ERROR","time":"2024-11-25T09:30:00+01:00"}
	// {"a":"firstEntry","b":true,"c":1.2,"message":"Fatal will be printed","severity":"FATAL","time":"2024-11-25T09:30:00+01:00"}
}

func ExampleLog_jsonFormatCustomKeys() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	os.Setenv("TYPEWRITER_LOG_LEVEL", "WARN")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "JSON")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_1", "the_time")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_2", "log_level")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_3", "correlation_id")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_4", "msg")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_5", "my_values")
	// Parameter 6 creates a subelement for the custom values
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_6", "true")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_7", time.RFC822)

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
	// {"log_level":"WARN","msg":"Warning will be printed","the_time":"25 Nov 24 09:30 CET"}
	// {"correlation_id":"CorrelationId123","log_level":"ERROR","msg":"Error will be printed","the_time":"25 Nov 24 09:30 CET"}
	// {"log_level":"FATAL","msg":"Fatal will be printed","my_values":{"a":"firstEntry","b":true,"c":1.2},"the_time":"25 Nov 24 09:30 CET"}
}

func ExampleLog_templateFormatDefaultTemplates() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
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
	// [2024-11-25T09:30:00+01:00] WARN : Warning will be printed
	// [2024-11-25T09:30:00+01:00] ERROR CorrelationId123: Error will be printed
	// [2024-11-25T09:30:00+01:00] FATAL: Fatal will be printed [a]: firstEntry [b]: true [c]: 1.2
}

func ExampleLog_templateFormatCustomTemplates() {
	// Clear test environment
	os.Clearenv()
	Reset()
	// use fixed time for output compare
	mockTime := time.Date(2024, 11, 25, 8, 30, 0, 0, time.FixedZone("DE", 0))
	format.SetFormatterMockTime(&mockTime)

	// Example beginn
	os.Setenv("TYPEWRITER_LOG_LEVEL", "WARN")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_TYPE", "TEMPLATE")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_1", "time=[%s], severity=[%s], msg=[%s]")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_2", "time=[%[1]s], severity=[%[2]s], msg=[%[4]s] correlation[%[3]s]")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_3", "time=[%s], severity=[%s], msg=[%s], %s=[%s], %s=[%t], %s=[%g]")
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_4", time.RFC822)
	os.Setenv("TYPEWRITER_LOG_FORMATTER_PARAMETER_5", "true")

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
	// time=[25 Nov 24 09:30 CET], severity=[WARN], msg=[Warning will be printed]
	// time=[25 Nov 24 09:30 CET], severity=[ERROR], msg=[Error will be printed] correlation[CorrelationId123]
	// time=[25 Nov 24 09:30 CET], severity=[FATAL], msg=[Fatal will be printed], a=[firstEntry], b=[true], c=[1.2]
}
