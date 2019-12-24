package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Name int
}

func NewStack(a []Card) []Card {
	r := make([]Card, len(a))
	j := 0
	for i := len(a) - 1; i >= 0; i-- {
		r[j] = a[i]
		j++
	}
	return r
}

func CutN(a []Card, n int) []Card {
	if n > 0 {
		return append(a[n:], a[:n]...)
	}
	return append(a[len(a)+n:], a[:len(a)+n]...)
}

func IncN(a []Card, n int) []Card {
	r := make([]Card, len(a))
	r[0] = a[0] // first item does not move
	var pos int
	for i := 1; i < len(a); i++ {
		pos = i * n % len(a)
		r[pos] = a[i]
	}
	return r
}

var numCards = 10007

// var numCards = 10

func main() {
	cards := make([]Card, numCards)
	for i := 0; i < len(cards); i++ {
		cards[i].Name = i
	}

	f, err := os.Open(os.Args[1])
	defer f.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	var instructions []string
	for scanner.Scan() {
		instructions = append(instructions, scanner.Text())
	}

	for _, inst := range instructions {
		if strings.Contains(inst, "deal into new stack") {
			cards = NewStack(cards)
			fmt.Println("NEW")
			continue
		}
		if strings.Contains(inst, "cut") {
			parts := strings.Split(inst, " ")
			n, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				panic(err)
			}
			cards = CutN(cards, n)
			fmt.Println("CUT", n)
			continue
		}
		if strings.Contains(inst, "increment") {
			parts := strings.Split(inst, " ")
			n, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				panic(err)
			}
			cards = IncN(cards, n)
			fmt.Println("INCREMENT", n)
			continue
		}
		panic(inst)
	}

	for pos, c := range cards {
		if c.Name == 2019 {
			fmt.Println("2019 at", pos)
		}
	}
}
