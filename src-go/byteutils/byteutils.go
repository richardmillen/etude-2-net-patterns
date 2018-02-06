package byteutils

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

// JoinBytes combines multiple byte slices into a single slice.
func JoinBytes(v ...[]byte) []byte {
	if v == nil {
		return nil
	}
	if len(v) == 0 {
		return []byte{}
	}

	joined := make([]byte, 0, len(v[0])*len(v))

	for _, slice := range v {
		joined = append(joined, slice...)
	}

	return joined
}
