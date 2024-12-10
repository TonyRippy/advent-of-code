package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

type Node struct {
	Value int
	Edges []*Node
}

func parseInput(filename string) ([]*Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var trailheads []*Node
	var nodes [][]*Node

	// Build nodes
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		ns := make([]*Node, len(line))
		for i, c := range line {
			if !unicode.IsDigit(c) {
				continue
			}
			d := int(c - '0')
			ns[i] = &Node{Value: d}
		}
		nodes = append(nodes, ns)
	}

	// Build edges
	for i, ns := range nodes {
		for j, n := range ns {
			if n == nil {
				continue
			}
			if n.Value == 0 {
				trailheads = append(trailheads, n)
			}
			if i > 0 {
				n2 := nodes[i-1][j]
				if n2 != nil && n2.Value == n.Value+1 {
					n.Edges = append(n.Edges, n2)
				}
			}
			if i < len(nodes)-1 {
				n2 := nodes[i+1][j]
				if n2 != nil && n2.Value == n.Value+1 {
					n.Edges = append(n.Edges, n2)
				}
			}
			if j > 0 {
				n2 := ns[j-1]
				if n2 != nil && n2.Value == n.Value+1 {
					n.Edges = append(n.Edges, n2)
				}
			}
			if j < len(ns)-1 {
				n2 := ns[j+1]
				if n2 != nil && n2.Value == n.Value+1 {
					n.Edges = append(n.Edges, ns[j+1])
				}
			}
		}
	}
	return trailheads, scanner.Err()
}

func push(stack []*Node, n *Node) []*Node {
	return append(stack, n)
}

func pop(stack []*Node) ([]*Node, *Node) {
	last := len(stack) - 1
	if last < 0 {
		return nil, nil
	}
	return stack[:last], stack[last]
}

func score(n *Node) int {
	// fmt.Printf("stack: n=%d, %v\n", len(stack), stack)
	// Walk the graph
	stack := make([]*Node, 0)
	visits := make(map[*Node]bool)
	nines := make(map[*Node]bool)
	for n != nil {
		//fmt.Printf("visit: %v (%d) -> %v\n", n, n.Value, n.Edges)
		visits[n] = true
		if n.Value == 9 {
			nines[n] = true
		}
		for _, e := range n.Edges {
			if _, ok := visits[e]; !ok {
				stack = push(stack, e)
			}
		}
		stack, n = pop(stack)
	}
	return len(nines)
}

func Part1(trailheads []*Node) int {
	var sum int
	for _, n := range trailheads {
		sum += score(n)
	}
	return sum
}

func rating(n *Node, ratings map[*Node]int) int {
	if r, ok := ratings[n]; ok {
		return r
	}
	if n.Value == 9 {
		ratings[n] = 1
		return 1
	}
	r := 0
	for _, e := range n.Edges {
		r += rating(e, ratings)
	}
	ratings[n] = r
	return r
}

func Part2(trailheads []*Node) int {
	ratings := make(map[*Node]int)
	var sum int
	for _, n := range trailheads {
		sum += rating(n, ratings)
	}
	return sum
}

func main() {
	_, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
