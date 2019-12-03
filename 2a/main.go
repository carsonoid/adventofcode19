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

func main() {
	data := getData(os.Args[1])
	fmt.Printf("Start: %v\n", data)

	// restore gravity assist
	data[1] = 12
	data[2] = 2

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

	fmt.Printf("End: %v\n", data)
}
