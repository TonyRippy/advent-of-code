package main

import (
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Button struct {
	dx, dy int
}

type Prize struct {
	x, y int
}

type Machine struct {
	a, b Button
	p    Prize
}

func (m *Machine) findCost(p Prize) int {
	// Solve the system of equations that will give the answer
	a1 := m.a.dx * m.b.dy
	p1 := p.x * m.b.dy

	a2 := m.a.dy * m.b.dx
	p2 := p.y * m.b.dx

	n := p1 - p2
	d := a1 - a2
	if n%d != 0 {
		return 0
	}
	a := n / d

	n = p.x - a*m.a.dx
	d = m.b.dx
	if n%d != 0 {
		return 0
	}
	b := n / d
	return a*3 + b
}

func parseInput(filename string) ([]*Machine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	input, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var machines []*Machine
	re := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+)`)
	ms := re.FindAllStringSubmatch(string(input), -1)
	for _, m := range ms {
		ax, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, nil
		}
		ay, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, nil
		}
		bx, err := strconv.Atoi(m[3])
		if err != nil {
			return nil, nil
		}
		by, err := strconv.Atoi(m[4])
		if err != nil {
			return nil, nil
		}
		px, err := strconv.Atoi(m[5])
		if err != nil {
			return nil, nil
		}
		py, err := strconv.Atoi(m[6])
		if err != nil {
			return nil, nil
		}
		machines = append(machines, &Machine{
			a: Button{ax, ay},
			b: Button{bx, by},
			p: Prize{px, py},
		})
	}
	return machines, nil
}

func Part1(machines []*Machine) int {
	var total int
	for _, m := range machines {
		cost := m.findCost(m.p)
		total += cost
	}
	return total
}

func Part2(machines []*Machine) int {
	var total int
	for _, m := range machines {
		p := Prize{
			m.p.x + 10000000000000,
			m.p.y + 10000000000000,
		}
		cost := m.findCost(p)
		total += cost
	}
	return total
}

func main() {
	_, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
