package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part2test.txt", "File path")

// it's a pair for left/right destinations, but I wanted to minimize branches for no reason.
type Node []string

func parseTravelPath(contents string) []int {
	var path []int
	for _, c := range strings.TrimSpace(contents) {
		if c == 'L' {
			path = append(path, 0)
		} else {
			path = append(path, 1)
		}
	}
	return path
}

func parseMap(lines []string) map[string]Node {
	m := map[string]Node{}
	for _, line := range lines {
		node := line[0:3]
		left := line[7:10]
		right := line[12:15]
		m[node] = []string{left, right}
	}
	return m
}

func traverse(startNode string, nodeMap map[string]Node, path []int) int {
	var steps int
	node := startNode
	for {
		nextPath := steps % len(path)
		next := nodeMap[node][path[nextPath]]
		node = next

		steps++
		if node[2] == 'Z' {
			return steps
		}

		if steps > 1000000 {
			return -1
		}
	}
}

func traverseAll(nodeMap map[string]Node, path []int) int64 {
	var s []int64
	for n := range nodeMap {
		if n[2] == 'A' {
			steps := traverse(n, nodeMap, path)
			log.Printf("%s made it in %d", n, steps)
			s = append(s, int64(steps))
		}
	}
	return lcmAll(s)
}

func lcmAll(steps []int64) int64 {
	a := steps[0]
	for _, b := range steps[1:] {
		next := lcmTwo(a, b)
		a = next
	}
	return a
}

func lcmTwo(a, b int64) int64 {
	return (a * b) / gcd(a, b)
}

func gcd(a, b int64) int64 {
	var rem int64
	for {
		rem = a % b
		a = b
		b = rem
		if b == 0 {
			break
		}
	}
	return a
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

	lines := strings.Split(string(contents), "\n")

	path := parseTravelPath(lines[0])
	nodeMap := parseMap(lines[2:])

	common := traverseAll(nodeMap, path)

	log.Printf("Made it in steps: %d", common)
}
