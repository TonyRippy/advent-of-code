package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
			line = line[:i] + "X" + line[i+1:]
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
	if b, ok := f.get(p.x, p.y); ok && b == '#' {
		return true
	}
	return false
}

func (f *floor) walk(guard point, d direction) int {
	marked := 1
	for {
		next := guard.next(d)
		if f.isExit(next) {
			return marked
		}
		for f.isObstruction(next) {
			d = d.Turn90()
			next = guard.next(d)
		}
		b, ok := f.get(next.x, next.y)
		if !ok { 
			panic(fmt.Sprintf("Unexpected position %v", next))
		}
		if b != 'X' {
			if b != '.' { panic(fmt.Sprintf("Unexpected value at position %v: %q", next, b)) } 
			f.set(next.x, next.y, 'X')
			marked += 1
		}
		guard = next
	}
}

func main() {
	floor, err := readFloor(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	marked := floor.walk(floor.guard, N)
	fmt.Printf("Part 1: walked %d steps\n", marked)
}
