package main

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
)

type direction struct {
	name   string
	dx, dy int
}

func (d direction) String() string {
	return d.name
}

var (
	N  = direction{"N", 0, -1}
	S  = direction{"S", 0, 1}
	E  = direction{"E", 1, 0}
	W  = direction{"W", -1, 0}
	NE = direction{"NE", 1, -1}
	NW = direction{"NW", -1, -1}
	SE = direction{"SE", 1, 1}
	SW = direction{"SW", -1, 1}
)

type puzzle struct {
	lines []string
}

func newPuzzle(lines []string) *puzzle {
	return &puzzle{lines}
}

func readPuzzle(filename string) (*puzzle, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return newPuzzle(lines), nil
}

func (p *puzzle) get(x, y int) (byte, bool) {
	if y < 0 || y >= len(p.lines) {
		return 0, false
	}
	line := p.lines[y]
	if x < 0 || x >= len(line) {
		return 0, false
	}
	return line[x], true
}

func (p *puzzle) walk(x, y int, d direction) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		for {
			x += d.dx
			y += d.dy
			r, ok := p.get(x, y)
			if !ok {
				return
			}
			if !yield(r) {
				return
			}
		}
	}
}

func (p *puzzle) findString(s iter.Seq[byte], target string) bool {
	var i int
	for r := range s {
		if r != target[i] {
			break
		}
		i++
		if i == len(target) {
			return true
		}
	}
	return false
}

func (p *puzzle) countXmasAt(x, y int) int {
	var count int
	for _, d := range []direction{N, S, E, W, NE, NW, SE, SW} {
		seq := p.walk(x, y, d)
		if p.findString(seq, "MAS") {
			count++
		}
	}
	return count
}

func (p *puzzle) countXmas() int {
	var count int
	// Search entire puzzle for all 'X' runes
	for y, line := range p.lines {
		for x, r := range line {
			if r == 'X' {
				count += p.countXmasAt(x, y)
			}
		}
	}
	return count
}

func (p *puzzle) charAtOffset(x, y int, d direction, target byte) bool {
	c, ok := p.get(x+d.dx, y+d.dy)
	return ok && c == target
}

func (p *puzzle) findCrossMasAt(x, y int) bool {
	for _, cross := range []struct {
		m1, s1, m2, s2 direction
	}{
		// This was a bug! I has mistakenly assumed up-down-left-right was also a valid cross.
		// {N, S, E, W},
		// {N, S, W, E},
		// {S, N, E, W},
		// {S, N, W, E},
		{NW, SE, NE, SW},
		{NW, SE, SW, NE},
		{SE, NW, NE, SW},
		{SE, NW, SW, NE},
	} {
		if p.charAtOffset(x, y, cross.m1, 'M') &&
			p.charAtOffset(x, y, cross.s1, 'S') &&
			p.charAtOffset(x, y, cross.m2, 'M') &&
			p.charAtOffset(x, y, cross.s2, 'S') {
			return true
		}
	}
	return false
}

func (p *puzzle) countCrossMas() int {
	var count int
	// Search entire puzzle for all 'A' runes
	for y, line := range p.lines {
		for x, r := range line {
			if r == 'A' && p.findCrossMasAt(x, y) {
				count++
			}
		}
	}
	return count
}

func main() {
	p, err := readPuzzle(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1: found %d XMAS\n", p.countXmas())
	fmt.Printf("Part 2: found %d X-MAS\n", p.countCrossMas())
}
