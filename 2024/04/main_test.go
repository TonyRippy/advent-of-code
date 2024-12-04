package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart2(t *testing.T) {
	p, err := readPuzzle("test.txt")
	require.NoError(t, err)
	fmt.Println(p.lines[2])
	c, ok := p.get(6,2)
	assert.True(t, ok)
	assert.Equal(t, byte('A'), c)
	assert.True(t, p.findCrossMasAt(6,2))
}
