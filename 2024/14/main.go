package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

type Robot struct {
	px, py int
	vx, vy int
}

func parseInput(filename string) ([]*Robot, error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var robots []*Robot
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	ms := re.FindAllStringSubmatch(string(input), -1)
	for _, m := range ms {
		r := &Robot{}
		r.px, err = strconv.Atoi(m[1])
		if err != nil {
			return nil, nil
		}
		r.py, err = strconv.Atoi(m[2])
		if err != nil {
			return nil, nil
		}
		r.vx, err = strconv.Atoi(m[3])
		if err != nil {
			return nil, nil
		}
		r.vy, err = strconv.Atoi(m[4])
		if err != nil {
			return nil, nil
		}
		robots = append(robots, r)
	}
	return robots, nil
}

func step(robots []*Robot, w, h int) {
	for _, r := range robots {
		r.px += r.vx
		for r.px < 0 {
			r.px += w
		}
		r.px %= w
		r.py += r.vy
		for r.py < 0 {
			r.py += h
		}
		r.py %= h
	}
}

func Part1(robots []*Robot, seconds, w, h int) int {
	for range seconds {
		step(robots, w, h)
	}
	if w%2 == 0 || h%2 == 0 {
		panic("width/height must be odd")
	}
	// count robots in quadrants
	mx := w / 2
	my := h / 2
	var quadrants [4]int
	for _, r := range robots {
		if r.px < mx {
			if r.py < my {
				quadrants[0]++
			} else if r.py > my {
				quadrants[2]++
			}
		} else if r.px > mx {
			if r.py < my {
				quadrants[1]++
			} else if r.py > my {
				quadrants[3]++
			}
		}
	}
	// Calculate safety factor
	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

func makeMap(w, h int) [][]byte {
	m := make([][]byte, h)
	for i := range h {
		m[i] = make([]byte, w)
	}
	return m
}

func setMap(m [][]byte, robots []*Robot) {
	// initialize the map
	for _, row := range m {
		for i, _ := range row {
			row[i] = '.'
		}
	}
	// mark the robots
	for _, r := range robots {
		m[r.py][r.px] = '#'
	}
}

func printMap(s int, m [][]byte) {
	for _, row := range m {
		fmt.Println(string(row))
	}
	fmt.Printf("============================================ %d\n", s)
}

func isSymmetric(m [][]byte) bool {
	// check to see if the map is symmetric accross the x-axis
	mx := len(m[0]) / 2
	for _, s := range m {
		h1 := s[:mx]
		h2 := s[mx+1:]
		slices.Reverse(h2)
		if !slices.Equal(h1, h2) {
			return false
		}
	}
	return true
}
 func isDenseX(m [][]byte, threshold int) bool {
	for _, row := range m {
		count := 0
		for _, c := range row {
			if c == '#' {
				count++
			}
		}
		if count > threshold {
			return true
		}
	}
	return false
}

func isDenseY(m [][]byte, threshold int) bool {
	for i := range m[0] {
		count := 0
		for _, row := range m {
			if row[i] == '#' {
				count++
			}
		}
		if count > threshold {
			return true
		}
	}
	return false
}

// Printf problem solving....
func Part2(robots []*Robot)  {
	m := makeMap(101, 103)
	var s int
	for {
		fmt.Printf("%d\r", s)
		setMap(m, robots)
		if isDenseX(m, 20) || isDenseY(m, 30) {
			printMap(s, m)
			fmt.Scanln()
		}
		step(robots, 101, 103)
		s++
	}
}

func main() {
	robots, err := parseInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	Part2(robots)
}
