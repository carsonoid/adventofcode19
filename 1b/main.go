package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func getFuelForMass(mass uint64) uint64 {
	answer := math.Floor(float64(mass)/3) - 2
	if answer <= 0.0 {
		return 0
	}
	return uint64(answer)
}

func getFuelForModuleMass(mass uint64) uint64 {
	var total uint64
	for {
		mass = getFuelForMass(mass)
		if mass > 0 {
			total += mass
		} else {
			break
		}
	}
	return total
}

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var totalFuel uint64
	var line string
	for {
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")

		if line != "" {
			moduleMass, err := strconv.ParseUint(line, 10, 64)
			if err != nil {
				panic(err)
			}

			totalFuel += getFuelForModuleMass(moduleMass)
		}

		if err != nil {
			break
		}
	}

	fmt.Printf("Total Fuel: %d\n", totalFuel)
}
