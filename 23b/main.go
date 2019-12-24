package main

import (
	"fmt"
	"os"
	"time"
)

type packet struct {
	Addr *int
	X, Y *int
}

var numComputers = 50
var computers = make([]*computer, numComputers)
var natPacket = []int{}
var natYs = map[int]int{}

func main() {
	code := getData(os.Args[1])

	for id := range computers {
		computers[id] = newComputer(id, code,
			func(c *computer) int { // Input
				if len(c.inBuffer) == 0 {
					// fmt.Println(c.id, "INPUT NONE")
					return -1 // empty buffer
				}
				var r int
				r, c.inBuffer = c.inBuffer[0], c.inBuffer[1:] // send first
				fmt.Println(c.id, "INPUT IS", r)
				return r
			},
			func(c *computer, o int) { // Output
				fmt.Println(c.id, "GETTING OUT LOCK")
				c.outLock.Lock()
				defer c.outLock.Unlock()
				c.outBuffer = append(c.outBuffer, o) // Append
				if len(c.outBuffer) == 3 {           // send and reset every 3 outs
					fmt.Println(c.id, "SEND PACKET", c.outBuffer)
					addr := c.outBuffer[0]
					x := c.outBuffer[1]
					y := c.outBuffer[2]
					if addr == 255 { // send to nat
						fmt.Println("NATLOG:", x, y)
						natPacket = []int{x, y}
					} else {
						fmt.Println(c.id, "GETTING IN LOCK FOR", addr)
						computers[addr].inLock.Lock()
						defer computers[addr].inLock.Unlock()
						computers[addr].BufferInput(x, y)
					}
					c.outBuffer = []int{}
				}
			},
		)
		computers[id].BufferInput(id) // first input is the comp id
	}

	for id := range computers {
		go computers[id].start()
	}

	go func() { // start "nat"
		idleCycles := 0
		for {
			time.Sleep(time.Second)
			if allComputersIdle() {
				idleCycles++
			} else {
				idleCycles = 0
			}
			if idleCycles >= 2 && len(natPacket) != 0 {
				fmt.Println("NAT SENDING", natPacket)
				y := natPacket[1]
				natYs[y]++
				if count, ok := natYs[y]; ok && count >= 2 {
					fmt.Println("DOUBLE NAT Y:", y, "TIMES", count)
					os.Exit(1)
				}
				computers[0].BufferInput(natPacket...)
				idleCycles = 0
			} else {
				fmt.Println("NAT IDLE CYCLES:", idleCycles)
			}
		}
	}()

	for i := range computers {
		<-computers[i].quit
	}
}

func allComputersIdle() bool {
	for id := range computers { // if any computer
		if len(computers[id].outBuffer) > 0 || len(computers[id].inBuffer) > 0 { // has anything pending input/output
			fmt.Println("NAT IGNORE", id, "IS BUSY. NETWORK NOT IDLE")
			return false
		}
	}
	return true
}
