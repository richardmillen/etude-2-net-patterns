package frames

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// ReadUInt8 reads a single octet/byte and returns the uint8 value.
func ReadUInt8(r io.Reader) (uint8, error) {
	n, err := readUInt(r, 1)
	return uint8(n), err
}

// ReadUInt16 reads two octets/bytes and returns the uint16 value.
func ReadUInt16(r io.Reader) (uint16, error) {
	n, err := readUInt(r, 2)
	return uint16(n), err
}

// ReadSig returns a protocol signature as a byte array.
func ReadSig(r io.Reader) ([2]byte, error) {
	b, err := ReadBytes(r, 2)
	if err != nil {
		return [2]byte{}, err
	}
	return [...]byte{b[0], b[1]}, nil
}

// ReadBytes returns a byte slice.
func ReadBytes(r io.Reader, numBytes int64) ([]byte, error) {
	reader := io.LimitReader(r, numBytes)
	frame := make([]byte, numBytes)

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

// readInt read 'n' bytes from a message frame and returns the value as a uint64.
func readUInt(r io.Reader, numBytes uint8) (uint64, error) {
	buf := make([]byte, numBytes)
	lenReader := io.LimitReader(r, int64(numBytes))

	_, err := lenReader.Read(buf)
	if err != nil {
		return 0, err
	}

	switch numBytes {
	case 1:
		return uint64(buf[0]), nil
	case 2:
		return uint64(binary.BigEndian.Uint16(buf)), nil
	case 3, 4:
		return uint64(binary.BigEndian.Uint32(buf)), nil
	case 5, 6, 7, 8:
		return binary.BigEndian.Uint64(buf), nil
	default:
		return 0, fmt.Errorf("readUInt cannot read %d bytes", numBytes)
	}
}
