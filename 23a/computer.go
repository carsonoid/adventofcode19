package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type inputFunc func(c *computer) int
type outputFunc func(c *computer, o int)

type computer struct {
	id        int
	memory    []int
	input     inputFunc
	output    outputFunc
	result    int
	quit      chan struct{}
	relPos    int
	inBuffer  []int
	outBuffer []int
}

func newComputer(id int, code []int, input inputFunc, output outputFunc) *computer {
	c := computer{
		id:        id,
		memory:    make([]int, len(code)),
		input:     input,
		output:    output,
		quit:      make(chan struct{}),
		inBuffer:  []int{},
		outBuffer: []int{},
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
	params []int
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getOpCode(m int) OpCode {
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

func getMode(m int) Mode {
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

func getOperation(code int, workSet []int) operation {
	op := operation{}

	// Get Opcode
	opStr := strconv.Itoa(code)
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
	opCode, err := strconv.Atoi(opCodeStr)
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
			op.modes[i] = getMode(int(c - '0'))
		}
	}
	if op.opCode == OpCodeJumpIfFalse || op.opCode == OpCodeJumpIfTrue {
		op.modes = make([]Mode, 2)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int(c - '0'))
		}
	}
	if op.opCode == OpCodeInput || op.opCode == OpCodeOutput || op.opCode == OpCodeRelPos {
		op.modes = make([]Mode, 1)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int(c - '0'))
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

func (c *computer) BufferInput(ins ...int) {
	c.inBuffer = append(c.inBuffer, ins...)
}

func (c *computer) shutdown() {
	close(c.quit)
}

func (c *computer) setMemory(pos, val int) {
	// fmt.Printf("SET POS %v TO: %v\n", pos, val)

	// Grow mem if needed
	memSize := int(len(c.memory))
	if pos >= memSize {
		newMem := make([]int, pos-memSize+1)
		c.memory = append(c.memory, newMem...)
	}

	c.memory[pos] = val
}

func (c *computer) getMemory(pos int) int {
	// fmt.Printf("GET POS %v\n", pos)

	// Grow mem if needed
	memSize := int(len(c.memory))
	if pos >= memSize {
		newMem := make([]int, pos-memSize+1)
		c.memory = append(c.memory, newMem...)
	}

	return c.memory[pos]
}

func (c *computer) doOperation(op operation) *int {
	// fmt.Printf("COMP %d DO OP: %v\n", c.id, op)
	switch op.opCode {
	case OpCodeAdd:
		var loc, v1, v2 int
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
		var loc, v1, v2 int
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
		var loc int
		// loc logic is reversed
		if op.modes[0] == ModePosition {
			loc = op.params[0]
		} else if op.modes[0] == ModeRelative {
			loc = op.params[0] + c.relPos
		} else {
			loc = c.getMemory(op.params[0])
		}
		c.setMemory(loc, c.input(c))
	case OpCodeOutput:
		if op.modes[0] == ModePosition {
			c.result = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			c.result = c.getMemory(c.relPos + op.params[0])
		} else {
			c.result = op.params[0]
		}
		// fmt.Printf("COMP %d OUTPUT TO CHAIN:\t%d\n", c.id, c.result)
		// fmt.Println("OUTPUT", c.result)
		c.output(c, c.result) // send result to output func
	case OpCodeJumpIfTrue:
		var v1, v2 int
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
		var v1, v2 int
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
		var loc, v1, v2 int
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
		var loc, v1, v2 int
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
		var v1 int
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

func getData(filePath string) []int {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var data []int
	for {
		dataPoint, err := reader.ReadString(',')
		dataPoint = strings.Trim(dataPoint, ",\n")

		u, err := strconv.Atoi(dataPoint)
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
