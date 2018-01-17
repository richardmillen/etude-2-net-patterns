package frames

import (
	"encoding/binary"
)

// WriteBytes writes a slice of bytes to another slice
// and returns the new position in the destination slice.
func WriteBytes(dst []byte, src []byte) []byte {
	copy(dst, src)
	return dst[len(src):]
}

// WriteUInt8 writes an uint8 to a slice and
// returns the new position in the destination slice.
func WriteUInt8(dst []byte, n uint8) []byte {
	dst[0] = byte(n)
	return dst[1:]
}

// WriteUInt16 writes an uint16 to a slice and
// returns the new position in the destination slice.
func WriteUInt16(dst []byte, n uint16) []byte {
	binary.BigEndian.PutUint16(dst, n)
	return dst[2:]
}
