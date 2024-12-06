package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type direction int

const (
	N direction = iota
	S
	E
	W
)

type point struct {
	x, y int
}

func (p point) next(d direction) point {
	switch d {
	case N:
		return point{p.x, p.y - 1}
	case S:
		return point{p.x, p.y + 1}
	case E:
		return point{p.x + 1, p.y}
	case W:
		return point{p.x - 1, p.y}
	}
	panic("invalid direction")
}

func (d direction) Turn90() direction {
	switch d {
	case N:
		return E
	case E:
		return S
	case S:
		return W
	case W:
		return N
	}
	panic("invalid direction")
}

type floor struct {
	lines [][]byte
	guard point
}

func readFloor(filename string) (*floor, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var guard point
	var lines [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexByte(line, '^')
		if i >= 0 {
			guard = point{i, len(lines)}
		}
		lines = append(lines, []byte(line))
	}
	return &floor{lines, guard}, nil
}

func (f *floor) set(x, y int, b byte) {
	f.lines[y][x] = b
}

func (f *floor) isExit(p point) bool {
	return p.x < 0 || p.y < 0 || p.y >= len(f.lines) || p.x >= len(f.lines[p.y])
}

func (f *floor) isObstruction(p point) bool {
	b := f.lines[p.y][p.x]
	return b == '#' || b == 'O'
}

func (f *floor) part1() (map[point][]direction, bool) {
	guard := f.guard
	d := N
	visited := make(map[point][]direction)
	visited[guard] = []direction{N}
	for {
		next := guard.next(d)
		if f.isExit(next) {
			return visited, true
		}
		for f.isObstruction(next) {
			d = d.Turn90()
			next = guard.next(d)
		}
		ds, ok := visited[next]
		if ok && slices.Contains(ds, d) {
			return visited, false
		}
		visited[next] = append(ds, d)
		guard = next
	}
}

func (f *floor) part2(visited map[point][]direction) int {
	// Loop through all positions the guard would visit, replacing open space ('.') with an obstacle. ('O')
	var loops int
	for p := range visited {
		if f.lines[p.y][p.x] != '.' {
			continue
		}
		f.lines[p.y][p.x] = 'O'
		if _, ok := f.part1(); !ok {
			loops++
		}
		f.lines[p.y][p.x] = '.'
	}
	return loops
}

func main() {
	floor, err := readFloor(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	visited, _ := floor.part1()
	fmt.Printf("Part 1: walked %d steps\n", len(visited))
	fmt.Printf("Part 2: %d loops\n", floor.part2(visited))
}
