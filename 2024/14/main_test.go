package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSymmetric(t *testing.T) {
	s := "..#.......#.."
	mx := len(s) / 2
	h1 := s[:mx]
	h2 := s[mx+1:]
	assert.NotEqual(t, h1, h2)

	// reverse the second half
	bs := []byte(h2)
	slices.Reverse(bs)
	h2 = string(bs)

	// They should be the same if they're symmetric
	assert.Equal(t, h1, h2)
}

func TestPart1(t *testing.T) {
	for _, tc := range []struct {
		filename string
		w, h     int
		want     int
	}{
		{"test.txt", 11, 7, 12},
		{"input.txt", 101, 103, 232253028},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			robots, err := parseInput(tc.filename)
			require.NoError(t, err)
			got := Part1(robots, 100, tc.w, tc.h)
			assert.Equal(t, tc.want, got)
		})
	}
}
