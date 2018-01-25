package utils

import "bytes"

// IsAt returns true if 'b' is a subset of 'a' at the specified position.
func IsAt(a, b []byte, at int) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	if at < 0 || at >= (len(a)-len(b))+1 {
		return false
	}
	return bytes.Equal(a[at:at+len(b)], b)
}
