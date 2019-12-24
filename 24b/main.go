package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile struct {
	x, y      int
	cur, next rune
}

type Grid struct {
	depth int
	tiles [][]Tile
}

func (t *Tile) Render() string {
	return string(t.cur)
}

func getGrid(path string) *Grid {
	var grid Grid
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]Tile, len(line))
		x := 0
		for i, c := range line {
			row[i] = Tile{x: x, y: y, cur: c}
			x++
		}
		grid.tiles = append(grid.tiles, row)
		y++
	}
	return &grid
}

func (g Grid) Render() {
	for _, row := range g.tiles {
		for _, t := range row {
			fmt.Printf("%s", t.Render())
		}
		fmt.Println()
	}
}

func (g Grid) GetNeighbors(t Tile) []Tile {
	x := t.x
	y := t.y
	var neighbors []Tile

	// get north neighbor(s)
	if y-1 == 2 && x == 2 {
		if inner, ok := dimensions[g.depth+1]; ok {
			neighbors = append(neighbors, inner.GetEdgeTiles('S')...)
		}
	} else {
		if y-1 >= 0 {
			neighbors = append(neighbors, g.tiles[y-1][x])
		} else {
			if outer, ok := dimensions[g.depth-1]; ok {
				neighbors = append(neighbors, outer.tiles[1][2])
			} else {
				neighbors = append(neighbors, Tile{cur: '.'})
			}
		}
	}

	// get south neighbor(s)
	if y+1 == 2 && x == 2 {
		if inner, ok := dimensions[g.depth+1]; ok {
			neighbors = append(neighbors, inner.GetEdgeTiles('N')...)
		}
	} else {
		if y+1 < len(g.tiles) {
			neighbors = append(neighbors, g.tiles[y+1][x])
		} else {
			if outer, ok := dimensions[g.depth-1]; ok {
				neighbors = append(neighbors, outer.tiles[3][2])
			} else {
				neighbors = append(neighbors, Tile{cur: '.'})
			}
		}
	}

	// get west neighbor(s)
	if y == 2 && x-1 == 2 {
		if inner, ok := dimensions[g.depth+1]; ok {
			neighbors = append(neighbors, inner.GetEdgeTiles('E')...)
		}
	} else {
		if x-1 >= 0 {
			neighbors = append(neighbors, g.tiles[y][x-1])
		} else {
			if outer, ok := dimensions[g.depth-1]; ok {
				neighbors = append(neighbors, outer.tiles[2][1])
			} else {
				neighbors = append(neighbors, Tile{cur: '.'})
			}
		}
	}

	// get east neighbor(s)
	if y == 2 && x+1 == 2 {
		if inner, ok := dimensions[g.depth+1]; ok {
			neighbors = append(neighbors, inner.GetEdgeTiles('W')...)
		}
	} else {
		if x+1 < len(g.tiles[x]) {
			neighbors = append(neighbors, g.tiles[y][x+1])
		} else {
			if outer, ok := dimensions[g.depth-1]; ok {
				neighbors = append(neighbors, outer.tiles[2][3])
			} else {
				neighbors = append(neighbors, Tile{cur: '.'})
			}
		}
	}

	return neighbors
}

func (g Grid) GetEdgeTiles(edge rune) []Tile {
	edgeTiles := make([]Tile, 5)
	switch edge {
	case 'E':
		x := len(g.tiles[0]) - 1
		for y := range g.tiles {
			edgeTiles[y] = g.tiles[y][x]
		}
	case 'W':
		x := 0
		for y := range g.tiles {
			edgeTiles[y] = g.tiles[y][x]
		}
	case 'S':
		y := len(g.tiles) - 1
		for x := range g.tiles[y] {
			edgeTiles[x] = g.tiles[y][x]
		}
	case 'N':
		y := 0
		for x := range g.tiles[y] {
			edgeTiles[x] = g.tiles[y][x]
		}
	default:
		panic(fmt.Sprintf("Invalid edge: %s", string(edge)))
	}
	return edgeTiles
}

func (g Grid) GetBugCount(t Tile) int {
	count := 0
	for _, n := range g.GetNeighbors(t) {
		if n.cur == '#' {
			count++
		}
	}
	return count
}

func (g Grid) GetNextState(t Tile) rune {
	numAdjBugs := g.GetBugCount(t)
	if t.cur == '#' && numAdjBugs != 1 { // A bug dies (becoming an empty space) unless there is exactly one bug adjacent to it.
		return '.'
	}
	if t.cur == '.' && numAdjBugs >= 1 && numAdjBugs <= 2 { //An empty space becomes infested with a bug if exactly one or two bugs are adjacent to it.
		return '#'
	}
	return t.cur // Otherwise, a bug or empty space remains the same.
}

func (g Grid) Tick() {
	for y, row := range g.tiles {
		for x, t := range row {
			g.tiles[y][x].next = g.GetNextState(t)
		}
	}
}

func (g Grid) Update() {
	for y, row := range g.tiles {
		for x, t := range row {
			g.tiles[y][x].cur = t.next
		}
	}
}

var dimensions = make(map[int]*Grid)
var minutes = 10

func main() {
	dimensions[-1] = getGrid(os.Args[1])
	dimensions[-1].depth = -1
	dimensions[0] = getGrid(os.Args[1])
	dimensions[0].depth = 0
	dimensions[1] = getGrid(os.Args[1])
	dimensions[1].depth = 1
	fmt.Println("INITIAL")
	dimensions[0].Render()
	fmt.Printf("%v\n", len(dimensions[1].GetNeighbors(dimensions[1].tiles[0][0])))
	// for i := 1; i < minutes/2; i++ { // go forward
	// 	fmt.Printf("\nAFTER %d MINUTES\n", i)
	// 	dimensions[0].Tick()
	// 	dimensions[0].Update()
	// 	dimensions[0].Render()
	// }
}
