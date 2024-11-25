[![Build and test](https://github.com/Ma-Vin/typewriter/actions/workflows/go-build.yaml/badge.svg)](https://github.com/Ma-Vin/typewriter/actions/workflows/go-build.yaml)

# Typewriter

This repository provides a GoLang logging package, which are inspired by some Log4j features like package specific enablement

:warning: :construction: **Not released yet**

## Usage

### Logging

#### Access

The logger can be accessed via importing `github.com/ma-vin/typewriter/logger` and using the `Log()` method at `logger` package. This will return a struct implementing the `Logger` interface.

#### Log Methods

There are permutations of the following categories for each severity

* log without additional parameters, e.g. `Information(args ...any)`
* log with correlation, e.g. `InformationWithCorrelation(correlationId string, args ...any)`
* log with custom values map, e.g. `InformationCustom(customValues map[string]any, args ...any)`
* message format with args, e.g. `Informationf(format string, args ...any)`
* panic (only for warning, error and fatal) e.g. `ErrorWithPanic(args ...any)`
* exit (only for fatal) e.g. `FatalWithExit(args ...any)`

Arguments are handled in the manner of `fmt.Sprint`/`fmt.Sprintf`

In addition, there exists an indicator if the severity is enabled or not.

### Configuration

Configuration is currently possible via environment variables or an properties file.

#### Configauration by file

To use a configuration file by setting an environment variable `TYPEWRITER_CONFIG_FILE` with a path to a properties file.
If this environment variable is set, all others will be ignored. Nevertheless, the same property names are used by both.
It is possible to use single line comments `#`,`//` or `--` and multiline commets with begin at `/*` and end at `*/`.

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
  * The target file of the appender has to be configured by setting the path to `TYPEWRITER_LOG_APPENDER_PARAMETER`  

The default appender is the standard output one

#### Formatter

There are three types of formatter which provides the texts to log for the appender. They can be configured by setting `TYPEWRITER_LOG_FORMATTER_TYPE` with one of the following values:

* `DELIMITER`: The record information will be appended and delimited with a given sign
  * The default value of the delimiter is ` - `. It can be configured by setting `TYPEWRITER_LOG_FORMATTER_PARAMETER`.
* `TEMPLATE`: The records will be derived from three templates and time layout. They can be set `TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>` where `<x>`has to be replaced by the following integers.
  1. template for writing time, severity and the message. Default `[%s] %s: %s`
  2. template for writing time, severity, correlationID and the message. Default `[%s] %s %s: %s`
  3. template for writing time, severity, message and custom value map. Default `[%s] %s: %s`
  4. time layout. Default value of `time.RFC3339`
  5. indicator wether to trim severity text or to add space at warn and info to algin following elements. Default `false`

  It is possible to reorder parameter by argument indices. The *custom value map* at *3* will be appended as key-value-pairs sorted by key (e.g. a custom map with three entries of string, number and boolean format at indices from 4 to 9: `severity: %[2]s message: %[3]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %[9]t time: %[1]s`)
* `JSON`: The records will be logged as *JSON*. It is possible to define key names, the time layout and if the custom value map should be a sub element or not. The keys of the custom value map will be used 1:1. The properties can be set via `TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>` where `<x>` has to be replaced by the following integers.
  1. key of time. Default `time`
  2. key of severity. Default: `severity`
  3. key of correlationID. Default: `correlation`
  4. key of message. Default:  `message`
  5. key of custom values if used as sub elements. Default: `custom`
  6. indicator to add custom value map as sub element: Default: `false`
  7. time layout. Default value of `time.RFC3339`

The default formatter is the delimiter one.

#### Package specific configuration

Each configuration element can be declared package specific by replacing `TYPEWRITER` with `TYPEWRITER_PACKAGE` and adding the package as postfix at environment variable names.

The table shows the corresponding variable names for the example package *logger*

| general                                | package specific                                      |
|----------------------------------------|-------------------------------------------------------|
| `TYPEWRITER_LOG_LEVEL`                 | `TYPEWRITER_PACKAGE_LOG_LEVEL_LOGGER`                 |
| `TYPEWRITER_LOG_APPENDER_TYPE`         | `TYPEWRITER_PACKAGE_LOG_APPENDER_TYPE_LOGGER`         |
| `TYPEWRITER_LOG_APPENDER_PARAMETER`    | `TYPEWRITER_PACKAGE_LOG_APPENDER_PARAMETER_LOGGER`    |
| `TYPEWRITER_LOG_FORMATTER_TYPE`        | `TYPEWRITER_PACKAGE_LOG_FORMATTER_TYPE_LOGGER`        |
| `TYPEWRITER_LOG_FORMATTER_PARAMETER`   | `TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_LOGGER`   |
| `TYPEWRITER_LOG_FORMATTER_PARAMETER_1` | `TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_LOGGER_1` |

## Examples

Examples are avialable at godoc [logger/example_test.go](logger/example_test.go)
