package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type computer struct {
	id     int
	memory []int64
	input  chan int64
	output chan int64
	result int64
	quit   chan struct{}
	relPos int64
}

func newComputer(id int, code []int64) *computer {
	c := computer{
		id:     id,
		memory: make([]int64, len(code)),
		input:  make(chan int64, 1),
		output: make(chan int64, 1),
		quit:   make(chan struct{}),
	}
	copy(c.memory, code)
	return &c
}

type OpCode int

const (
	OpCodeAdd         OpCode = 1
	OpCodeMultiply    OpCode = 2
	OpCodeInput       OpCode = 3
	OpCodeOutput      OpCode = 4
	OpCodeJumpIfTrue  OpCode = 5
	OpCodeJumpIfFalse OpCode = 6
	OpCodeLessThan    OpCode = 7
	OpCodeEquals      OpCode = 8
	OpCodeRelPos      OpCode = 9
	OpCodeQuit        OpCode = 99
)

type Mode int

const (
	ModePosition  Mode = 0
	ModeImmediate Mode = 1
	ModeRelative  Mode = 2
)

type operation struct {
	opCode OpCode
	modes  []Mode
	params []int64
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getOpCode(m int64) OpCode {
	switch m {
	case 1:
		return OpCodeAdd
	case 2:
		return OpCodeMultiply
	case 3:
		return OpCodeInput
	case 4:
		return OpCodeOutput
	case 5:
		return OpCodeJumpIfTrue
	case 6:
		return OpCodeJumpIfFalse
	case 7:
		return OpCodeLessThan
	case 8:
		return OpCodeEquals
	case 9:
		return OpCodeRelPos
	case 99:
		return OpCodeQuit
	default:
		panic(fmt.Sprintf("INVALID OPCODE: %v", m))
	}
}

func getMode(m int64) Mode {
	switch m {
	case 0:
		return ModePosition
	case 1:
		return ModeImmediate
	case 2:
		return ModeRelative
	default:
		panic(fmt.Sprintf("INVALID MODE: %v", m))
	}
}

func getOperation(code int64, workSet []int64) operation {
	op := operation{}

	// Get Opcode
	opStr := strconv.FormatInt(code, 10)
	// fmt.Printf("OPSTR: %v\tWORKSET: %v\n", opStr, workSet)
	var opCodeStr, modeStr string
	if len(opStr) > 2 {
		opCodeStr = opStr[len(opStr)-2:]
		modeStr = opStr[:len(opStr)-2]
	} else {
		opCodeStr = opStr
		modeStr = ""
	}
	// fmt.Printf("%v\n", opCodeStr)
	opCode, err := strconv.ParseInt(opCodeStr, 10, 64)
	if err != nil {
		panic(err)
	}
	op.opCode = getOpCode(opCode)

	// Get modes
	if op.opCode == OpCodeAdd || op.opCode == OpCodeMultiply || op.opCode == OpCodeLessThan || op.opCode == OpCodeEquals {
		op.modes = make([]Mode, 3)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int64(c - '0'))
		}
	}
	if op.opCode == OpCodeJumpIfFalse || op.opCode == OpCodeJumpIfTrue {
		op.modes = make([]Mode, 2)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int64(c - '0'))
		}
	}
	if op.opCode == OpCodeInput || op.opCode == OpCodeOutput || op.opCode == OpCodeRelPos {
		op.modes = make([]Mode, 1)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int64(c - '0'))
		}
	}

	switch op.opCode {
	case OpCodeAdd:
		fallthrough
	case OpCodeMultiply:
		fallthrough
	case OpCodeLessThan:
		fallthrough
	case OpCodeEquals:
		op.params = workSet[:3]
	case OpCodeInput:
		fallthrough
	case OpCodeOutput:
		fallthrough
	case OpCodeRelPos:
		op.params = workSet[:1]
	case OpCodeJumpIfFalse:
		fallthrough
	case OpCodeJumpIfTrue:
		op.params = workSet[:2]
	}

	// fmt.Printf("OP: %v\n", op)
	return op
}

func (c *computer) shutdown() {
	close(c.output)
	close(c.quit)
}

func (c *computer) getInput() int64 {
	i := <-c.input
	return i
}

func (c *computer) setMemory(pos, val int64) {
	// fmt.Printf("SET POS %v TO: %v\n", pos, val)

	// Grow mem if needed
	memSize := int64(len(c.memory))
	if pos >= memSize {
		newMem := make([]int64, pos-memSize+1)
		c.memory = append(c.memory, newMem...)
	}

	c.memory[pos] = val
}

func (c *computer) getMemory(pos int64) int64 {
	// fmt.Printf("GET POS %v\n", pos)

	// Grow mem if needed
	memSize := int64(len(c.memory))
	if pos >= memSize {
		newMem := make([]int64, pos-memSize+1)
		c.memory = append(c.memory, newMem...)
	}

	return c.memory[pos]
}

