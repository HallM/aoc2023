package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part2test.txt", "File path")

const (
	PIPE_GROUND = iota
	PIPE_START
	PIPE_NS
	PIPE_EW
	PIPE_NE
	PIPE_NW
	PIPE_SW
	PIPE_SE

	CONNECT_OUTSIDE
)

const (
	FROM_NORTH = iota
	FROM_EAST
	FROM_SOUTH
	FROM_WEST
)

var charToPipe = map[rune]int {
	'.': PIPE_GROUND,
	'S': PIPE_START,
	'|': PIPE_NS,
	'-': PIPE_EW,
	'L': PIPE_NE,
	'J': PIPE_NW,
	'7': PIPE_SW,
	'F': PIPE_SE,
}
var pipeToChar = map[int]rune {
	PIPE_GROUND: ' ',
	PIPE_START: 'S',
	PIPE_NS: '|',
	PIPE_EW: '-',
	PIPE_NE: 'L',
	PIPE_NW: 'J',
	PIPE_SW: '7',
	PIPE_SE: 'F',
	CONNECT_OUTSIDE: 'O',
}

type Seeker struct {
	id int
	x int
	y int
	arrivedFrom int
	distance int
	indices []int
}

type PipeMap struct {
	cells []int
	width int
	height int
	startX int
	startY int
}

func (m *PipeMap) print() {
	for y := 0; y < m.height; y++ {
		var row []rune
		for x := 0; x < m.width; x++ {
			row = append(row, pipeToChar[m.value(x, y)])
		}
		log.Printf("%d: %s", y, string(row))
	}
}

func (m *PipeMap) indexFor(x, y int) int {
	return y * m.width + x
}

func (m *PipeMap) value(x, y int) int {
	return m.cells[m.indexFor(x, y)]
}

func (m *PipeMap) set(x, y int, value int) {
	m.cells[m.indexFor(x, y)] = value
}


func (m *PipeMap) canMakeSeeker(x, y int) bool {
	if x < 0 || y < 0 || x >= m.width || y >= m.height {
		return false
	}
	if m.value(x, y) == PIPE_GROUND {
		return false
	}
	return true
}

type indexForFn func(x, y int) int
func (s *Seeker) makeNorth(indexer indexForFn) *Seeker {
	return &Seeker{s.id, s.x, s.y-1, FROM_SOUTH, s.distance+1, append(s.indices, indexer(s.x, s.y-1))}
}
func (s *Seeker) makeEast(indexer indexForFn) *Seeker {
	return &Seeker{s.id, s.x+1, s.y, FROM_WEST, s.distance+1, append(s.indices, indexer(s.x+1, s.y))}
}
func (s *Seeker) makeSouth(indexer indexForFn) *Seeker {
	return &Seeker{s.id, s.x, s.y+1, FROM_NORTH, s.distance+1, append(s.indices, indexer(s.x, s.y+1))}
}
func (s *Seeker) makeWest(indexer indexForFn) *Seeker {
	return &Seeker{s.id, s.x-1, s.y, FROM_EAST, s.distance+1, append(s.indices, indexer(s.x-1, s.y))}
}

func (pipeMap *PipeMap) generateLoopMap() *PipeMap {
	seekers := pipeMap.findLoopingSeekers()

	cells := make([]int, pipeMap.width * pipeMap.height)
	for _, s := range seekers {
		for _, i := range s.indices {
			cells[i] = pipeMap.cells[i]
		}
	}
	return &PipeMap{cells, pipeMap.width, pipeMap.height, pipeMap.startX, pipeMap.startY}
}

