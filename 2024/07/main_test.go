package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuffix(t *testing.T) {
	for _, tc := range []struct {
		value  int
		suffix int
		prefix int
		ok     bool
	}{
		{123, 3, 12, true},
		{123, 23, 1, true},
		// {123, 123, 0, true}, should never happen in practice
		{123, 1, 0, false},
		{123, 12, 0, false},
	} {
		prefix, ok := isSuffixOf(tc.suffix, tc.value)
		if tc.ok {
			assert.True(t, ok)
			assert.Equal(t, tc.prefix, prefix)
		} else {
			assert.False(t, ok)
		}
	}
}

func TestPart1(t *testing.T) {
	eqs, err := parseInput("test.txt")
	require.NoError(t, err)
	assert.Equal(t, 3749, Part1(eqs))

	eqs, err = parseInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 3312271365652, Part1(eqs))
}

func TestCheckPart2(t *testing.T) {
	// 7290: 6 8 6 15
	// 7290 = ((6 * 8) || 6) * 15.
	//        *
	//       / \
	//     ||   15
	//    /  \
	//   *    6
	//  / \
	// 6   8
	raw := []int{6, 8, 6, 15}
	slices.Reverse(raw)
	assert.True(t, CheckPart2(7290, raw))
}

func TestPart2(t *testing.T) {
	eqs, err := parseInput("test.txt")
	require.NoError(t, err)
	assert.Equal(t, 11387, Part2(eqs))

	eqs, err = parseInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 509463489296712, Part2(eqs))
}
