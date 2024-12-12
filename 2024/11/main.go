package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

func parseInput(filename string) (list []int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	all, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(strings.TrimSpace(string(all)))
	list = make([]int, len(fields))
	for i, field := range fields {
		list[i], err = strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
	}
	return
}

func pow10(exp int) int {
	pow := 1
	for ; exp > 0; exp-- {
		pow *= 10
	}
	return pow
}

type cacheKey struct {
	stone, blinks int
}

var cache sync.Map

func blink(stone int, blinks int) int {
	if blinks == 0 {
		return 1
	}
	blinks--
	if stone == 0 {
		return cachedBlink(1, blinks)
	}
	digits := int(math.Floor(math.Log10(float64(stone)))) + 1
	if digits%2 == 0 {
		pow := pow10(digits / 2)
		return cachedBlink(stone%pow, blinks) + cachedBlink(stone/pow, blinks)
	}
	return cachedBlink(stone*2024, blinks)
}

func cachedBlink(stone int, blinks int) int {
	key := cacheKey{stone, blinks}
	if out, ok := cache.Load(key); ok {
		return out.(int)
	}
	out := blink(stone, blinks)
	cache.Store(key, out)
	return out
}

func Part1(stones []int, blinks int) int {
	var sum int
	for _, stone := range stones {
		sum += blink(stone, blinks)
	}
	return sum
}

func Part2(stones []int) int {
	in := make(chan int, 1000)
	out := make(chan int, 1000)

	// Create workers
	var wg sync.WaitGroup
	for range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range in {
				out <- cachedBlink(n, 75)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	// Send input to workers
	for _, n := range stones {
		in <- n
	}
	close(in)

	// Collect results
	var sum int
	for n := range out {
		sum += n
	}
	return sum
}

func main() {
	stones, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 2: %d\n", Part2(stones))
}
