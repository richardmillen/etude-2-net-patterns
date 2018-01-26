package utils_test

import (
	"bytes"
	"testing"

	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

var isAtTestCases = []struct {
	name     string
	a        []byte
	b        []byte
	index    int
	expected bool
}{
	{
		name:     "BothEmpty",
		a:        []byte{},
		b:        []byte{},
		expected: false,
		index:    0,
	},
	{
		name:     "BothEmptyWithIndexOutOfRange",
		a:        []byte{},
		b:        []byte{},
		expected: false,
		index:    999,
	},
	{
		name:     "BothEmptyWithNegativeIndex",
		a:        []byte{},
		b:        []byte{},
		expected: false,
		index:    -1,
	},
	{
		name:     "BIsEmpty",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{},
		expected: false,
		index:    0,
	},
	{
		name:     "BNotInA",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{8, 9},
		expected: false,
		index:    0,
	},
	{
		name:     "BNotInAWithIndexOutOfRange",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{8, 9},
		expected: false,
		index:    999,
	},
	{
		name:     "AStartsWithSingularBWithNegativeIndex",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1},
		expected: false,
		index:    -1,
	},
	{
		name:     "AStartsWithSingularBWithIndexOutOfRange",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1},
		expected: false,
		index:    999,
	},
	{
		name:     "AStartsWithSingularB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1},
		expected: true,
		index:    0,
	},
	{
		name:     "AContainsSingularB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{2},
		expected: true,
		index:    1,
	},
	{
		name:     "AEndsWithSingularB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{4},
		expected: true,
		index:    3,
	},
	{
		name:     "AStartsWithB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 2},
		expected: true,
		index:    0,
	},
	{
		name:     "AContainsB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{2, 3},
		expected: true,
		index:    1,
	},
	{
		name:     "AEndsWithB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{3, 4},
		expected: true,
		index:    2,
	},
	{
		name:     "AStartsWithPartialB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 9},
		expected: false,
		index:    3,
	},
	{
		name:     "AEndsWithPartialB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{4, 5},
		expected: false,
		index:    3,
	},
	{
		name:     "AAndBMatch",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 2, 3, 4},
		expected: true,
		index:    0,
	},
	{
		name:     "AIsSubsetOfB",
		a:        []byte{1, 2, 3, 4},
		b:        []byte{1, 2, 3, 4, 5, 6, 7},
		expected: false,
		index:    0,
	},
}

func TestIsAt(t *testing.T) {
	for _, tc := range isAtTestCases {
		t.Run(tc.name, func(*testing.T) {
			actual := utils.IsAt(tc.a, tc.b, tc.index)
			if actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

var joinBytesTestCases = []struct {
	name     string
	b        [][]byte
	expected []byte
}{
	{
		name:     "Nil",
		b:        nil,
		expected: nil,
	},
	{
		name:     "Empty",
		b:        [][]byte{},
		expected: []byte{},
	},
	{
		name: "EmptySlices",
		b: [][]byte{
			{},
			{},
		},
		expected: []byte{},
	},
	{
		name: "EmptyAndNilSlices",
		b: [][]byte{
			{},
			nil,
			{},
		},
		expected: []byte{},
	},
	{
		name: "SingleSliceWithSingleElement",
		b: [][]byte{
			{1},
		},
		expected: []byte{1},
	},
	{
		name: "SingleSliceWithMultipleElements",
		b: [][]byte{
			{1, 2, 3},
		},
		expected: []byte{1, 2, 3},
	},
	{
		name: "MultipleSlicesWithSingleElement",
		b: [][]byte{
			{1},
			{2},
		},
		expected: []byte{1, 2},
	},
	{
		name: "MultipleSlicesWithMultipleElements",
		b: [][]byte{
			{1, 2, 3},
			{4, 5, 6},
		},
		expected: []byte{1, 2, 3, 4, 5, 6},
	},
	{
		name: "MultipleSlicesIncludeNilAndEmpty",
		b: [][]byte{
			{1, 2, 3},
			{},
			{4, 5, 6},
			nil,
			{7, 8, 9},
		},
		expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
	},
}

func TestJoinBytes(t *testing.T) {
	for _, tc := range joinBytesTestCases {
		t.Run(tc.name, func(*testing.T) {
			actual := utils.JoinBytes(tc.b...)
			if !bytes.Equal(actual, tc.expected) {
				t.Errorf("expectd %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestJoinBytesMultipleParams(t *testing.T) {
	b1 := []byte{1, 2, 3}
	b2 := []byte{4, 5, 6}
	b3 := []byte{7, 8, 9}
	expected := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

	actual := utils.JoinBytes(b1, b2, b3)
	if !bytes.Equal(actual, expected) {
		t.Errorf("expectd %v, got %v", expected, actual)
	}
}
