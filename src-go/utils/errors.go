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

// LogIfError writes to log and returns true if an error is provided.
func LogIfError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
