package main

import (
	"testing"

	"github.com/rensr/spack/solidity"
)

// The findOptimalPacking function doesn't use the fields in the
// slots, so we can leave them empty.
func TestFindOptimalPacking(t *testing.T) {
	tests := []struct {
		name     string
		options  [][]solidity.StorageSlot
		expected []solidity.StorageSlot
	}{
		{
			name: "Single best answer",
			options: [][]solidity.StorageSlot{
				{
					{Offset: 32}, {Offset: 16}, {Offset: 18}, {Offset: 27},
				},
				{
					{Offset: 32}, {Offset: 16}, {Offset: 13}, {Offset: 32},
				},
				{
					{Offset: 30}, {Offset: 15}, {Offset: 20}, {Offset: 30},
				},
			},
			expected: []solidity.StorageSlot{
				{Offset: 32}, {Offset: 16}, {Offset: 13}, {Offset: 32},
			},
		},
		{
			name: "Multiple best answers - pick the first",
			options: [][]solidity.StorageSlot{
				{
					{Offset: 32}, {Offset: 32}, {Offset: 17}, {Offset: 17},
				},
				{
					{Offset: 32}, {Offset: 32}, {Offset: 18}, {Offset: 16},
				},
				{
					{Offset: 32}, {Offset: 32}, {Offset: 10}, {Offset: 24},
				},
			},
			expected: []solidity.StorageSlot{
				{Offset: 32}, {Offset: 32}, {Offset: 17}, {Offset: 17},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := findOptimalPacking(tt.options)
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, actual)
			}
			for i := range actual {
				if actual[i].Offset != tt.expected[i].Offset {
					t.Errorf("Expected %v, got %v", tt.expected, actual)
				}
			}
		})
	}
}
