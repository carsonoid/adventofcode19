package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func calculateModuleFuel(mass uint64) uint64 {
	answer := math.Floor(float64(mass)/3) - 2
	return uint64(answer)
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

			totalFuel += calculateModuleFuel(moduleMass)
		}

		if err != nil {
			break
		}
	}

	fmt.Printf("Total Fuel: %d\n", totalFuel)
}
