package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
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
}
