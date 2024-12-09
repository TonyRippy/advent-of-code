package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"unicode"
)

type Range struct {
	Free bool
	File int
	Pos  int
	Size int
}

func parseInput(filename string) ([]*Range, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ranges []*Range
	pos := 0
	fid := 0
	free := false
	reader := bufio.NewReader(file)
	for {
		c, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !unicode.IsDigit(c) {
			return nil, fmt.Errorf("invalid character: %q", c)
		}
		d := int(c - '0')
		r := &Range{free, 0, pos, d}
		if !free {
			r.File = fid
			fid++
		}
		pos += d
		free = !free
		ranges = append(ranges, r)
	}
	return ranges, nil
}

func printRanges(ranges []*Range) {
	var s strings.Builder
	for _, r := range ranges {
		var c rune
		if r.Free {
			c = '.'
		} else {
			c = '0' + rune(r.File)
		}
		for range r.Size {
			s.WriteRune(c)
		}
	}
	fmt.Println(s.String())
}

func defrag1(ranges []*Range) []*Range {
	i := 0               // first free range
	j := len(ranges) - 1 // last file
	for i < j {
		free := ranges[i]
		if !free.Free {
			i++
			continue
		}
		f := ranges[j]
		if f.Free {
			j--
			continue
		}
		if free.Size == f.Size {
			// The file fits perfectly, swap them
			free.Free = false
			free.File = f.File
			f.Free = true
			f.File = -1
			continue
		}
		if free.Size < f.Size {
			// The file is too big, split it
			remainder := &Range{
				Free: false,
				File: f.File,
				Pos:  f.Pos + free.Size,
				Size: f.Size - free.Size,
			}
			free.Free = false
			free.File = f.File
			// Free the first part of the file
			f.Size = free.Size
			f.Free = true
			f.File = -1
			f.Size = free.Size
			// Insert remainder after f
			ranges = slices.Insert(ranges, j+1, remainder)
			j++
			continue
		}
		// else the file is too small, split the free range
		remainder := &Range{
			Free: true,
			File: -1,
			Pos:  free.Pos + f.Size,
			Size: free.Size - f.Size,
		}
		free.Free = false
		free.File = f.File
		free.Size = f.Size
		f.Free = true
		f.File = -1
		// Insert remainder after free
		ranges = slices.Insert(ranges, i+1, remainder)
	}
	return ranges
}

func calcChecksum(ranges []*Range) int {
	var sum int
	var pos int
	for _, r := range ranges {
		if r.Free {
			pos += r.Size
			continue
		}
		end := r.Pos + r.Size
		for pos < end {
			sum += pos * r.File
			pos++
		}
	}
	return sum
}

func Part1(ranges []*Range) int {
	return calcChecksum(defrag1(ranges))
}

func defrag2(ranges []*Range) []*Range {
	j := len(ranges) - 1 // last file
	for j >= 0 {
		f := ranges[j]
		if f.Free {
			j--
			continue
		}
		i := 0 // first free range
		for {
			if i >= j {
				j--
				break
			}
			free := ranges[i]
			if !free.Free || free.Size < f.Size {
				i++
				continue
			}
			if free.Size == f.Size {
				// The file fits perfectly, swap them
				free.Free = false
				free.File = f.File
				f.Free = true
				f.File = -1
				break
			}
			// else the file is too small, split the free range
			remainder := &Range{
				Free: true,
				File: -1,
				Pos:  free.Pos + f.Size,
				Size: free.Size - f.Size,
			}
			free.Free = false
			free.File = f.File
			free.Size = f.Size
			f.Free = true
			f.File = -1
			// Insert remainder after free
			ranges = slices.Insert(ranges, i+1, remainder)
			break
		}
	}
	return ranges
}

func Part2(ranges []*Range) int {
	ranges = defrag2(ranges)
	return calcChecksum(ranges)
}

func main() {
	ranges, err := parseInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	printRanges(ranges)
	ranges = defrag2(ranges)
	printRanges(ranges)
}
