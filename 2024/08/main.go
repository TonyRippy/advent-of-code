package main

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"
)

type Point struct {
	Y, X int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.Y, p.X)
}

type AntennaLists map[rune][]Point

type Map struct {
	Lines    []string
	Antennas AntennaLists
	W, H     int
	debug    bool
}

func (m *Map) InBounds(p Point) bool {
	return p.X >= 0 && p.X < m.W && p.Y >= 0 && p.Y < m.H
}

func (m *Map) parseAntennas(line string, y int) {
	for x, c := range line {
		if c == '.' || c == '#' {
			continue
		}
		p := Point{y, x}
		if alist, ok := m.Antennas[c]; ok {
			m.Antennas[c] = append(alist, p)
		} else {
			m.Antennas[c] = []Point{p}
		}
	}
}

func parseInput(filename string) (*Map, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("missing first line")
	}

	line := strings.TrimSpace(scanner.Text())
	out := &Map{
		Antennas: make(AntennaLists),
		W:        len(line),
	}
	y := 0
	out.Lines = append(out.Lines, line)
	out.parseAntennas(line, y)

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		out.Lines = append(out.Lines, line)
		y++
		out.parseAntennas(line, y)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	out.H = y + 1
	return out, nil
}

func combinations(n int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if j == i {
					continue
				}
				if !yield(i, j) {
					return
				}
			}
		}
	}
}

func antinodes1(_ *Map, p1 Point, p2 Point) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		dx := p2.X - p1.X
		dy := p2.Y - p1.Y
		antinode := Point{p1.Y - dy, p1.X - dx}
		if !yield(antinode) {
			return
		}
		antinode = Point{p2.Y + dy, p2.X + dx}
		yield(antinode)
	}
}

func (m *Map) Part1() int {
	var count int
	seen := make(map[Point]bool)
	for _, alist := range m.Antennas {
		for i, j := range combinations(len(alist)) {
			p1 := alist[i]
			p2 := alist[j]
			for antinode := range antinodes1(m, p1, p2) {
				// Is this point on the map?
				if !m.InBounds(antinode) {
					continue
				}
				// Check if this point is already an antinode
				if _, ok := seen[antinode]; ok {
					continue
				}
				seen[antinode] = true
				count++
			}
		}
	}
	return count
}

func antinodes2(m *Map, p1 Point, p2 Point) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		yield(p1)
		yield(p2)
		dx := p2.X - p1.X
		dy := p2.Y - p1.Y
		for {
			p1 = Point{p1.Y - dy, p1.X - dx}
			if !m.InBounds(p1) {
				break
			}
			if !yield(p1) {
				return
			}
		}
		for {
			p2 = Point{p2.Y + dy, p2.X + dx}
			if !m.InBounds(p2) {
				break
			}
			if !yield(p2) {
				return
			}
		}
	}
}

func (m *Map) Part2() int {
	var count int
	seen := make(map[Point]bool)
	for _, alist := range m.Antennas {
		for i, j := range combinations(len(alist)) {
			p1 := alist[i]
			p2 := alist[j]
			for antinode := range antinodes2(m, p1, p2) {
				// Is this point on the map?
				if !m.InBounds(antinode) {
					continue
				}
				if _, ok := seen[antinode]; ok {
					continue
				}
				seen[antinode] = true
				count++
			}
		}
	}
	return count
}

func main() {
	m, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1: %d\n", m.Part1())
	fmt.Printf("Part 2: %d\n", m.Part2())
}
