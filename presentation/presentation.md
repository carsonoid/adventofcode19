---
marp: true
theme: default
paginate: true
#footer: ![width:1em](img/github.png) carsonoid | ![width:1em](img/twit.png) carson_ops
---
<!-- 
_paginate: false
_footer: ""
-->

<!-- Scoped style -->
<style scoped>
h1,h2 {
	text-align: center;
}
h1 {
  color: black;
  text-decoration: underline;
  font-size: 3em;
}
h2 {
  color: black;
  font-size: 2em;
}
</style>

# Advent of Code 19
## _Takeaways and Go Behaviors_
## Carson Anderson - Weave
## ![width:1em](img/github.png) carsonoid &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; ![width:1em](img/twit.png) carson_ops

---
# Day 1 - Integer division

> Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for > a module, take its mass, divide by three, **round down**, and subtract 2.

---

```golang
var mass = 298

func main() {
	fmt.Println("// Non-even division")
	fmt.Println(float64(mass) / 3)

	fmt.Println("// Round down via math.Floor")
	fmt.Println(math.Floor(float64(mass) / 3))

	fmt.Println("// Oh... int division in go rounds down by default :D")
	fmt.Println(mass / 3)
}

/* OUTPUT:
// Non-even division
99.33333333333333
// Round down via math.Floor
99
// Oh... int division in go rounds down by default :D
99
*/
```

Playground: https://play.golang.org/p/JfTkBHYdxeO

---

# Day 2 - Intcode: Scanners vs Readers

> Intcode programs are given as a list of integers

```
1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,13,1,19,1,10,19,23,1,23,9,27,1,5,27,31,2,31,13,35,1,35,5,39,1,39,5,43,2,...
```

The input for this day was a single line. And could potentially get very, very long. So it was important to understand the difference between readers and scanners.

---

### Scanners

```golang
// Make a test reader with a bad line
file := strings.NewReader("" + // Empty add to make fmt clean
    "1\n2\n3\n" +
    strings.Repeat("4", bufio.MaxScanTokenSize+1) +
    "\n5\n6\n",
)

// Scanners: 
// * Read up to a newline
// * Don't include the newline in the text
// * Have a default line length limit of 64K (Can be increased manually)
//    * Will return an empty line if it's too long
//    * Check for scanner.Err() == bufio.ErrTooLong to catch this
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}
fmt.Println(scanner.Err())
```

Playground: https://play.golang.org/p/KOgOzt_-ROV

---

### Readers

```golang
// Make a test reader with a bad line
file := strings.NewReader("" + // Empty add to make fmt clean
    "1\n2\n3\n" +
    strings.Repeat("4", bufio.MaxScanTokenSize+1) +
    "\n5\n6",
)

// Readers:
// * Can read many ways: all, byte, bytes, rune, string, slice, line
// * May or may not include delimiters, depeding on read method so check docs
reader := bufio.NewReader(file)
for {
    text, err := reader.ReadString(byte('\n'))
    text = strings.TrimRight(text, "\n")
    fmt.Println(text)

    if err != nil {
        if errors.Is(err, io.EOF) {
            break
        }
        panic(err)
    }
}
```

Playground: https://play.golang.org/p/fcaBp7G8L5D

---

### Result

Final intcode parsing function

```golang
func getIntCode(filePath string) []uint64 {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var data []uint64
	for {
		dataPoint, err := reader.ReadString(',')
		dataPoint = strings.Trim(dataPoint, ",")

		u, err := strconv.ParseUint(dataPoint, 10, 64)
		if err != nil {
			break
		}
		data = append(data, u)
	}

	return data
}
```

---

# Day 3

> For example, if the first wire's path is R8,U5,L5,D3,
> then starting from the central port (o), it goes right 8, up 5, left 5, and finally down 3:

More parsing fun! The first step was to convert lines of instruction strings to slices of instructions. Reading the input by lines and further splitting by strings was easy! But how to convert `R75` to a struct? Runes to the rescue!

---

```golang
type instruction struct {
	direction rune
	distance  int64
}

func getInstructionsFromString(in string) []instruction {
	stringInstructions := strings.Split(in, ",")
	var instructions []instruction
	for _, sInst := range stringInstructions {
		runes := []rune(sInst)
		direction := runes[0]
		distance, err := strconv.ParseInt(sInst[1:], 10, 64)
		if err != nil {
			break
		}
		instructions = append(instructions, instruction{
			direction: direction,
			distance:  distance,
		})
	}
	return instructions
}
```

---

Be wary of the default stringification of runes though...

