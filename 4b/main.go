package main

import (
	"fmt"
	"strconv"
)

func main() {
	pws := []int{}

	for x := 136760; x < 595730; x++ {
		pw := strconv.Itoa(x)

		// check doubles
		hasDoubles := false
		var last, numRepetitions int
		for j := 0; j <= 9; j++ {
			numRepetitions = 0
			for _, c := range pw {
				cur := int(c - '0')
				if cur == j {
					numRepetitions++
				}
				last = cur
			}
			if numRepetitions == 2 {
				hasDoubles = true
				break
			}
		}

		// check increasing numbers
		hasIncreasingOnly := true
		last = 1
		for _, c := range pw {
			cur := int(c)
			if cur < last {
				hasIncreasingOnly = false
				break
			}
			last = cur
		}

		if hasDoubles && hasIncreasingOnly {
			pws = append(pws, x)
		}

	}

	// fmt.Printf("%v\n", pws)
	fmt.Printf("%v\n", len(pws))
}
