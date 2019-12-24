package main

import (
	"fmt"
	"os"
)

type packet struct {
	Addr *int
	X, Y *int
}

var numComputers = 50
var computers = make([]*computer, numComputers)

func main() {
	code := getData(os.Args[1])

	for id := range computers {
		computers[id] = newComputer(id, code,
			func(c *computer) int { // Input
				if len(c.inBuffer) == 0 {
					fmt.Println(c.id, "INPUT NONE")
					return -1 // empty buffer
				}
				var r int
				r, c.inBuffer = c.inBuffer[0], c.inBuffer[1:] // send first
				fmt.Println(c.id, "INPUT IS", r)
				return r
			},
			func(c *computer, o int) { // Output
				c.outBuffer = append(c.outBuffer, o) // Append
				if len(c.outBuffer) == 3 {           // send and reset every 3 outs
					fmt.Println(c.id, "SEND PACKET", c.outBuffer)
					computers[c.outBuffer[0]].BufferInput(c.outBuffer[1:]...)
					c.outBuffer = []int{}
				}
			},
		)
		computers[id].BufferInput(id) // first input is the comp id
		go computers[id].start()
	}

	for i := range computers {
		<-computers[i].quit
	}
}
