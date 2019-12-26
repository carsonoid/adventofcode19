package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var inputInts = []int{}

var xDim = 20
var yDim = 20
var grid = make([][]rune, yDim)
var x = xDim / 2
var y = yDim / 2
var lastIn string
var lastOut string

func draw() {
	for i := 0; i < yDim; i++ {
		for j := 0; j < xDim; j++ {
			if i == y && j == x {
				fmt.Printf("D")
			} else {
				fmt.Printf("%s", string(grid[i][j]))
			}
		}
		fmt.Println()
	}
}

func input(c *computer) int {
	if len(inputInts) == 0 {
		if strings.Contains(lastOut, "You can't go that way.") || strings.Contains(lastOut, "you are ejected back") {
			switch lastIn {
			case "north\n":
				grid[y-1][x] = '#'
			case "south\n":
				grid[y+1][x] = '#'
			case "west\n":
				grid[y][x-1] = '#'
			case "east\n":
				grid[y][x+1] = '#'
			}
		} else {
			switch lastIn {
			case "north\n":
				y--
			case "south\n":
				y++
			case "west\n":
				x--
			case "east\n":
				x++
			}
			grid[y][x] = ' '
		}
		draw()

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		lastIn = text
		if err != nil {
			panic(err)
		}
		fmt.Println()

		for _, r := range text {
			inputInts = append(inputInts, int(r))
		}
	}
	lastOut = ""
	return autoInput()
}

func autoInput() int {
	var i int
	i, inputInts = inputInts[0], inputInts[1:]
	return i
}

func main() {
	code := getData(os.Args[1])

	for i := 0; i < yDim; i++ {
		grid[i] = make([]rune, xDim)
		for j := 0; j < xDim; j++ {
			grid[i][j] = '?'
		}
	}

	// set start room
	grid[y][x] = '#'

	c := newComputer(0, code,
		input,
		func(c *computer, o int) { // Output
			lastOut += string(rune(o))
			fmt.Printf("%s", string(rune(o)))
		},
	)
	go c.start()
	<-c.quit
}
