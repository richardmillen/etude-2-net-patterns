package frames

import (
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

// readUInt read 'n' bytes from a message frame and returns the value as a uint64.
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
