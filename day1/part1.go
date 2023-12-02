package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var docFilePath = flag.String("doc", "part1test.txt", "File path to the calibration doc")

var runeDigitMap = map[rune]int {
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

func computeCalibration(doc string) int {
	lines := strings.Split(doc, "\n")

	var total int
	for _, line := range lines {
		total += computeCalibrationLine(line)
	}
	return total
}

func computeCalibrationLine(line string) int {
	var first, last int
	var hasFirst bool

	for _, c := range line {
		if d, ok := runeDigitMap[c]; ok {
			if !hasFirst {
				hasFirst = true
				first = d
			}
			last = d
		}
	}
	return (10 * first) + last
}

func main() {
	flag.Parse()

	if *docFilePath == "" {
		log.Fatalf("Must specify the calibration doc!")
	}

	doc, err := os.ReadFile(*docFilePath)
	if err != nil {
		log.Fatal(err)
	}
	calibration := computeCalibration(string(doc))
	log.Printf("Calibration: %d", calibration)
}
