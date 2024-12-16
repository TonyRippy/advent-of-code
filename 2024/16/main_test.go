package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart1(t *testing.T) {
	for _, tc := range []struct {
		filename string
		expected int
	}{
		{"test1a.txt", 7036},
		{"test1b.txt", 11048},
		{"input.txt", 72400},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			assert.Equal(t, tc.expected, Part1(tc.filename))
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range []struct {
		filename string
		expected int
	}{
		{"test1a.txt", 45},
		{"test1b.txt", 64},
		{"input.txt", 435},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			assert.Equal(t, tc.expected, Part2(tc.filename))
		})
	}
}
