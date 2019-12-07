package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// file, err := os.Open("in-test.txt")
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
	parents := make(map[string]string)

	for _, line := range lines {
		objs := strings.SplitN(line, ")", 2)
		o1 := objs[0]
		o2 := objs[1]
		parents[o2] = o1

		if _, ok := orbits[o1]; ok {
			orbits[o1] = append(orbits[o1], o2)
		} else {
			orbits[o1] = []string{o2}
		}
	}
	// fmt.Printf("%v\n", orbits)

	count := getOrbitCount(orbits, "COM", 0)
	fmt.Printf("COUNT: %v\n", count)

	yps := make(map[string]string)
	yPath := []string{}
	o := "YOU"
	for {
		if p, ok := parents[o]; ok {
			fmt.Printf("%s PARENT = %v\n", o, p)
			yPath = append(yPath, p)
			yps[o] = p
			o = p
		} else {
			break
		}
	}

	fmt.Printf("YOU PARENTS: %v\n", yps)
	fmt.Printf("YOU PATH: %v\n", yPath)

	sps := make(map[string]string)
	sPath := []string{}
	o = "SAN"
	for {
		if p, ok := parents[o]; ok {
			fmt.Printf("%s PARENT = %v\n", o, p)
			sps[o] = p
			sPath = append(sPath, p)
			o = p
		} else {
			break
		}
	}

	fmt.Printf("SAN PARENTS: %v\n", yps)
	fmt.Printf("SAN PATH: %v\n", sPath)

	for _, v := range yPath {
		fmt.Printf("CHECK SAN PARENT: %v\n", v)
		if _, ok := sps[v]; ok {
			fmt.Printf("FIRST SHARED PARENT: %v\n", v)
			var youDist int
			getDist(orbits, yps, v, "YOU", 0, &youDist)
			var sanDist int
			getDist(orbits, sps, v, "SAN", 0, &sanDist)
			fmt.Printf("FSP to YOU: %v\n", youDist)
			fmt.Printf("FSP to SAN: %v\n", sanDist)
			fmt.Printf("NUM XFERS: %v\n", youDist+sanDist)
			break
		}
	}
}

func getDist(orbits map[string][]string, lineage map[string]string, start string, end string, depth int, result *int) int {
	count := 0
	if children, ok := orbits[start]; ok {

		fmt.Printf("%v ADD CURRENT DEPTH %d\n", start, depth)
		count += depth
		depth = depth + 1
		fmt.Printf("%v CHILDREN %v\n", start, children)
		for _, child := range children {
			if _, ok := lineage[child]; ok {
				if child == end {
					fmt.Printf("REACHED END: %s DEPTH: %v\n", child, depth)
					*result = depth - 1
				}
				fmt.Printf("%v IS A PARENT OF %v\n", child, start)
				getDist(orbits, lineage, child, end, depth, result)
			}
		}
	}
	return depth
}

func getOrbitCount(orbits map[string][]string, k string, depth int) int {
	count := 0
	if children, ok := orbits[k]; ok {
		// fmt.Printf("%v ADD CURRENT DEPTH %d\n", k, depth)
		count += depth
		depth = depth + 1

		// fmt.Printf("%v CHILDREN %v\n", k, children)
		for _, child := range children {
			count += getOrbitCount(orbits, child, depth)
		}
	} else {
		// fmt.Printf("%s NO CHILDREN\n", k)
		count += depth
	}
	return count
}
