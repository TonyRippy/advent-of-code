package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func parseInput(filename string) (list1, list2 []int, err error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, nil, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    fields := strings.Fields(line)
    if len(fields) != 2 {
      return nil, nil, fmt.Errorf("invalid input: %q", line)
    }
    n, err := strconv.Atoi(fields[0])
    if err != nil {
      return nil, nil, fmt.Errorf("invalid list item: %q", fields[0])
    }
    list1 = append(list1, n)
    n, err = strconv.Atoi(fields[1])
    if err != nil {
      return nil, nil, fmt.Errorf("invalid list item: %q", fields[0])
    }
    list2 = append(list2, n)
  }
  if err := scanner.Err(); err != nil {
    return nil, nil, err
  }

  return list1, list2, nil
}

func Part1(list1, list2 []int) {
  var distance int
  for i := 0; i < len(list1); i++ {
    d := list1[i] - list2[i]
    if d < 0 {
      d = -d
    }
    distance += d
  }
  fmt.Printf("Total distance: %d\n", distance)
}

func findCount(list []int, n int) int {
  i, found := slices.BinarySearch(list, n)
  if !found {
    return 0
  }
  count := 1
  for i += 1; i < len(list) && list[i] == n; i++ {
    count++
  }
  return count
}

func Part2(list1, list2 []int) {
  var similarity int
  for _, n := range list1 {
    count := findCount(list2, n)
    similarity += n * count
  }
  fmt.Printf("Similarity: %d\n", similarity)
}

func main() {
  list1, list2, err := parseInput(os.Args[1])
  if err != nil {
    log.Fatalf("failed to parse input: %v", err)
  }
  
  slices.Sort(list1)
  slices.Sort(list2)

  Part1(list1, list2)
  Part2(list1, list2)
}
