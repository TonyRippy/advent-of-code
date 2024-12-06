package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type point struct {
	x, y int
}

func (p point) next(d direction) point {
	return point{p.x + d.dx, p.y + d.dy}
}

type direction struct {
	name   string
	dx, dy int
}

func (d direction) String() string {
	return d.name
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

var (
	N  = direction{"N", 0, -1}
	S  = direction{"S", 0, 1}
	E  = direction{"E", 1, 0}
	W  = direction{"W", -1, 0}
)

type floor struct {
	lines []string
	guard point
}

func readFloor(filename string) (*floor, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var guard point
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, '^')
		if i >= 0 {
			guard = point{i, len(lines)}
		}
		lines = append(lines, line)
	}
	return &floor{lines, guard}, nil
}

func (p *floor) get(x, y int) (byte, bool) {
	if y < 0 || y >= len(p.lines) {
		return 0, false
	}
	line := p.lines[y]
	if x < 0 || x >= len(line) {
		return 0, false
	}
	return line[x], true
}

func (f *floor) set(x, y int, b byte) {
	f.lines[y] = f.lines[y][:x] + string(b) + f.lines[y][x+1:]
}

func (f *floor) isExit(p point) bool {
	return p.x < 0 || p.y < 0 || p.y >= len(f.lines) || p.x >= len(f.lines[p.y])
}

func (f *floor) isObstruction(p point) bool {
	if b, ok := f.get(p.x, p.y); ok && (b == '#' || b == 'O') {
		return true
	}
	return false
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
		if b, _ := f.get(p.x, p.y); b != '.' {
			continue
		}
		f.set(p.x, p.y, 'O')
		if _, ok := f.part1(); !ok {
			loops++
		}
		f.set(p.x, p.y, '.')
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
