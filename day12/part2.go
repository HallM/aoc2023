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

const MULTIPLIER = 5

const (
	PART_UNKNOWN = iota
	PART_OPERATIONAL
	PART_BROKEN
)

var cache = map[string]int64{}

type Row struct {
	line string
	parts []int
	checksum []int
}

func makeRow(parts []int, checksum []int) *Row {
	var partString []rune
	for _, p := range parts {
		if p == PART_UNKNOWN {
			partString = append(partString, '?')
		} else if p == PART_OPERATIONAL {
			partString = append(partString, '.')
		} else if p == PART_BROKEN {
			partString = append(partString, '#')
		}
	}
	partString = append(partString, ' ')

	var cs []string
	for _, c := range checksum {
		cs = append(cs, fmt.Sprintf("%d", c))
	}

	line := string(partString) + strings.Join(cs, ",")
	return &Row{line: line, parts: parts, checksum: checksum}
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

	var realparts []int
	var realchk []int
	for i := 0; i < MULTIPLIER; i++ {
		if i > 0 {
			realparts = append(realparts, PART_UNKNOWN)
		}
		realparts = append(realparts, parts...)
		realchk = append(realchk, checksum...)
	}

	return makeRow(realparts, realchk)
}

func computePossibles(row *Row) int64 {
	if c, ok := cache[row.line]; ok {
		return c
	}

	if !row.isPossible() {
		cache[row.line] = 0
		return 0
	}
	if len(row.parts) == 0 {
		cache[row.line] = 1
		return 1
	}

	if row.parts[0] == PART_OPERATIONAL {
		next := 1
		for i, p := range row.parts {
			if p != PART_OPERATIONAL {
				next = i
				break
			}
		}
		v := computePossibles(makeRow(row.parts[next:], row.checksum))
		cache[row.line] = v
		return v
	}
	if row.parts[0] == PART_BROKEN {
		if len(row.checksum) == 0 {
			cache[row.line] = 0
			return 0
		}
		required := row.checksum[0]
		if len(row.parts) < required {
			cache[row.line] = 0
			return 0
		}
		for _, p := range row.parts[0:required] {
			if p == PART_OPERATIONAL {
				cache[row.line] = 0
				return 0
			}
		}
		next := required
		if len(row.parts) > required {
			if row.parts[required] == PART_BROKEN {
				cache[row.line] = 0
				return 0
			}
			next++
		}
		v := computePossibles(makeRow(row.parts[next:], row.checksum[1:]))
		cache[row.line] = v
		return v
	}

	ifOperational := append([]int{PART_OPERATIONAL}, row.parts[1:]...)
	ifBroken := append([]int{PART_BROKEN}, row.parts[1:]...)

	a := computePossibles(makeRow(ifOperational, row.checksum))
	b := computePossibles(makeRow(ifBroken, row.checksum))

	cache[row.line] = a + b
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

	var total int64
	for _, line := range lines {
		row := parseRow(line)
		possibles := computePossibles(row)
		log.Printf("Line %s has %d possibles", line, possibles)
		total += possibles
	}

	log.Printf("Total: %d", total)
}
