package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var basePattern = []int{-1, 0, 1, 0} // in reverse order for easy append

var inputReps = 10000
var phaseReps = 100

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

	msgOffsetInts := ints[:7]
	msgOffset := ""
	for _, i := range msgOffsetInts {
		msgOffset += strconv.Itoa(i)
	}
	offset, err := strconv.Atoi(msgOffset)
	if err != nil {
		panic(err)
	}
	fmt.Println("Getting msg at offset", msgOffset)

	allInts := []int{}
	for i := 1; i <= inputReps; i++ {
		allInts = append(allInts, ints...)
	}

	input := allInts[offset:]

	fmt.Println(len(input))

	fmt.Println("Starting phase calculations")
	for i := 1; i <= phaseReps; i++ {
		sum := 0
		for i := len(input) - 1; i >= 0; i-- {
			sum += input[i]
			input[i] = sum % 10
		}
	}
	for _, d := range input[:8] {
		fmt.Print(d)
	}
	fmt.Println()
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
