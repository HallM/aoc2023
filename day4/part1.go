package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

var filePath = flag.String("file", "part1test.txt", "File path to the scratch offs file")

func computeSum(contents string) (int, error) {
	lines := strings.Split(contents, "\n")

	var total int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		score, err := computeScratchoffScore(line)
		if err != nil {
			return 0, err
		}
		total += score
	}
	return total, nil
}

func computeScratchoffScore(contents string) (int, error) {
	// Skip over the "Card #:"
	segments := strings.Split(contents[strings.Index(contents, ":")+1:], "|")
	if len(segments) != 2 {
		return 0, fmt.Errorf("Unknown format for %q", contents)
	}

	winningNumbers := map[int]bool{}
	for _, s := range strings.Split(strings.TrimSpace(segments[0]), " ") {
		// Note that splitting like "52  1" will create 52,"",1
		if len(s) == 0 || s == " " {
			continue
		}

		val, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("Cannot parse number from %q as int: %w", s, err)
		}

		winningNumbers[int(val)] = true
	}
	log.Printf("All winning numbers %v", winningNumbers)

	score := 0
	for _, s := range strings.Split(strings.TrimSpace(segments[1]), " ") {
		if len(s) == 0 || s == " " {
			continue
		}

		val, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("Cannot parse number from %q as int: %w", s, err)
		}

		if winningNumbers[int(val)] {
			log.Printf("Have matching number %d", val)
			if score == 0 {
				score = 1
			} else {
				score *= 2
			}
		}
	}

	return score, nil
}

func main() {
	flag.Parse()

	if *filePath == "" {
		log.Fatalf("Must specify the scratch off file!")
	}

	contents, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	score, err := computeSum(string(contents))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total: %d", score)
}
