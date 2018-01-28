package check

import (
	"log"
)

// MustFalse logs and panics if v isn't false.
func MustFalse(v bool, name string) {
	if v {
		log.Panicf("%s is true. expected value false\n", name)
	}
}

// MustNotEqual logs and panics if a is equal to b.
func MustNotEqual(a int, b int, hint string) {
	if a == b {
		log.Panicf("%s equal to %d. expected not equal to %d\n", hint, a, b)
	}
}

// MustEqual logs and panics if a isn't equal to b.
func MustEqual(a int, b int, hint string) {
	if a != b {
		log.Panicf("%s equal to %d. expected equal to %d\n", hint, a, b)
	}
}

// MustGreater logs and panics if a isn't greater than b.
func MustGreater(a int, b int, hint string) {
	if a <= b {
		log.Panicf("%s equal to %d. expected greater than %d\n", hint, a, b)
	}
}

// MustGreaterEqual logs and panics if a isn't greater than or equal to b.
func MustGreaterEqual(a int, b int, hint string) {
	if a < b {
		log.Panicf("%s equal to %d. expected greater than or equal to %d\n", hint, a, b)
	}
}

// MustInRange logs and panics if v isn't between min and max inclusive.
func MustInRange(v int, min int, max int, hint string) {
	if v < min || v > max {
		log.Panicf("%s equal to %d. expected between minimum %d and maximum %d\n", hint, v, min, max)
	}
}

// MustNotNil logs and panics if v is nil.
func MustNotNil(v interface{}, hint string) {
	if v == nil {
		log.Panicf("%s is nil. expected not nil\n", hint)
	}
}
