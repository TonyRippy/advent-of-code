package main

import (
	"bufio"
	"os"
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

func (d direction) Turn() direction {
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

type node struct {
	edges   [4]*node
	end     bool
	visited bool
}

func parseInput(filename string) (*node, *node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var lines [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines = append(lines, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	// Build nodes
	var start *node
	var end *node
	var nodes [][]*node
	for _, line := range lines {
		var row []*node
		for _, c := range line {
			var n *node
			switch c {
			case '#':
				n = nil
			case '.':
				n = &node{}
			case 'S':
				n = &node{}
				start = n
			case 'E':
				n = &node{end: true}
				end = n
			}
			row = append(row, n)
		}
		nodes = append(nodes, row)
	}

	// Connect nodes
	for y, row := range nodes {
		for x, n := range row {
			if n == nil {
				continue
			}
			if y > 0 && nodes[y-1][x] != nil {
				n.edges[N] = nodes[y-1][x]
				nodes[y-1][x].edges[S] = n
			}
			if x > 0 && nodes[y][x-1] != nil {
				n.edges[W] = nodes[y][x-1]
				nodes[y][x-1].edges[E] = n
			}
		}
	}
	return start, end, nil
}

type key struct {
	n *node
	d direction
}

type cache map[key]int

const moveCost = 1
const turnCost = 1000

func traverse(n *node, d direction, cost int, best cache) {
	if n.end {
		key := key{n, N}
		if c, ok := best[key]; !ok || c > cost {
			best[key] = cost
		}
		return
	}
	key := key{n, d}
	if c, ok := best[key]; ok && c <= cost {
		return
	}
	best[key] = cost
	n2 := n.edges[d]
	if n2 != nil {
		traverse(n2, d, cost+moveCost, best)
	}
	d2 := d.Turn() // 90
	n2 = n.edges[d2]
	if n2 != nil {
		traverse(n2, d2, cost+turnCost+moveCost, best)
	}
	d2 = d2.Turn() // 180
	n2 = n.edges[d2]
	if n2 != nil {
		traverse(n2, d2, cost+turnCost+turnCost+moveCost, best)
	}
	d2 = d2.Turn() // -90
	n2 = n.edges[d2]
	if n2 != nil {
		traverse(n2, d2, cost+turnCost+moveCost, best)
	}
}

func Part1(filename string) int {
	start, end, err := parseInput(filename)
	if err != nil {
		panic(err)
	}
	best := make(cache)
	traverse(start, E, 0, best)
	return best[key{end, N}]
}

var visited = make(map[*node]bool)

func traverse2(n *node, d direction, steps []key, cost int, best cache) {
	var k key
	if n.end {
		k = key{n, N}
	} else {
		k = key{n, d}
	}
	if c, ok := best[k]; !ok || c < cost {
		return
	}
	if n.end {
		visited[n] = true
		for _, k2 := range steps {
			visited[k2.n] = true
		}
		return
	}
	steps = append(steps, k)
	n2 := n.edges[d]
	if n2 != nil {
		traverse2(n2, d, steps, cost+moveCost, best)
	}
	d2 := d.Turn() // 90
	n2 = n.edges[d2]
	if n2 != nil {
		traverse2(n2, d2, steps, cost+turnCost+moveCost, best)
	}
	d2 = d2.Turn() // 180
	n2 = n.edges[d2]
	if n2 != nil {
		traverse2(n2, d2, steps, cost+turnCost+turnCost+moveCost, best)
	}
	d2 = d2.Turn() // -90
	n2 = n.edges[d2]
	if n2 != nil {
		traverse2(n2, d2, steps, cost+turnCost+moveCost, best)
	}
}

func Part2(filename string) int {
	start, _, err := parseInput(filename)
	if err != nil {
		panic(err)
	}
	// Traverse once to find the best paths
	best := make(cache)
	traverse(start, E, 0, best)
	// Traverse again to find the visited nodes
	visited = make(map[*node]bool)
	traverse2(start, E, nil, 0, best)
	return len(visited)
}
