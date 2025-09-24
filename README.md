[![Build and test](https://github.com/Ma-Vin/typewriter/actions/workflows/go-build.yaml/badge.svg)](https://github.com/Ma-Vin/typewriter/actions/workflows/go-build.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ma-vin/typewriter.svg)](https://pkg.go.dev/github.com/ma-vin/typewriter)
[![Go Report Card](https://goreportcard.com/badge/github.com/ma-vin/typewriter)](https://goreportcard.com/report/github.com/ma-vin/typewriter)

# Typewriter

This repository provides a GoLang logging package, which are inspired by some Log4j features like package specific enablement.

## Sonar analysis

* [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)
* [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)  [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)
* [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)  [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)
* [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)  [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)  [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)
* [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)
* [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)  [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ma-vin%3Atypewriter&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ma-vin%3Atypewriter)

## Usage

### Logging

#### Access

The logger can be accessed via importing `github.com/ma-vin/typewriter/logger` and using the `Log()` method at `logger` package. This will return a struct implementing the `Logger` interface.
There are also the same methods directly in the `logger` package, but without returning an interface.
If you use `Log()` in combination with `Close()` afterwards, keep in mind that the appenders used by the logger returned by `Log()` might referenced a closed appender (A closed appender will not write any log entries).

#### Log Methods

There are permutations of the following categories for each severity

* log without additional parameters, e.g. `Information(args ...any)`
* log with correlation id, e.g. `InformationWithCorrelation(correlationId string, args ...any)`
* log with custom values map, e.g. `InformationCustom(customValues map[string]any, args ...any)`
* log with correlation id from context, e.g. `InformationCtx(context context.Context, args ...any)`
* message format with args, e.g. `Informationf(format string, args ...any)`
* panic (only for warning, error and fatal) e.g. `ErrorWithPanic(args ...any)`
* exit (only for fatal) e.g. `FatalWithExit(args ...any)`

Arguments are handled in the manner of `fmt.Sprint`/`fmt.Sprintf`

In addition, there exists an indicator if the severity is enabled or not.

### Configuration

Configuration is currently possible via environment variables or an properties file.

#### Configuration by file

To use a configuration file by setting an environment variable `TYPEWRITER_CONFIG_FILE` with a path to a properties file.
If this environment variable is set, all others will be ignored. Nevertheless, the same property names are used by both.
It is possible to use single line comments `#`,`//` or `--` and multiline comments with begin at `/*` and end at `*/`.

#### Log levels

The log level can be set via `TYPEWRITER_LOG_LEVEL` with one of the following values:

* `DEBUG`
* `INFO` or `INFORMATION`
* `WARN` or `WARNING`
* `ERROR`
* `FATAL`

The default log level is `ERROR`

#### Appender

There are two types of appender at the moment. They can be configured by setting `TYPEWRITER_LOG_APPENDER_TYPE` with one of the following values:

* `STDOUT`: Standard output appender
* `FILE`: File appender which writes to a file
  * The target file of the appender has to be configured by setting the path to `TYPEWRITER_LOG_APPENDER_FILE`
  * Renaming of the log file with existing entries at a specific point in time can be configured with a cron expression at `TYPEWRITER_LOG_APPENDER_CRON_RENAMING` (allowed characters are `* , - /`. The day of week is zero based)
  * Renaming of the log file with existing entries, if a limit of bytes will be exceeded by next log entry, can be configured by an integer at `TYPEWRITER_LOG_APPENDER_SIZE_RENAMING` (`kb`/`KB` or `mb`/`MB` can be added). This value will be ignored if `TYPEWRITER_LOG_APPENDER_CRON_RENAMING` is set also.
* Your id of a registered custom appender
* It is possible to use multiple appenders in parallel by setting them as comma separated list

The default appender is the standard output one

#### Formatter

There are three types of formatter which provides the texts to log for the appender. They can be configured by setting `TYPEWRITER_LOG_FORMATTER_TYPE` with one of the following values:

* `DELIMITER`: The record information will be appended and delimited with a given sign. They can be set `TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>` where `<x>` has to be replaced by the following values:
  1. `DELIMITER` the delimiting signs. The default value of the delimiter is ` - `.
  2. `TIME_LAYOUT` time layout. Default value of `time.RFC3339`
* `TEMPLATE`: The records will be derived from three templates and time layout. They can be set `TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>` where `<x>` has to be replaced by the following values:
  1. `TEMPLATE` template for writing time, severity and the message. Default `[%s] %s: %s`
  2. `TEMPLATE_CORRELATION` template for writing time, severity, correlationID and the message. Default `[%s] %s %s: %s`
  3. `TEMPLATE_CUSTOM` template for writing time, severity, message and custom value map. Default `[%s] %s: %s`
  4. `TIME_LAYOUT` time layout. Default value of `time.RFC3339`
  5. `TEMPLATE_TRIM_SEVERITY` indicator whether to trim severity text or to add space at warn and info to algin following elements. Default `false`
  6. `TEMPLATE_CALLER` like 1. with caller function, file and line placed in front of message. Default `[%s] %s %s(%s.%d): %s`
  7. `TEMPLATE_CALLER_CORRELATION` like 2. with caller function, file and line placed in front of message. Default `[%s] %s %s %s(%s.%d): %s`
  8. `TEMPLATE_CALLER_CUSTOM` like 3. with caller function, file and line placed in front of message. Default `[%s] %s %s(%s.%d): %s`

  It is possible to reorder parameter by argument indices. The *custom value map* at *3* will be appended as key-value-pairs sorted by key (e.g. a custom map with three entries of string, number and boolean format at indices from 4 to 9: `severity: %[2]s message: %[3]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %[9]t time: %[1]s`)
* `JSON`: The records will be logged as *JSON*. It is possible to define key names, the time layout and if the custom value map should be a sub element or not. The keys of the custom value map will be used 1:1. The properties can be set via `TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>` where `<x>` has to be replaced by the following values:
  1. `JSON_TIME_KEY` key of time. Default `time`
  2. `JSON_SEVERITY_KEY` key of severity. Default: `severity`
  3. `JSON_CORRELATION_KEY` key of correlationID. Default: `correlation`
  4. `JSON_MESSAGE_KEY` key of message. Default:  `message`
  5. `JSON_CUSTOM_VALUES_KEY` key of custom values if used as sub elements. Default: `custom`
  6. `JSON_CUSTOM_VALUES_SUB` indicator to add custom value map as sub element: Default: `false`
  7. `TIME_LAYOUT` time layout. Default value of `time.RFC3339`
  8. `JSON_CALLER_FUNCTION_KEY` key of the caller function. Default: `caller`
  9. `JSON_CALLER_FILE_KEY` key of the caller file. Default: `file`
  10. `JSON_CALLER_LINE_KEY` key of the caller file line. Default: `line`

The default formatter is the delimiter one.

#### Caller

The logging of the caller function, file and line can be activated by setting `true` at `TYPEWRITER_LOG_CALLER`. The default is `false`.

#### Context

The key, which is used to get the correlation id from `context.Context.Value(key)`, can be defined by setting `TYPEWRITER_CONTEXT_CORRELATION_ID_KEY`.
The default value is `correlationId`. If `context.Context.Value(key)` returns `nil`, a log statement without correlation id will be used.

#### Package specific configuration

Each configuration element can be declared package specific by replacing `TYPEWRITER` with `TYPEWRITER_PACKAGE` and adding an package identifier as postfix at environment variable names. The identifier is mapped to the package by setting the package name at `TYPEWRITER_PACKAGE_LOG_PACKAGE_<identifier>`. If the package name should be interpreted as a full qualified name, the variable `TYPEWRITER_PACKAGE_FULL_QUALIFIED` is to set with value `true`.

The table shows the corresponding variable names for the example package *logger*

| general                                        | package specific                                              |
|------------------------------------------------|---------------------------------------------------------------|
| `TYPEWRITER_LOG_LEVEL`                         | `TYPEWRITER_PACKAGE_LOG_LEVEL_LOGGER`                         |
| `TYPEWRITER_LOG_APPENDER_TYPE`                 | `TYPEWRITER_PACKAGE_LOG_APPENDER_TYPE_LOGGER`                 |
| `TYPEWRITER_LOG_APPENDER_FILE`                 | `TYPEWRITER_PACKAGE_LOG_APPENDER_FILE_LOGGER`                 |
| `TYPEWRITER_LOG_FORMATTER_TYPE`                | `TYPEWRITER_PACKAGE_LOG_FORMATTER_TYPE_LOGGER`                |
| `TYPEWRITER_LOG_FORMATTER_PARAMETER`           | `TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_LOGGER`           |
| `TYPEWRITER_LOG_FORMATTER_PARAMETER_DELIMITER` | `TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_LOGGER_DELIMITER` |

If `TYPEWRITER_PACKAGE_LOG_PACKAGE_<identifier>` is not defined but an other package specific element. The `identifier` in lower case will be used as package name instead.

### Custom appender and formatter

It is possible to register or deregister custom appender and formatter together with their configuration provider.

* To register a custom appender use `logger.RegisterAppenderWithConfig`, where to pass
  * the name of the appender to access it via property `TYPEWRITER_LOG_APPENDER_TYPE` and its package variant (You are not able to override the build in ones)
  * a slice of key prefixes for properties that you want to map to your configuration implementation (It will be used to filter values of environment or property file before passing them to your configuration constructor)
  * your constructor function for your implementation of the `appender.Appender` interface at [appender/appender.go](appender/appender.go)
  * your constructor function for your implementation of the corresponding `config.AppenderConfig` interface  at [config/appender_configuration.go](config/appender_configuration.go)
* To deregister a custom appender use `logger.DeregisterAppenderTogetherWithConfig`, where to pass the name of your appender (You are not able to remove the build in ones)
* To register a custom formatter use `logger.RegisterFormatterWithConfig`, where to pass
  * the name of the formatter to access it via property `TYPEWRITER_LOG_FORMATTER_TYPE` and its package variant (You are not able to override the build in ones)
  * a slice of key prefixes for properties that you want to map to your configuration implementation (It will be used to filter values of environment or property file before passing them to your configuration constructor)
  * your constructor function for your implementation of the `format.Formatter` interface at [format/formatter.go](format/formatter.go)
  * your constructor function for your implementation of the corresponding `config.FormatterConfig` interface at [config/formatter_configuration.go](config/formatter_configuration.go)
* To deregister a custom formatter use `logger.DeregisterFormatterTogetherWithConfig`, where to pass the name of your formatter (You are not able to remove the build in ones)

## Examples

* Examples are available at godoc [logger/example_test.go](logger/example_test.go)
* Some multi threading file appender test to verify renaming of log files [logger/test/integration_test.go](logger/test/integration_test.go)
