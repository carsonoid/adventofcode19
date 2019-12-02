package main

import "testing"

func Test_calculateModuleFuel(t *testing.T) {
	tests := []struct {
		mass uint64
		want uint64
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, tt := range tests {
		t.Run(string(tt.mass), func(t *testing.T) {
			if got := calculateModuleFuel(tt.mass); got != tt.want {
				t.Errorf("calculateModuleFuel() = %v, want %v", got, tt.want)
			}
		})
	}
}
