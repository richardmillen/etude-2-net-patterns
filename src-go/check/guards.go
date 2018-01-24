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

// IsNotEqual panics if a is equal to b.
func IsNotEqual(a int, b int, hint string) {
	if a == b {
		log.Fatalf("%s equal to %d. expected not equal to %d\n", hint, a, b)
	}
}

// IsEqual panics if a isn't equal to b.
func IsEqual(a int, b int, hint string) {
	if a != b {
		log.Fatalf("%s equal to %d. expected equal to %d\n", hint, a, b)
	}
}

// IsGreater panics if a isn't greater than b.
func IsGreater(a int, b int, hint string) {
	if a <= b {
		log.Fatalf("%s equal to %d. expected greater than %d\n", hint, a, b)
	}
}

// IsGreaterEqual panics if a isn't greater than or equal to b.
func IsGreaterEqual(a int, b int, hint string) {
	if a < b {
		log.Fatalf("%s equal to %d. expected greater than or equal to %d\n", hint, a, b)
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
