package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"strconv"
)

var filePath = flag.String("file", "part1test.txt", "File path")

type RangeMap struct {
	ranges []*Range
}

func NewRangeMap(ranges []*Range) *RangeMap {
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].src < ranges[j].src
	})
	return &RangeMap{ranges: ranges}
}

func (m *RangeMap) valueOf(src int64) int64 {
	for _, r := range m.ranges {
		if src >= r.src {
			d := src - r.src
			if d < r.size {
				return r.dest + d
			}
		}
	}
	return src
}

type Range struct {
	dest int64
	src int64
	size int64
}

func computeClosestLocation(contents string) (int64, error) {
	// cause windows
	contents = strings.ReplaceAll(contents, "\r\n", "\n")
	blocks := strings.Split(contents, "\n\n")

	values, err := parseSeeds(blocks[0])
	if err != nil {
		return 0, err
	}

	for bn, block := range blocks[1:] {
		lines := strings.Split(block, "\n")
		rm, err := parseRangemap(lines[1:])
		if err != nil {
			return 0, err
		}

		for i, v := range values {
			after := rm.valueOf(v)
			values[i] = after
		}
	}

	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min, nil
}

func parseSeeds(line string) ([]int64, error) {
	var seeds []int64
	for _, s := range strings.Split(line[7:], " ") {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse number from %q as int: %w", s, err)
		}
		seeds = append(seeds, id)
	}
	return seeds, nil
}

func parseRangemap(lines []string) (*RangeMap, error) {
	var ranges []*Range

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var src, dest, size int64
		var err error

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("Range line expecting 3 parts (dest, src, size), but got %d from %q", len(parts), line)
		}

		dest, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse dest from %q as int: %w", parts[0], err)
		}

		src, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse src from %q as int: %w", parts[1], err)
		}

		size, err = strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse size from %q as int: %w", parts[2], err)
		}

		ranges = append(ranges, &Range{src: src, dest: dest, size: size})
	}
	return NewRangeMap(ranges), nil
}

func main() {
	flag.Parse()

	if *filePath == "" {
		log.Fatalf("Must specify the file!")
	}

	contents, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	score, err := computeClosestLocation(string(contents))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Number: %d", score)
}
