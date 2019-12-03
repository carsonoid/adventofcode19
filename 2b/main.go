package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func operate(data, workSet []uint64) error {
	opCode := workSet[0]
	inPos1 := workSet[1]
	inPos2 := workSet[2]
	outPos := workSet[3]
	switch opCode {
	case 1:
		data[outPos] = data[inPos1] + data[inPos2]
	case 2:
		data[outPos] = data[inPos1] * data[inPos2]
	default:
		return fmt.Errorf("invalid opCode: %d", opCode)
	}
	return nil
}

func getData(filePath string) []uint64 {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var data []uint64
	for {
		dataPoint, err := reader.ReadString(',')
		dataPoint = strings.Trim(dataPoint, ",")

		u, err := strconv.ParseUint(dataPoint, 10, 64)
		if err != nil {
			break
		}
		data = append(data, u)
	}

	return data
}

func test(data []uint64, noun, verb uint64) {
	// fmt.Printf("Start: %v\n", data)

	// restore gravity assist
	data[1] = noun
	data[2] = verb

	var startPos uint64
	for {
		workSet := data[startPos : startPos+4]
		if workSet[0] == 99 {
			break
		}
		if err := operate(data, workSet); err != nil {
			panic(err)
		}
		startPos = startPos + 4
	}

	// fmt.Printf("End: %v\n", data)
}

func main() {
	memory := getData(os.Args[1])

	var noun, verb uint64
	for noun = 0; noun < 99; noun++ {
		for verb = 0; verb < 99; verb++ {
			var data []uint64
			data = make([]uint64, len(memory))
			copy(data, memory)
			test(data, noun, verb)

			if data[0] == 19690720 {
				fmt.Printf("MATCH! noun: %d, verb: %d\n", noun, verb)
				fmt.Printf("Answer: %d\n", 100*noun+verb)
			}
		}
	}
}
