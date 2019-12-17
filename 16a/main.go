package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var basePattern = []int{-1, 0, 1, 0} // in reverse order for easy append

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	ints := []int{}
	for {
		b, err := reader.ReadByte()
		if err != nil || b == byte('\n') {
			break
		}
		i, err := strconv.Atoi(string(b))
		if err != nil {
			panic(err)
		}
		ints = append(ints, i)
	}

	// fmt.Println(ints)
	for i := 1; i <= 100; i++ {
		ints = runPhase(ints)
	}
	fmt.Println(ints)
}

func getModList(pos int) []int {
	ret := []int{}
	for j := len(basePattern) - 1; j >= 0; j-- {
		for i := 0; i < pos; i++ { // repeat each char pos number of times
			ret = append(ret, basePattern[j])
		}
	}
	return ret
}

func runPhase(inputs []int) []int {
	result := []int{}
	for i := range inputs {
		modResults := []int{}
		i++ // index starts at one
		modList := getModList(i)
		// fmt.Println(modList)

		maxIndex := len(modList) - 1
		curIndex := 1 // skip first element on first loop
		for _, input := range inputs {
			// Multiply each pos by it's mod
			mod := modList[curIndex]
			// fmt.Println("operation", input, "*", mod)
			modResults = append(modResults, input*mod)
			// fmt.Println(modResults)

			curIndex++
			if curIndex > maxIndex {
				curIndex = 0
			}

		}
		out := 0
		for _, r := range modResults {
			out += r
		}
		outStr := strconv.Itoa(out)
		lastChar := string(outStr[len(outStr)-1])
		// fmt.Println(lastChar)
		lastInt, err := strconv.Atoi(lastChar)
		if err != nil {
			panic(err)
		}
		result = append(result, lastInt)
	}
	return result
}
