package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

const (
	PART_UNKNOWN = iota
	PART_OPERATIONAL
	PART_BROKEN
)

type Row struct {
	operational int
	broken int
	unknown int
}

func (r *Row) isValid(potential int) bool {
	// Every bit in the operational should be set in potential
	if (potential & r.operational) != r.operational {
		return false
	}
	// Every bit set in the broken should be NOT set in potential
	if ((^potential) & r.broken) != r.broken {
		return false
	}
	// Every bit set in the unknown doesn't matter
}

func (r *Row) minRequired() {
	if len(r.checksum) == 0 {
		return 0
	}

	var val int
	for _, c := range r.checksum {
		val += c
	}
	return val + len(r.checksum) - 1
}

func parseRow(line string) *Row {
	line = strings.TrimSpace(line)

	var parts []int
	var checksumStart int
	for i, r := range line {
		if r == '.' {
			parts = append(parts, PART_OPERATIONAL)
		} else if r == '#' {
			parts = append(parts, PART_BROKEN)
		} else if r == '?' {
			parts = append(parts, PART_UNKNOWN)
		} else if r == ' ' {
			checksumStart = i+1
		}
	}

	var checksum []int
	checkParts := strings.Split(line[checksumStart:], ",")
	for _, s := range checkParts {
		value, _ := strconv.ParseInt(s, 10, 32)
		checksum = append(checksum, int(value))
	}
	return &Row{parts: parts, checksum: checksum}
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

	lines := strings.Split(string(contents), "\n")

	var total int
	for _, line := range lines {
		row := parseRow(line)
		possibles := computePossibles(row)
		log.Printf("Line %s has %d possibles", line, possibles)
		total += possibles
	}

	log.Printf("Total: %d", total)
}
