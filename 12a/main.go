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
	moons = []*moon{
		{pos: vector{-14, -4, -11}, vel: vector{0, 0, 0}},
		{pos: vector{-9, 6, -7}, vel: vector{0, 0, 0}},
		{pos: vector{4, 1, 4}, vel: vector{0, 0, 0}},
		{pos: vector{2, -14, -9}, vel: vector{0, 0, 0}},
	}
)

func main() {
	var step int
	for step = 0; step < 1000; step++ {
		fmt.Printf("After %v steps:\n", step)
		for _, moon := range moons {
			fmt.Printf("%v\n", *moon)
		}
		fmt.Println()
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
