package check

import (
	"log"
)

// IsFalse panics if v isn't false.
func IsFalse(v bool, name string) {
	if v {
		log.Fatalf("%s is true. expected value false\n", name)
	}
}

// IsGreaterEqual panics if v1 isn't greater than or equal to v2.
func IsGreaterEqual(v1 int, v2 int, hint string) {
	if v1 < v2 {
		log.Fatalf("%s equal to %d. expected greater than or equal to %d\n", hint, v1, v2)
	}
}

// IsInRange panics if v isn't between min and max inclusive.
func IsInRange(v int, min int, max int, hint string) {
	if v < min || v > max {
		log.Fatalf("%s equal to %d. expected between minimum %d and maximum %d\n", hint, v, min, max)
	}
}

// IsNotNil panics if v is nil.
func IsNotNil(v interface{}, hint string) {
	if v == nil {
		log.Fatalf("%s is nil. expected not nil\n", hint)
	}
}
