package main

import (
	"container/heap"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

var filePath = flag.String("file", "part1test.txt", "File path")

var parseNumber = map[rune]int {
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
}

const (
	DIRECTION_NORTH = iota
	DIRECTION_EAST
	DIRECTION_SOUTH
	DIRECTION_WEST
)
const (
	MOVE_LEFT = iota
	MOVE_RIGHT
	MOVE_AHEAD
)

type Vertex struct {
	x int
	y int
}

var changeDirection = map[int]map[int]int {
	DIRECTION_NORTH: map[int]int{
		MOVE_LEFT: DIRECTION_WEST,
		MOVE_RIGHT: DIRECTION_EAST,
		MOVE_AHEAD: DIRECTION_NORTH,
	},
	DIRECTION_EAST: map[int]int{
		MOVE_LEFT: DIRECTION_NORTH,
		MOVE_RIGHT: DIRECTION_SOUTH,
		MOVE_AHEAD: DIRECTION_EAST,
	},
	DIRECTION_SOUTH: map[int]int{
		MOVE_LEFT: DIRECTION_EAST,
		MOVE_RIGHT: DIRECTION_WEST,
		MOVE_AHEAD: DIRECTION_SOUTH,
	},
	DIRECTION_WEST: map[int]int{
		MOVE_LEFT: DIRECTION_SOUTH,
		MOVE_RIGHT: DIRECTION_NORTH,
		MOVE_AHEAD: DIRECTION_WEST,
	},
}

var directionOffsets = map[int]Vertex {
	DIRECTION_NORTH: Vertex{x: 0, y: -1},
	DIRECTION_EAST: Vertex{x: 1, y: 0},
	DIRECTION_SOUTH: Vertex{x: 0, y: 1},
	DIRECTION_WEST: Vertex{x: -1, y: 0},
}

type Graph struct {
	weights [][]int
	width int
	height int
}

type Path struct {
	location Vertex
	blocksMoved int
	direction int
}

type Searcher struct {
	path Path
	heat int
}

type SearchQueue []*Searcher

func (sq SearchQueue) Len() int {
	return len(sq)
}
func (sq SearchQueue) Less(i, j int) bool {
	if sq[i].heat != sq[j].heat {
		return sq[i].heat < sq[j].heat
	}
	if sq[i].path.location.y != sq[j].path.location.y {
		return sq[i].path.location.y < sq[j].path.location.y
	}
	return sq[i].path.location.x < sq[j].path.location.x
}
func (sq SearchQueue) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}
func (sq *SearchQueue) Push(s interface{}) {
	*sq = append(*sq, s.(*Searcher))
}
func (sq *SearchQueue) Pop() interface{} {
	start := *sq
	size := len(start)
	searcher := start[size-1]
	start[size-1] = nil
	*sq = start[0:size-1]
	return searcher
}

func (g *Graph) canMoveTo(loc Vertex) bool {
	if loc.x < 0 || loc.y < 0 || loc.x >= g.width || loc.y >= g.height {
		return false
	}
	return true
}

func (s *Searcher) makeMoves(g *Graph) []*Searcher {
	var searchers []*Searcher

	moves := make([]int, 0, 3)
	moves = append(moves, MOVE_LEFT)
	moves = append(moves, MOVE_RIGHT)
	if s.path.blocksMoved < 3 {
		moves = append(moves, MOVE_AHEAD)
	}

	for _, move := range moves {
		blocksMoved := 1
		if move == MOVE_AHEAD {
			blocksMoved = s.path.blocksMoved + 1
		}

		newDirection := changeDirection[s.path.direction][move]
		offset := directionOffsets[newDirection]
		coord := Vertex{x: s.path.location.x+offset.x, y: s.path.location.y+offset.y}
		if g.canMoveTo(coord) {
			searchers = append(searchers, &Searcher{
				path: Path{
					location: coord,
					blocksMoved: blocksMoved,
					direction: newDirection,
				},
				heat: s.heat + g.weights[coord.y][coord.x],
			})
		}
	}

	return searchers
}

func (g *Graph) pathToTarget(start, end Vertex) *Searcher {
	searchers := SearchQueue{
		&Searcher{
			path: Path{location: Vertex{x: start.x, y: start.y+1}, direction: DIRECTION_SOUTH, blocksMoved: 1},
			heat: g.weights[start.y+1][start.x],
		},
		&Searcher{
			path: Path{location: Vertex{x: start.x+1, y: start.y}, direction: DIRECTION_EAST, blocksMoved: 1},
			heat: g.weights[start.y][start.x+1],
		},
	}
	heap.Init(&searchers)

	visited := map[Path]bool{}

	startTime := time.Now()
	defer func() {
		t := time.Now()
		d := t.Sub(startTime).Microseconds()
		if d > 0 {
			log.Printf("Elapsed %v", d)
		}
	}()

	for len(searchers) > 0 {
		s := heap.Pop(&searchers).(*Searcher)
		if visited[s.path] {
			continue
		}

		if s.path.location.x == end.x && s.path.location.y == end.y {
			return s
		}
		visited[s.path] = true

		for _, newS := range s.makeMoves(g) {
			heap.Push(&searchers, newS)
		}
	}
	return nil
}

func parseGraph(contents string) *Graph {
	lines := strings.Split(contents, "\n")
	height := len(lines)
	width := len(lines[0])

	weights := make([][]int, height)

	for y, l := range lines {
		weights[y] = make([]int, width)
		for x, r := range l {
			weights[y][x] = parseNumber[r]
		}
	}

	return &Graph{weights: weights, width: width, height: height}
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

	graph := parseGraph(str)
	searcher := graph.pathToTarget(Vertex{x: 0, y: 0}, Vertex{x: graph.width-1, y: graph.height-1})

	log.Printf("Heat cost: %d", searcher.heat)
}
