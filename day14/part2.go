package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

const cycles = 1000000000

type Slot struct {
	x int
	y int
	isRound bool
	isSquare bool
}

type Platform struct {
	grid []Slot
	width int
	height int
}

func (p *Platform) computeHash() string {
	hash := make([]rune, len(p.grid))
	for i, r := range p.grid {
		if r.isRound {
			hash[i] = 'O'
		} else if r.isSquare {
			hash[i] = '#'
		} else {
			hash[i] = '.'
		}
	}
	return string(hash)
}

func (p *Platform) computeLoad() int {
	var load int
	for i, r := range p.grid {
		if r.isRound {
			y := i / p.width
			load += p.height - y
		}
	}
	return load
}

func (p *Platform) rotateNorth() {
	colNextY := map[int]int{}

	for i := range p.grid {
		x := i % p.width
		y := i / p.width

		if p.grid[i].isRound {
			moveY := colNextY[x]
			p.grid[i].isRound = false
			p.grid[moveY*p.width + x].isRound = true
			colNextY[x]++
		} else if p.grid[i].isSquare {
			colNextY[x] = y+1
		}
	}
}

func (p *Platform) rotateSouth() {
	colNextY := map[int]int{}
	for x := 0; x < p.width; x++ {
		colNextY[x] = p.height - 1
	}

	for i := len(p.grid)-1; i >= 0; i-- {
		x := i % p.width
		y := i / p.width

		if p.grid[i].isRound {
			moveY := colNextY[x]
			p.grid[i].isRound = false
			p.grid[moveY*p.width + x].isRound = true
			colNextY[x]--
		} else if p.grid[i].isSquare {
			colNextY[x] = y-1
		}
	}
}


func (p *Platform) rotateWest() {
	rowNextX := 0

	for i := range p.grid {
		x := i % p.width
		y := i / p.width
		if x == 0 {
			rowNextX = 0
		}

		if p.grid[i].isRound {
			moveX := rowNextX
			p.grid[i].isRound = false
			p.grid[y*p.width + moveX].isRound = true
			rowNextX++
		} else if p.grid[i].isSquare {
			rowNextX = x+1
		}
	}
}

func (p *Platform) rotateEast() {
	rowNextX := p.width - 1

	for i := len(p.grid)-1; i >= 0; i-- {
		x := i % p.width
		y := i / p.width
		if x == p.width - 1 {
			rowNextX = p.width - 1
		}

		if p.grid[i].isRound {
			moveX := rowNextX
			p.grid[i].isRound = false
			p.grid[y*p.width + moveX].isRound = true
			rowNextX--
		} else if p.grid[i].isSquare {
			rowNextX = x-1
		}
	}
}

func parsePlatform(contents string) *Platform {
	lines := strings.Split(contents, "\n")
	height := len(lines)
	width := len(lines[0])

	grid := make([]Slot, width*height)
	i := 0

	for _, line := range lines {
		for _, r := range line {
			if r == 'O' {
				grid[i].isRound = true
			} else if r == '#' {
				grid[i].isSquare = true
			}
			i++
		}
	}
	return &Platform{grid: grid, width: width, height: height}
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

	seen := map[string]int{}

	var offset, loopLength int
	for i := 0; i < cycles; i++ {
		platform.rotateNorth()
		platform.rotateWest()
		platform.rotateSouth()
		platform.rotateEast()

		hash := platform.computeHash()
		if prev, ok := seen[hash]; ok {
			offset = prev
			loopLength = i - prev
			log.Printf("Seen previous cycle(%d) after %d cycles", prev, i)
			break
		}
		seen[hash] = i
	}

	loops := (cycles - offset) / loopLength
	midLoop := (cycles - offset) % loopLength
	log.Printf("I think offset=%d looplength=%d so %d cycles hits in %d loops at %d cycles inside that loop", offset, loopLength, cycles, loops, midLoop)

	platform = parsePlatform(str)
	for i := 0; i < (offset + midLoop); i++ {
		platform.rotateNorth()
		platform.rotateWest()
		platform.rotateSouth()
		platform.rotateEast()
	}

	total := platform.computeLoad()

	log.Printf("Sum: %d", total)
}
