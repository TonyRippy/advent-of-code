package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart2(t *testing.T) {
	rules, update, err := parseInput("test.txt")
	require.NoError(t, err)
	var sum PageNumber
	for _, u := range update {
		if u.Check(rules) {
			continue
		}
		u.CorrectOrder(rules)
		sum += u.MiddlePage()
	}
	assert.Equal(t, PageNumber(123), sum)
}
