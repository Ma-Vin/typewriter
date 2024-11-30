// This package provides appender to delegate messages to their target
package appender

import "github.com/ma-vin/typewriter/common"

// Interface to delegate messages to their target
type Appender interface {
	// Writes the given logValues
	Write(logValues *common.LogValues)
	// Closes writer if necessary
	Close()
}
