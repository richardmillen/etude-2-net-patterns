// Package uuid may be used to generate random UUID's according to RFC4122.
//
// see sections 4.1.1 & 4.1.3 for lines 74 & 75 respectively.
//
// slightly modified version of go playground version (link below).
// https://play.golang.org/p/4FkNSiUDMg
//
// error handling changed slightly from the playground version.
// this implementation writes to log (and panics) rather than
// returning the error because:
// a) the returned error pollutes the interface
// b) if an error occurs then there's something fundamentally wrong
package uuid

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// Size is the byte length of a raw UUID.
const Size = 16

// Bytes represents the raw (16) bytes of a UUID.
type Bytes []byte

// New generates a random UUID.
func New() Bytes {
	uuid := make([]byte, Size)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(err)
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return uuid
}

// NewFrom constructs a new UUID from a byte slice.
func NewFrom(src []byte) Bytes {
	check.IsEqual(len(src), Size, "uuid byte length")

	b := make([]byte, Size)
	copy(b, src)
	return Bytes(b)
}

// Equal returns true if two UUIDs are identical.
func Equal(a, b Bytes) bool {
	return bytes.Equal([]byte(a), []byte(b))
}

// String returns the UUID bytes as a string.
func (b Bytes) String() string {
	b2 := []byte(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b2[0:4], b2[4:6], b2[6:8], b2[8:10], b2[10:])
}
