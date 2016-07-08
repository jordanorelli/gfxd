package main

import (
	"testing"
)

func TestNorm(t *testing.T) {
	type test struct {
		i   int
		min int
		max int
		n   float64
	}
	tests := []test{
		{1, 2, 3, 0},
		{3, 2, 1, 1},
		{1, 0, 10, 0.1},
		{5, 0, 10, 0.5},
		{125, 100, 200, 0.25},
		{-125, -200, -100, 0.75},
	}

	for _, test := range tests {
		n := norm(test.i, test.min, test.max)
		if n == test.n {
			t.Logf("norm(%d, %d, %d) == %f", test.i, test.min, test.max, n)
		} else {
			t.Errorf("norm(%d, %d, %d) is %f, expected %f", test.i, test.min, test.max, n, test.n)
		}
	}
}
