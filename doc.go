/*
This repository provides a GoLang logging package, which are inspired by some Log4j features like package specific enablement

# 1. Logging

# 1.1 Access Log Methods

The logger can be accessed via the [pkg/github.com/ma-vin/typewriter/logger.Log]. This will return a struct implementing the [pkg/github.com/ma-vin/typewriter/logger.Logger] interface.
There are also the same methods directly in the “logger” package, but without returning an interface.

# 1.2 Log Methods

There are permutations of the following categories for each severity
  - log without additional parameters, e.g. “Information(args ...any)”
  - log with correlation, e.g. “InformationWithCorrelation(correlationId string, args ...any)”
  - log with custom values map, e.g. “InformationCustom(customValues map[string]any, args ...any)”
  - message format with args, e.g. “Informationf(format string, args ...any)”
  - panic (only for warning, error and fatal) e.g. “ErrorWithPanic(args ...any)”
  - exit (only for fatal) e.g. “FatalWithExit(args ...any)”

Arguments are handled in the manner of “fmt.Sprint”/“fmt.Sprintf”
In addition, there exists an indicator if the severity is enabled or not.

# 2. Configuration

Configuration is currently possible via environment variables or an properties file.

# 2.1 Configuration by file

To use a configuration file by setting an environment variable “TYPEWRITER_CONFIG_FILE” with a path to a properties file.
If this environment variable is set, all others will be ignored. Nevertheless, the same property names are used by both.
It is possible to use single line comments “#”,“//” or “--” and multiline comments with begin at “/*” and end at “* /” (Without blank, because of unknown escaping at doc.go).

# 2.2 Log levels

The log level can be set via “TYPEWRITER_LOG_LEVEL” with one of the following values:
  - DEBUG
  - INFO or INFORMATION
  - WARN or WARNING
  - ERROR (The default log level)
  - FATAL

# 2.3 Appender

There are two types of appender at the moment. They can be configured by setting “TYPEWRITER_LOG_APPENDER_TYPE” with one of the following values:

  - STDOUT: Standard output appender (The default appender)
  - FILE: File appender which writes to a file. The target file of the appender has to be configured by setting the path to “TYPEWRITER_LOG_APPENDER_FILE”.
    Renaming of the log file with existing entries at a specific point in time can be configured with a cron expression at “TYPEWRITER_LOG_APPENDER_CRON_RENAMING” (allowed characters are “* , - /”. The day of week is zero based)
    Renaming of the log file with existing entries, if a limit of bytes will be exceeded by next log entry, can be configured by an integer at “TYPEWRITER_LOG_APPENDER_SIZE_RENAMING”. This value will be ignored if “TYPEWRITER_LOG_APPENDER_CRON_RENAMING” is set also.

# 2.4 Formatter

There are three types of formatter which provides the texts to log for the appender. They can be configured by setting “TYPEWRITER_LOG_FORMATTER_TYPE” with one of the following values:

DELIMITER: The record information will be appended and delimited with a given sign. They can be set “TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>” where <x>has to be replaced by the following values:
 1. “DELIMITER” the delimiting signs. The default value of the delimiter is “ - ”.
 2. “TIME_LAYOUT” time layout. Default value of [time.RFC3339]

TEMPLATE: The records will be derived from three templates and time layout. They can be set “TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>” where <x>has to be replaced by the following values:
 1. “TEMPLATE” template for writing time, severity and the message. Default “[%s] %s: %s”
 2. “TEMPLATE_CORRELATION” template for writing time, severity, correlationID and the message. Default “[%s] %s %s: %s”
 3. “TEMPLATE_CUSTOM”template for writing time, severity, message and custom value map. Default “[%s] %s: %s”
 4. “TIME_LAYOUT” time layout. Default value of [time.RFC3339]
 5. “TEMPLATE_TRIM_SEVERITY” indicator wether to trim severity text or to add space at warn and info to algin following elements. Default “false”
 6. “TEMPLATE_CALLER” like 1. with caller function, file and line placed in front of message. Default “[%s] %s %s(%s.%d): %s”
 7. “TEMPLATE_CALLER_CORRELATION” like 2. with caller function, file and line placed in front of message. Default “[%s] %s %s %s(%s.%d): %s”
 8. “TEMPLATE_CALLER_CUSTOM” like 3. with caller function, file and line placed in front of message. Default “[%s] %s %s(%s.%d): %s”

It is possible to reorder parameter by argument indices. The custom value map at 3 will be appended as key-value-pairs sorted by key
(e.g. a custom map with three entries of string, number and boolean format at indices from 4 to 9: “severity: %[2]s message: %[3]s %[4]s: %[5]s %[6]s: %[7]d %[8]s: %[9]t time: %[1]s”)

JSON: The records will be logged as JSON. It is possible to define key names, the time layout and if the custom value map should be a sub element or not.
The keys of the custom value map will be used 1:1. The properties can be set via “TYPEWRITER_LOG_FORMATTER_PARAMETER_<x>” where <x> has to be replaced by the following values:
 1. “JSON_TIME_KEY” key of time. Default “time”
 2. “JSON_SEVERITY_KEY” key of severity. Default: “severity”
 3. “JSON_CORRELATION_KEY” key of correlationID. Default: “correlation”
 4. “JSON_MESSAGE_KEY” key of message. Default:  “message”
 5. “JSON_CUSTOM_VALUES_KEY” key of custom values if used as sub elements. Default: “custom”
 6. “JSON_CUSTOM_VALUES_SUB” indicator to add custom value map as sub element: Default: “false”
 7. “TIME_LAYOUT” time layout. Default value of [time.RFC3339]
 8. “JSON_CALLER_FUNCTION_KEY” key of the caller function. Default: “caller”
 9. “JSON_CALLER_FILE_KEY” key of the caller file. Default: “file”
 10. “JSON_CALLER_LINE_KEY” key of the caller file line. Default: “line”

The default formatter is the delimiter one.

# 2.5 Package specific configuration

Each configuration element can be declared package specific by replacing “TYPEWRITER” with “TYPEWRITER_PACKAGE” and adding the package identifier as postfix at environment variable names.
The identifier is mapped to the package by setting the package name at “TYPEWRITER_PACKAGE_LOG_PACKAGE_<identifier>”.
If the package name should be interpreted as a full qualified name, the variable vTYPEWRITER_PACKAGE_FULL_QUALIFIED” is to set with value “true”.
If “TYPEWRITER_PACKAGE_LOG_PACKAGE_<identifier>” is not defined but an other package specific element. The “identifier” in lower case will be used as package name instead.

  - “TYPEWRITER_LOG_LEVEL” -> “TYPEWRITER_PACKAGE_LOG_LEVEL_LOGGER”
  - “TYPEWRITER_LOG_APPENDER_TYPE” -> “TYPEWRITER_PACKAGE_LOG_APPENDER_TYPE_LOGGER”
  - “TYPEWRITER_LOG_APPENDER_FILE” -> “TYPEWRITER_PACKAGE_LOG_APPENDER_FILE_LOGGER”
  - “TYPEWRITER_LOG_FORMATTER_TYPE” -> “TYPEWRITER_PACKAGE_LOG_FORMATTER_TYPE_LOGGER”
  - “TYPEWRITER_LOG_FORMATTER_PARAMETER” -> “TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_LOGGER”
  - “TYPEWRITER_LOG_FORMATTER_PARAMETER_DELIMITER” -> “TYPEWRITER_PACKAGE_LOG_FORMATTER_PARAMETER_LOGGER_DELIMITER”
*/
package typewriter
