package frames

import (
	"bytes"
	"io"
	"log"
)

var (
	kvSep     = byte(':')
	spaceChar = byte(' ')
	propSep   = byte(0)

	kvSepBytes   = []byte{kvSep, spaceChar}
	propSepBytes = []byte{propSep}
)

// PropsToBytes turns a property map into a byte slice.
func PropsToBytes(props map[string][]byte) []byte {
	var buf bytes.Buffer
	for key, value := range props {
		buf.WriteString(key)
		buf.Write(kvSepBytes)
		buf.Write(value)
		buf.Write(propSepBytes)
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

	log.Println("reading props len:", propsLen)

	readBytes := 0
	for readBytes < len(buf) {
		n, err := reader.Read(buf[readBytes:])
		if err != nil {
			return nil, err
		}
		readBytes += n
	}

	log.Println("read props into buf.")

	props := make(map[string][]byte)
	var pair []*bytes.Buffer

	beginPair := func() {
		pair = make([]*bytes.Buffer, 1, 2)
		pair[0] = &bytes.Buffer{}
	}

	log.Println("building props map...")

	beginPair()
	for n := 0; n < len(buf); n++ {
		if isChar(buf, n, propSep) {
			props[pair[0].String()] = pair[1].Bytes()

			// walk past any trailing delimiter chars:
			for isChar(buf, n+1, propSep) {
				n++
			}
			beginPair()
		} else if isChar(buf, n, kvSep) && isChar(buf, n+1, spaceChar) {
			pair = append(pair, &bytes.Buffer{})

			// walk past any trailing spaces:
			n += 2
			for isChar(buf, n, spaceChar) {
				n++
			}
		} else {
			pair[len(pair)-1].WriteByte(buf[n])
		}
	}

	log.Println("adding final prop pair if exists...")

	if pair[0].Len() > 0 {
		if len(pair) == 2 {
			props[pair[0].String()] = pair[1].Bytes()
		} else if pair[0].Len() > 0 {
			props[pair[0].String()] = make([]byte, 0)
		}
	}

	log.Println("props map built.")
	log.Println("uuid len:", len(props["uuid"]))

	return props, nil
}

func isChar(buf []byte, n int, char byte) bool {
	if n >= len(buf) {
		return false
	}
	return buf[n] == char
}
