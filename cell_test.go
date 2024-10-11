package main

import "testing"

func TestGetAliveNeighborCells(t *testing.T) {
	// Create a test grid
	grid := &Grid{
		x: 5,
		y: 5,
		Cells: [][]Cell{
			{Alive, Dead, Alive, Dead, Dead},
			{Dead, Dead, Dead, Alive, Dead},
			{Alive, Dead, Alive, Dead, Dead},
			{Dead, Alive, Dead, Dead, Dead},
			{Dead, Dead, Dead, Dead, Dead},
		},
	}

	testCases := []struct {
		name     string
		x, y     int
		expected int
	}{
		{"Diagonals", 1, 1, 4},
		{"Corner cell", 0, 0, 0},
		{"Edge cell", 2, 0, 2},
		{"All dead neighbors", 5, 4, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grid.GetAliveNeighborCells(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("GetAliveNeighborCells(%d, %d) = %d; want %d", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}
