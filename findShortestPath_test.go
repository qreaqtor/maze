package main

import (
	"reflect"
	"testing"
)

type TestCase struct {
	testName string
	graph    [][]int
	start    point
	end      point
	expected []point
}

func TestFindShortestPath(t *testing.T) {
	tests := []TestCase{
		{
			graph: [][]int{
				{1, 2, 0},
				{2, 0, 1},
				{9, 1, 0},
			},
			start:    point{0, 0},
			end:      point{2, 1},
			expected: []point{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
			testName: "Simple",
		},
		{
			graph: [][]int{
				{1, 1, 1, 1},
				{1, 0, 1, 1},
				{9, 0, 9, 1},
				{1, 1, 1, 1},
			},
			start:    point{0, 0},
			end:      point{3, 0},
			expected: []point{{0, 0}, {0, 1}, {0, 2}, {1, 2}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}, {3, 0}},
			testName: "Shortest path contains more cells",
		},
		{
			graph: [][]int{
				{1, 1, 1},
				{1, 1, 0},
				{1, 1, 1},
			},
			start:    point{0, 0},
			end:      point{2, 2},
			expected: []point{{0, 0}, {1, 0}, {1, 1}, {2, 1}, {2, 2}},
			testName: "Save came from, not came to",
		},
		{
			graph: [][]int{
				{1, 1, 1, 0, 1},
				{0, 0, 1, 0, 1},
				{1, 1, 1, 1, 1},
				{1, 0, 0, 0, 1},
				{1, 1, 1, 1, 1},
			},
			start:    point{0, 0},
			end:      point{4, 4},
			expected: []point{{0, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {2, 3}, {2, 4}, {3, 4}, {4, 4}},
			testName: "Multiple paths",
		},
		{
			graph: [][]int{
				{1, 9, 9, 9, 9, 9, 9},
				{9, 1, 9, 1, 1, 1, 1},
				{1, 1, 9, 1, 9, 9, 9},
				{1, 1, 9, 1, 1, 1, 1},
				{9, 9, 9, 9, 9, 9, 1},
				{1, 1, 1, 1, 1, 1, 1},
			},
			start:    point{0, 0},
			end:      point{5, 6},
			expected: []point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {5, 6}},
			testName: "Maze with long detours",
		},
	}

	for _, test := range tests {
		actual := findShortestPath(test.graph, test.start, test.end)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("test failed: %s\nexpected:\t%v\nactual: \t%v", test.testName, test.expected, actual)
		}
	}
}
