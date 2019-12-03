package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type collision struct {
	p     position
	steps int
}

type position struct {
	x, y int64
}

// Abs returns the absolute value of x.
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (p *position) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p *position) Distance() int64 {
	return Abs(p.x) + Abs(p.y)
}

type instruction struct {
	direction rune
	distance  int64
}

func getCrossedPositions(a, b []position) []collision {
	collisions := []collision{}

	// Create a map of all items in a
	m := make(map[string]int)
	for aSteps, item := range a {
		m[item.String()] = aSteps + 1 // Steps start with 1 even though index starts with 0
	}

	// Find any occurances in map of items in b add ot collisions if true
	for bSteps, item := range b {
		if aSteps, ok := m[item.String()]; ok {
			bSteps++ // Steps start with 1 even though index starts with 0
			collisions = append(collisions, collision{
				p:     item,
				steps: aSteps + bSteps,
			})
		}
	}

	return collisions
}

func getClosestCollision(collisions []collision) collision {
	closest := collision{}
	var min int64
	for _, col := range collisions {
		distance := col.p.Distance()
		if min == 0 {
			min = distance
			closest = col
		} else if distance < min {
			min = distance
			closest = col
		}
	}
	return closest
}

func getFastestCollision(collisions []collision) collision {
	fastest := collision{}
	var min int
	for _, col := range collisions {
		steps := col.steps
		if min == 0 {
			min = steps
			fastest = col
		} else if steps < min {
			min = steps
			fastest = col
		}
	}
	return fastest
}

func getClosestPosition(collisions []collision) position {
	closest := position{}
	var min int64
	for _, c := range collisions {
		distance := c.p.Distance()
		if min == 0 {
			min = distance
			closest = c.p
		} else if distance < min {
			min = distance
			closest = c.p
		}
	}
	return closest
}

func getPassedPositions(instructions []instruction) []position {
	pos := position{}
	positions := []position{}

	for _, inst := range instructions {
		switch inst.direction {
		case 'R':
			for j := inst.distance; j > 0; j-- {
				pos.x++
				positions = append(positions, pos)
			}
		case 'L':
			for j := inst.distance; j > 0; j-- {
				pos.x--
				positions = append(positions, pos)
			}
		case 'U':
			for j := inst.distance; j > 0; j-- {
				pos.y++
				positions = append(positions, pos)
			}
		case 'D':
			for j := inst.distance; j > 0; j-- {
				pos.y--
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

func getFewestWireCrossSteps(inst1, inst2 []instruction) int {

	pos1 := getPassedPositions(inst1)
	// fmt.Printf("%v\n", pos1)

	pos2 := getPassedPositions(inst2)
	// fmt.Printf("%v\n", pos2)

	crossed := getCrossedPositions(pos1, pos2)
	// fmt.Printf("%v\n", crossed)

	c := getFastestCollision(crossed)
	fmt.Printf("fastest position: %v\n", c.p)
	return c.steps
}

func getClosestWireCrossDistance(inst1, inst2 []instruction) int64 {

	pos1 := getPassedPositions(inst1)
	// fmt.Printf("%v\n", pos1)

	pos2 := getPassedPositions(inst2)
	// fmt.Printf("%v\n", pos2)

	crossed := getCrossedPositions(pos1, pos2)
	// fmt.Printf("%v\n", crossed)

	p := getClosestPosition(crossed)
	fmt.Printf("Closest position: %v\n", p)
	return p.Distance()
}

func getInstructions(filePath string) ([]instruction, []instruction) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	inst1 := getInstructionsFromString(lines[0])
	inst2 := getInstructionsFromString(lines[1])

	return inst1, inst2
}

func getInstructionsFromString(in string) []instruction {
	stringInstructions := strings.Split(in, ",")
	var instructions []instruction
	for _, sInst := range stringInstructions {
		runes := []rune(sInst)
		direction := runes[0]
		distance, err := strconv.ParseInt(sInst[1:], 10, 64)
		if err != nil {
			break
		}
		instructions = append(instructions, instruction{
			direction: direction,
			distance:  distance,
		})
	}
	return instructions
}

func main() {
	// line1Instructions := []instruction{{'R', 75}, {'D', 30}, {'R', 83}, {'U', 83}, {'L', 12}, {'D', 49}, {'R', 71}, {'U', 7}, {'L', 72}}
	// line2Instructions := []instruction{{'U', 62}, {'R', 66}, {'U', 55}, {'R', 34}, {'D', 71}, {'R', 55}, {'D', 58}, {'R', 83}}

	inst1, inst2 := getInstructions(os.Args[1])

	closest := getClosestWireCrossDistance(inst1, inst2)
	fmt.Printf("Closest Distance: %v\n", closest)

	fastest := getFewestWireCrossSteps(inst1, inst2)
	fmt.Printf("Fastest Steps: %v\n", fastest)
}
