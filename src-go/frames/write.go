package frames

import (
	"encoding/binary"
	"fmt"
)

// WriteBytes writes a slice of bytes to another slice
// and returns the new position in the destination slice.
func WriteBytes(dst []byte, src []byte) []byte {
	copy(dst, src)
	return dst[len(src):]
}

// WriteInt64 writes an int64 to a slice and
// returns the new position in the destination slice.
func WriteInt64(dst []byte, n int64, octets int) []byte {
	switch octets {
	case 1:
		dst[0] = byte(n)
	case 2:
		binary.BigEndian.PutUint16(dst, uint16(n))
	default:
		panic(fmt.Sprintf("invalid octet count of %d passed to WriteInt64", octets))
	}
	return dst[octets:]
}
