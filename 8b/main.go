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

		if curWidth == imgSize {
			layers = append(layers, layer)
			layer = []int{}
			curWidth = 0
		}
		if len(layers) == imgHeight && curWidth == imgSize {
			break
		}
	}

	for i := range layers {
		fmt.Printf("%v\n", layers[i])
	}

	// Get the rendered layer. The first non-transparent pixel
	// for a position is the rendered image
	rendered := make([]int, imgSize)
	for i := range rendered {
		rendered[i] = 2
	}
	for _, layer := range layers {
		for pos, pixel := range layer {
			if pixel != 2 && rendered[pos] == 2 {
				rendered[pos] = pixel
			}
		}
	}
	fmt.Printf("\n%v\n", rendered)

	curWidth = 0
	for _, pixel := range rendered {
		curWidth++
		switch pixel {
		case 0:
			fmt.Printf("\u2591")
		case 1:
			fmt.Printf("\u2588")
		}
		// fmt.Printf("%d", pixel)
		if curWidth%imgWidth == 0 {
			fmt.Printf("\n")
		}
	}
}
