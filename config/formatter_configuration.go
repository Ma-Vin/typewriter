package config

import "slices"

// common properties of all formatter configurations
type CommonFormatterConfig struct {
	Id               string
	FormatterType    string
	IsDefault        bool
	PackageParameter string
	TimeLayout       string
	IsSequenceActive bool
	EnvNamesToLog    []string
}

// Checks whether an other common config equals with this one with respect to type and time layout
func (c *CommonFormatterConfig) Equals(other *CommonFormatterConfig) bool {
	return c.FormatterType == other.FormatterType && c.TimeLayout == other.TimeLayout &&
		c.IsSequenceActive == other.IsSequenceActive && slices.Compare(c.EnvNamesToLog, other.EnvNamesToLog) == 0
}

// Checks whether the current common config is less than a given other one
func (c *CommonFormatterConfig) LessCompareForSort(other *CommonFormatterConfig) bool {
	return (c.IsDefault && !other.IsDefault) || (c.IsDefault == other.IsDefault && c.PackageParameter < other.PackageParameter)
}

// config of a formatter
type FormatterConfig interface {
	// The id from common element
	Id() string
	// The formatter type from common element
	FormatterType() string
	// The default indicator from common element
	IsDefault() bool
	// The package from common element
	PackageParameter() string
	// The time layout from common element
	TimeLayout() string
	// Pointer to the common element
	GetCommon() *CommonFormatterConfig
	// Checks whether two formatter config equals without regarding pointers to formatter, id, defaultType or package
	Equals(other *FormatterConfig) bool
	// creates a copy of the current struct. Any Pointer has to address a new copy also
	CreateFullCopy() FormatterConfig
}

// Configuration representation of a delimiter formatter
type DelimiterFormatterConfig struct {
	Common    *CommonFormatterConfig
	Delimiter string
}

func (c DelimiterFormatterConfig) Id() string {
	return c.Common.Id
}

func (c DelimiterFormatterConfig) FormatterType() string {
	return c.Common.FormatterType
}

func (c DelimiterFormatterConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c DelimiterFormatterConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c DelimiterFormatterConfig) TimeLayout() string {
	return c.Common.TimeLayout
}

func (c DelimiterFormatterConfig) GetCommon() *CommonFormatterConfig {
	return c.Common
}

func (c DelimiterFormatterConfig) Equals(other *FormatterConfig) bool {
	return c.Common.Equals((*other).GetCommon()) && c.Delimiter == (*other).(DelimiterFormatterConfig).Delimiter
}

func (c DelimiterFormatterConfig) CreateFullCopy() FormatterConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	return c
}

// Configuration representation of a template formatter
type TemplateFormatterConfig struct {
	Common                               *CommonFormatterConfig
	Template                             string
	IsDefaultTemplate                    bool
	CallerTemplate                       string
	IsDefaultCallerTemplate              bool
	CorrelationIdTemplate                string
	IsDefaultCorrelationIdTemplate       bool
	CallerCorrelationIdTemplate          string
	IsDefaultCallerCorrelationIdTemplate bool
	CustomTemplate                       string
	IsDefaultCustomTemplate              bool
	CallerCustomTemplate                 string
	IsDefaultCallerCustomTemplate        bool
	TrimSeverityText                     bool
}

func (c TemplateFormatterConfig) Id() string {
	return c.Common.Id
}

func (c TemplateFormatterConfig) FormatterType() string {
	return c.Common.FormatterType
}

func (c TemplateFormatterConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c TemplateFormatterConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c TemplateFormatterConfig) TimeLayout() string {
	return c.Common.TimeLayout
}

func (c TemplateFormatterConfig) GetCommon() *CommonFormatterConfig {
	return c.Common
}

func (c TemplateFormatterConfig) Equals(other *FormatterConfig) bool {
	return c.Common.Equals((*other).GetCommon()) &&
		c.Template == (*other).(TemplateFormatterConfig).Template &&
		c.IsDefaultTemplate == (*other).(TemplateFormatterConfig).IsDefaultTemplate &&
		c.CorrelationIdTemplate == (*other).(TemplateFormatterConfig).CorrelationIdTemplate &&
		c.IsDefaultCorrelationIdTemplate == (*other).(TemplateFormatterConfig).IsDefaultCorrelationIdTemplate &&
		c.CustomTemplate == (*other).(TemplateFormatterConfig).CustomTemplate &&
		c.IsDefaultCustomTemplate == (*other).(TemplateFormatterConfig).IsDefaultCustomTemplate &&
		c.CallerTemplate == (*other).(TemplateFormatterConfig).CallerTemplate &&
		c.IsDefaultCallerTemplate == (*other).(TemplateFormatterConfig).IsDefaultCallerTemplate &&
		c.CallerCorrelationIdTemplate == (*other).(TemplateFormatterConfig).CallerCorrelationIdTemplate &&
		c.IsDefaultCallerCorrelationIdTemplate == (*other).(TemplateFormatterConfig).IsDefaultCallerCorrelationIdTemplate &&
		c.CallerCustomTemplate == (*other).(TemplateFormatterConfig).CallerCustomTemplate &&
		c.IsDefaultCallerCustomTemplate == (*other).(TemplateFormatterConfig).IsDefaultCallerCustomTemplate &&
		c.TrimSeverityText == (*other).(TemplateFormatterConfig).TrimSeverityText
}

func (c TemplateFormatterConfig) CreateFullCopy() FormatterConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	return c
}

// Configuration representation of a json formatter
type JsonFormatterConfig struct {
	Common                   *CommonFormatterConfig
	TimeKey                  string
	SequenceKey              string
	SeverityKey              string
	MessageKey               string
	CorrelationKey           string
	CustomValuesKey          string
	CustomValuesAsSubElement bool
	CallerFunctionKey        string
	CallerFileKey            string
	CallerFileLineKey        string
}

func (c JsonFormatterConfig) Id() string {
	return c.Common.Id
}

func (c JsonFormatterConfig) FormatterType() string {
	return c.Common.FormatterType
}

func (c JsonFormatterConfig) IsDefault() bool {
	return c.Common.IsDefault
}

func (c JsonFormatterConfig) PackageParameter() string {
	return c.Common.PackageParameter
}

func (c JsonFormatterConfig) TimeLayout() string {
	return c.Common.TimeLayout
}

func (c JsonFormatterConfig) GetCommon() *CommonFormatterConfig {
	return c.Common
}

func (c JsonFormatterConfig) Equals(other *FormatterConfig) bool {
	return c.Common.Equals((*other).GetCommon()) &&
		c.TimeKey == (*other).(JsonFormatterConfig).TimeKey &&
		c.SequenceKey == (*other).(JsonFormatterConfig).SequenceKey &&
		c.SeverityKey == (*other).(JsonFormatterConfig).SeverityKey &&
		c.MessageKey == (*other).(JsonFormatterConfig).MessageKey &&
		c.CorrelationKey == (*other).(JsonFormatterConfig).CorrelationKey &&
		c.CustomValuesKey == (*other).(JsonFormatterConfig).CustomValuesKey &&
		c.CustomValuesAsSubElement == (*other).(JsonFormatterConfig).CustomValuesAsSubElement &&
		c.CallerFunctionKey == (*other).(JsonFormatterConfig).CallerFunctionKey &&
		c.CallerFileKey == (*other).(JsonFormatterConfig).CallerFileKey &&
		c.CallerFileLineKey == (*other).(JsonFormatterConfig).CallerFileLineKey
}

func (c JsonFormatterConfig) CreateFullCopy() FormatterConfig {
	commonConfig := *c.Common
	c.Common = &commonConfig
	return c
}
