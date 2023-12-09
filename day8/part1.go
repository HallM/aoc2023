package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var filePath = flag.String("file", "part1test.txt", "File path")

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

func traverse(startNode, targetNode string, nodeMap map[string]Node, path []int) int {
	var steps int
	node := startNode
	for {
		nextPath := steps % len(path)
		next := nodeMap[node][path[nextPath]]
		wentdir := "left"
		if path[nextPath] == 1 {
			wentdir = "right"
		}
		log.Printf("step %d went from %s -> %s along %s(%d)", steps, node, next, wentdir, nextPath)
		node = next

		steps++
		if node == targetNode {
			return steps
		}

		if steps > 1000000 {
			return -1
		}
	}
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

	steps := traverse("AAA", "ZZZ", nodeMap, path)

	log.Printf("Made it in steps: %d", steps)
}
