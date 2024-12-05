package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type PageNumber int

func parsePageNumber(s string) (PageNumber, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return PageNumber(n), nil
}

// A map of "key comes before values"
type Rules map[PageNumber][]PageNumber

type Update struct {
	Valid bool
	Pages []PageNumber
	Index map[PageNumber]int
}

func newUpdate() *Update {
	return &Update{
		Valid: true,
		Pages: nil,
		Index: make(map[PageNumber]int),
	}
}

func parseUpdate(s string) (*Update, error) {
	parts := strings.Split(s, ",")
	out := newUpdate()
	out.Pages = make([]PageNumber, 0, len(parts))
	for i, part := range parts {
		pn, err := parsePageNumber(part)
		if err != nil {
			return nil, err
		}
		out.Pages = append(out.Pages, pn)
		if _, ok := out.Index[pn]; ok {
			out.Valid = false
		} else {
			out.Index[pn] = i
		}
	}
	return out, nil
}

func (u *Update) Check(r Rules) bool {
	for i, before := range u.Pages {
		rules, ok := r[before]
		if !ok {
			continue // nothing to check
		}
		for _, after := range rules {
			if j, ok := u.Index[after]; ok { 
				if j < i {
					return false
				}
			}
		}
	}
	return true
}

func (u *Update) CorrectOrder(r Rules) {
	pages := u.Pages
	swap:
	for {
		for i, before := range pages {
			rules, ok := r[before]
			if !ok {
				continue // nothing to check
			}
			for _, after := range rules {
				if j, ok := u.Index[after]; ok {
					if j < i {
						pages[i], pages[j] = pages[j], pages[i]
						u.Index[before] = j
						u.Index[after] = i
						continue swap
					}
				}
			}
		}
		return
	}
}

func (u *Update) MiddlePage() PageNumber {
	n := len(u.Pages)
	if n % 2 == 0 {
		panic(fmt.Sprintf("even number of pages: %v", u.Pages))
	}
	return u.Pages[n/2]
}

func parseInput(filename string) (Rules, []Update, error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, nil, err
  }
  defer file.Close()

	var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
		line := scanner.Text()
    lines = append(lines, strings.TrimSpace(line))
	}
  if err := scanner.Err(); err != nil {
    return nil, nil, err
  }
	rules := make(Rules)
	var next int
	for i, line := range lines {
		if len(line) == 0 {
			next = i + 1
			break
		}
		parts := strings.Split(line, "|")
		before, err := parsePageNumber(parts[0])
		if err != nil {
			return nil, nil, err
		}
		after, err := parsePageNumber(parts[1])
		if err != nil {
			return nil, nil, err
		}
		rules[before] = append(rules[before], after)
	}
	updates := make([]Update, 0, len(lines)-next)
	for _, line := range lines[next:] {
		update, err := parseUpdate(line)
		if err != nil {
			return nil, nil, err
		}
		updates = append(updates, *update)
	}
  return rules, updates, nil
}

func main() {
	rules, update, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var sum1, sum2 PageNumber
	for _, u := range update {
		if u.Check(rules) {
			sum1 += u.MiddlePage()
			continue
		}
		u.CorrectOrder(rules)
		sum2 += u.MiddlePage()
	}
	fmt.Printf("Part 1: sum = %d\n", sum1)
	fmt.Printf("Part 2: sum = %d\n", sum2)
}
