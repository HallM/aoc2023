package main

import (
	"flag"
	"log"
	"os"
)

var filePath = flag.String("file", "part1test.txt", "File path")

func computeHash(s []byte) int {
	var hash int
	for _, c := range s {
		hash += int(c)
		hash = hash * 17
		hash = (hash & 0xFF)
	}
	return hash
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

	total := 0
	start := 0
	for i, c := range contents {
		if c == ',' {
			s := contents[start:i]
			hash := computeHash(s)
			log.Printf("%s hash is %d", string(s), hash)
			start = i+1
			total += hash
		}
	}
	s := contents[start:]
	hash := computeHash(s)
	log.Printf("%s hash is %d", string(s), hash);
	total += hash

	log.Printf("Sum: %d", total)
}
