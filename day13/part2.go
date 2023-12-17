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
	IS_ASH = int64(iota)
	IS_ROCK
)

var cellTypeMap = map[rune]int64 {
	'.': IS_ASH,
	'#': IS_ROCK,
}

const (
	horizMultiplier = 100
	vertMultiplier = 1
)

type Block struct {
	lineNumber int

	// bitmasks where 0=ash, 1=rock
	// top/left is left-most bit, bottom/right is right-most bit
	cols []int64
	rows []int64
}

func (b *Block) print() {
	log.Printf("Block:")
	log.Printf("  Rows:")
	for y, r := range b.rows {
		log.Printf("    row %d: %s / %d", y+1, strconv.FormatInt(r, 2), r)
	}
	log.Printf("  Cols:")
	for x, c := range b.cols {
		log.Printf("    col %d: %s / %d", x+1, strconv.FormatInt(c, 2), c)
	}
}

func (b *Block) flipCell(x, y int) *Block {
	var cols []int64
	for i, c := range b.cols {
		if x == i {
			bitToFlip := int64(1) << (len(b.rows) - y - 1)
			cols = append(cols, c ^ bitToFlip)
		} else {
			cols = append(cols, c)
		}
	}
	var rows []int64
	for j, r := range b.rows {
		if y == j {
			bitToFlip := int64(1) << (len(b.cols) - x - 1)
			rows = append(rows, r ^ bitToFlip)
		} else {
			rows = append(rows, r)
		}
	}
	return &Block{lineNumber: b.lineNumber, cols: cols, rows: rows}
}

// returns 0 if no vertical reflection
// or number of columns to the left of the vertical reflection line
func (b *Block) findVertical(exclude int) int {
	return findReflection(b.cols, exclude)
}

func (b *Block) findHorizontal(exclude int) int {
	return findReflection(b.rows, exclude)
}

func (b *Block) findScore(excludeV int, excludeH int) (int, int) {
	return b.findVertical(excludeV), b.findHorizontal(excludeH)
}

func (b *Block) findSmidgeScore() int {
	excludeV, excludeH := b.findScore(-1, -1)

	for y := 0; y < len(b.rows); y++ {
		for x := 0; x < len(b.cols); x++ {
			smudged := b.flipCell(x, y)
			scoreV, scoreH := smudged.findScore(excludeV, excludeH)
			if scoreV > 0 {
				return vertMultiplier * scoreV
			}
			if scoreH > 0 {
				return horizMultiplier * scoreH
			}
		}
	}
	return 0
}

func findReflection(arr []int64, exclude int) int {
	for x, v := range arr[1:] {
		if x == exclude-1 {
			continue
		}
		if v == arr[x] {
			allMatches := true
			j := x+2
			for i := x-1; i >= 0 && j < len(arr); i-- {
				if arr[i] != arr[j] {
					allMatches = false
					break
				}
				j++
			}
			if allMatches {
				return x+1
			}
		}
	}
	return 0
}

func parseBlock(block string, lineNumber int) *Block {
	lines := strings.Split(block, "\n")
	height := len(lines)
	width := len(lines[0])

	rows := make([]int64, height)
	cols := make([]int64, width)

	for y := 0; y < height; y++ {
		rows[y] = 0
	}
	for x := 0; x < width; x++ {
		cols[x] = 0
	}

	for y, line := range lines {
		for x, c := range line {
			rows[y] = (rows[y] << 1) | cellTypeMap[c]
			cols[x] = (cols[x] << 1) | cellTypeMap[c]
		}
	}
	return &Block{lineNumber: lineNumber, rows: rows, cols: cols}
}

func parseInput(contents string) []*Block {
	blocks := strings.Split(contents, "\n\n")
	line := 1
	var ret []*Block
	for i, b := range blocks {
		log.Printf("parsing block %d at line %d", i+1, line)
		ret = append(ret, parseBlock(b, line))
		line += len(strings.Split(b, "\n")) + 1
	}
	return ret
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

	blocks := parseInput(str)
	var total int
	for i, b := range blocks {
		score := b.findSmidgeScore();
		// b.print()
		log.Printf("Block %d (line %d) => %d", (i+1), b.lineNumber, score)
		total += score
	}

	log.Printf("Sum: %d", total)
}
