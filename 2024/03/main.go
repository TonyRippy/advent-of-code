package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type state struct {
	sum     int
	enabled bool
}

func newState() *state {
	return &state{0, true}
}

type op interface {
	Apply(*state)
}

type mul struct {
	a, b int
}

func (m mul) Apply(s *state) {
	if s.enabled {
		s.sum += m.a * m.b
	}
}

type doOp struct {
}

func (doOp) Apply(s *state) {
	s.enabled = true
}

type dontOp struct {
}

func (dontOp) Apply(s *state) {
	s.enabled = false
}

var pattern1 = regexp.MustCompile(`mul\((\d+),(\d+)\)`)

func parsePart1(input string, c chan<- op) {
	matches := pattern1.FindAllStringSubmatch(string(input), -1)
	for _, m := range matches {
		a, err := strconv.Atoi(m[1])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}
		c <- mul{a, b}
	}
	close(c)
}

var pattern2 = regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d+),(\d+)\)`)

func parsePart2(input string, c chan<- op) {
	matches := pattern2.FindAllStringSubmatch(string(input), -1)
	for _, m := range matches {
		if m[0] == "do()" {
			c <- doOp{}
			continue
		}
		if m[0] == "don't()" {
			c <- dontOp{}
			continue
		}
		a, err := strconv.Atoi(m[1])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}
		c <- mul{a, b}
	}
	close(c)
}


func applyAll(c <-chan op) int {
	s := newState()
	for o := range c {
		o.Apply(s)
	}
	return s.sum
}

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan op)
	go parsePart1(string(input), c)
	sum := applyAll(c)
	fmt.Printf("Part 1: sum = %d\n", sum)

	c = make(chan op)
	go parsePart2(string(input), c)
	sum = applyAll(c)
	fmt.Printf("Part 2: sum = %d\n", sum)
}
