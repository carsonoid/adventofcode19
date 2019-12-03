package main

import "testing"

func Test_operate(t *testing.T) {
	tests := []struct {
		name    string
		data    []uint64
		wantErr bool
		want    uint64
	}{
		{
			name: "add",
			data: []uint64{1, 0, 0, 0, 99},
			want: 2,
		},
		{
			name: "multi",
			data: []uint64{2, 3, 0, 3, 99},
			want: 6,
		},
		{
			name: "multi2",
			data: []uint64{2, 4, 4, 5, 99, 0},
			want: 9801,
		},
		{
			name:    "invalidcode",
			data:    []uint64{123, 4, 4, 5, 99, 0},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workSet := tt.data[:4]
			resultPos := workSet[3]
			if err := operate(tt.data, workSet); (err != nil) != tt.wantErr {
				t.Errorf("operate() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.data[resultPos]
			if got != tt.want {
				t.Errorf("operate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
