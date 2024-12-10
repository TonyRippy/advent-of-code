package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart1(t *testing.T) {
	for _, tc := range []struct{
		filename string
		expected int
	} {
		{"test1a.txt", 1},
		{"test1b.txt", 2},
		{"test1c.txt", 4},
		{"test1d.txt", 3},
		{"test1e.txt", 36},
		{"input.txt", 644},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			v, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, Part1(v))
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range []struct{
		filename string
		expected int
	} {
		{"test1e.txt", 81},
		{"input.txt", 1366},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			v, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, Part2(v))
		})
	}
}
