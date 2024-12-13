package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart1(t *testing.T) {
	for _, tc := range []struct {
		filename string
		expected int
	}{
		{"test.txt", 480},
		{"input.txt", 33481},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			ms, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, Part1(ms))
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range []struct {
		filename string
		expected int
	}{
		{"input.txt", 92572057880885},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			ms, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, Part2(ms))
		})
	}
}
