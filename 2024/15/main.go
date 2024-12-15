package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Map struct {
	lines  [][]byte
	moves  []byte
	rx, ry int // robot position
}

func (m *Map) canPush(x, y, dx, dy int) bool {
	x += dx
	y += dy
	switch m.lines[y][x] {
	case '#':
		return false
	case '.':
		return true
	case 'O':
		return m.canPush(x, y, dx, dy)
	case '[':
		if dy != 0 {
			return m.canPush(x, y, dx, dy) && m.canPush(x+1, y, dx, dy)
		} else {
			return m.canPush(x, y, dx, dy)
		}
	case ']':
		if dy != 0 {
			return m.canPush(x, y, dx, dy) && m.canPush(x-1, y, dx, dy)
		} else {
			return m.canPush(x, y, dx, dy)
		}
	default:
		panic("invalid character " + string(m.lines[y][x]))
	}
}

// split into pushX and pushY?
// keep a range of x values that need to be pushed?

func (m *Map) doPushX(x, y, dx int, last byte) {
	x += dx
	switch m.lines[y][x] {
	case '#':
		if last != '.' {
			panic("still holding a box")
		}
	case '.':
		m.lines[y][x] = last
	case 'O', '[', ']':
		m.lines[y][x], last = last, m.lines[y][x]
		m.doPushX(x, y, dx, last)
	default:
		panic("invalid character " + string(m.lines[y][x]))
	}
}

func (m *Map) doPushY(x1, x2, y, dy int, last ...byte) {
	y += dy
	line := m.lines[y]
	c := rune(line[x1])
	c2 := rune(line[x2])
	if c == '#' || c2 == '#' {
		if last[x1] != '.' || last[x2] != '.' {
			panic("still holding a box")
		}
		return
	}
	line[x1] = last[0]
	switch c {
	case 'O':
		m.doPushY(x1, x1, y, dy, 'O')
	case '[':
		m.doPushY(x1, x1+1, y, dy, '[', ']')
		line[x1+1] = '.'
	case ']':
		m.doPushY(x1-1, x1, y, dy, '[', ']')
		line[x1-1] = '.'
	}
	if x2 == x1 {
		return
	}
	line[x2] = last[1]
	switch c2 {
	case 'O':
		m.doPushY(x2, x2, y, dy, 'O')
	case '[':
		m.doPushY(x2, x2+1, y, dy, '[', ']')
		line[x2+1] = '.'
	}
}

func (m *Map) Run() {
	// var step int
	for _, move := range m.moves {
		// fmt.Printf("Step %d - move %c\n", step, move)
		// m.Print()
		// step++

		var dx, dy int
		switch move {
		case '^':
			dx = 0
			dy = -1
		case 'v':
			dx = 0
			dy = 1
		case '<':
			dx = -1
			dy = 0
		case '>':
			dx = 1
			dy = 0
		}
		if m.canPush(m.rx, m.ry, dx, dy) {
			if dy != 0 {
				m.doPushY(m.rx, m.rx, m.ry, dy, '.')
			} else {
				m.doPushX(m.rx, m.ry, dx, '.')
			}
			m.rx += dx
			m.ry += dy
		}
	}
}

func (m *Map) Sum() int {
	var sum int
	for y, line := range m.lines {
		for x, c := range line {
			if c == 'O' || c == '[' {
				sum += 100*y + x
			}
		}
	}
	return sum
}

func (m *Map) Print() {
	for y, line := range m.lines {
		for x, c := range line {
			if x == m.rx && y == m.ry {
				fmt.Print("@")
			} else {
				fmt.Print(string(c))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (m *Map) findRobot() {
	for y, line := range m.lines {
		for x, c := range line {
			if c == '@' {
				line[x] = '.'
				m.rx, m.ry = x, y
				return
			}
		}
	}
	panic("robot not found")
}

func parseInput(filename string, parseLine func(string) []byte) (*Map, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := &Map{}
	scanner := bufio.NewScanner(file)

	// Scan the map
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			break
		}
		m.lines = append(m.lines, parseLine(line))
	}
	m.findRobot()

	// Scan the moves
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		m.moves = append(m.moves, []byte(line)...)
	}
	return m, scanner.Err()
}

func parseInput1(filename string) (*Map, error) {
	return parseInput(filename, func(line string) []byte {
		return []byte(line)
	})
}

func parseInput2(filename string) (*Map, error) {
	return parseInput(filename, func(line string) []byte {
		out := make([]byte, 0, len(line)*2)
		for _, c := range line {
			switch c {
			case '#':
				out = append(out, '#', '#')
			case '@':
				out = append(out, '@', '.')
			case '.':
				out = append(out, '.', '.')
			case 'O':
				out = append(out, '[', ']')
			}
		}
		return out
	})
}

func Part1(filename string) int {
	m, err := parseInput1(filename)
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
	return m.Sum()
}

func Part2(filename string) int {
	m, err := parseInput2(filename)
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
	return m.Sum()
}
