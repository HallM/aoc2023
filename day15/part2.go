package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"os"
)

var filePath = flag.String("file", "part1test.txt", "File path")

type Hashmap struct {
	boxes []*List
}

type Instruction struct {
	label string
	hash int
	removeLens bool
	setFocal int // 0 if not setting
}

type LabeledLens struct {
	label string
	focal int
}

type Node struct {
	data *LabeledLens
	next *Node
}

type List struct {
	first *Node
}

func makeHashmap() *Hashmap {
	boxes := make([]*List, 256)
	for i := range boxes {
		boxes[i] = &List{}
	}
	return &Hashmap{
		boxes: boxes,
	}
}

func (h *Hashmap) print() {
	for i, box := range h.boxes {
		if box.first != nil {
			var contents []string
			node := box.first
			for node != nil {
				contents = append(contents, fmt.Sprintf("[%s %d]", node.data.label, node.data.focal))
				node = node.next
			}
			log.Printf("Box %d: %s", i, strings.Join(contents, " "))
		}
	}
}

func (h *Hashmap) runInstruction(s []byte) {
	instr := computeInstruction(s)
	log.Printf("After %s = %s hash is %d - %v / %v", string(s), instr.label, instr.hash, instr.removeLens, instr.setFocal)

	box := h.boxes[instr.hash]

	if instr.removeLens {
		// log.Printf("Removing a lens")
		var lastNode *Node
		node := box.first
		for node != nil {
			if node.data.label == instr.label {
				if lastNode == nil {
					// log.Printf("Removing the only lens in the box")
					box.first = node.next
				} else {
					// log.Printf("Removing one lens and shifting up")
					lastNode.next = node.next
				}
				break
			}
			lastNode = node
			node = node.next
		}
	} else {
		// log.Printf("Set a lens")
		if box.first == nil {
			box.first = &Node{data: &LabeledLens{label: instr.label, focal: instr.setFocal}}
		} else {
			node := box.first
			for node != nil {
				if node.data.label == instr.label {
					// log.Printf("Replace a lens (previously %d)", node.data.focal)
					node.data.focal = instr.setFocal
					break
				}
				if node.next == nil {
					// log.Printf("Adding a new lens to the box")
					node.next = &Node{data: &LabeledLens{label: instr.label, focal: instr.setFocal}}
					break
				}
				node = node.next
			}
		}
	}

	// h.print()
}

func (h *Hashmap) computeFocusPower() int {
	power := 0
	for boxNum, box := range h.boxes {
		if box.first != nil {
			slot := 1
			node := box.first
			for node != nil {
				p := (boxNum + 1) * slot * node.data.focal
				log.Printf("%s: %d * %d * %d = %d", node.data.label, boxNum + 1, slot, node.data.focal, p)
				power += p
				slot++
				node = node.next
			}
		}
	}
	return power
}

func computeInstruction(s []byte) *Instruction {
	var hash int
	for i, c := range s {
		if c == 45 {
			// "-"
			return &Instruction{label: string(s[0:i]), hash: hash, removeLens: true, setFocal: 0}
		} else if c == 61 {
			// "="
			focal := int(s[i+1]) - 48
			return &Instruction{label: string(s[0:i]), hash: hash, removeLens: false, setFocal: focal}
		}

		hash += int(c)
		hash = hash * 17
		hash = (hash & 0xFF)
	}
	return nil
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

	hashmap := makeHashmap()

	start := 0
	for i, c := range contents {
		if c == ',' {
			hashmap.runInstruction(contents[start:i])
			start = i+1
		}
	}
	hashmap.runInstruction(contents[start:])

	total := hashmap.computeFocusPower()
	log.Printf("Sum: %d", total)
}
