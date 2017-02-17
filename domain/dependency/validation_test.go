package dependency

import (
	"testing"

)


func TestCycleDetection(t *testing.T) {
	testCases := []struct{
		Edges [][]int
		Expected bool
	}{
		{
			Edges: [][]int{{1,2}, {2,3}, {3,4}},
			Expected: false,
		}, {
			Edges: [][]int{{1,2}, {2,3}, {3,1}},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			res := detectCycle(tc.Edges)
			if res != tc.Expected {
				t.Error("invalid result")
			}
		})
	}
}