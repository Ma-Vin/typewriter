package config

import (
	"testing"

	"github.com/ma-vin/testutil-go"
)

func TestEqualsStdOutAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: "AppenderType"}
	var config AppenderConfig = StdOutAppenderConfig{Common: &commonConfig}

	testutil.AssertTrue(config.Equals(&config), t, "stdout same instance")
}

func TestNotEqualsStdOutAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: "AppenderType"}
	var config AppenderConfig = StdOutAppenderConfig{Common: &commonConfig}

	otherConfig := config.CreateFullCopy().(StdOutAppenderConfig)
	otherConfig.Common.AppenderType = "other"
	var castOtherConfig AppenderConfig = otherConfig

	testutil.AssertFalse(config.Equals(&castOtherConfig), t, "stdout AppenderType diff")
}

func TestEqualsFileAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: "AppenderType"}
	var config AppenderConfig = FileAppenderConfig{
		Common:         &commonConfig,
		PathToLogFile:  "PathToLogFile",
		CronExpression: "CronExpression",
		LimitByteSize:  "LimitByteSize",
	}

	testutil.AssertTrue(config.Equals(&config), t, "file same instance")
}

func TestNotEqualsFileAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: "AppenderType"}
	config := FileAppenderConfig{
		Common:         &commonConfig,
		PathToLogFile:  "PathToLogFile",
		CronExpression: "CronExpression",
		LimitByteSize:  "LimitByteSize",
	}

	checkNotEqualsFileAppenderConfig(&config, func(otherConfig *FileAppenderConfig) { otherConfig.Common.AppenderType = "other" }, "file AppenderType diff", t)
	checkNotEqualsFileAppenderConfig(&config, func(otherConfig *FileAppenderConfig) { otherConfig.PathToLogFile = "Other" }, "file PathToLogFile diff", t)
	checkNotEqualsFileAppenderConfig(&config, func(otherConfig *FileAppenderConfig) { otherConfig.CronExpression = "Other" }, "file CronExpression diff", t)
	checkNotEqualsFileAppenderConfig(&config, func(otherConfig *FileAppenderConfig) { otherConfig.LimitByteSize = "Other" }, "file LimitByteSize diff", t)
}

func checkNotEqualsFileAppenderConfig(config *FileAppenderConfig, modifier func(otherConfig *FileAppenderConfig), message string, t *testing.T) {
	otherConfig := config.CreateFullCopy().(FileAppenderConfig)

	modifier(&otherConfig)

	var castOtherConfig AppenderConfig = otherConfig
	testutil.AssertFalse((*config).Equals(&castOtherConfig), t, message)
}
