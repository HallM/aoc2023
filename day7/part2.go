package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"strconv"
)

var filePath = flag.String("file", "part1test.txt", "File path")

var cardValueMap = map[rune]int {
	'J': 1,

	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 11,
	'K': 12,
	'A': 13,
}

const (
	FIVE_KIND_HAND = 7
	FOUR_KIND_HAND = 6
	FULL_HOUSE_HAND = 5
	THREE_KIND_HAND = 4
	TWO_PAIR_HAND = 3
	ONE_PAIR_HAND = 2
	HIGH_CARD_HAND = 1
)

type PokerHand struct {
	hand string
	cards []int
	handType int
	bid int64
}

func computeWinnings(hands []*PokerHand) int64 {
	sort.Slice(hands, func(a, b int) bool {
		if hands[a].handType != hands[b].handType {
			return hands[a].handType < hands[b].handType
		}
		for i, x := range hands[a].cards {
			y := hands[b].cards[i]
			if x != y {
				return x < y
			}
		}
		return false
	})

	var winnings int64
	for m, hand := range hands {
		handWinning := (int64(m)+1) * hand.bid
		log.Printf("%d place is hand %q with bid %d won %d", m+1, hand.hand, hand.bid, handWinning)
		winnings += handWinning
	}
	return winnings
}

func parseHands(contents string) ([]*PokerHand, error) {
	var hands []*PokerHand
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		hand, err := parseHand(line)
		if err != nil {
			return nil, err
		}
		hands = append(hands, hand)
	}
	return hands, nil
}

func parseHand(line string) (*PokerHand, error) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Expected 2 numbers separated by a space, but got %q", line)
	}

	var cards []int
	for _, c := range parts[0] {
		cards = append(cards, cardValueMap[c])
	}

	bid, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse number from %q as int: %w", parts[1], err)
	}

	return &PokerHand{parts[0], cards, handType(cards), bid}, nil
}

type cardCounts struct {
	card int
	count int
}

func handType(cards []int) int {
	counts, jokers := countOfCards(cards)

	if len(counts) <= 1 {
		return FIVE_KIND_HAND
	}
	if len(counts) == 2 {
		if counts[0].count + jokers.count == 4 {
			return FOUR_KIND_HAND
		}
		// else then we must have counts[0].count == 3 && counts[1].count == 2
		return FULL_HOUSE_HAND
	}
	if len(counts) == 3 {
		if counts[0].count + jokers.count == 3 {
			return THREE_KIND_HAND
		}
		// else then we have counts[0].count == 2 && counts[1].count == 2
		return TWO_PAIR_HAND
	}
	if len(counts) == 4 {
		return ONE_PAIR_HAND
	}
	return HIGH_CARD_HAND
}

// Returns the counts of all non-jokers and then the jokers separately
func countOfCards(cards []int) ([]cardCounts, cardCounts) {
	counts := map[int]int{}
	for _, c := range cards {
		counts[c]++
	}

	var c []cardCounts
	for k, v := range counts {
		// jokers are returned separately
		if k != 1 {
			c = append(c, cardCounts{k, v})
		}
	}
	sort.Slice(c, func(a, b int) bool { return c[a].count > c[b].count })
	return c, cardCounts{1, counts[1]}
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

	hands, err := parseHands(string(contents))
	if err != nil {
		log.Fatal(err)
	}
	winnings := computeWinnings(hands)
	log.Printf("All winnings: %d", winnings)
}
