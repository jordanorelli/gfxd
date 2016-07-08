package main

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}

	eq := func(c1, c2 color.Color) bool {
		r1, g1, b1, a1 := c1.RGBA()
		r2, g2, b2, a2 := c2.RGBA()
		return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
	}

	type test struct {
		in  string
		out color.RGBA
	}

	var tests = []test{
		{"", white},
		{"0000ff", blue},
		{"0000ffff", blue},
		{"00ff00", green},
		{"00ff00ff", green},
		{"ff0000", red},
		{"ff0000ff", red},
		{"FF0000", red},
	}

	for _, tt := range tests {
		c := parseColor(tt.in, white)
		if eq(c, tt.out) {
			t.Logf("ok: '%s' == %v", tt.in, c)
		} else {
			t.Errorf("parse color failed: '%s' yielded %v, expected %v", tt.in, c, tt.out)
		}
	}
}
