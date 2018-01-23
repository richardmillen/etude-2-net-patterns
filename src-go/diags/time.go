package diags

import (
	"log"
	"time"
)

// StopFunc is called after Start() in order to stop a timed operation.
type StopFunc func()

// Start is called to time an operation and get the
func Start(message string) StopFunc {
	log.Println("started:", message)
	started := time.Now()
	return func() {
		log.Printf("elapsed: %v.\n", time.Since(started))
	}
}
