package config

// common properties of all logger configurations
type CommonLoggerConfig struct {
	Id               string
	LoggerType       string
	IsDefault        bool
	PackageParameter string
	PackageName      string
}

// Checks whether an other common config equals with this one with respect to type
func (c *CommonLoggerConfig) Equals(other *CommonLoggerConfig) bool {
	return c.LoggerType == other.LoggerType
}

// Checks whether the current common config is less than a given other one
func (c *CommonLoggerConfig) LessCompareForSort(other *CommonLoggerConfig) bool {
	return (c.IsDefault && !other.IsDefault) || (c.IsDefault == other.IsDefault && c.PackageParameter < other.PackageParameter)
}

// config of a logger
type LoggerConfig interface {
	// The id from common element
	Id() string
	// The default indicator from common element
	IsDefault() bool
	// The package from common element
	PackageParameter() string
	// The package name from common element
	PackageName() string
	// Pointer to the common element
	GetCommon() *CommonLoggerConfig
	// Checks whether two formatter config equals without regarding id, defaultType or package
	Equals(other *LoggerConfig) bool
	// creates a copy of the current struct. Any Pointer has to address a new copy also
	CreateFullCopy() LoggerConfig
}

// common properties of general logger configurations
type GeneralLoggerConfig struct {
	Common        *CommonLoggerConfig
	Severity      int
	IsCallerToSet bool
}

func (c GeneralLoggerConfig) Id() string {
	return c.Common.Id
}

func (c GeneralLoggerConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c GeneralLoggerConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c GeneralLoggerConfig) PackageName() string {
	return c.Common.PackageName
}

func (c GeneralLoggerConfig) GetCommon() *CommonLoggerConfig {
	return c.Common
}

func (c GeneralLoggerConfig) Equals(other *LoggerConfig) bool {
	return c.Common.Equals((*other).GetCommon()) &&
		c.Severity == (*other).(GeneralLoggerConfig).Severity && c.IsCallerToSet == (*other).(GeneralLoggerConfig).IsCallerToSet
}

func (c GeneralLoggerConfig) CreateFullCopy() LoggerConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	return c
}
