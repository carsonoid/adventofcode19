package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	X, Y     int
	Tile     rune
	Children []Node
}

var Nodes = []Node{}

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var screen [][]rune
	for {
		line, err := reader.ReadString('\n')

		row := []rune{}
		if line != "" {
			for _, c := range line {
				if c == '\n' {
					screen = append(screen, row)
					row = []rune{}
				} else {
					row = append(row, c)
				}
			}
		}

		if err != nil {
			break
		}
	}

	var start Node
	for y, _ := range screen {
		for x := range screen[y] {
			fmt.Printf(string(screen[y][x]))
			if screen[y][x] == '@' {
				start = Node{
					X:    x,
					Y:    y,
					Tile: '@',
				}
			}
		}
		fmt.Println()
	}

	start.setChildren(screen)
	fmt.Println("START", start, len(start.Children))
	for _, child := range start.Children {
		fmt.Printf("%d,%d %s\n", child.X, child.Y, string(child.Tile))
	}
}

func (n *Node) setChildren(screen [][]rune) {
	if n.X-1 >= 0 {
		n.Children = append(n.Children, Node{
			X:    n.X - 1,
			Y:    n.Y,
			Tile: screen[n.Y][n.X-1],
		})
	}
	if n.X+1 <= len(screen[n.Y]) {
		n.Children = append(n.Children, Node{
			X:    n.X + 1,
			Y:    n.Y,
			Tile: screen[n.Y][n.X+1],
		})
	}

	if n.Y-1 >= 0 {
		n.Children = append(n.Children, Node{
			X:    n.X,
			Y:    n.Y - 1,
			Tile: screen[n.Y-1][n.X],
		})
	}
	if n.Y+1 <= len(screen) {
		n.Children = append(n.Children, Node{
			X:    n.X,
			Y:    n.Y + 1,
			Tile: screen[n.Y+1][n.X],
		})
	}
}
