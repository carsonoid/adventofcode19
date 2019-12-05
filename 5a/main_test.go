package main

import (
	"reflect"
	"testing"
)

// func Test_operate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    []int64
// 		wantErr bool
// 		want    int64
// 	}{
// 		{
// 			name: "add",
// 			data: []int64{1, 0, 0, 0, 99},
// 			want: 2,
// 		},
// 		{
// 			name: "multi",
// 			data: []int64{2, 3, 0, 3, 99},
// 			want: 6,
// 		},
// 		{
// 			name: "multi2",
// 			data: []int64{2, 4, 4, 5, 99, 0},
// 			want: 9801,
// 		},
// 		{
// 			name:    "invalidcode",
// 			data:    []int64{123, 4, 4, 5, 99, 0},
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
		code    int64
		workSet []int64
		want    operation
	}{
		{
			"default",
			1002,
			[]int64{4, 3, 4, 33},
			operation{
				opCode: OpCodeMultiply,
				modes: []Mode{
					ModePosition,
					ModeImmediate,
					ModePosition,
				},
				params: []int64{4, 3, 4},
			},
		},
		{
			"short",
			02,
			[]int64{4, 3, 4, 33},
			operation{
				opCode: OpCodeMultiply,
				modes: []Mode{
					ModePosition,
					ModePosition,
					ModePosition,
				},
				params: []int64{4, 3, 4},
			},
		},
		{
			"in",
			03,
			[]int64{4, 3, 4, 33},
			operation{
				opCode: OpCodeInput,
				modes: []Mode{
					ModePosition,
				},
				params: []int64{4},
			},
		},
		{
			"out",
			104,
			[]int64{4, 3, 4, 33},
			operation{
				opCode: OpCodeOutput,
				modes: []Mode{
					ModeImmediate,
				},
				params: []int64{4},
			},
		},
		// {
		// 	"quit",
		// 	99,
		// 	[]int64{4, 3, 4, 33},
		// 	operation{
		// 		opCode: OpCodeQuit,
		// 		modes:  []Mode{},
		// 		params: []int64{},
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
