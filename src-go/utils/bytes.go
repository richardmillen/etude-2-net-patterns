package utils

import "bytes"

// IsAt returns true if b exists in a at a specified position.
func IsAt(a, b []byte, at int) bool {
	if at < 0 || at >= (len(a)-len(b))+1 {
		return false
	}
	return bytes.Equal(a[at:at+len(b)], b)
}
