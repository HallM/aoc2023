package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

const (
	PIPE_GROUND = iota
	PIPE_START
	PIPE_NS
	PIPE_EW
	PIPE_NE
	PIPE_NW
	PIPE_SW
	PIPE_SE
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

type Seeker struct {
	id int
	x int
	y int
	arrivedFrom int
	distance int
}

type PipeMap struct {
	cells []int
	width int
	height int
	startX int
	startY int
}

func (m *PipeMap) indexFor(x, y int) int {
	return y * m.width + x
}

func (m *PipeMap) value(x, y int) int {
	return m.cells[m.indexFor(x, y)]
}

func findMaxDistanceLoop(pipeMap *PipeMap) int {
	var seekers []*Seeker
	if pipeMap.canMakeSeeker(pipeMap.startX-1, pipeMap.startY) {
		seekers = append(seekers, &Seeker{1, pipeMap.startX-1, pipeMap.startY, FROM_EAST, 1})
	}
	if pipeMap.canMakeSeeker(pipeMap.startX, pipeMap.startY-1) {
		seekers = append(seekers, &Seeker{2, pipeMap.startX, pipeMap.startY-1, FROM_SOUTH, 1})
	}
	if pipeMap.canMakeSeeker(pipeMap.startX+1, pipeMap.startY) {
		seekers = append(seekers, &Seeker{3, pipeMap.startX+1, pipeMap.startY, FROM_WEST, 1})
	}
	if pipeMap.canMakeSeeker(pipeMap.startX, pipeMap.startY+1) {
		seekers = append(seekers, &Seeker{4, pipeMap.startX, pipeMap.startY+1, FROM_NORTH, 1})
	}

	for len(seekers) > 0 {
		var next []*Seeker

		for _, s := range seekers {
			t := pipeMap.value(s.x, s.y)

			var newSeeker *Seeker
			if t == PIPE_NS {
				if s.arrivedFrom == FROM_NORTH {
					newSeeker = s.makeSouth()
				} else if s.arrivedFrom == FROM_SOUTH {
					newSeeker = s.makeNorth()
				}
			} else if t == PIPE_EW {
				if s.arrivedFrom == FROM_EAST {
					newSeeker = s.makeWest()
				} else if s.arrivedFrom == FROM_WEST {
					newSeeker = s.makeEast()
				}
			} else if t == PIPE_NE {
				if s.arrivedFrom == FROM_NORTH {
					newSeeker = s.makeEast()
				} else if s.arrivedFrom == FROM_EAST {
					newSeeker = s.makeNorth()
				}
			} else if t == PIPE_NW {
				if s.arrivedFrom == FROM_NORTH {
					newSeeker = s.makeWest()
				} else if s.arrivedFrom == FROM_WEST {
					newSeeker = s.makeNorth()
				}
			} else if t == PIPE_SW {
				if s.arrivedFrom == FROM_SOUTH {
					newSeeker = s.makeWest()
				} else if s.arrivedFrom == FROM_WEST {
					newSeeker = s.makeSouth()
				}
			} else if t == PIPE_SE {
				if s.arrivedFrom == FROM_SOUTH {
					newSeeker = s.makeEast()
				} else if s.arrivedFrom == FROM_EAST {
					newSeeker = s.makeSouth()
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
					return newSeeker.distance
				}
			}
			next = append(next, newSeeker)
		}
		seekers = next
	}
	return 0
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

func (s *Seeker) makeNorth() *Seeker {
	return &Seeker{s.id, s.x, s.y-1, FROM_SOUTH, s.distance+1}
}
func (s *Seeker) makeEast() *Seeker {
	return &Seeker{s.id, s.x+1, s.y, FROM_WEST, s.distance+1}
}
func (s *Seeker) makeSouth() *Seeker {
	return &Seeker{s.id, s.x, s.y+1, FROM_NORTH, s.distance+1}
}
func (s *Seeker) makeWest() *Seeker {
	return &Seeker{s.id, s.x-1, s.y, FROM_EAST, s.distance+1}
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

	distance := findMaxDistanceLoop(pipeMap)

	log.Printf("Max Distance: %d", distance)
}
