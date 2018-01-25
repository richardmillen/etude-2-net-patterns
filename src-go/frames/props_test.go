package frames_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
)

var testCases = []struct {
	name     string
	props    map[string][]byte
	bytesStr string
	err      error
}{
	{
		name:     "Empty",
		err:      nil,
		bytesStr: "",
		props:    map[string][]byte{},
	},
	{
		name:     "SingleEmptyProperty",
		err:      nil,
		bytesStr: "",
		props: map[string][]byte{
			"": []byte{},
		},
	},
	{
		name:     "SinglePropertyEmptyValue",
		err:      nil,
		bytesStr: fmt.Sprintf("abc%s%s", frames.KeyValueSep, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{},
		},
	},
	{
		name:     "SinglePropertyWithValue",
		err:      nil,
		bytesStr: fmt.Sprintf("abc%s%s%s", frames.KeyValueSep, []byte{1, 2, 3}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, 2, 3},
		},
	},
	{
		name: "MultiplePropertiesEmptyValues",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%sdef%s%s",
			frames.KeyValueSep, frames.PropTerm,
			frames.KeyValueSep, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{},
			"def": []byte{},
		},
	},
	{
		name: "MultiplePropertiesFirstWithEmptyValue",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%sdef%s%s%s",
			frames.KeyValueSep, frames.PropTerm,
			frames.KeyValueSep, []byte{4, 5, 6}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{},
			"def": []byte{4, 5, 6},
		},
	},
	{
		name: "MultiplePropertiesLastWithEmptyValue",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%sdef%s%s",
			frames.KeyValueSep, []byte{1, 2, 3}, frames.PropTerm,
			frames.KeyValueSep, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, 2, 3},
			"def": []byte{},
		},
	},
	{
		name: "MultiplePropertiesWithValues",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%sdef%s%s%s",
			frames.KeyValueSep, []byte{1, 2, 3}, frames.PropTerm,
			frames.KeyValueSep, []byte{4, 5, 6}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, 2, 3},
			"def": []byte{4, 5, 6},
		},
	},
	{
		name: "SinglePropertyValueContainsSeparator",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%s",
			frames.KeyValueSep, []byte{1, frames.KeyValueSep[0], frames.KeyValueSep[1], 2, 3}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, frames.KeyValueSep[0], frames.KeyValueSep[1], 2, 3},
		},
	},
	{
		name: "SinglePropertyValueStartsWithSeparator",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%s",
			frames.KeyValueSep, []byte{frames.KeyValueSep[0], frames.KeyValueSep[1], 1, 2, 3}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{frames.KeyValueSep[0], frames.KeyValueSep[1], 1, 2, 3},
		},
	},
	{
		name: "SinglePropertyValueEndsWithSeparator",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%s",
			frames.KeyValueSep, []byte{1, 2, 3, frames.KeyValueSep[0], frames.KeyValueSep[1]}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, 2, 3, frames.KeyValueSep[0], frames.KeyValueSep[1]},
		},
	},
	{
		name: "SinglePropertyValueContainsTerminator",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%s",
			frames.KeyValueSep, []byte{1, frames.PropTerm[0], frames.PropTerm[1], frames.PropTerm[2], 2, 3}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, frames.PropTerm[0], frames.PropTerm[1], frames.PropTerm[2], 2, 3},
		},
	},
	{
		name: "SinglePropertyValueStartsWithTerminator",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%s",
			frames.KeyValueSep, []byte{frames.PropTerm[0], frames.PropTerm[1], frames.PropTerm[2], 1, 2, 3}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{frames.PropTerm[0], frames.PropTerm[1], frames.PropTerm[2], 1, 2, 3},
		},
	},
	{
		name: "SinglePropertyValueEndsWithTerminator",
		err:  nil,
		bytesStr: fmt.Sprintf("abc%s%s%s",
			frames.KeyValueSep, []byte{1, 2, 3, frames.PropTerm[0], frames.PropTerm[1], frames.PropTerm[2]}, frames.PropTerm),
		props: map[string][]byte{
			"abc": []byte{1, 2, 3, frames.PropTerm[0], frames.PropTerm[1], frames.PropTerm[2]},
		},
	},
}

func TestProps(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("PropsToBytes_%s", tc.name), func(t *testing.T) {
			expected := []byte(tc.bytesStr)
			actual := frames.PropsToBytes(tc.props)

			if !bytes.Equal(actual, expected) {
				t.Errorf("expected %b, got %b", expected, actual)
			}
		})
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ReadProps_%s", tc.name), func(t *testing.T) {
			r := bytes.NewReader([]byte(tc.bytesStr))
			actual, err := frames.ReadProps(r, int64(r.Len()))

			if !reflect.DeepEqual(actual, tc.props) {
				t.Errorf("expected %v,\ngot      %v", tc.props, actual)
			}
			if err != tc.err {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}
		})
	}
}