package main

import (
	"fmt"
)

type vector struct {
	x, y, z int
}

type moon struct {
	pos, vel vector
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (v *vector) Add(o vector) {
	v.x += o.x
	v.y += o.y
	v.z += o.z
}

func (m *moon) Move() {
	m.pos.Add(m.vel)
}

func (m *moon) PotentialEnergy() int {
	return abs(m.pos.x) + abs(m.pos.y) + abs(m.pos.z)
}

func (m *moon) KineticEnergy() int {
	return abs(m.vel.x) + abs(m.vel.y) + abs(m.vel.z)
}

func ApplyGravity(m1, m2 *moon) {
	if m1.pos.x < m2.pos.x {
		m1.vel.x++
		m2.vel.x--
	} else if m1.pos.x > m2.pos.x {
		m1.vel.x--
		m2.vel.x++
	}
	if m1.pos.y < m2.pos.y {
		m1.vel.y++
		m2.vel.y--
	} else if m1.pos.y > m2.pos.y {
		m1.vel.y--
		m2.vel.y++
	}
	if m1.pos.z < m2.pos.z {
		m1.vel.z++
		m2.vel.z--
	} else if m1.pos.z > m2.pos.z {
		m1.vel.z--
		m2.vel.z++
	}
}

var (
	initialState = []*moon{
		{pos: vector{-14, -4, -11}, vel: vector{0, 0, 0}},
		{pos: vector{-9, 6, -7}, vel: vector{0, 0, 0}},
		{pos: vector{4, 1, 4}, vel: vector{0, 0, 0}},
		{pos: vector{2, -14, -9}, vel: vector{0, 0, 0}},
	}
	moons = []*moon{
		{pos: vector{-14, -4, -11}, vel: vector{0, 0, 0}},
		{pos: vector{-9, 6, -7}, vel: vector{0, 0, 0}},
		{pos: vector{4, 1, 4}, vel: vector{0, 0, 0}},
		{pos: vector{2, -14, -9}, vel: vector{0, 0, 0}},
	}
)

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	var ttx, tty, ttz int

	var step int
	for step = 0; ; step++ {
		// fmt.Printf("After %v steps:\n", step)
		// for _, moon := range moons {
		// 	fmt.Printf("%v\n", *moon)
		// }
		// fmt.Println()
		for i := range moons {
			for j := i; j < len(moons); j++ {
				if i != j {
					// fmt.Printf("Calc Gravity from %d to %d\n", i, j)
					ApplyGravity(moons[i], moons[j])
				}
			}
		}
		for i := range moons {
			moons[i].Move()
		}
		if step != 0 {
			xAligned := false
			for i := range moons {
				if moons[i].pos.x == initialState[i].pos.x && moons[i].vel.x == initialState[i].vel.x {
					xAligned = true
				} else {
					xAligned = false
					break
				}
			}
			if xAligned && ttx == 0 {
				ttx = step + 1
			}

			yAligned := false
			for i := range moons {
				if moons[i].pos.y == initialState[i].pos.y && moons[i].vel.y == initialState[i].vel.y {
					yAligned = true
				} else {
					yAligned = false
					break
				}
			}
			if yAligned && tty == 0 {
				tty = step + 1
			}

			zAligned := false
			for i := range moons {
				if moons[i].pos.z == initialState[i].pos.z && moons[i].vel.z == initialState[i].vel.z {
					zAligned = true
				} else {
					zAligned = false
					break
				}
			}
			if zAligned && ttz == 0 {
				ttz = step + 1
			}

			if ttx != 0 && tty != 0 && ttz != 0 {
				fmt.Printf("TIME TO REPEAT: %d, %d, %d\n", ttx, tty, ttz)
				fmt.Printf("LCM: %d\n", LCM(ttx, tty, ttz))
				break
			}
		}
	}
	fmt.Printf("After %v steps:\n", step)
	for _, moon := range moons {
		fmt.Printf("%v\n", *moon)
	}
	var total int
	for i := range moons {
		pe := moons[i].PotentialEnergy()
		ke := moons[i].KineticEnergy()
		total += pe * ke
	}
	fmt.Printf("Total: %d\n", total)
}
