package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart1(t *testing.T) {
	v, err := parseInput("test.txt")
	require.NoError(t, err)
	assert.Equal(t, 1928, Part1(v))

	v, err = parseInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 6200294120911, Part1(v))
}

func TestPart2(t *testing.T) {
	v, err := parseInput("test.txt")
	require.NoError(t, err)
	assert.Equal(t, 2858, Part2(v))

	v, err = parseInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 6227018762750, Part2(v))
}
