package frames

import (
	"bytes"
	"io"
)

var (
	// KeyValueSepChar is a colon (:) character which is used in a message property key/value separator.
	KeyValueSepChar = byte(':')
	// SpaceChar is a space character.
	SpaceChar = byte(' ')
	// PropTermChar is a zero value byte which is used in a message property terminator.
	PropTermChar = byte(0)
)

var (
	// KeyValueSep is the complete property key/value separator as a byte slice.
	KeyValueSep = []byte{KeyValueSepChar, SpaceChar}
	// PropTerm is the complete property terminator as a byte slice.
	PropTerm = []byte{PropTermChar, PropTermChar, PropTermChar}
)

// PropsToBytes turns a property map into a byte slice.
func PropsToBytes(props map[string][]byte) []byte {
	var buf bytes.Buffer
	for key, value := range props {
		buf.WriteString(key)
		buf.Write(KeyValueSep)
		buf.Write(value)
		buf.Write(PropTerm)
	}
	return buf.Bytes()
}

// ReadProps returns a map containing all property name/value pairs.
//
// This function assumes that the buffer doesn't start with 'delimChar's,
// if it does then you'll get a panic (index out of range).
func ReadProps(r io.Reader, propsLen int64) (map[string][]byte, error) {
	reader := io.LimitReader(r, propsLen)
	buf := make([]byte, propsLen)

	readBytes := 0
	for readBytes < len(buf) {
		n, err := reader.Read(buf[readBytes:])
		if err != nil {
			return nil, err
		}
		readBytes += n
	}

	props := make(map[string][]byte)
	var pair []*bytes.Buffer

	beginPair := func() {
		pair = make([]*bytes.Buffer, 1, 2)
		pair[0] = &bytes.Buffer{}
	}

	beginPair()
	for n := 0; n < len(buf); n++ {
		if isChar(buf, n, PropTerm[0]) &&
			isChar(buf, n, PropTerm[1]) &&
			isChar(buf, n, PropTerm[2]) {

			props[pair[0].String()] = pair[1].Bytes()

			// walk past any trailing delimiter chars:
			n += 2
			for isChar(buf, n+1, PropTermChar) {
				n++
			}
			beginPair()
		} else if isChar(buf, n, KeyValueSep[0]) &&
			isChar(buf, n+1, KeyValueSep[1]) {

			pair = append(pair, &bytes.Buffer{})

			// walk past any trailing spaces:
			n++
			for isChar(buf, n+1, SpaceChar) {
				n++
			}
		} else {
			pair[len(pair)-1].WriteByte(buf[n])
		}
	}

	if pair[0].Len() > 0 {
		if len(pair) == 2 {
			props[pair[0].String()] = pair[1].Bytes()
		} else if pair[0].Len() > 0 {
			props[pair[0].String()] = make([]byte, 0)
		}
	}

	return props, nil
}

func isChar(buf []byte, n int, char byte) bool {
	if n >= len(buf) {
		return false
	}
	return buf[n] == char
}
