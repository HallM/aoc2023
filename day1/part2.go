package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

var docFilePath = flag.String("doc", "part2test.txt", "File path to the calibration doc")

var matchDigitMap = map[string]int {
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,

	"zero": 0,
	"one": 1,
	"two": 2,
	"three": 3,
	"four": 4,
	"five": 5,
	"six": 6,
	"seven": 7,
	"eight": 8,
	"nine": 9,
}

func computeCalibration(doc string) int {
	lines := strings.Split(doc, "\n")

	var total int
	for _, line := range lines {
		total += computeCalibrationLine(strings.TrimSpace(line))
	}
	return total
}

func computeCalibrationLine(line string) int {
	var first, last int
	var hasFirst bool

	re := regexp.MustCompile("(one|two|three|four|five|six|seven|eight|nine|[0-9])")

	// FindAllStrings will not find overlapping like oneight.
	for i := 0; i < len(line); {
		m := re.FindStringIndex(line[i:])
		if m == nil {
			break
		}

		match := line[i+m[0]:i+m[1]]
		if hasFirst == false {
			hasFirst = true
			first = matchDigitMap[match]
		}
		last = matchDigitMap[match]
		i += m[0] + 1
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
