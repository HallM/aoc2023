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

func (m *RangeMap) rangesOf(rng *Range) []*Range {
	var ret []*Range
	for ri, r := range m.ranges {
		// We may have to split ranges

		// if the range size matches OR map's range is bigger, return 1 range
		// if the range is larger, then create a new range for the next

		if rng.src >= r.src {
			d := rng.src - r.src
			// Is the range big enough to contain this src?
			if d >= r.size {
				continue
			}
			endA := rng.src + rng.size
			endB := r.src + r.size

			// If we have remaining, then keep processing
			if endA > endB {
				ret = append(ret, &Range{
					src: r.dest + d,
					size: r.size - d,
				})
				// create a new range to process the remainder
				rng = &Range{src: endB, size: endA - endB}
			} else {
				ret = append(ret, &Range{
					src: r.dest + d,
					size: rng.size,
				})
				// no more range to process
				rng = &Range{src: endB, size: 0}
				break
			}
		}
	}
	if rng.size > 0 {
		ret = append(ret, rng)
	}
	return ret
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

		var after []*Range
		for x, v := range values {
			after = append(after, rm.rangesOf(v)...)
		}
		values = after
	}

	min := values[0].src
	for _, v := range values[1:] {
		if v.src < min {
			min = v.src
		}
	}
	return min, nil
}

func parseSeeds(line string) ([]*Range, error) {
	var seeds []*Range
	parts := strings.Split(line[7:], " ")
	if len(parts) % 2 != 0 {
		return nil, fmt.Errorf("Seed numbers come in pairs (start, length)")
	}
	for i := 0; i < len(parts); i+=2 {
		start, err := strconv.ParseInt(parts[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse number from %q as int: %w", parts[i], err)
		}

		size, err := strconv.ParseInt(parts[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse number from %q as int: %w", parts[i+1], err)
		}

		seeds = append(seeds, &Range{src: start, size: size})
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
