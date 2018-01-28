package check

import (
	"fmt"
	"log"
)

// Must logs an error and panics.
func Must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Error logs an error and panics.
func Error(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Log returns true if a valid error was logged.
func Log(err error) bool {
	if err != nil {
		log.Printf("error logged: %s\n", err)
		return true
	}
	return false
}

// NotEqual returns an error if a equals b.
func NotEqual(a int, b int, hint string) error {
	if a == b {
		return fmt.Errorf("%s equal to %d. expected not equal to %d", hint, a, b)
	}
	return nil
}

// NotNil returns an error if v is nil.
func NotNil(v interface{}, hint string) error {
	if v == nil {
		return fmt.Errorf("%s is nil. expected not nil", hint)
	}
	return nil
}
