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
		{"test1a.txt", 10092},
		{"test1b.txt", 2028},
		{"input.txt", 1412971},
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
		{"test1a.txt", 9021},
		{"input.txt", 1429299},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			assert.Equal(t, tc.expected, Part2(tc.filename))
		})
	}
}
