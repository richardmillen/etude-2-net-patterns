package utils

import (
	"log"
)

// CheckError logs an error and panics.
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// LogError returns true if an error was logged.
func LogError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
