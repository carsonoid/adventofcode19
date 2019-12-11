package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
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
	numAsteroids := 0
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
				lineData = append(lineData, point{x, y, false, 0.0})
			case '#':
				numAsteroids++
				lineData = append(lineData, point{x, y, true, 0.0})
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

	monCounts := map[point]map[float64][]point{}
	for _, lineData := range mapData {
		for _, p := range lineData {
			if p.hasAsteroid == true {
				ma := getMonitorableAsteroids(p, mapData)
				monCounts[p] = ma
			}
		}
	}
	// fmt.Printf("%v\n", monCounts)

	winner := point{}
	maxCount := 0
	for pt, ma := range monCounts {
		if len(ma) > maxCount {
			maxCount = len(ma)
			winner = pt
		}
	}

	fmt.Printf("WINNER: %v. NUM: %d\n", winner, maxCount)

	// Part 2
	winnerData := monCounts[winner]

	// Fill all points with distance from winner
	for angle, asteroids := range winnerData {
		for i, asteroid := range asteroids {
			winnerData[angle][i].distance = getDistance(winner, asteroid)
		}
	}

	// Sort each slice of points by distance from winner
	for angle, asteroids := range winnerData {
		sort.Slice(asteroids, func(i, j int) bool {
			return asteroids[i].distance < asteroids[j].distance
		})
		winnerData[angle] = asteroids
	}

	// Get Angles
	angles := []float64{}
	for angle := range winnerData {
		angles = append(angles, angle)
	}
	fmt.Printf("ANGLES: %v\n", angles)

	// Sort
	sort.Float64s(angles)
	fmt.Printf("SORTED ANGLES: %v\n", angles)

	// Loop through sorted angles
	explodableAsteroids := numAsteroids - 1 // exclude self
	explodedAsteroids := 0
TOP:
	for {
		for _, angle := range angles {
			fmt.Printf("LAZOR AT ANGLE: %v\n", angle)
			if len(winnerData[angle]) > 0 {
				// shift off front
				var p point
				p, winnerData[angle] = winnerData[angle][0], winnerData[angle][1:]
				fmt.Printf("%v EXPLODES!!!\n", p.String())
				explodedAsteroids++

				if explodedAsteroids == 200 {
					fmt.Printf("BET ON %v\n", p.String())
				}

				if explodedAsteroids == explodableAsteroids { // exclude self from explosion
					fmt.Printf("EXPLODED %d/%d ASTEROIDS\n", explodedAsteroids, explodableAsteroids)
					fmt.Printf("SPACE IS CLEAR\n")
					break TOP
				}
			}
		}
	}
}

func getMonitorableAsteroids(source point, mapData [][]point) map[float64][]point {
	asteroids := map[float64][]point{}
	// fmt.Printf("%s\n", source.String())
	for _, lineData := range mapData {
		for _, p := range lineData {
			if p.hasAsteroid && !(p.x == source.x && p.y == source.y) { // has asteroid and not self
				a := getAngleBetweenPoints(source, p)
				asteroids[a] = append(asteroids[a], p)
			}
		}
	}
	// fmt.Printf("%v\n", asteroids)
	return asteroids
}

func getAngleBetweenPoints(o, t point) float64 {
	x := float64(t.x) - float64(o.x)
	y := float64(t.y) - float64(o.y)

	var degrees float64

	degrees = (math.Atan2(x, y) - math.Atan2(1, 0)) * 180 / math.Pi //shift
	degrees = math.Abs(degrees - 90)                                // invert

	return degrees
}

// (x2 - x1)^2 + (y2 - y1)^2
func getDistance(p1, p2 point) float64 {
	first := math.Pow(float64(p2.x-p1.x), 2)
	second := math.Pow(float64(p2.y-p1.y), 2)
	return math.Sqrt(first + second)
}

type point struct {
	x, y        int
	hasAsteroid bool
	distance    float64
}

func (p *point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
