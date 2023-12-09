package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

type Line struct {
	history []int64
}

func parseFile(contents string) ([]*Line, error) {
	var lines []*Line
	for _, l := range strings.Split(contents, "\n") {
		var history []int64
		for _, n := range strings.Split(strings.TrimSpace(l), " ") {
			v, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Cannot parse number from %q as int: %w", n, err)
			}
			history = append(history, v)
		}
		lines = append(lines, &Line{history})
	}
	return lines, nil
}

func (l *Line) extrapolate() int64 {
	pyramid := [][]int64{l.history}
	var i int
	for {
		var diffs []int64
		for j, x := range pyramid[i][1:] {
			// because of the slice, "j" is 1 less than the real index
			d := x - pyramid[i][j]
			diffs = append(diffs, d)
		}
		pyramid = append(pyramid, diffs)
		i++

		allzero := true
		for _, x := range pyramid[i] {
			allzero = allzero && x == 0
		}
		if allzero {
			break
		}
	}

	var extraps []int64
	for _ = range pyramid {
		extraps = append(extraps, 0)
	}

	// the last row of the pyramid just is 0 anyway
	for index := len(pyramid)-1; index > 0; index-- {
		i := index-1
		extraps[i] = extraps[index] + pyramid[i][len(pyramid[i])-1]
	}
	return extraps[0]
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

	lines, err := parseFile(string(contents))
	if err != nil {
		log.Fatal(err)
	}

	var total int64
	for i, l := range lines {
		extrap := l.extrapolate()
		log.Printf("Row %d, extrapolated %d", i+1, extrap)
		total += extrap
	}

	log.Printf("Sum: %d", total)
}
