package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type amplifier struct {
	id       int
	memory   []int64
	input    chan int64
	output   chan int64
	phase    int64
	phaseSet bool
	result   int64
	quit     chan struct{}
}

func newAmplifier(id int, code []int64, phase int64) *amplifier {
	a := amplifier{
		id:     id,
		memory: make([]int64, len(code)),
		phase:  phase,
		output: make(chan int64, 1),
		quit:   make(chan struct{}),
	}
	copy(a.memory, code)
	return &a
}

const numAmps = 5

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

func (a *amplifier) shutdown() {
	close(a.output)
	close(a.quit)
}

func (a *amplifier) getInput() int64 {
	// first input request on run returns the phase, otherwise return the input
	if !a.phaseSet {
		a.phaseSet = true
		// fmt.Printf("AMP %d INPUT PHASE: %v\n", a.id, a.phase)
		return a.phase
	}
	i := <-a.input
	// fmt.Printf("AMP %d INPUT FROM CHAIN:\t%v\n", a.id, i)
	return i
}

func (a *amplifier) doOperation(op operation, memory []int64) *int64 {
	// fmt.Printf("AMP %d DO OP: %v\n", a.id, op)
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
		memory[loc] = a.getInput()
		// memory[loc] = INPUT
	case OpCodeOutput:
		var loc int64
		// loc logic is reversed
		if op.modes[0] == ModePosition {
			loc = op.params[0]
		} else {
			loc = memory[op.params[0]]
		}
		a.result = memory[loc] // store result
		// fmt.Printf("AMP %d OUTPUT TO CHAIN:\t%d\n", a.id, a.result)
		a.output <- a.result // send result to output chan
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

func (a *amplifier) start() {
	var op operation
	var pos int
	var quit bool
	for {
		// fmt.Printf("AMP %v POS: %d, len: %d\n", a.id, pos, len(a.memory))
		if pos == len(a.memory) {
			quit = true
		} else if pos+1 == len(a.memory) {
			op = getOperation(a.memory[pos], a.memory[pos+1:pos+2])
			quit = true
		} else if pos+2 == len(a.memory) {
			op = getOperation(a.memory[pos], a.memory[pos+1:pos+2])
			quit = true
		} else if pos+3 == len(a.memory) {
			op = getOperation(a.memory[pos], a.memory[pos+1:pos+3])
			quit = true
		} else if pos+4 == len(a.memory) {
			op = getOperation(a.memory[pos], a.memory[pos+1:pos+4])
			quit = true
		} else {
			op = getOperation(a.memory[pos], a.memory[pos+1:pos+4])
		}
		pos += len(op.params) + 1

		if quit { // do final operation
			// fmt.Printf("AMP %v HIT END OF CODE\n", a.id)
			// fmt.Printf("AMP %v LAST OP:%#v\n", a.id, op)
			a.doOperation(op, a.memory)
			a.shutdown()
			return
		}

		if op.opCode == OpCodeQuit { // Quit immediately
			// fmt.Printf("AMP %v QUIT CODE DURING RUN!\n", a.id)
			a.shutdown()
			return
		}

		newPos := a.doOperation(op, a.memory)
		if newPos != nil {
			pos = int(*newPos)
			// fmt.Printf("MOVED PTR to %d\n", pos)
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
	// fmt.Printf("CODE: %v\n", code)

	allPhases := permutations([]int64{5, 6, 7, 8, 9})
	// allPhases := [][]int64{{5, 7, 6, 9, 8}}
	// fmt.Printf("phase perms: %v\n", allPhases)

	outputs := []int64{}

	for _, phasePerm := range allPhases {
		// fmt.Printf("CUR_PHASE_PERM: %v\n", phasePerm)

		// Setup amplifiers
		amps := make([]*amplifier, numAmps)
		for i := 0; i < numAmps; i++ {
			amps[i] = newAmplifier(i, code, phasePerm[i])
		}

		// Link inputs to outputs
		for i := 0; i < numAmps; i++ {
			if i == 0 { // link first amp input to last amp output
				// fmt.Printf("AMP %v input tied to AMP %v output\n", i, len(amps)-1)
				amps[i].input = amps[len(amps)-1].output
			} else { // link to previous amp input to current amp output
				// fmt.Printf("AMP %v input tied to AMP %v output\n", i, i-1)
				amps[i].input = amps[i-1].output
			}
		}

		// Start amps
		for i := 0; i < numAmps; i++ {
			go amps[i].start()
		}

		// Send start signal to first amp to begin processing
		amps[0].input <- 0

		// time.Sleep(10 * time.Second)
		// Wait for amps to close
		for i := 0; i < numAmps; i++ {
			// fmt.Printf("WAIT AMP %v to quit\n", i)
			<-amps[i].quit
		}

		// Process result
		output := amps[len(amps)-1].result
		outputs = append(outputs, output)
	}

	// fmt.Printf("OUTPUTS: %v\n", outputs)
	fmt.Printf("MAX OUTPUT: %v\n", max(outputs))
}
