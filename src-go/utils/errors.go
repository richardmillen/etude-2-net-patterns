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
