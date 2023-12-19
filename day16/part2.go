package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

const (
	CELL_EMPTY = iota
	CELL_SLASH
	CELL_BACKSLASH
	CELL_SPLITVERT
	CELL_SPLITHORIZ
)

const (
	DIRECTION_RIGHT = iota
	DIRECTION_UP
	DIRECTION_LEFT
	DIRECTION_DOWN
)

var celltypes = map[rune]int {
	'.': CELL_EMPTY,
	'/': CELL_SLASH,
	'\\': CELL_BACKSLASH,
	'|': CELL_SPLITVERT,
	'-': CELL_SPLITHORIZ,
}

type Cell struct {
	kind int
	isEnergized bool
}

type Room struct {
	grid []Cell
	width int
	height int
}

type Laser struct {
	x int
	y int
	direction int
}

// mapping kind > input direction > output directions
var laserMovement = map[int]map[int][]int {
	CELL_EMPTY: map[int][]int {
		DIRECTION_RIGHT: []int{ DIRECTION_RIGHT },
		DIRECTION_UP: []int{ DIRECTION_UP },
		DIRECTION_LEFT: []int{ DIRECTION_LEFT },
		DIRECTION_DOWN: []int{ DIRECTION_DOWN },
	},
	CELL_SLASH: map[int][]int {
		DIRECTION_RIGHT: []int{ DIRECTION_UP },
		DIRECTION_UP: []int{ DIRECTION_RIGHT },
		DIRECTION_LEFT: []int{ DIRECTION_DOWN },
		DIRECTION_DOWN: []int{ DIRECTION_LEFT },
	},
	CELL_BACKSLASH: map[int][]int {
		DIRECTION_RIGHT: []int{ DIRECTION_DOWN },
		DIRECTION_UP: []int{ DIRECTION_LEFT },
		DIRECTION_LEFT: []int{ DIRECTION_UP },
		DIRECTION_DOWN: []int{ DIRECTION_RIGHT },
	},
	CELL_SPLITVERT: map[int][]int {
		DIRECTION_RIGHT: []int{ DIRECTION_UP, DIRECTION_DOWN },
		DIRECTION_UP: []int{ DIRECTION_UP },
		DIRECTION_LEFT: []int{ DIRECTION_UP, DIRECTION_DOWN },
		DIRECTION_DOWN: []int{ DIRECTION_DOWN },
	},
	CELL_SPLITHORIZ: map[int][]int {
		DIRECTION_RIGHT: []int{ DIRECTION_RIGHT },
		DIRECTION_UP: []int{ DIRECTION_LEFT, DIRECTION_RIGHT },
		DIRECTION_LEFT: []int{ DIRECTION_LEFT },
		DIRECTION_DOWN: []int{ DIRECTION_LEFT, DIRECTION_RIGHT },
	},
}

func (r *Room) energizeCell(x, y int) {
	r.grid[y*r.width+x].isEnergized = true
}

func (r *Room) shootLaser(startX, startY int, startDirection int) {
	var lasers []*Laser

	seen := map[int]map[int]bool{}
	for i := range r.grid {
		seen[i] = map[int]bool{}
	}

	nextDirections := laserMovement[r.grid[startY*r.width+startX].kind][startDirection]
	for _, d := range nextDirections {
		lasers = append(lasers, &Laser{ x: startX, y: startY, direction: d })
	}

	step := 0
	for len(lasers) > 0 {
		step++
		var next []*Laser
		for _, l := range lasers {
			r.energizeCell(l.x, l.y)

			nextX := l.x
			nextY := l.y
			if l.direction == DIRECTION_RIGHT {
				nextX = l.x + 1
			} else if l.direction == DIRECTION_LEFT {
				nextX = l.x - 1
			} else if l.direction == DIRECTION_DOWN {
				nextY = l.y + 1
			} else if l.direction == DIRECTION_UP {
				nextY = l.y - 1
			}
			gridcoords := nextY*r.width + nextX
			if nextX < 0 || nextY < 0 || nextX >= r.width || nextY >= r.height || seen[gridcoords][l.direction] {
				continue
			}
			seen[gridcoords][l.direction] = true
			nextDirections := laserMovement[r.grid[gridcoords].kind][l.direction]
			for _, d := range nextDirections {
				next = append(next, &Laser{ x: nextX, y: nextY, direction: d })
			}
		}
		lasers = next
	}
}

func (r *Room) energizedCount() int {
	total := 0
	for _, c := range r.grid {
		if c.isEnergized {
			total++
		}
	}
	return total
}

func (r *Room) reset() {
	for i := range r.grid {
		r.grid[i].isEnergized = false
	}
}

func (r *Room) print() {
	i := 0
	for y := 0; y < r.height; y++ {
		var row []rune
		for x := 0; x < r.width; x++ {
			if r.grid[i].isEnergized {
				row = append(row, '#')
			} else {
				row = append(row, '.')
			}
			i++
		}
		log.Printf("%s", string(row))
	}
}

func parseRoom(contents string) *Room {
	lines := strings.Split(contents, "\n")
	height := len(lines)
	width := len(lines[0])

	grid := make([]Cell, width*height)

	i := 0
	for _, line := range lines {
		for _, r := range line {
			kind := celltypes[r]
			grid[i].kind = kind
			grid[i].isEnergized = false
			i++
		}
	}
	return &Room{grid: grid, width: width, height: height}
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

	room := parseRoom(str)
	var maxTotal int

	// do left to right along top side pointing down
	for x := 0; x < room.width; x++ {
		room.shootLaser(x, 0, DIRECTION_DOWN)
		total := room.energizedCount()
		if total > maxTotal {
			maxTotal = total
		}
		room.reset()
	}
	// do left to right along bottom side pointing up
	for x := 0; x < room.width; x++ {
		room.shootLaser(x, room.height-1, DIRECTION_UP)
		total := room.energizedCount()
		if total > maxTotal {
			maxTotal = total
		}
		room.reset()
	}
	// do top to bottom along left side pointing right
	for y := 0; y < room.height; y++ {
		room.shootLaser(0, y, DIRECTION_RIGHT)
		total := room.energizedCount()
		if total > maxTotal {
			maxTotal = total
		}
		room.reset()
	}
	// do top to bottom along right side pointing left
	for y := 0; y < room.height; y++ {
		room.shootLaser(room.width-1, y, DIRECTION_LEFT)
		total := room.energizedCount()
		if total > maxTotal {
			maxTotal = total
		}
		room.reset()
	}

	log.Printf("Sum: %d", maxTotal)
}