```golang
// R75,D30,R83,U83,L12,D49,R71,U7,L72
[]main.instruction{
		main.instruction{direction: 82, distance: 75},
		main.instruction{direction: 68, distance: 30},
		main.instruction{direction: 82, distance: 83},
		main.instruction{direction: 85, distance: 83},
		main.instruction{direction: 76, distance: 12},
		main.instruction{direction: 68, distance: 49},
		main.instruction{direction: 82, distance: 71},
		main.instruction{direction: 85, distance: 7},
		main.instruction{direction: 76, distance: 72},
	}
```

Playground: https://play.golang.org/p/CEtf-SBKppL

---

# Day 4 - Irregular Regex

Find the password that fits a given set of rules. Most  were easy, but...

> Two adjacent digits are the same (like 22 in 122345).

I tried to use a regex like `(\d)\1` to test for doubles. Suddenly, no compilation:
```
unknown escape sequence
```

---

It turns out that Go uses the RE2 engine, which doesn't support backreferences. So you are not able to to do anything like: `(\d)\1` to test for doubles.

Reasoning from https://swtch.com/~rsc/regexp/regexp3.html. Emhpasis added.

> RE2 disallows PCRE features that cannot be implemented efficiently using automata. **(The most notable such feature is backreferences.)** In return for giving up these difficult to implement (and often incorrectly used) features, RE2 can provably analyze the regular expressions or the automata. We've already seen examples of analysis for use in RE2 itself, in the DFA's use of memchr and in the analysis of whether a regular expression is one-pass. **RE2 can also provide analyses that let higher-level applications speed searches.**

**Translation: No backreferences or other complex things so we can ensure a linear parse time.**

---

# Day 7 - Amplifiers In A Row

More use of the intcode computer. But this time the problem required more than one to run. And for some fun connections:

> There are five amplifiers connected in series; each one receives an input signal and produces an output signal. They are connected such that the first amplifier's output leads to the second amplifier's input, ...

```
    O-------O  O-------O  O-------O  O-------O  O-------O
0 ->| Amp A |->| Amp B |->| Amp C |->| Amp D |->| Amp E |-> (to thrusters)
    O-------O  O-------O  O-------O  O-------O  O-------O
```

So run 5 computers in a row, passing the output of one into the input of another. Easy!

---

## A Wild Part 2 Appears:

> The Elves quickly talk you through rewiring the amplifiers into a feedback loop:

```
      O-------O  O-------O  O-------O  O-------O  O-------O
0 -+->| Amp A |->| Amp B |->| Amp C |->| Amp D |->| Amp E |-.
   |  O-------O  O-------O  O-------O  O-------O  O-------O |
   |                                                        |
   '--------------------------------------------------------+
                                                            |
                                                            v
                                                     (to thrusters)
```

---

# Channels!

With channels, we can actually "wire" the amps together in code and let them run at their own pace concurrently. Not a requirement, but pretty cool!

---

```golang
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
```

---

Make Amps:

```golang
// Setup amplifiers
amps := make([]*amplifier, numAmps)
for i := 0; i < numAmps; i++ {
    amps[i] = newAmplifier(i, code, phasePerm[i])
}

```

---

"Wire" amps together

```golang
// Link inputs to outputs
for i := 0; i < numAmps; i++ {
    if i == 0 { // link first amp input to last amp output
        amps[i].input = amps[len(amps)-1].output
    } else { // link to previous amp input to current amp output
        amps[i].input = amps[i-1].output
    }
}
```

---

Start the amps and wait for completion

```golang
// Start amps
for i := 0; i < numAmps; i++ {
    go amps[i].start()
}

// Send start signal to first amp to begin processing
amps[0].input <- 0

// Wait for amps to close
for i := 0; i < numAmps; i++ {
    <-amps[i].quit
}
```

---

# Day 8 - Runes to ints and Unicode!

Render an image sent to you.

> Images are sent as a series of digits that each represent the color of a single pixel. The digits fill each row of the image left-to-right, then move downward to the next row, filling rows top-to-bottom until every pixel of the image is filled.

### Runes

Another case where inputs could be very long lines. Time to use a reader! One rune at a time. (Hindsight: probably should have used ReadByte or converted to a string and then int for proper error checking)


---

This code converts a rune for 0-9 to an int. But it will fail silently for runes out of range.

```golang
r, _, err := reader.ReadRune()
if err != nil {
    panic(err)
} else if r == '\n' {
    break
}
pixel := int(r - '0') 
```


---

```golang
pixel := int(r - '0') 
```

Uhh. what is that thing? Turns out it is just convenient ascii math.


CODE | CHAR |   | CODE | CHAR
-----|------|---|------|-----
48   | 0    |   | 53   | 5
49   | 1    |   | 54   | 6
50   | 2    |   | 55   | 7
51   | 3    |   | 56   | 8
52   | 4    |   | 57   | 9

---


