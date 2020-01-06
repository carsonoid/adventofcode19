package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func main() {
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
}
