package logger

import (
	"os"
	"strings"

	"github.com/ma-vin/typewriter/appender"
)

func CreatePackageLoggers(appender *appender.Appender) map[string]*CommonLogger {
	result := map[string]*CommonLogger{}

	for _, envEntry := range os.Environ() {

		keyValue := strings.SplitN(strings.ToUpper(envEntry), "=", 2)

		if len(keyValue) == 2 {

			packageName, found := strings.CutPrefix(keyValue[0], DEFAULT_LOG_LEVEL_ENV_NAME+"_")

			if found {
				logger := CommonLogger{appender: appender}
				determineSeverity(strings.TrimSpace(keyValue[1]), &logger)
				result[strings.ToLower(packageName)] = &logger
			}
		}
	}

	return result
}
