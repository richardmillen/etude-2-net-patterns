package frames

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

const (
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
	PropTerm = []byte{PropTermChar}
)

var (
	// ErrNoPropKey occurs when an empty property key is found.
	ErrNoPropKey = errors.New("empty property key")
	// ErrInvPropKey occurs when a property key contains invalid characters.
	ErrInvPropKey = errors.New("invalid property key")
)

// PropsToBytes validates and turns a property map into a byte slice.
func PropsToBytes(props map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer

	for key, value := range props {
		err := checkKey(key)
		if err != nil {
			return nil, err
		}

		buf.WriteString(key)
		buf.Write(KeyValueSep)
		buf.Write(value)
		buf.Write(PropTerm)
	}

	return buf.Bytes(), nil
}

// ReadProps returns a map containing all property name/value pairs.
//
// This function assumes that the buffer doesn't start with 'delimChar's,
// if it does then you'll get a panic (index out of range).
func ReadProps(r io.Reader, propsLen int64) (map[string][]byte, error) {
	if propsLen == 0 {
		return nil, nil
	}

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

	var pair []*bytes.Buffer
	var prev []*bytes.Buffer

	beginPair := func() {
		prev = pair
		pair = make([]*bytes.Buffer, 1, 2)
		pair[0] = &bytes.Buffer{}
	}

	props := make(map[string][]byte)

	beginPair()
	for n := 0; n < len(buf); n++ {
		if len(pair) == 1 && utils.IsAt(buf, KeyValueSep, n) {
			err := checkKey(pair[0].String())
			if err != nil {
				return nil, err
			}

			pair = append(pair, &bytes.Buffer{})
			n += (len(KeyValueSep) - 1)
		} else if utils.IsAt(buf, PropTerm, n) {
			if len(pair) == 1 && prev != nil {
				// we still have more data for the previous property value:
				prev[1].Write(PropTerm)
				prev[1].Write(pair[0].Bytes())
				props[prev[0].String()] = prev[1].Bytes()
				pair[0].Reset()
			} else {
				props[pair[0].String()] = pair[1].Bytes()
				beginPair()
			}
			n += (len(PropTerm) - 1)
		} else {
			pair[len(pair)-1].WriteByte(buf[n])
		}
	}

	if pair[0].Len() > 0 {
		if len(pair) == 2 {
			props[pair[0].String()] = pair[1].Bytes()
		} else if pair[0].Len() > 0 {
			props[pair[0].String()] = []byte{}
		}
	}

	return props, nil
}

func checkKey(key string) (err error) {
	if key == "" {
		return ErrNoPropKey
	}
	if strings.Contains(key, string(PropTerm)) {
		return ErrInvPropKey
	}
	return
}
