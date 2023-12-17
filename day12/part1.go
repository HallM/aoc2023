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
	parts []int
	checksum []int
}

func (r *Row) isPossible() bool {
	if len(r.parts) == 0 && len(r.checksum) == 0 {
		return true
	}

	req := r.minRequired()
	if (req + len(r.checksum) - 1) > len(r.parts) {
		return false
	}
	numBroke := r.numMaybeDamaged()
	if numBroke < req {
		return false
	}
	return true
}

func (r *Row) minRequired() int {
	if len(r.checksum) == 0 {
		return 0
	}

	var val int
	for _, c := range r.checksum {
		val += c
	}
	return val
}

func (r *Row) numMaybeDamaged() int {
	var op int
	for _, p := range r.parts {
		if p == PART_BROKEN || p == PART_UNKNOWN {
			op++
		}
	}
	return op
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

func computePossibles(row *Row, indent string) int {
	nextIndent := indent + "    "

	if !row.isPossible() {
		return 0
	}
	if len(row.parts) == 0 {
		return 1
	}

	if row.parts[0] == PART_OPERATIONAL {
		return computePossibles(&Row{parts: row.parts[1:], checksum: row.checksum}, nextIndent)
	}
	if row.parts[0] == PART_BROKEN {
		if len(row.checksum) == 0 {
			return 0
		}
		required := row.checksum[0]
		if len(row.parts) < required {
			return 0
		}
		for _, p := range row.parts[0:required] {
			if p == PART_OPERATIONAL {
				return 0
			}
		}
		next := required
		if len(row.parts) > required {
			if row.parts[required] == PART_BROKEN {
				return 0
			}
			next++
		}
		return computePossibles(&Row{parts: row.parts[next:], checksum: row.checksum[1:]}, nextIndent)
	}

	ifOperational := append([]int{PART_OPERATIONAL}, row.parts[1:]...)
	ifBroken := append([]int{PART_BROKEN}, row.parts[1:]...)

	a := computePossibles(&Row{parts: ifOperational, checksum: row.checksum}, nextIndent)
	b := computePossibles(&Row{parts: ifBroken, checksum: row.checksum}, nextIndent)
	return a + b
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
		possibles := computePossibles(row, "")
		log.Printf("Line %s has %d possibles", line, possibles)
		total += possibles
	}

	log.Printf("Total: %d", total)
}