```golang
fmt.Printf("%d\n", '1')                              // Prints ascii code for '1'
fmt.Printf("%d - %d = %d\n", '1', '0', int('1'-'0')) // Prints expected int
fmt.Printf("%d - %d = %d\n", 'A', '0', int('A'-'0')) // Silent failure of logic

// convert allowed runes
digit, err := strconv.Atoi(string('1'))
if err != nil {
    panic(err)
}
fmt.Printf("%d\n", digit)

// properly catch bad runes
digit, err = strconv.Atoi(string('A'))
if err != nil {
    fmt.Printf("error: '%v' cannot be converted to a digit", string('A'))
}

/* OUTPUT:
49
49 - 48 = 1
65 - 48 = 17
1
error: 'A' cannot be converted to a digit
*/
```

Playground: https://play.golang.org/p/mKCo7SsOF2l

---

### Also: Unicode

Cleaner visuals!

```golang
	for _, pixel := range rendered { // loop over a slice of ints.
		curWidth++
		switch pixel {
		case 0:
			fmt.Printf("\u2591") // Shaded block
		case 1:
			fmt.Printf("\u2588") // Full Block
		}
		if curWidth%imgWidth == 0 {
			fmt.Printf("\n")
		}
    }

/*
████░█░░░░███░░░░██░████░
░░░█░█░░░░█░░█░░░░█░█░░░░
░░█░░█░░░░███░░░░░█░███░░
░█░░░█░░░░█░░█░░░░█░█░░░░
█░░░░█░░░░█░░█░█░░█░█░░░░
████░████░███░░░██░░█░░░░
*/
```


---

### About that output....

```
████░█░░░░███░░░░██░████░
░░░█░█░░░░█░░█░░░░█░█░░░░
░░█░░█░░░░███░░░░░█░███░░
░█░░░█░░░░█░░█░░░░█░█░░░░
█░░░░█░░░░█░░█░█░░█░█░░░░
████░████░███░░░██░░█░░░░
```

* It wasn't until a later day that I finally realized why I couldn't find a unicode character for an empty block...

---

```
███
██  <--- It's a called a "space". I am dumb.
███
```

---

# Day 5 - Custom Types for Compilation Safety

> ... you'll need to add support for parameter modes:
> 
> Each parameter of an instruction is handled based on its parameter mode. ...
> Parameter modes are stored in the same value as the instruction's opcode. The opcode is a two-digit number based only on the ones and tens digit of the value, that is, the opcode is the rightmost two digits of the first value in an instruction. Parameter modes are single digits, one per parameter, read right-to-left from the opcode:

---

> The first instruction, 1002,4,3,4, is a multiply instruction - the rightmost two digits of the first value, 02, indicate opcode 2, multiplication. Then, going right to left, the parameter modes are 0 (hundreds digit), 1 (thousands digit), and 0 (ten-thousands digit, not present and therefore zero):

```
ABCDE
 1002

DE - two-digit opcode,      02 == opcode 2
 C - mode of 1st parameter,  0 == position mode
 B - mode of 2nd parameter,  1 == immediate mode
 A - mode of 3rd parameter,  0 == position mode,
                                  omitted due to being a leading zero
```

---

So now intcode has 3 different "kinds" of `int`: 

- values
- operational codes
- parameter modes.

But they are all `int` to the compiler by default. That means we open ourselves up to cases where the default type checking in the compiler will fail to find issues at compile time. Let's fix it with custom types!

---

```golang
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

func getOpCode(m int64) OpCode {
	switch m {
	case 1:
		return OpCodeAdd
	case 2:
		return OpCodeMultiply
	...
}
```
---

```golang
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

func getMode(m int64) Mode {
	switch m {
	case 0:
		return ModePosition
	case 1:
		return ModeImmediate
	...
}
```

Now we can catch more errors at compile time and get more readable code to boot!

---

### Side Note

 You can define methods for custom types

```golang
func (m Mode) Print() {
	fmt.Println(m)
}

func main() {
	m := Mode(2)
	m.Print()
}
/* OUTPUT:
main.Mode
2
*/
```

Playground: https://play.golang.org/p/eXchfaKqcBg

---

# Day 11 - First "Screen" to draw

Graphics libraries have 0,0 as the top,left corner for a reason! It's hard to print rows of pixels starting from the center... So you have to store your "screen" as a slice of slice and print top-to-bottom and left-to-right from zero both ways.

```golang
func print(screen [][]int) {
	for y := range screen {
		for x := range screen[row] {
			fmt.println(screen[y][x])
		}
	}
}
```

---

# Day 13 - No good graphical rendering in go.

Lots of printing out the whole map and then clearing the screen with some special chars.

```golang
var screen [][]int64

func Draw(c *computer) {
	frames := 0
	for range c.print { // print on demand until channel is closed
		os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
		frames++
		fmt.Println("FRAMES", frames)
		for y := range screen {
			for _, tileID := range screen[y] {
				fmt.Printf(GetTile(tileID))
			}
			fmt.Println()
		}
		fmt.Println("SCORE:", score)
	}
}
```

