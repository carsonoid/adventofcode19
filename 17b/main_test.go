package main

import (
	"reflect"
	"testing"
)

// func Test_operate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    []int
// 		wantErr bool
// 		want    int
// 	}{
// 		{
// 			name: "add",
// 			data: []int{1, 0, 0, 0, 99},
// 			want: 2,
// 		},
// 		{
// 			name: "multi",
// 			data: []int{2, 3, 0, 3, 99},
// 			want: 6,
// 		},
// 		{
// 			name: "multi2",
// 			data: []int{2, 4, 4, 5, 99, 0},
// 			want: 9801,
// 		},
// 		{
// 			name:    "invalidcode",
// 			data:    []int{123, 4, 4, 5, 99, 0},
// 			want:    0,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			workSet := tt.data[:4]
// 			resultPos := workSet[3]
// 			if err := operate(tt.data, workSet); (err != nil) != tt.wantErr {
// 				t.Errorf("operate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			got := tt.data[resultPos]
// 			if got != tt.want {
// 				t.Errorf("operate() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_getOperation(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		workSet []int
		want    operation
	}{
		{
			"default",
			1002,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeMultiply,
				modes: []Mode{
					ModePosition,
					ModeImmediate,
					ModePosition,
				},
				params: []int{4, 3, 4},
			},
		},
		{
			"short",
			02,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeMultiply,
				modes: []Mode{
					ModePosition,
					ModePosition,
					ModePosition,
				},
				params: []int{4, 3, 4},
			},
		},
		{
			"in",
			03,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeInput,
				modes: []Mode{
					ModePosition,
				},
				params: []int{4},
			},
		},
		{
			"out",
			104,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeOutput,
				modes: []Mode{
					ModeImmediate,
				},
				params: []int{4},
			},
		},
		{
			"jit",
			05,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeJumpIfTrue,
				modes: []Mode{
					ModePosition,
					ModePosition,
				},
				params: []int{4, 3},
			},
		},
		{
			"jif",
			6,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeJumpIfFalse,
				modes: []Mode{
					ModePosition,
					ModePosition,
				},
				params: []int{4, 3},
			},
		},
		{
			"lt",
			7,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeLessThan,
				modes: []Mode{
					ModePosition,
					ModePosition,
					ModePosition,
				},
				params: []int{4, 3, 4},
			},
		},
		{
			"gt",
			8,
			[]int{4, 3, 4, 33},
			operation{
				opCode: OpCodeEquals,
				modes: []Mode{
					ModePosition,
					ModePosition,
					ModePosition,
				},
				params: []int{4, 3, 4},
			},
		},
		// {
		// 	"quit",
		// 	99,
		// 	[]int{4, 3, 4, 33},
		// 	operation{
		// 		opCode: OpCodeQuit,
		// 		modes:  []Mode{},
		// 		params: []int{},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOperation(tt.code, tt.workSet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
