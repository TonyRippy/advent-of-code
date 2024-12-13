package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Square struct {
	id int
	c  rune
}

type Map struct {
	Width   int
	Height  int
	Squares [][]Square
}

func (m *Map) assignIds() {
	changed := true
	for changed {
		changed = false
		for y, row := range m.Squares {
			for x, me := range row {
				min := me.id
				if x > 0 && m.Squares[y][x-1].c == me.c {
					id := m.Squares[y][x-1].id
					if id < min {
						min = id
					}
				}
				if x < m.Width-1 && m.Squares[y][x+1].c == me.c {
					id := m.Squares[y][x+1].id
					if id < min {
						min = id
					}
				}
				if y > 0 && m.Squares[y-1][x].c == me.c {
					id := m.Squares[y-1][x].id
					if id < min {
						min = id
					}
				}
				if y < m.Height-1 && m.Squares[y+1][x].c == me.c {
					id := m.Squares[y+1][x].id
					if id < min {
						min = id
					}
				}
				if min != me.id {
					m.Squares[y][x].id = min
					changed = true
				}
			}
		}
	}
}

func (m *Map) Part1() int {
	perimeter := make(map[int]int)
	area := make(map[int]int)
	for y, row := range m.Squares {
		for x, square := range row {
			if a, ok := area[square.id]; !ok {
				area[square.id] = 1
			} else {
				area[square.id] = a + 1
			}
			var boundaries int
			if x == 0 || row[x-1].id != square.id {
				boundaries += 1
			}
			if x == (m.Width-1) || row[x+1].id != square.id {
				boundaries += 1
			}
			if y == 0 || m.Squares[y-1][x].id != square.id {
				boundaries += 1
			}
			if y == (m.Height-1) || m.Squares[y+1][x].id != square.id {
				boundaries += 1
			}
			if p, ok := perimeter[square.id]; !ok {
				perimeter[square.id] = boundaries
			} else {
				perimeter[square.id] = p + boundaries
			}
		}
	}
	var total int
	for id, area := range area {
		p := perimeter[id]
		price := area * p
		total += price
	}
	return total
}

func (m *Map) corners(x, y int) int {
	square := m.Squares[y][x]

	//   [ 1] [ 2] [4]
	//   [ 8]   X  [16]
	//   [32] [64] [128]
	var boundaries int
	if y > 0 {
		line := m.Squares[y-1]
		if x > 0 && line[x-1].id == square.id {
			boundaries |= 1
		}
		if line[x].id == square.id {
			boundaries |= 2
		}
		if x < (m.Width-1) && line[x+1].id == square.id {
			boundaries |= 4
		}
	}
	if x > 0 && m.Squares[y][x-1].id == square.id {
		boundaries |= 8
	}
	if x < (m.Width-1) && m.Squares[y][x+1].id == square.id {
		boundaries |= 16
	}
	if y < (m.Height - 1) {
		line := m.Squares[y+1]
		if x > 0 && line[x-1].id == square.id {
			boundaries |= 32
		}
		if line[x].id == square.id {
			boundaries |= 64
		}
		if x < (m.Width-1) && line[x+1].id == square.id {
			boundaries |= 128
		}
	}
	//   [ 1] [ 2] [4]
	//   [ 8]   X  [16]
	//   [32] [64] [128]
	var corner int
	mask := boundaries & (1 | 2 | 8)
	if mask == 0 || mask == 1 || mask == (2|8) {
		// NW
		corner += 1
	}
	mask = boundaries & (2 | 4 | 16)
	if mask == 0 || mask == 4 || mask == (2|16) {
		// NE
		corner += 1
	}
	mask = boundaries & (8 | 32 | 64)
	if mask == 0 || mask == 32 || mask == (8|64) {
		// NE
		corner += 1
	}
	mask = boundaries & (16 | 64 | 128)
	if mask == 0 || mask == 128 || mask == (16|64) {
		// SE
		corner += 1
	}
	return corner
}

func (m *Map) Part2() int {
	corners := make(map[int]int)
	area := make(map[int]int)
	for y, row := range m.Squares {
		for x, square := range row {
			if a, ok := area[square.id]; !ok {
				area[square.id] = 1
			} else {
				area[square.id] = a + 1
			}
			c := m.corners(x, y)
			if cs, ok := corners[square.id]; !ok {
				corners[square.id] = c
			} else {
				corners[square.id] = cs + c
			}
		}
	}
	var total int
	for id, area := range area {
		sides := corners[id]
		price := area * sides
		total += price
	}
	return total
}

func parseInput(filename string) (*Map, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var squares [][]Square
	nextid := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var line []Square
		for _, c := range strings.TrimSpace(scanner.Text()) {
			line = append(line, Square{nextid, c})
			nextid++
		}
		squares = append(squares, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	m := &Map{
		Width:   len(squares[0]),
		Height:  len(squares),
		Squares: squares,
	}
	m.assignIds()
	return m, nil
}

func main() {
	_, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
