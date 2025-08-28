package config

import (
	"testing"

	"github.com/ma-vin/testutil-go"
	"github.com/ma-vin/typewriter/common"
)

func TestEqualsGeneralLoggerConfig(t *testing.T) {
	commonConfig := CommonLoggerConfig{LoggerType: LOGGER_GENERAL}
	var config LoggerConfig = GeneralLoggerConfig{Common: &commonConfig, Severity: common.ERROR_SEVERITY, IsCallerToSet: true}

	testutil.AssertTrue(config.Equals(&config), t, "general same instance")
}

func TestNotEqualsGeneralLoggerConfig(t *testing.T) {
	commonConfig := CommonLoggerConfig{LoggerType: LOGGER_GENERAL}
	config := GeneralLoggerConfig{Common: &commonConfig, Severity: common.ERROR_SEVERITY, IsCallerToSet: true}

	checkNotEqualsGeneralLoggerConfig(&config, func(otherConfig *GeneralLoggerConfig) { otherConfig.Common.LoggerType = "other" }, "general LoggerType diff", t)
	checkNotEqualsGeneralLoggerConfig(&config, func(otherConfig *GeneralLoggerConfig) { otherConfig.Severity = common.DEBUG_SEVERITY }, "general Severity diff", t)
	checkNotEqualsGeneralLoggerConfig(&config, func(otherConfig *GeneralLoggerConfig) { otherConfig.IsCallerToSet = false }, "general IsCallerToSet diff", t)
}

func checkNotEqualsGeneralLoggerConfig(config *GeneralLoggerConfig, modifier func(otherConfig *GeneralLoggerConfig), message string, t *testing.T) {
	otherConfig := config.CreateFullCopy().(GeneralLoggerConfig)

	modifier(&otherConfig)

	var castOtherConfig LoggerConfig = otherConfig
	testutil.AssertFalse((*config).Equals(&castOtherConfig), t, message)
}