func (pipeMap *PipeMap) findLoopingSeekers() []*Seeker {
	startIndex := pipeMap.indexFor(pipeMap.startX, pipeMap.startY)

	var seekers []*Seeker
	if pipeMap.canMakeSeeker(pipeMap.startX-1, pipeMap.startY) {
		index := pipeMap.indexFor(pipeMap.startX-1, pipeMap.startY)
		seekers = append(seekers, &Seeker{1, pipeMap.startX-1, pipeMap.startY, FROM_EAST, 1, []int{startIndex, index}})
	}
	if pipeMap.canMakeSeeker(pipeMap.startX, pipeMap.startY-1) {
		index := pipeMap.indexFor(pipeMap.startX, pipeMap.startY-1)
		seekers = append(seekers, &Seeker{2, pipeMap.startX, pipeMap.startY-1, FROM_SOUTH, 1, []int{startIndex, index}})
	}
	if pipeMap.canMakeSeeker(pipeMap.startX+1, pipeMap.startY) {
		index := pipeMap.indexFor(pipeMap.startX+1, pipeMap.startY)
		seekers = append(seekers, &Seeker{3, pipeMap.startX+1, pipeMap.startY, FROM_WEST, 1, []int{startIndex, index}})
	}
	if pipeMap.canMakeSeeker(pipeMap.startX, pipeMap.startY+1) {
		index := pipeMap.indexFor(pipeMap.startX, pipeMap.startY+1)
		seekers = append(seekers, &Seeker{4, pipeMap.startX, pipeMap.startY+1, FROM_NORTH, 1, []int{startIndex, index}})
	}

	for len(seekers) > 0 {
		var next []*Seeker

		for _, s := range seekers {
			t := pipeMap.value(s.x, s.y)

			var newSeeker *Seeker
			if t == PIPE_NS {
				if s.arrivedFrom == FROM_NORTH {
					newSeeker = s.makeSouth(pipeMap.indexFor)
				} else if s.arrivedFrom == FROM_SOUTH {
					newSeeker = s.makeNorth(pipeMap.indexFor)
				}
			} else if t == PIPE_EW {
				if s.arrivedFrom == FROM_EAST {
					newSeeker = s.makeWest(pipeMap.indexFor)
				} else if s.arrivedFrom == FROM_WEST {
					newSeeker = s.makeEast(pipeMap.indexFor)
				}
			} else if t == PIPE_NE {
				if s.arrivedFrom == FROM_NORTH {
					newSeeker = s.makeEast(pipeMap.indexFor)
				} else if s.arrivedFrom == FROM_EAST {
					newSeeker = s.makeNorth(pipeMap.indexFor)
				}
			} else if t == PIPE_NW {
				if s.arrivedFrom == FROM_NORTH {
					newSeeker = s.makeWest(pipeMap.indexFor)
				} else if s.arrivedFrom == FROM_WEST {
					newSeeker = s.makeNorth(pipeMap.indexFor)
				}
			} else if t == PIPE_SW {
				if s.arrivedFrom == FROM_SOUTH {
					newSeeker = s.makeWest(pipeMap.indexFor)
				} else if s.arrivedFrom == FROM_WEST {
					newSeeker = s.makeSouth(pipeMap.indexFor)
				}
			} else if t == PIPE_SE {
				if s.arrivedFrom == FROM_SOUTH {
					newSeeker = s.makeEast(pipeMap.indexFor)
				} else if s.arrivedFrom == FROM_EAST {
					newSeeker = s.makeSouth(pipeMap.indexFor)
				}
			}
			if newSeeker == nil {
				// probably hit a wall / invalid pipe
				continue
			}

			if !pipeMap.canMakeSeeker(newSeeker.x, newSeeker.y) {
				continue
			}

			for _, s2 := range next {
				if newSeeker.x == s2.x && newSeeker.y == s2.y {
					return []*Seeker{newSeeker, s2}
				}
			}
			next = append(next, newSeeker)
		}
		seekers = next
	}
	return nil
}

func (pipeMap *PipeMap) computeInside() *PipeMap {
	filledMap := &PipeMap{
		make([]int, pipeMap.width * pipeMap.height),
		pipeMap.width,
		pipeMap.height,
		pipeMap.startX,
		pipeMap.startY,
	}

	for y := 0; y < pipeMap.height; y++ {
		isInside := false
		for x := 0; x < pipeMap.width; x++ {
			t := pipeMap.value(x, y)
			if t == PIPE_NE || t == PIPE_NW || t == PIPE_NS || t == PIPE_START {
				isInside = !isInside
			}
			if !isInside {
				filledMap.set(x, y, CONNECT_OUTSIDE)
			}
		}
	}
	return filledMap
}

func (pipeMap *PipeMap) markOccupiedFrom(other *PipeMap) {
	for index, cell := range other.cells {
		if cell != PIPE_GROUND {
			pipeMap.cells[index] = cell
		}
	}
}

func (pipeMap *PipeMap) countGround() int {
	var count int
	for _, cell := range pipeMap.cells {
		if cell == PIPE_GROUND {
			count++
		}
	}
	return count
}

func (pipeMap *PipeMap) printOnlyGround() {
	filledMap := &PipeMap{
		make([]int, pipeMap.width * pipeMap.height),
		pipeMap.width,
		pipeMap.height,
		pipeMap.startX,
		pipeMap.startY,
	}

	for index, cell := range pipeMap.cells {
		if cell == PIPE_GROUND {
			filledMap.cells[index] = CONNECT_OUTSIDE
		}
	}
	filledMap.print()
}

func parseMap(contents string) *PipeMap {
	lines := strings.Split(strings.TrimSpace(contents), "\n")
	if len(lines) == 0 || len(lines[0]) == 0 {
		return &PipeMap{}
	}

	height := len(lines)
	width := len(strings.TrimSpace(lines[0]))

	var startX, startY int
	cells := make([]int, 0, width * height)
	for y, l := range lines {
		for x, c := range strings.TrimSpace(l) {
			ctype := charToPipe[c]
			if ctype == PIPE_START {
				startX = x
				startY = y
			}
			cells = append(cells, ctype)
		}
	}

	return &PipeMap{cells, width, height, startX, startY}
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

	pipeMap := parseMap(string(contents))

	loopMap := pipeMap.generateLoopMap()
	loopMap.print()
	insideMap := loopMap.computeInside()
	insideMap.markOccupiedFrom(loopMap)
	insideMap.print()
	insideMap.printOnlyGround()

	log.Printf("Number contained: %d", insideMap.countGround())
}
