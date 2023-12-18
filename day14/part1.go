package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

type Rock struct {
	y int
	isRound bool
}

type Platform struct {
	cols [][]Rock
	height int
}

func (p *Platform) computeLoad() int {
	var load int
	for _, rocks := range p.cols {
		for _, r := range rocks {
			if r.isRound {
				load += p.height - r.y
			}
		}
	}
	return load
}

func parsePlatform(contents string) *Platform {
	lines := strings.Split(contents, "\n")
	height := len(lines)
	width := len(lines[0])

	cols := make([][]Rock, width)
	colNextY := map[int]int{}

	for y, line := range lines {
		for x, r := range line {
			if r == 'O' {
				cols[x] = append(cols[x], Rock{y: colNextY[x], isRound: true})
				colNextY[x]++
			} else if r == '#' {
				colNextY[x] = y
				cols[x] = append(cols[x], Rock{y: colNextY[x], isRound: false})
				colNextY[x]++
			}
		}
	}
	return &Platform{cols: cols, height: height}
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
	str := strings.ReplaceAll(string(contents), "\r", "")

	platform := parsePlatform(str)
	total := platform.computeLoad()

	log.Printf("Sum: %d", total)
}
