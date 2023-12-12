package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

type Galaxy struct {
	id int
	x int
	y int
}

type Universe struct {
	galaxies []*Galaxy
}

func parseMap(contents string) *Universe {
	var galaxies []*Galaxy

	nextId := 1
	yExpand := 0

	columnCounts := map[int]int{}

	lines := strings.Split(contents, "\n")
	if len(lines) == 1 {
		return &Universe{}
	}

	width := len(lines[0])

	for y, line := range lines {
		line = strings.TrimSpace(line)

		hadOne := false
		for x, c := range line {
			if c == '#' {
				hadOne = true
				columnCounts[x] = columnCounts[x] + 1
				galaxies = append(galaxies, &Galaxy{id: nextId, x: x, y: y + yExpand})
				nextId++
			}
		}
		if !hadOne {
			yExpand++
		}
	}

	xMapping := map[int]int{}
	xExpand := 0
	for x := 0; x < width; x++ {
		xMapping[x] = x + xExpand
		if columnCounts[x] == 0 {
			xExpand++
		}
	}

	log.Printf("%d rows, %d cols", yExpand, xExpand)

	for _, g := range galaxies {
		g.x = xMapping[g.x]
	}

	return &Universe{galaxies: galaxies}
}

func (u *Universe) sumDistances() int {
	total := 0
	for i, a := range u.galaxies {
		for j, b := range u.galaxies {
			if i == j {
				break
			}
			x := a.x - b.x
			if x < 0 {
				x = -x
			}
			y := a.y - b.y
			if x < 0 {
				y = -y
			}
			path := x + y
			// log.Printf("%d -> %d is %d units", b.id, a.id, path)
			total += path
		}
	}
	return total
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

	universe := parseMap(string(contents))
	total := universe.sumDistances()

	log.Printf("Sum: %d", total)
}
