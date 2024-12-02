package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseInput(filename string) (input [][]int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		values := make([]int, len(fields))
		for i, f := range fields {
			values[i], err = strconv.Atoi(f)
			if err != nil {
				return nil, err
			}
		}
		input = append(input, values)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

func IsGraduallyIncreasing(report []int) bool {
	for i := 1; i < len(report); i++ {
		if report[i] <= report[i-1] {
			return false
		}
	}
	return true
}

func IsGraduallyDecreasing(report []int) bool {
	for i := 1; i < len(report); i++ {
		if report[i] >= report[i-1] {
			return false
		}
	}
	return true
}

func IsSafeReport1(report []int) bool {
	if !IsGraduallyIncreasing(report) && !IsGraduallyDecreasing(report) {
		return false
	}
	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]
		if diff < 0 {
			diff = -diff
		}
		// diff must be at least 1, because the increase/decrease check fails if zero
		if diff > 3 {
			return false
		}
	}
	return true
}

func IsSafeReport2(report []int) bool {
	if IsSafeReport1(report) {
		return true
	}
	n := len(report)
	for i := 0; i < n; i++ {
		newReport := make([]int, 0, n-1)
		newReport = append(newReport, report[:i]...)
		newReport = append(newReport, report[i+1:]...)
		if IsSafeReport1(newReport) {
			return true
		}
	}
	return false
}

func main() {
	reports, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var safe int
	for _, report := range reports {
		if IsSafeReport1(report) {
			safe++
		}
	}
	fmt.Printf("Part 1: %d reports are safe\n", safe)

	safe = 0
	for _, report := range reports {
		if IsSafeReport2(report) {
			safe++
		}
	}
	fmt.Printf("Part 2: %d reports are safe\n", safe)
}
