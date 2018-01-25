package utils_test

import (
	"testing"

	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

var testCases = []struct {
	desc     string
	a        []byte
	b        []byte
	index    int
	expected bool
}{
	{
		desc:     "both empty",
		a:        []byte{},
		b:        []byte{},
		expected: false,
		index:    0,
	},
	{
		desc:     "both empty, index out of range",
		a:        []byte{},
		b:        []byte{},
		expected: false,
		index:    999,
	},
	{
		desc:     "both empty, negative index",
		a:        []byte{},
		b:        []byte{},
		expected: false,
		index:    -1,
	},
	{
		desc:     "b is empty",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{},
		expected: false,
		index:    0,
	},
	{
		desc:     "b not in a",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{8, 9},
		expected: false,
		index:    0,
	},
	{
		desc:     "b not in a, index out of range",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{8, 9},
		expected: false,
		index:    999,
	},
	{
		desc:     "a starts with singular b, negative index",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1},
		expected: false,
		index:    -1,
	},
	{
		desc:     "a starts with singular b, index out of range",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1},
		expected: false,
		index:    999,
	},
	{
		desc:     "a starts with singular b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1},
		expected: false,
		index:    0,
	},
	{
		desc:     "a contains singular b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{2},
		expected: false,
		index:    1,
	},
	{
		desc:     "a ends with singular b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{4},
		expected: false,
		index:    3,
	},
	{
		desc:     "a starts with b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 2},
		expected: false,
		index:    0,
	},
	{
		desc:     "a contains b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{2, 3},
		expected: false,
		index:    1,
	},
	{
		desc:     "a ends with b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{3, 4},
		expected: false,
		index:    2,
	},
	{
		desc:     "a starts with partial b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 9},
		expected: false,
		index:    3,
	},
	{
		desc:     "a ends with partial b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{4, 5},
		expected: false,
		index:    3,
	},
	{
		desc:     "a and b match",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 2, 3, 4},
		expected: false,
		index:    0,
	},
	{
		desc:     "a is a subset of b",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 2, 3, 4, 5, 6, 7},
		expected: false,
		index:    0,
	},
}

func TestIsAt(t *testing.T) {
	for _, tc := range testCases {
		actual := utils.IsAt(tc.index, tc.a, tc.b...)
		if actual != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.desc, tc.expected, actual)
		}
	}
}
