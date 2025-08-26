package config

import (
	"testing"

	"github.com/ma-vin/testutil-go"
)

func TestEqualsDelimiterFormatterConfig(t *testing.T) {
	commonConfig := CommonFormatterConfig{FormatterType: "FormatterType", TimeLayout: "TimeLayout"}
	var config FormatterConfig = DelimiterFormatterConfig{Common: &commonConfig, Delimiter: "Delimiter"}

	testutil.AssertTrue(config.Equals(&config), t, "delimiter same instance")
}

func TestNotEqualsDelimiterFormatterConfig(t *testing.T) {
	commonConfig := CommonFormatterConfig{FormatterType: "FormatterType", TimeLayout: "TimeLayout"}
	config := DelimiterFormatterConfig{Common: &commonConfig, Delimiter: "Delimiter"}

	checkNotEqualsDelimiterFormatterConfig(&config, func(otherConfig *DelimiterFormatterConfig) { otherConfig.Common.FormatterType = "other" }, "delimiter FormatterType diff", t)
	checkNotEqualsDelimiterFormatterConfig(&config, func(otherConfig *DelimiterFormatterConfig) { otherConfig.Common.TimeLayout = "other" }, "delimiter TimeLayout diff", t)
	checkNotEqualsDelimiterFormatterConfig(&config, func(otherConfig *DelimiterFormatterConfig) { otherConfig.Delimiter = "Other" }, "delimiter Delimiter diff", t)
}

func checkNotEqualsDelimiterFormatterConfig(config *DelimiterFormatterConfig, modifier func(otherConfig *DelimiterFormatterConfig), message string, t *testing.T) {
	otherConfig := config.CreateFullCopy().(DelimiterFormatterConfig)

	modifier(&otherConfig)

	var castOtherConfig FormatterConfig = otherConfig
	testutil.AssertFalse((*config).Equals(&castOtherConfig), t, message)
}

func TestEqualsTemplateFormatterConfig(t *testing.T) {
	commonConfig := CommonFormatterConfig{FormatterType: "FormatterType", TimeLayout: "TimeLayout"}
	var config FormatterConfig = TemplateFormatterConfig{
		Common:                      &commonConfig,
		Template:                    "Template",
		CorrelationIdTemplate:       "CorrelationIdTemplate",
		CustomTemplate:              "CustomTemplate",
		CallerTemplate:              "CallerTemplate",
		CallerCorrelationIdTemplate: "CallerCorrelationIdTemplate",
		CallerCustomTemplate:        "CallerCustomTemplate",
		TrimSeverityText:            true,
	}

	testutil.AssertTrue(config.Equals(&config), t, "template same instance")
}

func TestNotEqualsTemplateFormatterConfig(t *testing.T) {
	commonConfig := CommonFormatterConfig{FormatterType: "FormatterType", TimeLayout: "TimeLayout"}
	config := TemplateFormatterConfig{
		Common:                      &commonConfig,
		Template:                    "Template",
		CorrelationIdTemplate:       "CorrelationIdTemplate",
		CustomTemplate:              "CustomTemplate",
		CallerTemplate:              "CallerTemplate",
		CallerCorrelationIdTemplate: "CallerCorrelationIdTemplate",
		CallerCustomTemplate:        "CallerCustomTemplate",
		TrimSeverityText:            true,
	}

	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.Common.FormatterType = "other" }, "template FormatterType diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.Common.TimeLayout = "other" }, "template TimeLayout diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.Template = "other" }, "template Template diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.CorrelationIdTemplate = "other" }, "template CorrelationIdTemplate diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.CustomTemplate = "other" }, "template CustomTemplate diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.CallerTemplate = "other" }, "template CallerTemplate diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.CallerCorrelationIdTemplate = "other" }, "template CallerCorrelationIdTemplate diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.CallerCustomTemplate = "other" }, "template CallerCustomTemplate diff", t)
	checkNotEqualsTemplateFormatterConfig(&config, func(otherConfig *TemplateFormatterConfig) { otherConfig.TrimSeverityText = false }, "template TrimSeverityText diff", t)
}

func checkNotEqualsTemplateFormatterConfig(config *TemplateFormatterConfig, modifier func(otherConfig *TemplateFormatterConfig), message string, t *testing.T) {
	otherConfig := config.CreateFullCopy().(TemplateFormatterConfig)

	modifier(&otherConfig)

	var castOtherConfig FormatterConfig = otherConfig
	testutil.AssertFalse((*config).Equals(&castOtherConfig), t, message)
}

func TestEqualsJsonFormatterConfig(t *testing.T) {
	commonConfig := CommonFormatterConfig{FormatterType: "FormatterType", TimeLayout: "TimeLayout"}
	var config FormatterConfig = JsonFormatterConfig{
		Common:                   &commonConfig,
		TimeKey:                  "TimeKey",
		SeverityKey:              "SeverityKey",
		MessageKey:               "MessageKey",
		CorrelationKey:           "CorrelationKey",
		CustomValuesKey:          "CustomValuesKey",
		CustomValuesAsSubElement: true,
		CallerFunctionKey:        "CallerFunctionKey",
		CallerFileKey:            "CallerFileKey",
		CallerFileLineKey:        "CallerFileLineKey",
	}

	testutil.AssertTrue(config.Equals(&config), t, "json same instance")
}

func TestNotEqualsJsonFormatterConfig(t *testing.T) {
	commonConfig := CommonFormatterConfig{FormatterType: "FormatterType", TimeLayout: "TimeLayout"}
	config := JsonFormatterConfig{
		Common:                   &commonConfig,
		TimeKey:                  "TimeKey",
		SeverityKey:              "SeverityKey",
		MessageKey:               "MessageKey",
		CorrelationKey:           "CorrelationKey",
		CustomValuesKey:          "CustomValuesKey",
		CustomValuesAsSubElement: true,
		CallerFunctionKey:        "CallerFunctionKey",
		CallerFileKey:            "CallerFileKey",
		CallerFileLineKey:        "CallerFileLineKey",
	}

	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.Common.FormatterType = "other" }, "json FormatterType diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.Common.TimeLayout = "other" }, "json TimeLayout diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.TimeKey = "other" }, "json TimeKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.SeverityKey = "other" }, "json SeverityKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.MessageKey = "other" }, "json MessageKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.CorrelationKey = "other" }, "json CorrelationKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.CustomValuesKey = "other" }, "json CustomValuesKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.CustomValuesAsSubElement = false }, "json CustomValuesAsSubElement diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.CallerFunctionKey = "other" }, "json CallerFunctionKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.CallerFileKey = "other" }, "json CallerFileKey diff", t)
	checkNotEqualsJsonFormatterConfig(&config, func(otherConfig *JsonFormatterConfig) { otherConfig.CallerFileLineKey = "other" }, "json CallerFileLineKey diff", t)
}

func checkNotEqualsJsonFormatterConfig(config *JsonFormatterConfig, modifier func(otherConfig *JsonFormatterConfig), message string, t *testing.T) {
	otherConfig := config.CreateFullCopy().(JsonFormatterConfig)

	modifier(&otherConfig)

	var castOtherConfig FormatterConfig = otherConfig
	testutil.AssertFalse((*config).Equals(&castOtherConfig), t, message)
}
