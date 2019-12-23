package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
		fmt.Printf("%v", o)
	}
}

func main() {
	code := getData(os.Args[1])
	inputInts = readSpringscript(os.Args[2])
	c := newComputer(0, code, autoInput, output)
	go c.start()
	<-c.quit
}
