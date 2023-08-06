package main

import "testing"

// The findOptimalPacking function doesn't use the fields in the
// slots, so we can leave them empty.
func TestFindOptimalPacking(t *testing.T) {
	tests := []struct {
		name     string
		options  [][]StorageSlot
		expected []StorageSlot
	}{
		{
			name: "Single best answer",
			options: [][]StorageSlot{
				{
					{Offset: 32}, {Offset: 16}, {Offset: 18}, {Offset: 27},
					{Offset: 32}, {Offset: 16}, {Offset: 13}, {Offset: 32},
					{Offset: 30}, {Offset: 15}, {Offset: 20}, {Offset: 30},
				},
			},
			expected: []StorageSlot{
				{Offset: 32}, {Offset: 16}, {Offset: 13}, {Offset: 32},
			},
		},
		{
			name: "Multiple best answers - pick the first",
			options: [][]StorageSlot{
				{
					{Offset: 32}, {Offset: 32}, {Offset: 17}, {Offset: 17},
					{Offset: 32}, {Offset: 32}, {Offset: 18}, {Offset: 16},
					{Offset: 32}, {Offset: 32}, {Offset: 10}, {Offset: 24},
				},
			},
			expected: []StorageSlot{
				{Offset: 32}, {Offset: 32}, {Offset: 17}, {Offset: 17},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findOptimalPacking(tt.options)
		})
	}
}
