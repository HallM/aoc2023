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

	scratch := &Scratchoffs{copies: map[int64]int{}}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		err := scratch.computeScratchoffCopies(line)
		if err != nil {
			return 0, err
		}
	}

	total := 0
	for ID, copies := range scratch.copies {
		// Do not count copies of cards we do not have
		if ID > scratch.maxID {
			break
		}
		log.Printf("Card %d has %d copies", ID, copies)
		total += copies
	}
	return total, nil
}

type Scratchoffs struct {
	copies map[int64]int
	maxID int64
}

func (scratch *Scratchoffs) computeScratchoffCopies(contents string) error {
	index := strings.Index(contents, ":")
	cardID, err := strconv.ParseInt(strings.TrimSpace(contents[5:index]), 10, 32)
	if err != nil {
		return fmt.Errorf("Cannot parse number from %q as int: %w", strings.TrimSpace(contents[5:index]), err)
	}

	scratch.maxID = cardID
	scratch.copies[cardID]++

	segments := strings.Split(contents[index+1:], "|")
	if len(segments) != 2 {
		return fmt.Errorf("Unknown format for %q", contents)
	}

	winningNumbers := map[int]bool{}
	for _, s := range strings.Split(strings.TrimSpace(segments[0]), " ") {
		// Note that splitting like "52  1" will create 52,"",1
		if len(s) == 0 || s == " " {
			continue
		}

		val, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("Cannot parse number from %q as int: %w", s, err)
		}

		winningNumbers[int(val)] = true
	}
	log.Printf("All winning numbers %v", winningNumbers)

	next := int64(1)
	for _, s := range strings.Split(strings.TrimSpace(segments[1]), " ") {
		if len(s) == 0 || s == " " {
			continue
		}

		val, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("Cannot parse number from %q as int: %w", s, err)
		}

		if winningNumbers[int(val)] {
			log.Printf("Have matching number %d, adding a copy of %d", val, cardID + next)
			scratch.copies[cardID + next] += scratch.copies[cardID]
			next++
		}
	}

	return nil
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
