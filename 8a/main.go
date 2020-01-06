package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	imgWidth, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	imgHeight, _ := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}
	imgSize := imgHeight * imgWidth

	var layers [][]int

	layer := []int{}
	curWidth := 0
	for {
		curWidth++

		r, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		} else if r == '\n' {
			break
		}
		pixel := int(r - '0')
		layer = append(layer, pixel)

		if len(layers) == imgHeight && curWidth == imgSize {
			break
		}
		if curWidth == imgSize {
			layers = append(layers, layer)
			layer = []int{}
			curWidth = 0
		}
	}

	for i := range layers {
		fmt.Printf("Layer %v: %v\n", i, layers[i])
	}

	checksums := map[int]map[int]int{}
	for i, layer := range layers {
		checksums[i] = map[int]int{}
		for _, pixel := range layer {
			checksums[i][pixel]++
		}
	}
	fmt.Printf("Checksums %v\n", checksums)

	winner := 0
	lastCount := -1
	for lnum, counts := range checksums {
		count := counts[0]
		if count < lastCount || lastCount < 0 {
			winner = lnum
			lastCount = count
		}
	}

	fmt.Printf("WINNER! %v\n", winner)

	pixelCounts := map[int]int{}
	for _, pixel := range layers[winner] {
		pixelCounts[pixel]++
	}
	fmt.Printf("ANSWER! %v\n", pixelCounts[1]*pixelCounts[2])
}
