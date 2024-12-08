package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCombinations(t *testing.T) {
	type pair struct{ i, j int }
	for _, tc := range []struct {
		n   int
		out []pair
	}{
		{0, nil},
		{1, nil},
		{2, []pair{{0, 1}}},
		{3, []pair{{0, 1}, {0, 2}, {1, 2}}},
		{4, []pair{{0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 3}, {2, 3}}},
	} {
		t.Run(strconv.Itoa(tc.n), func(t *testing.T) {
			var got []pair
			for i, j := range combinations(tc.n) {
				got = append(got, pair{i, j})
			}
			assert.Equal(t, tc.out, got)
		})
	}
}
func TestPart1(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"test1a.txt", 2},
		{"test1b.txt", 4},
		{"test1c.txt", 4},
		{"test1d.txt", 14},
		{"input.txt", 361},
	} {
		t.Run(tc.in, func(t *testing.T) {
			m, err := parseInput(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, m.Part1())
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"test2a.txt", 9},
		{"input.txt", 1249},
	} {
		t.Run(tc.in, func(t *testing.T) {
			m, err := parseInput(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, m.Part2())
		})
	}
}
