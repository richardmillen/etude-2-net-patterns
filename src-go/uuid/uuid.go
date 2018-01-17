package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

// New generates a random UUID according to RFC4122
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
func New() string {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		log.Fatal(err)
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
