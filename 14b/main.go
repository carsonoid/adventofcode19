package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type chem struct {
	units int
	name  string
}

func (c chem) String() string {
	return fmt.Sprintf("%d %s", c.units, c.name)
}

type reaction struct {
	output chem
	inputs []chem
}

func (r reaction) String() string {
	instr := []string{}
	for _, in := range r.inputs {
		instr = append(instr, in.String())
	}
	return fmt.Sprintf("%v => %v", strings.Join(instr, ", "), r.output)
}

var reactions = map[string]reaction{}
var storage = make(map[string]int)
var oreCount int

func getChem(s string) chem {
	s = strings.Trim(s, " ")
	p := strings.SplitN(s, " ", 2)
	u, err := strconv.Atoi(p[0])
	if err != nil {
		panic(err)
	}
	return chem{
		units: u,
		name:  p[1],
	}
}

func getReactions() map[string]reaction {
	reactions := make(map[string]reaction)
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")

		if line != "" {
			r := reaction{}
			parts := strings.SplitN(line, "=>", 2)
			for _, x := range strings.Split(parts[0], ",") {
				r.inputs = append(r.inputs, getChem(x))
			}
			r.output = getChem(parts[1])
			reactions[r.output.name] = r
		}

		if err != nil {
			break
		}
	}
	return reactions
}

var fuel int
var oreStored = 1000000000000

func main() {
	reactions = getReactions()

	// Brute forced
	for {
		doReaction(-1, reactions["FUEL"], 1)
		fuel++
		if oreStored%1000 == 0 {
			fmt.Printf("FUEL: %d, ORE: %d\n", fuel, oreStored)
		}
	}
}

func doReaction(depth int, r reaction, need int) {
	depth++
	// fmt.Printf("%sNEED %v, DO REATION %v\n", strings.Repeat("\t", depth), r.output, r)
	// printStorage("BEFORE", depth)
	for _, input := range r.inputs {
		if input.name == "ORE" {
			// fmt.Printf("%sTAKE %d ORE", strings.Repeat("\t", depth), input.units)
			if oreStored < input.units {
				fmt.Printf("OUT OF ORE. FULE %d\n", fuel)
				os.Exit(1)
			}
			oreCount += input.units
			oreStored -= input.units
		} else {
			for {
				if storage[input.name] < input.units {
					doReaction(depth, reactions[input.name], input.units)
				} else {
					break
				}
			}
			// fmt.Printf("%sTAKE %d of %d %s; ", strings.Repeat("\t", depth), input.units, storage[input.name], input.name)
			storage[input.name] -= input.units
			// fmt.Printf("%v %s LEFT\n", storage[input.name], input.name)
		}
	}
	// fmt.Printf("%sGENERATED: %v\n", strings.Repeat("\t", depth), r.output)
	storage[r.output.name] += r.output.units
	// printStorage("AFTER", depth)
}

func printStorage(t string, d int) {
	fmt.Printf("%sSTORAGE %s: ", strings.Repeat("\t", d), t)
	if len(storage) == 0 {
		fmt.Printf("EMPTY")
	}
	for c, u := range storage {
		fmt.Printf(" %d %s,", u, c)
	}
	fmt.Println()
}
