package check

import (
	"log"
)

// IsFalse panics if v isn't false.
func IsFalse(v bool, name string) {
	if v {
		log.Fatalf("%s is true. expected value false", name)
	}
}

// IsGreaterEqual panics if v1 isn't greater than or equal to v2.
func IsGreaterEqual(v1 int, v2 int, message string) {
	if v1 < v2 {
		log.Fatalf("%s equal to %d. expected greater than or equal to %d", message, v1, v2)
	}
}

// IsInRange panics if v isn't between min and max inclusive.
func IsInRange(v int, min int, max int, message string) {
	if v < min || v > max {
		log.Fatalf("%s equal to %d. expected between minimum %d and maximum %d", message, v, min, max)
	}
}
