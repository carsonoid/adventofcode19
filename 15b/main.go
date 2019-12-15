package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jmcvetta/randutil"
)

type computer struct {
	id        int
	memory    []int64
	input     chan int64
	output    chan int64
	result    int64
	print     chan struct{}
	quit      chan struct{}
	relPos    int64
	color     int64
	cout      []int64
	lastMoved int64
	lastState State
	forceMove Move
	x, y      int64
	neighbors []point
	choices   []randutil.Choice
}

func newComputer(id int, code []int64) *computer {
	c := computer{
		id:        id,
		memory:    make([]int64, len(code)),
		input:     make(chan int64, 1),
		output:    make(chan int64, 1),
		print:     make(chan struct{}, 1),
		quit:      make(chan struct{}),
		neighbors: make([]point, 4),
		choices: []randutil.Choice{
			{Weight: 0, Item: 1}, //, "NORTH"
			{Weight: 0, Item: 2}, //, "SOUTH"
			{Weight: 0, Item: 3}, //, "WEST"
			{Weight: 0, Item: 4}, //, "EAST"
		},
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

func readChar() string {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	var b []byte = make([]byte, 1)

	os.Stdin.Read(b)
	return string(b)
}

func getConsoleInput() int64 {
	// reader := bufio.NewReader(os.Stdin)
	for {
		// b, _ := reader.ReadString('\n')
		// b = strings.Trim(b, "\n")
		b := readChar()
		switch b {
		case "w":
			return int64(MoveNorth)
		case "a":
			return int64(MoveWest)
		case "s":
			return int64(MoveSouth)
		case "d":
			return int64(MoveEast)
		}
	}
}

func (c *computer) Draw() {
	c.print <- struct{}{}
}

func (c *computer) getInput() int64 {
	if c.forceMove != 0 {
		return int64(c.forceMove)
	}
	// time.Sleep(time.Second / 10)
	// in := getConsoleInput()
	if c.lastState == StateMoved {
		c.choices[c.lastMoved-1].Weight += 10 // keep going on move
	}

	for i, a := range c.neighbors {
		switch a.tile {
		case TileWall: // Never go toward walls walls
			c.choices[i].Weight = 0
		case TileUnknown: // Always go to unknown
			c.lastMoved = int64(i + 1)
			return int64(i + 1)
		default:

			c.choices[i].Weight = 1 + int(a.weightAdj) // Randomize based on weight and history
		}
	}

	allZero := true
	for i := range c.choices {
		if c.choices[i].Weight <= 0 {
			c.choices[i].Weight = 0
		}
		if c.choices[i].Weight > 0 {
			allZero = false
		}
	}
	if allZero { // prevent all zero
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		c.choices[r1.Intn(4)].Weight = 1
	}

	result, err := randutil.WeightedChoice(c.choices)
	if err != nil {
		panic(err)
	}

	in := int64(result.Item.(int))

	c.lastMoved = in
	return in
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
		c.setMemory(loc, <-c.input)
	case OpCodeOutput:
		if op.modes[0] == ModePosition {
			c.result = c.getMemory(op.params[0])
		} else if op.modes[0] == ModeRelative {
			c.result = c.getMemory(c.relPos + op.params[0])
		} else {
			c.result = op.params[0]
		}
		c.output <- c.result
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
			fmt.Printf("COMP %v QUIT CODE DURING RUN!\n", c.id)
			c.shutdown()
			continue
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

var screen [][]point
var xDim int64 = 50
var yDim int64 = 50
var score int64
var crumbs = []point{}

type point struct {
	tile      Tile
	weightAdj int64
	x, y      int64
}

func main() {
	code := getData(os.Args[1])
	// fmt.Printf("CODE: %v\n", code)

	screen = make([][]point, yDim)
	for y := int64(0); y < yDim; y++ {
		screen[y] = make([]point, xDim)
		for x := int64(0); x < xDim; x++ {
			screen[y][x] = point{x: x, y: y}
		}
	}

	c := newComputer(0, code)
	c.x = int64(xDim / 2)
	c.y = int64(yDim / 2)

	screen[c.y][c.x].tile = TileDroid
	screen[37][39].tile = TileDest

	go c.start()
	go Draw(c)
	c.Draw()

	c.neighbors[0].tile = TileUnknown
	c.neighbors[1].tile = TileUnknown
	c.neighbors[2].tile = TileUnknown
	c.neighbors[3].tile = TileUnknown

	c.input <- c.getInput()
	for state := range c.output {
		screen[c.y][c.x].weightAdj--
		x2, y2 := getMovedCoords(c)
		switch State(state) {
		case StateHitWall: // Hit wall
			screen[y2][x2].tile = TileWall
			c.lastState = StateHitWall
			// Go backward in crumb
			if len(crumbs) > 0 {
				prev := crumbs[len(crumbs)-1]
				if c.x == prev.x && c.y == prev.y { // if on prev tile, pop off crumb
					c.forceMove = 0 // allow normal movement
				} else { // go backwards to previous
					if prev.x < c.x {
						c.forceMove = MoveWest
					}
					if prev.x > c.x {
						c.forceMove = MoveEast
					}
					if prev.y < c.y {
						c.forceMove = MoveNorth
					}
					if prev.y > c.y {
						c.forceMove = MoveSouth
					}
					crumbs = crumbs[:len(crumbs)-1]
				}
			}
		case StateMoved: // Moved
			screen[y2][x2].tile = TileDroid
			screen[c.y][c.x].tile = TileEmpty
			c.x = x2
			c.y = y2
			c.lastState = StateMoved
			crumbs = append(crumbs, screen[c.y][c.x])
		case StateAtDest: // At Dest
			screen[y2][x2].tile = TileDroid
			screen[c.y][c.x].tile = TileDest
			c.x = x2
			c.y = y2
			c.lastState = StateAtDest
			crumbs = append(crumbs, screen[c.y][c.x])
		}
		c.neighbors[0] = screen[c.y-1][c.x]
		c.neighbors[1] = screen[c.y+1][c.x]
		c.neighbors[2] = screen[c.y][c.x-1]
		c.neighbors[3] = screen[c.y][c.x+1]
		c.Draw()
		time.Sleep(time.Second / 200)
		c.input <- c.getInput()
	}
}

type State int

const (
	StateHitWall State = 0
	StateMoved   State = 1
	StateAtDest  State = 2
)

func getMovedCoords(c *computer) (int64, int64) {
	var x2, y2 int64
	switch Move(c.lastMoved) {
	case MoveNorth:
		x2 = c.x
		y2 = c.y - 1
	case MoveSouth:
		x2 = c.x
		y2 = c.y + 1
	case MoveWest:
		x2 = c.x - 1
		y2 = c.y
	case MoveEast:
		x2 = c.x + 1
		y2 = c.y
	}
	return x2, y2
}

func Draw(c *computer) {
	frames := 0
	for range c.print {
		os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
		frames++
		for y := range screen {
			for _, p := range screen[y] {
				fmt.Printf(GetTile(p.tile))
			}
			fmt.Println()
		}
		fmt.Println("POS:", c.x, c.y, screen[c.y][c.x].weightAdj)
		fmt.Println("LAST STATE", c.lastState)
		fmt.Println("LAST MOVED", c.lastMoved)
		fmt.Println("CHOICES", c.choices)
		// fmt.Println("CRUMBS", crumbs)
	}
}

type Move int64

const (
	MoveNorth Move = 1
	MoveSouth Move = 2
	MoveWest  Move = 3
	MoveEast  Move = 4
)

type Tile int64

const (
	TileUnknown Tile = 0
	TileWall    Tile = 1
	TileEmpty   Tile = 2
	TileDroid   Tile = 3
	TileDest    Tile = 4
)

func GetTile(v Tile) string {
	switch v {
	case TileUnknown:
		return fmt.Sprintf("?") // unknown
	case TileWall:
		return fmt.Sprintf("\u2588") // wall
	case TileEmpty:
		return fmt.Sprintf(" ") // empty
	case TileDroid:
		return fmt.Sprintf("D") // droid
	case TileDest:
		return fmt.Sprintf("X") // droid
	default:
		return "?"
	}
}
