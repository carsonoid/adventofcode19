package main

import (
	"reflect"
	"testing"
)

func Test_getPassedPositions(t *testing.T) {

	tests := []struct {
		name         string
		instructions []instruction
		want         []position
	}{
		{
			"right",
			[]instruction{{'R', 3}},
			[]position{{1, 0}, {2, 0}, {3, 0}},
		},
		{
			"left",
			[]instruction{{'L', 3}},
			[]position{{-1, 0}, {-2, 0}, {-3, 0}},
		},
		{
			"up",
			[]instruction{{'U', 3}},
			[]position{{0, 1}, {0, 2}, {0, 3}},
		},
		{
			"down",
			[]instruction{{'D', 3}},
			[]position{{0, -1}, {0, -2}, {0, -3}},
		},
		{
			"line1",
			[]instruction{{'R', 8}, {'U', 5}, {'L', 5}, {'D', 3}},
			[]position{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {7, 0}, {8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {7, 5}, {6, 5}, {5, 5}, {4, 5}, {3, 5}, {3, 4}, {3, 3}, {3, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPassedPositions(tt.instructions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPassedPositions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClosestWireCrossDistance(t *testing.T) {
	tests := []struct {
		name  string
		inst1 []instruction
		inst2 []instruction
		want  int64
	}{
		{
			"set1",
			[]instruction{{'R', 8}, {'U', 5}, {'L', 5}, {'D', 3}},
			[]instruction{{'U', 7}, {'R', 6}, {'D', 4}, {'L', 4}},
			6,
		},
		{
			"set2",
			[]instruction{{'R', 75}, {'D', 30}, {'R', 83}, {'U', 83}, {'L', 12}, {'D', 49}, {'R', 71}, {'U', 7}, {'L', 72}},
			[]instruction{{'U', 62}, {'R', 66}, {'U', 55}, {'R', 34}, {'D', 71}, {'R', 55}, {'D', 58}, {'R', 83}},
			159,
		},
		{
			"set3",
			[]instruction{{'R', 98}, {'U', 47}, {'R', 26}, {'D', 63}, {'R', 33}, {'U', 87}, {'L', 62}, {'D', 20}, {'R', 33}, {'U', 53}, {'R', 51}},
			[]instruction{{'U', 98}, {'R', 91}, {'D', 20}, {'R', 16}, {'D', 67}, {'R', 40}, {'U', 7}, {'R', 15}, {'U', 6}, {'R', 7}},
			135,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClosestWireCrossDistance(tt.inst1, tt.inst2); got != tt.want {
				t.Errorf("getClosestWireCrossDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}
