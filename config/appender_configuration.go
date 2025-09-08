package config

// common properties of all appender configurations
type CommonAppenderConfig struct {
	Id               string
	AppenderType     string
	IsDefault        bool
	PackageParameter string
}

// Checks whether an other common config equals with this one with respect to type
func (c *CommonAppenderConfig) Equals(other *CommonAppenderConfig) bool {
	return c.AppenderType == other.AppenderType
}

// Checks whether the current common config is less than a given other one
func (c *CommonAppenderConfig) LessCompareForSort(other *CommonAppenderConfig) bool {
	return (c.IsDefault && !other.IsDefault) || (c.IsDefault == other.IsDefault && c.PackageParameter < other.PackageParameter)
}

// config of an appender
type AppenderConfig interface {
	// The id from common element
	Id() string
	// The appender type from common element
	AppenderType() string
	// The default indicator from common element
	IsDefault() bool
	// The package from common element
	PackageParameter() string
	// Pointer to the common element
	GetCommon() *CommonAppenderConfig
	// Checks whether two appender config equals without regarding pointers to appender, id, defaultType or package
	Equals(other *AppenderConfig) bool
	// creates a copy of the current struct. Any Pointer has to address a new copy also
	CreateFullCopy() AppenderConfig
}

// Configuration representation of a file appender
type FileAppenderConfig struct {
	Common         *CommonAppenderConfig
	PathToLogFile  string
	CronExpression string
	LimitByteSize  string
}

func (c FileAppenderConfig) Id() string {
	return c.Common.Id
}

func (c FileAppenderConfig) AppenderType() string {
	return c.Common.AppenderType
}

func (c FileAppenderConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c FileAppenderConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c FileAppenderConfig) GetCommon() *CommonAppenderConfig {
	return c.Common
}

func (c FileAppenderConfig) Equals(other *AppenderConfig) bool {
	return c.Common.Equals((*other).GetCommon()) &&
		c.PathToLogFile == (*other).(FileAppenderConfig).PathToLogFile &&
		c.CronExpression == (*other).(FileAppenderConfig).CronExpression &&
		c.LimitByteSize == (*other).(FileAppenderConfig).LimitByteSize
}

func (c FileAppenderConfig) CreateFullCopy() AppenderConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	return c
}

// Configuration representation of a standard output appender
type StdOutAppenderConfig struct {
	Common *CommonAppenderConfig
}

func (c StdOutAppenderConfig) Id() string {
	return c.Common.Id
}

func (c StdOutAppenderConfig) AppenderType() string {
	return c.Common.AppenderType
}

func (c StdOutAppenderConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c StdOutAppenderConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c StdOutAppenderConfig) GetCommon() *CommonAppenderConfig {
	return c.Common
}

func (c StdOutAppenderConfig) Equals(other *AppenderConfig) bool {
	return c.Common.Equals((*other).GetCommon())
}

func (c StdOutAppenderConfig) CreateFullCopy() AppenderConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	return c
}

type MultiAppenderConfig struct {
	Common          *CommonAppenderConfig
	AppenderConfigs *[]AppenderConfig
}

func (c MultiAppenderConfig) Id() string {
	return c.Common.Id
}

func (c MultiAppenderConfig) AppenderType() string {
	return c.Common.AppenderType
}

func (c MultiAppenderConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c MultiAppenderConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c MultiAppenderConfig) GetCommon() *CommonAppenderConfig {
	return c.Common
}

func (c MultiAppenderConfig) Equals(other *AppenderConfig) bool {
	if !c.Common.Equals((*other).GetCommon()) || len(*c.AppenderConfigs) != len(*(*other).(MultiAppenderConfig).AppenderConfigs) {
		return false
	}
	for _, a1 := range *c.AppenderConfigs {
		exist := false
		for _, a2 := range *(*other).(MultiAppenderConfig).AppenderConfigs {
			if a1.Equals(&a2) {
				exist = true
				break
			}
		}
		if !exist {
			return false
		}
	}
	return true
}

func (c MultiAppenderConfig) CreateFullCopy() AppenderConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	appenderConfigs := make([]AppenderConfig, len(*c.AppenderConfigs))
	for i, a := range *c.AppenderConfigs {
		appenderConfigs[i] = a.CreateFullCopy()
	}
	c.AppenderConfigs = &appenderConfigs
	return c
}
