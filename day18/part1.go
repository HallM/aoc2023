package main

import (
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

type Vertex struct {
	x float64
	y float64
}

var directionOffsets = map[byte]Vertex {
	'U': Vertex{x: 0, y: -1},
	'D': Vertex{x: 0, y: 1},
	'L': Vertex{x: -1, y: 0},
	'R': Vertex{x: 1, y: 0},
}

type Polygon struct {
	vertices []Vertex
}

func (p *Polygon) area() float64 {
	if len(p.vertices) < 3 {
		log.Printf("not enough verts %d", len(p.vertices))
		return 0
	}
	first := p.vertices[0]
	a := first
	var tally float64
	for _, b := range p.vertices[1:] {
		tally = tally + (a.x * b.y) - (a.y * b.x)
		tally = tally + math.Sqrt(math.Pow(a.x - b.x, 2) + math.Pow(a.y - b.y, 2))
		a = b
	}
	tally = tally + (a.x * first.y) - (a.y * first.x)
	tally = tally + math.Sqrt(math.Pow(a.x - first.x, 2) + math.Pow(a.y - first.y, 2))
	return (tally / 2) + 1
}

func diggyDiggyHole(contents string) *Polygon {
	location := Vertex{x: 0, y: 0}
	vertices := []Vertex{}

	for _, line := range strings.Split(contents, "\n") {
		offset := directionOffsets[line[0]]
		endOfNumber := strings.IndexRune(line[2:], ' ')
		move, _ := strconv.ParseInt(line[2:endOfNumber+2], 10, 32)
		location.x += offset.x * float64(move)
		location.y += offset.y * float64(move)
		log.Printf("next up [%.0f, %.0f]", location.x, location.y)
		vertices = append(vertices, location)
	}

	return &Polygon{vertices: vertices}
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

	polygon := diggyDiggyHole(str)
	area := polygon.area()

	log.Printf("Area: %f", area)
}
