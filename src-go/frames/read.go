package frames

import (
	"bytes"
	"encoding/binary"
	"io"
)

// ReadInt returns an int64 from a message frame.
func ReadInt(r io.Reader, lenBytes int) (int64, error) {
	buf := make([]byte, lenBytes)
	lenReader := io.LimitReader(r, int64(lenBytes))

	_, err := lenReader.Read(buf)
	if err != nil {
		return 0, err
	}

	return int64(int64(binary.BigEndian.Uint16(buf))), nil
}

// ReadBytes returns a byte slice containing a message frame.
func ReadBytes(r io.Reader, frameLen int64) ([]byte, error) {
	reader := io.LimitReader(r, frameLen)
	frame := make([]byte, frameLen)

	readBytes := 0
	for readBytes < len(frame) {
		n, err := reader.Read(frame[readBytes:])
		if err != nil {
			return nil, err
		}
		readBytes += n
	}

	return frame, nil
}

// ReadProps returns a map containing all property name/value pairs.
func ReadProps(r io.Reader, propsLen int64) (map[string]string, error) {
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

	props := make(map[string]string)
	var pair []*bytes.Buffer

	beginPair := func() {
		pair = make([]*bytes.Buffer, 1, 2)
		pair[0] = &bytes.Buffer{}
	}

	beginPair()
	for n := 0; n < len(buf); n++ {
		if buf[n] == '\n' {
			props[pair[0].String()] = pair[1].String()
			if n == len(buf) {
				break
			}
			beginPair()
		} else if buf[n] == ':' {
			pair = append(pair, &bytes.Buffer{})
			// assume trailing space; jump over it:
			n++
		} else if len(pair) == 1 {
			pair[0].WriteByte(buf[n])
		} else if len(pair) == 2 {
			pair[1].WriteByte(buf[n])
		}
	}

	if pair[0].Len() > 0 {
		if len(pair) == 2 {
			props[pair[0].String()] = pair[1].String()
		} else {
			props[pair[0].String()] = ""
		}
	}

	return props, nil
}
