package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func printMap(mapData [][]point) {
	fmt.Printf("Map:\n")
	for _, lineData := range mapData {
		for _, p := range lineData {
			if p.hasAsteroid == true {
				fmt.Printf("\u00b7")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	mapData := [][]point{}
	var line string
	y := 0
	for {
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")

		if line == "" {
			break
		}

		lineData := []point{}
		for x, p := range line {
			switch p {
			case '.':
				lineData = append(lineData, point{x, y, false})
			case '#':
				lineData = append(lineData, point{x, y, true})
			default:
				panic(fmt.Sprintf("Invalid Map Point: %v", p))
			}
		}
		mapData = append(mapData, lineData)

		if err != nil {
			break
		}

		y++
	}

	printMap(mapData)

	monCounts := map[string]int{}
	for _, lineData := range mapData {
		for _, p := range lineData {
			if p.hasAsteroid == true {
				monCounts[p.String()] = getNumMonitorableAsteroids(p, mapData)
			}
		}
	}
	// fmt.Printf("%v\n", monCounts)

	winner := ""
	maxCount := 0
	for pt, count := range monCounts {
		if count > maxCount {
			maxCount = count
			winner = pt
		}
	}

	fmt.Printf("WINNER: %v. NUM: %d\n", winner, maxCount)
}

func getNumMonitorableAsteroids(source point, mapData [][]point) int {
	angles := map[string]int{}
	// fmt.Printf("%s\n", source.String())
	for _, lineData := range mapData {
		for _, p := range lineData {
			if p.hasAsteroid && !(p.x == source.x && p.y == source.y) { // has asteroid and not self
				a := getAngleBetweenPoints(source, p)
				angles[fmt.Sprintf("%f", a)]++
			}
		}
	}
	// fmt.Printf("%v\n", angles)
	return len(angles)
}

func getAngleBetweenPoints(p1, p2 point) float64 {
	dx := float64(p1.x) - float64(p2.x)
	dy := float64(p1.y) - float64(p2.y)
	return math.Atan2(dy, dx)
}

type point struct {
	x, y        int
	hasAsteroid bool
}

func (p *point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
