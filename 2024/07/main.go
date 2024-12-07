package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Equation struct {
	TestValue int
	Args []int
}

func (eq *Equation) Check(f func(int, []int) bool) bool {
	args := make([]int, len(eq.Args))
	copy(args, eq.Args)
	slices.Reverse(args)
	return f(eq.TestValue, args)
}

func parseInput(filename string) ([]*Equation, error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer file.Close()

	var eqs []*Equation
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i == -1 {
			return nil, errors.New("invalid input, colon not found")
		}
		eq := &Equation{}
		eq.TestValue, err = strconv.Atoi(line[:i])
		if err != nil {
			return nil, fmt.Errorf("invalid test value: %q", line[:i])
		}
		line = line[i+1:]
    fields := strings.Fields(line)
		if len(fields) == 0 {
			return nil, errors.New("invalid input, no list items")
		}
		eq.Args = make([]int, len(fields))
		for i, f := range fields {
			eq.Args[i], err = strconv.Atoi(strings.TrimSpace(f))
			if err != nil {
				return nil, fmt.Errorf("invalid list item: %q", f)
			}
		}
		slices.Reverse(eq.Args)
    eqs = append(eqs, eq)
  }
  if err := scanner.Err(); err != nil {
    return nil, err
  }
  return eqs, nil
}

func CheckPart1(tv int, args []int) bool {
	// if len(args) == 0 {
	// 	return false
	// }
	arg := args[0]
	if len(args) == 1 {
		return tv == arg
	}
	// Can multiplication be used?
	if tv % arg == 0 {
		// Try solution
		if CheckPart1(tv / arg, args[1:]) {
			return true
		}
	}
	// Can addition be used?
	return CheckPart1(tv - arg, args[1:])
}

func Part1(eqs []*Equation) int {
	var sum int
	for _, eq := range eqs {
		if CheckPart1(eq.TestValue, eq.Args) {
			sum += eq.TestValue
		}
	}
	return sum
}

// Returns the prefix of value after suffix removed
func isSuffixOf(suffix, value int) (int, bool) {
	suffixstr := strconv.Itoa(suffix)
	valuestr := strconv.Itoa(value)
	if !strings.HasSuffix(valuestr, suffixstr) {
		return 0, false
	}
	prefixstr := valuestr[:len(valuestr)-len(suffixstr)]
	prefix, _ := strconv.Atoi(prefixstr)
	return prefix, true
}

func CheckPart2(tv int, args []int) bool {
	// if len(args) == 0 {
	// 	return false
	// }
	arg := args[0]
	if len(args) == 1 {
		return tv == arg
	}
	// Can multiplication be used?
	if tv % arg == 0 {
		// Try solution
		if CheckPart2(tv / arg, args[1:]) {
			return true
		}
	}
	// Can addition be used?
	if CheckPart2(tv - arg, args[1:]) {
		return true
	}
	// Last resort, try concatenation
	if prefix, ok := isSuffixOf(arg, tv); ok {
		return CheckPart2(prefix, args[1:])
	}
	return false
}


func Part2(eqs []*Equation) int {
	var sum int
	for _, eq := range eqs {
		if CheckPart2(eq.TestValue, eq.Args) {
			sum += eq.TestValue
		}
	}
	return sum
}

func main() {
	eqs, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part1: %d\n", Part1(eqs))
	fmt.Printf("Part2: %d\n", Part2(eqs))
}
