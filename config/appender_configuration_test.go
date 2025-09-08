package config

import (
	"slices"
	"testing"

	"github.com/ma-vin/testutil-go"
)

func TestEqualsStdOutAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: APPENDER_STDOUT}
	var config AppenderConfig = StdOutAppenderConfig{Common: &commonConfig}

	testutil.AssertTrue(config.Equals(&config), t, "stdout same instance")
}

func TestNotEqualsStdOutAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: APPENDER_STDOUT}
	var config AppenderConfig = StdOutAppenderConfig{Common: &commonConfig}

	otherConfig := config.CreateFullCopy().(StdOutAppenderConfig)
	otherConfig.Common.AppenderType = "other"
	var castOtherConfig AppenderConfig = otherConfig

	testutil.AssertFalse(config.Equals(&castOtherConfig), t, "stdout AppenderType diff")
}

func TestEqualsFileAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: APPENDER_FILE}
	var config AppenderConfig = FileAppenderConfig{
		Common:         &commonConfig,
		PathToLogFile:  "PathToLogFile",
		CronExpression: "CronExpression",
		LimitByteSize:  "LimitByteSize",
	}

	testutil.AssertTrue(config.Equals(&config), t, "file same instance")
}

func TestNotEqualsFileAppenderConfig(t *testing.T) {
	commonConfig := CommonAppenderConfig{AppenderType: APPENDER_FILE}
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

func TestEqualsMultiAppenderConfig(t *testing.T) {
	stdOutCommonConfig := CommonAppenderConfig{AppenderType: APPENDER_STDOUT}
	var stdOutConfig AppenderConfig = StdOutAppenderConfig{Common: &stdOutCommonConfig}

	fileCommonConfig := CommonAppenderConfig{AppenderType: APPENDER_FILE}
	var fileConfig AppenderConfig = FileAppenderConfig{
		Common:         &fileCommonConfig,
		PathToLogFile:  "PathToLogFile",
		CronExpression: "CronExpression",
		LimitByteSize:  "LimitByteSize",
	}

	multiCommonConfig := CommonAppenderConfig{AppenderType: APPENDER_MULTIPLE}
	var multiConfig AppenderConfig = MultiAppenderConfig{
		Common:    &multiCommonConfig,
		AppenderConfigs: &[]AppenderConfig{stdOutConfig, fileConfig},
	}

	testutil.AssertTrue(multiConfig.Equals(&multiConfig), t, "multi same instance")
}

func TestNotEqualsMultiAppenderConfig(t *testing.T) {
	stdOutCommonConfig := CommonAppenderConfig{AppenderType: APPENDER_STDOUT}
	var stdOutConfig AppenderConfig = StdOutAppenderConfig{Common: &stdOutCommonConfig}

	fileCommonConfig := CommonAppenderConfig{AppenderType: APPENDER_FILE}
	var fileConfig AppenderConfig = FileAppenderConfig{
		Common:         &fileCommonConfig,
		PathToLogFile:  "PathToLogFile",
		CronExpression: "CronExpression",
		LimitByteSize:  "LimitByteSize",
	}

	multiCommonConfig := CommonAppenderConfig{AppenderType: APPENDER_MULTIPLE}
	var multiConfig = MultiAppenderConfig{
		Common:    &multiCommonConfig,
		AppenderConfigs: &[]AppenderConfig{stdOutConfig, fileConfig},
	}

	checkNotEqualsMultiAppenderConfig(&multiConfig, func(otherConfig *MultiAppenderConfig) { otherConfig.Common.AppenderType = "other" }, "file AppenderType diff", t)
	checkNotEqualsMultiAppenderConfig(&multiConfig, func(otherConfig *MultiAppenderConfig) {
		removed := slices.Delete(*otherConfig.AppenderConfigs, 0, 1)
		otherConfig.AppenderConfigs = &removed
	}, "first removed", t)
	checkNotEqualsMultiAppenderConfig(&multiConfig, func(otherConfig *MultiAppenderConfig) {
		removed := slices.Delete(*otherConfig.AppenderConfigs, 1, 2)
		otherConfig.AppenderConfigs = &removed
	}, "last removed", t)
	checkNotEqualsMultiAppenderConfig(&multiConfig, func(otherConfig *MultiAppenderConfig) {
		(*otherConfig.AppenderConfigs)[0].(StdOutAppenderConfig).Common.AppenderType = "other"
	}, "first diff", t)
	checkNotEqualsMultiAppenderConfig(&multiConfig, func(otherConfig *MultiAppenderConfig) {
		modified := (*otherConfig.AppenderConfigs)[1].(FileAppenderConfig)
		modified.PathToLogFile = "other"
		(*otherConfig.AppenderConfigs)[1] = modified
	}, "last diff", t)
}

func checkNotEqualsMultiAppenderConfig(config *MultiAppenderConfig, modifier func(otherConfig *MultiAppenderConfig), message string, t *testing.T) {
	otherConfig := config.CreateFullCopy().(MultiAppenderConfig)

	modifier(&otherConfig)

	var castOtherConfig AppenderConfig = otherConfig
	testutil.AssertFalse((*config).Equals(&castOtherConfig), t, message)
}
