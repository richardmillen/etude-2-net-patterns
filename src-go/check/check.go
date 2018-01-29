package check

import (
	"fmt"
	"log"
	"runtime/debug"
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
		debug.PrintStack()
		return true
	}
	return false
}

// NewFailedError constructs a new check FailedError.
func NewFailedError(format string, a ...interface{}) *FailedError {
	return &FailedError{
		err: fmt.Sprintf(format, a...),
	}
}

// FailedError is the error returned by all check functions that return an error.
type FailedError struct {
	err string
}

func (e FailedError) Error() string {
	return e.err
}

// NotEqual returns an error if a equals b.
func NotEqual(a int, b int, hint string) error {
	if a == b {
		return NewFailedError("%s equal to %d. expected not equal to %d", hint, a, b)
	}
	return nil
}

// NotNil returns an error if v is nil.
func NotNil(v interface{}, hint string) error {
	if v == nil {
		return NewFailedError("%s is nil. expected not nil", hint)
	}
	return nil
}
