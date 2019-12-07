package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var lines []string

	var line string
	for {
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		if len(line) > 0 {
			lines = append(lines, line)
		}

		if err != nil {
			break
		}
	}

	orbits := make(map[string][]string)

	for _, line := range lines {
		objs := strings.SplitN(line, ")", 2)
		o1 := objs[0]
		o2 := objs[1]

		if _, ok := orbits[o1]; ok {
			orbits[o1] = append(orbits[o1], o2)
		} else {
			orbits[o1] = []string{o2}
		}
	}
	fmt.Printf("%v\n", orbits)

	count := getOrbitCount(orbits, "COM", 0)

	fmt.Printf("COUNT: %v\n", count)
}

func getOrbitCount(orbits map[string][]string, k string, depth int) int {
	count := 0
	if children, ok := orbits[k]; ok {
		fmt.Printf("%v ADD CURRENT DEPTH %d\n", k, depth)
		count += depth
		depth = depth + 1

		fmt.Printf("%v CHILDREN %v\n", k, children)
		for _, child := range children {
			count += getOrbitCount(orbits, child, depth)
		}
	} else {
		fmt.Printf("%s NO CHILDREN\n", k)
		count += depth
	}
	return count
}
