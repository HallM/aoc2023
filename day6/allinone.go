package main

import (
	"log"
)

type Race struct {
	time int64
	distance int64
}

func main() {
	races := []*Race{
		// The sample from part 1
		// &Race{
		// 	time: 7,
		// 	distance: 9,
		// },
		// &Race{
		// 	time: 15,
		// 	distance: 40,
		// },
		// &Race{
		// 	time: 30,
		// 	distance: 200,
		// },
		// The sample from part 2
		&Race{
			time: 71530,
			distance: 940200,
		},
	}

	total := int64(1)
	for _, race := range races {
		var won, t int64

		for t = 0; t < race.time; t++ {
			remainingTime := race.time - t
			distance := remainingTime * t
			if distance > race.distance {
				won++
			}
		}
		total *= won
	}

	log.Printf("Number: %d", total)
}
