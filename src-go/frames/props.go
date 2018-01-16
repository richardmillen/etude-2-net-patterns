package frames

import "bytes"

// PropsToBytes turns a property map into a byte slice.
func PropsToBytes(props map[string]string) []byte {
	var buf bytes.Buffer
	for key, value := range props {
		buf.WriteString(key + ": " + value + "\n")
	}
	return buf.Bytes()
}
