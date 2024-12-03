package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePart1(t *testing.T) {
	const input = "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"
	c := make(chan op)
	go parsePart1(input, c)
	assert.Equal(t, 161, applyAll(c))
}

func TestParsePart2(t *testing.T) {
	const input = "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"
	c := make(chan op)
	go parsePart2(input, c)
	assert.Equal(t, 48, applyAll(c))
}
