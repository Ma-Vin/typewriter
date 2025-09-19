// This package provides appender to delegate messages to their target
package appender

import "github.com/ma-vin/typewriter/common"

// Interface to delegate messages to their target
type Appender interface {
	// Writes the given logValues. If the appender is closed, it should not write any write entries to its target anymore
	Write(logValues *common.LogValues)
	// Closes writer if necessary
	Close()
}
