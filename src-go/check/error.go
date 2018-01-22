package check

import "log"

// Must panics when passed a valid error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Error logs an error and panics.
func Error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Log returns true if a valid error was logged.
func Log(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
