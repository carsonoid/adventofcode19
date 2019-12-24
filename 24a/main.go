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

type Grid [][]Tile

func (t *Tile) Render() string {
	return string(t.cur)
}

func getGrid(path string) Grid {
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
		grid = append(grid, row)
		y++
	}
	return grid
}

func (g Grid) Render() {
	for _, row := range g {
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
	if y+1 < len(g) {
		neighbors = append(neighbors, g[y+1][x])
	} else { // Tiles on the edges of the grid have fewer than four adjacent tiles; the missing tiles count as empty space.)
		neighbors = append(neighbors, Tile{cur: '.'})
	}
	if y-1 >= 0 {
		neighbors = append(neighbors, g[y-1][x])
	} else {
		neighbors = append(neighbors, Tile{cur: '.'})
	}
	if x+1 < len(g[x]) {
		neighbors = append(neighbors, g[y][x+1])
	} else {
		neighbors = append(neighbors, Tile{cur: '.'})
	}
	if x-1 >= 0 {
		neighbors = append(neighbors, g[y][x-1])
	} else {
		neighbors = append(neighbors, Tile{cur: '.'})
	}
	return neighbors
}

func (g Grid) GetNextState(t Tile) rune {
	neighbors := g.GetNeighbors(t)
	numAdjBugs := 0
	for _, n := range neighbors {
		if n.cur == '#' {
			numAdjBugs++
		}
	}
	if t.cur == '#' && numAdjBugs != 1 { // A bug dies (becoming an empty space) unless there is exactly one bug adjacent to it.
		return '.'
	}
	if t.cur == '.' && numAdjBugs >= 1 && numAdjBugs <= 2 { //An empty space becomes infested with a bug if exactly one or two bugs are adjacent to it.
		return '#'
	}
	return t.cur // Otherwise, a bug or empty space remains the same.
}

func (g Grid) Tick() {
	for y, row := range g {
		for x, t := range row {
			g[y][x].next = g.GetNextState(t)
		}
	}
}

func (g Grid) Update() {
	for y, row := range g {
		for x, t := range row {
			g[y][x].cur = t.next
		}
	}
}

func (g Grid) GetBioRating() int {
	r := 0
	v := 1
	for _, row := range g {
		for _, t := range row {
			if t.cur == '#' {
				r += v
			}
			v = v * 2
		}
	}
	return r
}

func main() {

	grid := getGrid(os.Args[1])
	fmt.Println("INITIAL")
	bioRatings := map[int]int{
		0: grid.GetBioRating(),
	}
	grid.Render()
	for i := 1; ; i++ {
		// fmt.Printf("\nAFTER %d MINUTES\n", i)
		grid.Tick()
		grid.Update()
		// grid.Render()
		br := grid.GetBioRating()
		// fmt.Printf("BIO: %d\n", br)
		bioRatings[br]++
		if bioRatings[br] == 2 {
			grid.Render()
			fmt.Printf("%d MINUTES TO FIRST REPEATED RATING of %d\n", i, br)
			break
		}
	}
}
