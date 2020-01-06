package main

import (
	"fmt"
)

var iters = 10

func main() {
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
			fmt.Printf("i: %d, out: %d\n", i, out)
		}(i) // pass i as the parameter (will be passed by copy)
	}

	// wait for all funcs to complete
	for i := 0; i < iters; i++ {
		<-done
	}
}
