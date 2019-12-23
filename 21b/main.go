package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var inputInts = []int{}

func input() int {
	if len(inputInts) == 0 {
		fmt.Printf("INPUT: ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println()

		for _, r := range text {
			inputInts = append(inputInts, int(r))
		}
	}
	return autoInput()
}

func readSpringscript(path string) []int {
	ret := []int{}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	rawInput, err := ioutil.ReadAll(reader)
	for _, op := range rawInput {
		ret = append(ret, int(op))
	}
	return ret
}

func autoInput() int {
	var i int
	i, inputInts = inputInts[0], inputInts[1:]
	return i
}

func output(o int) {
	if 0 <= o && o <= 127 {
		fmt.Printf("%v", string(rune(o)))
	} else {
		success = true
		fmt.Printf("%v", o)
	}
}

var success = false
var ops = []string{"AND", "OR", "NOT"}
var roRegs = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
var rwRegs = []string{"T", "J"}
var allRegs = append(roRegs, rwRegs...)
var finalOp = "RUN"

func getInstructions() []string {
	instructions := []string{}
	for _, op := range ops {
		for _, reg1 := range allRegs {
			for _, rwReg := range rwRegs {
				if reg1 != rwReg {
					c := fmt.Sprintf("%s %s %s\n", op, reg1, rwReg)
					instructions = append(instructions, c)
				}
			}
		}
	}
	return instructions
}

func getProgramInts(in []string) []int {
	p := []int{}
	for _, r := range strings.Join(in, "") {
		p = append(p, int(r))
	}
	return p
}

func main() {
	code := getData(os.Args[1])
	instructions := getInstructions()
	for _, inst := range instructions {
		for _, inst2 := range instructions {
			for _, inst3 := range instructions {
				for _, inst4 := range instructions {
					for _, inst5 := range instructions {
						for _, inst6 := range instructions {
							for _, inst7 := range instructions {
								program := []string{
									inst,
									inst2,
									inst3,
									inst4,
									inst5,
									inst6,
									inst7,
									"RUN",
									"\n",
								}
								// fmt.Println(strings.Join(program, ""))
								inputInts = getProgramInts(program)
								runProgram(code)
							}
						}
					}
				}
			}
		}
	}
}

func runProgram(code []int) {
	c := newComputer(0, code, autoInput, func(o int) {
		if o > 127 {
			fmt.Printf("%v", o)
			success = true
		}
	})
	go c.start()
	<-c.quit
	if success {
		os.Exit(0)
	}
}
