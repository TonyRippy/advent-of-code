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
		{"test1a.txt", 140},
		{"test1b.txt", 772},
		{"test1c.txt", 1930},
		{"input.txt", 1370100},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			m, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, m.Part1())
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range []struct {
		filename string
		expected int
	}{
		{"test1a.txt", 80},
		{"test1b.txt", 436},
		{"test2a.txt", 236},
		{"test2b.txt", 368},
		{"input.txt", 818286},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			m, err := parseInput(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, m.Part2())
		})
	}
}