func (c *computer) doOperation(op operation) *int64 {
	// fmt.Printf("COMP %d DO OP: %v\n", c.id, op)
	switch op.opCode {
	case OpCodeAdd:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = c.getMemory(op.params[1])
		} else if op.modes[1] == ModeRelative {
			v2 = c.getMemory(c.relPos + op.params[1])
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else if op.modes[2] == ModeRelative {
			loc = c.relPos + op.params[2]
		} else {
			loc = c.getMemory(op.params[2])
		}
		c.setMemory(loc, v1+v2)
	case OpCodeMultiply:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = c.getMemory(op.params[1])
		} else if op.modes[1] == ModeRelative {
			v2 = c.getMemory(c.relPos + op.params[1])
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else if op.modes[2] == ModeRelative {
			loc = c.relPos + op.params[2]
		} else {
			loc = c.getMemory(op.params[2])
		}
		c.setMemory(loc, v1*v2)
	case OpCodeInput:
		var loc int64
		// loc logic is reversed
		if op.modes[0] == ModePosition {
			loc = op.params[0]
		} else if op.modes[0] == ModeRelative {
			loc = op.params[0] + c.relPos
		} else {
			loc = c.getMemory(op.params[0])
		}
		c.setMemory(loc, c.getInput())
	case OpCodeOutput:
		if op.modes[0] == ModePosition {
			c.result = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			c.result = c.getMemory(c.relPos + op.params[0])
		} else {
			c.result = op.params[0]
		}
		// fmt.Printf("COMP %d OUTPUT TO CHAIN:\t%d\n", c.id, c.result)
		c.output <- c.result // send result to output chan
	case OpCodeJumpIfTrue:
		var v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = c.getMemory(op.params[1])
		} else if op.modes[1] == ModeRelative {
			v2 = c.getMemory(c.relPos + op.params[1])
		} else {
			v2 = op.params[1]
		}
		if v1 != 0 {
			return &v2
		}
	case OpCodeJumpIfFalse:
		var v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = c.getMemory(op.params[1])
		} else if op.modes[1] == ModeRelative {
			v2 = c.getMemory(c.relPos + op.params[1])
		} else {
			v2 = op.params[1]
		}
		if v1 == 0 {
			return &v2
		}
	case OpCodeLessThan:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = c.getMemory(op.params[1])
		} else if op.modes[1] == ModeRelative {
			v2 = c.getMemory(c.relPos + op.params[1])
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else if op.modes[2] == ModeRelative {
			loc = c.relPos + op.params[2]
		} else {
			loc = c.getMemory(op.params[2])
		}
		if v1 < v2 {
			c.setMemory(loc, 1)
		} else {
			c.setMemory(loc, 0)
		}
	case OpCodeEquals:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = c.getMemory(op.params[1])
		} else if op.modes[1] == ModeRelative {
			v2 = c.getMemory(c.relPos + op.params[1])
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else if op.modes[2] == ModeRelative {
			loc = c.relPos + op.params[2]
		} else {
			loc = c.getMemory(op.params[2])
		}
		if v1 == v2 {
			c.setMemory(loc, 1)
		} else {
			c.setMemory(loc, 0)
		}
	case OpCodeRelPos:
		var v1 int64
		if op.modes[0] == ModePosition {
			v1 = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			v1 = c.getMemory(c.relPos + op.params[0])
		} else {
			v1 = op.params[0]
		}
		c.relPos += v1
	}
	return nil
}

func getData(filePath string) []int64 {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var data []int64
	for {
		dataPoint, err := reader.ReadString(',')
		dataPoint = strings.Trim(dataPoint, ",\n")

		u, err := strconv.ParseInt(dataPoint, 10, 64)
		if err != nil {
			break
		}
		data = append(data, u)
	}

	return data
}

func (c *computer) start() {
	var op operation
	var pos int
	var quit bool
	for {
		// fmt.Printf("COMP %v POS: %d, len: %d\n", c.id, pos, len(c.memory))
		if pos == len(c.memory) {
			quit = true
		} else if pos+1 == len(c.memory) {
			op = getOperation(c.memory[pos], c.memory[pos+1:pos+2])
			quit = true
		} else if pos+2 == len(c.memory) {
			op = getOperation(c.memory[pos], c.memory[pos+1:pos+2])
			quit = true
		} else if pos+3 == len(c.memory) {
			op = getOperation(c.memory[pos], c.memory[pos+1:pos+3])
			quit = true
		} else if pos+4 == len(c.memory) {
			op = getOperation(c.memory[pos], c.memory[pos+1:pos+4])
			quit = true
		} else {
			op = getOperation(c.memory[pos], c.memory[pos+1:pos+4])
		}
		pos += len(op.params) + 1

		if quit { // do final operation
			// fmt.Printf("COMP %v HIT END OF CODE\n", c.id)
			// fmt.Printf("COMP %v LAST OP:%#v\n", c.id, op)
			c.doOperation(op)
			c.shutdown()
			return
		}

		if op.opCode == OpCodeQuit { // Quit immediately
			// fmt.Printf("COMP %v QUIT CODE DURING RUN!\n", c.id)
			c.shutdown()
			return
		}

		newPos := c.doOperation(op)
		if newPos != nil {
			pos = int(*newPos)
			// fmt.Printf("MOVED PTR to %d\n", pos)
		}
	}
}

func permutations(arr []int64) [][]int64 {
	var helper func([]int64, int64)
	res := [][]int64{}

	helper = func(arr []int64, n int64) {
		if n == 1 {
			tmp := make([]int64, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := int64(0); i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, int64(len(arr)))
	return res
}

func max(in []int64) int64 {
	m := int64(0)
	for i, e := range in {
		if i == 0 || int64(e) > m {
			m = int64(e)
		}
	}
	return m
}

func main() {
	code := getData(os.Args[1])
	// fmt.Printf("CODE: %v\n", code)

	c := newComputer(0, code)
	c.input <- 2

	go c.start()
	for out := range c.output {
		fmt.Printf("OUTPUT: %v\n", out)
	}
	<-c.quit

}
