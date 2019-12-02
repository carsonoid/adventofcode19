package main

import "testing"

func Test_getFuelForMass(t *testing.T) {
	tests := []struct {
		mass uint64
		want uint64
	}{
		{0, 0},
		{1, 0},
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, tt := range tests {
		t.Run(string(tt.mass), func(t *testing.T) {
			if got := getFuelForMass(tt.mass); got != tt.want {
				t.Errorf("getFuelForMass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFuelForModuleMass(t *testing.T) {
	tests := []struct {
		mass uint64
		want uint64
	}{
		{2, 0},
		{12, 2},
		{14, 2},
		{1969, 966},
		{100756, 50346},
	}
	for _, tt := range tests {
		t.Run(string(tt.mass), func(t *testing.T) {
			if got := getFuelForModuleMass(tt.mass); got != tt.want {
				t.Errorf("getFuelForModuleMass() = %v, want %v", got, tt.want)
			}
		})
	}
}
