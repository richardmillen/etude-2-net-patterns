package frames

import (
	"bytes"
)

var endOfLine = []byte{'\n'}

// PropsToBytes turns a property map into a byte slice.
func PropsToBytes(props map[string][]byte) []byte {
	var buf bytes.Buffer
	for key, value := range props {
		buf.WriteString(key + ": ")
		buf.Write(value)
		buf.Write(endOfLine)
	}
	return buf.Bytes()
}