---

# Day 21 & 23 - Fun With Functions

Refactored computer to use funcs instead of channels for input/ouput.

```golang
type inputFunc func() int
type outputFunc func(o int)

type computer struct {
	input     inputFunc
	output    outputFunc
	...
}

func newComputer(id int, code []int, input inputFunc, output outputFunc) *computer {...}

c := newComputer(0, code,
	func() int {
		return 0
	},
	func(o int) {
		fmt.Println(o)
	},
)

go c.start()
```

---

Or what my actual Day 23 code did:

```golang
computers := make([]*computer, numComputers)

for id := range computers {
	computers[id] = newComputer(id, code,
			func(c *computer) int { // Input
				...
				// Problems came here
			},
			func(c *computer, o int) { // Output
				...
				// Problems came here
			},
		)
}

for id := range computers {
	go computers[id].start()
}
```

---

Contrived example to illustrate problem:

```golang
iters := 10

for i := 0; i < iters; i++ {
	fmt.Printf("i: %d\n", i)
}
/*
i: 0
i: 1
i: 2
i: 3
i: 4
i: 5
i: 6
i: 7
i: 8
i: 9
*/
```

Playground: https://play.golang.org/p/LHExek7g1E_g

---

```golang
iters := 10

for i := 0; i < iters; i++ {
	go func() {
		fmt.Printf("i: %d\n", i)
	}()
}
/*
Program exited.
*/
```

Playground: https://play.golang.org/p/3e-vEr1ZTIe

---


```golang
iters := 10

// Make a channel to wait for all goroutines to complete
done := make(chan struct{})

for i := 0; i < iters; i++ {
	go func() {
		 // Defer a send to the done chanel to make sure panics don't cause problems
        defer func() {
            done <- struct{}{}
		}()

		fmt.Printf("i: %d\n", i)
	}()
}
for i := 0; i < iters; i++ {
    <-done
}
/*
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
*/
```

Playground: https://play.golang.org/p/3e-vEr1ZTIe

---

Final, working example:

```golang
iters := 10

// Make a channel to wait for all goroutines to complete
done := make(chan struct{}, iters) // fully buffer done channel for zero memory cost

for i := 0; i < iters; i++ {
    // Start the goroutine
    go func(out int) {

        // Defer a send to the done chanel to make sure panics don't cause problems
        defer func() {
            done <- struct{}{}
        }()

        // Print the "i" loop variable, and the "out" function parameter value
        fmt.Printf("out: %d\n", out)
    }(i) // pass i as the parameter (will be passed by copy)
}

// wait for all funcs to complete. 
// Note that it's 100% possible that all funcs are done before we even get to 
// this part of the code. Thanks to the fully buffered channel.
for i := 0; i < iters; i++ {
    <-done
}
```

Playground: https://play.golang.org/p/h50Vu83DZNk

---

After I did all the work with channels and returns.  I learned about waitgroups... https://gobyexample.com/waitgroups

```golang
func worker(id int, wg *sync.WaitGroup) {
    fmt.Printf("Worker %d starting\n", id)

    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)

    wg.Done()
}

func main() {

    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
}
```

---

## Language Obscura - Empty structs are neat!

Empty structs take up **zero** memory. And are great for buffered channels!

```golang
var iters = 10

// Chans of empty structs take up zero bytes of mem!
done := make(chan struct{}, iters) // buffer iters
elemSize := unsafe.Sizeof(struct{}{})
fmt.Printf("done size in memory: %d\n", unsafe.Sizeof(done))
fmt.Printf("done buffer size in memory: %d\n", uint64(iters)*uint64(elemSize))

// Compare to a chan of booleans
doneBool := make(chan bool, iters) // buffer iters as boolean
elemSize = unsafe.Sizeof(true)
fmt.Printf("doneBool size in memory: %d\n", unsafe.Sizeof(doneBool))
fmt.Printf("doneBool buffer size in memory: %d\n", uint64(iters)*uint64(elemSize))

// done size in memory: 4
// done buffer size in memory: 0
// doneBool size in memory: 4
// doneBool buffer size in memory: 10
```

Playground: https://play.golang.org/p/SQjL0fH63wi

---
<!-- 
_paginate: false
_footer: ""
-->

<!-- Scoped style -->
<style scoped>
h1,h2,h3 {
	text-align: center;
}
h1 {
  color: black;
  text-decoration: underline;
  font-size: 3em;
}
h2 {
  color: black;
  font-size: 2em;
}
h3 {
  color: black;
  font-size: 1.5em;
}
</style>

# Questions?
## Carson Anderson - Weave
## ![width:1em](img/github.png) carsonoid &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; ![width:1em](img/twit.png) carson_ops
## Markdown To Slides: https://marp.app
### https://github.com/carsonoid/adventofcode19
