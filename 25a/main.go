package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	combinations "github.com/mxschmitt/golang-combinations"
)

var inputInts = []int{}

var xDim = 20
var yDim = 20
var grid = make([][]rune, yDim)
var x = xDim / 2
var y = yDim / 2
var lastIn string
var lastOut string
var items []string

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

		fmt.Println("text", text, "lastin", lastIn)
		if lastIn == "auto\n" {
			items = getItems(lastOut)
			text = getPermutationTestCode(items)
		}

		for _, r := range text {
			inputInts = append(inputInts, int(r))
		}
	}
	lastOut = ""
	return autoInput()
}

func permutations(arr []string) [][]string {
	var helper func([]string, int)
	res := [][]string{}

	helper = func(arr []string, n int) {
		if n == 1 {
			tmp := make([]string, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

var direction string = "west"

func getPermutationTestCode(items []string) string {
	var code []string
	// first drop all
	for _, item := range items {
		code = append(code, fmt.Sprintf("drop %s", item))
	}
	// then try all combos
	for _, comb := range combinations.All(items) {
		for _, item := range comb {
			code = append(code, fmt.Sprintf("take %s", item))
		}
		code = append(code, "west")
		for _, item := range comb {
			code = append(code, fmt.Sprintf("drop %s", item))
		}
	}
	fmt.Println(strings.Join(code, "\n"))
	return strings.Join(code, "\n")
}

func getItems(s string) []string {
	fmt.Println("RAW", s)
	items = []string{}
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "- ") {
			items = append(items, strings.TrimLeft(line, "- "))
		}
	}
	fmt.Println("ITEMS", items)
	return items
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
