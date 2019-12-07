package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputs []int64
var lastOut int64

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
	OpCodeQuit        OpCode = 99
)

type Mode int

const (
	ModePosition  Mode = 0
	ModeImmediate Mode = 1
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
	if op.opCode == OpCodeAdd || op.opCode == OpCodeMultiply {
		op.modes = make([]Mode, 3)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int64(c - '0'))
		}
	}
	if op.opCode == OpCodeInput || op.opCode == OpCodeOutput {
		op.modes = make([]Mode, 1)
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
	if op.opCode == OpCodeLessThan || op.opCode == OpCodeEquals {
		op.modes = make([]Mode, 3)
		// fmt.Printf("%v\n", modeStr)
		// fmt.Printf("%v\n", reverse(modeStr))
		for i, c := range reverse(modeStr) {
			op.modes[i] = getMode(int64(c - '0'))
		}
	}

	switch op.opCode {
	case OpCodeAdd:
		op.params = workSet[:3]
	case OpCodeMultiply:
		op.params = workSet[:3]
	case OpCodeInput:
		op.params = workSet[:1]
	case OpCodeOutput:
		op.params = workSet[:1]
	case OpCodeJumpIfFalse:
		fallthrough
	case OpCodeJumpIfTrue:
		op.params = workSet[:2]
	case OpCodeLessThan:
		fallthrough
	case OpCodeEquals:
		op.params = workSet[:3]
	}

	// fmt.Printf("OP: %v\n", op)
	return op
}

// func getInput() int64 {
// 	var i int64
// 	fmt.Printf("INPUT: ")
// 	fmt.Scan(&i)
// 	fmt.Printf("\n")
// 	return i
// }

func getInput() int64 {
	var x int64
	x, inputs = inputs[0], inputs[1:]
	fmt.Printf("AUTO INPUT: %v\n", x)
	return x
}

func doOperation(op operation, memory []int64) *int64 {
	switch op.opCode {
	case OpCodeAdd:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = memory[op.params[0]]
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = memory[op.params[1]]
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else {
			loc = memory[op.params[2]]
		}
		memory[loc] = v1 + v2
	case OpCodeMultiply:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = memory[op.params[0]]
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = memory[op.params[1]]
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else {
			loc = memory[op.params[2]]
		}
		memory[loc] = v1 * v2
	case OpCodeInput:
		var loc int64
		// loc logic is reversed
		if op.modes[0] == ModePosition {
			loc = op.params[0]
		} else {
			loc = memory[op.params[0]]
		}
		memory[loc] = getInput()
		// memory[loc] = INPUT
	case OpCodeOutput:
		var loc int64
		// loc logic is reversed
		if op.modes[0] == ModePosition {
			loc = op.params[0]
		} else {
			loc = memory[op.params[0]]
		}
		lastOut = memory[loc]
		fmt.Printf("OUTPUT: %d\n", lastOut)
	case OpCodeJumpIfTrue:
		var v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = memory[op.params[0]]
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = memory[op.params[1]]
		} else {
			v2 = op.params[1]
		}
		if v1 != 0 {
			return &v2
		}
	case OpCodeJumpIfFalse:
		var v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = memory[op.params[0]]
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = memory[op.params[1]]
		} else {
			v2 = op.params[1]
		}
		if v1 == 0 {
			return &v2
		}
	case OpCodeLessThan:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = memory[op.params[0]]
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = memory[op.params[1]]
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else {
			loc = memory[op.params[2]]
		}
		if v1 < v2 {
			memory[loc] = 1
		} else {
			memory[loc] = 0
		}
	case OpCodeEquals:
		var loc, v1, v2 int64
		if op.modes[0] == ModePosition {
			v1 = memory[op.params[0]]
		} else {
			v1 = op.params[0]
		}
		if op.modes[1] == ModePosition {
			v2 = memory[op.params[1]]
		} else {
			v2 = op.params[1]
		}
		// loc logic is reversed
		if op.modes[2] == ModePosition {
			loc = op.params[2]
		} else {
			loc = memory[op.params[2]]
		}
		if v1 == v2 {
			memory[loc] = 1
		} else {
			memory[loc] = 0
		}
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

func run(memory []int64) {
	var op operation
	var pos int
	var exit bool
	for {
		fmt.Printf("POS: %d, len: %d\n", pos, len(memory))
		if pos == len(memory) {
			exit = true
		} else if pos+1 == len(memory) {
			op = getOperation(memory[pos], memory[pos+1:pos+2])
			exit = true
		} else if pos+2 == len(memory) {
			op = getOperation(memory[pos], memory[pos+1:pos+2])
			exit = true
		} else if pos+3 == len(memory) {
			op = getOperation(memory[pos], memory[pos+1:pos+3])
			exit = true
		} else if pos+4 == len(memory) {
			op = getOperation(memory[pos], memory[pos+1:pos+4])
			exit = true
		} else {
			op = getOperation(memory[pos], memory[pos+1:pos+4])
		}
		pos += len(op.params) + 1

		if op.opCode == OpCodeQuit {
			fmt.Printf("DONE!\n")
			return
		}

		newPos := doOperation(op, memory)
		if newPos != nil {
			pos = int(*newPos)
			// fmt.Printf("MOVED PTR to %d\n", pos)
		}
		if exit {
			break
		}
	}
	return
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
	fmt.Printf("CODE: %v\n", code)

	outputs := []int64{}

	allPhases := permutations([]int64{0, 1, 2, 3, 4})
	// allPhases := [][]int64{{4, 3, 2, 1, 0}}
	fmt.Printf("phase perms: %v\n", allPhases)

	for _, phasePerm := range allPhases {

		lastOut = 0
		for _, phaseInput := range phasePerm {
			memory := make([]int64, len(code))
			copy(memory, code)
			fmt.Printf("MEMORY: %v\n", memory)

			inputs = []int64{phaseInput, lastOut}
			fmt.Printf("INPUTS: %v\n", inputs)

			run(memory)
		}
		outputs = append(outputs, lastOut)
	}
	fmt.Printf("OUTPUTS: %v\n", outputs)
	fmt.Printf("MAX OUTPUT: %v\n", max(outputs))
}
