package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart1(t *testing.T) {
	for _, tc := range []struct {
		filename string
		blinks   int
		expected int
	}{
		{"test.txt", 6, 22},
		{"test.txt", 25, 55312},
		{"input.txt", 25, 217812},
	} {
		t.Run(fmt.Sprintf("%s/%d", tc.filename, tc.blinks), func(t *testing.T) {
			stones, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, Part1(stones, tc.blinks))
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range []struct {
		filename string
		expected int
	}{
		{"input.txt", 259112729857522},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			v, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, Part2(v))
		})
	}
}
