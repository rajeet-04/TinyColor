package tinycolor

import (
	"math"
	"testing"
)

func TestModifiers(t *testing.T) {
	red, _ := FromCompat("red", false)
	if returned := red.Lighten(10); returned != &red {
		t.Fatal("Lighten must return its receiver")
	}
	if got := red.ToHexString(); got != "#ff3333" {
		t.Fatalf("Lighten(10) = %s", got)
	}

	for _, modifier := range []struct {
		name  string
		apply func(*Color) *Color
	}{
		{"lighten", func(c *Color) *Color { return c.Lighten(0) }},
		{"darken", func(c *Color) *Color { return c.Darken(0) }},
		{"saturate", func(c *Color) *Color { return c.Saturate(0) }},
		{"desaturate", func(c *Color) *Color { return c.Desaturate(0) }},
		{"brighten", func(c *Color) *Color { return c.Brighten(0) }},
		{"spin", func(c *Color) *Color { return c.Spin(0) }},
	} {
		t.Run(modifier.name+" zero", func(t *testing.T) {
			color, _ := FromCompat("#336699", false)
			if modifier.apply(&color) != &color || color.ToHexString() != "#336699" {
				t.Fatalf("%s(0) changed %#v", modifier.name, color)
			}
		})
	}

	white, _ := FromCompat("#fff", false)
	white.Lighten(100)
	if got := white.ToHexString(); got != "#ffffff" {
		t.Fatalf("Lighten clamp = %s", got)
	}
	black, _ := FromCompat("#000", false)
	black.Darken(100)
	if got := black.ToHexString(); got != "#000000" {
		t.Fatalf("Darken clamp = %s", got)
	}
	red, _ = FromCompat("red", false)
	red.Desaturate(200)
	if got := red.ToHexString(); got != "#808080" {
		t.Fatalf("Desaturate clamp = %s", got)
	}

	black, _ = FromCompat("#000", false)
	black.Brighten(.2)
	if got := black.ToHexString(); got != "#010101" {
		t.Fatalf("Brighten pre-clamp rounding = %s", got)
	}

	red, _ = FromCompat("red", false)
	red.Spin(-120)
	if got := red.ToHexString(); got != "#0000ff" {
		t.Fatalf("Spin(-120) = %s", got)
	}
	red, _ = FromCompat("red", false)
	red.Spin(480)
	if got := red.ToHexString(); got != "#00ff00" {
		t.Fatalf("Spin(480) = %s", got)
	}

	color, _ := FromCompatWithOptions("rgba(255, 0, 0, .4)", false, CompatOptions{GradientType: true})
	original, format, alpha := color.Original(), color.Format(), color.Alpha()
	if color.Spin(math.NaN()) != &color {
		t.Fatal("Spin must return its receiver")
	}
	if color.ToHexString() != "#000000" || !color.Valid() || color.Format() != format || color.Original() != original || color.Alpha() != alpha {
		t.Fatalf("Spin(NaN) did not retain metadata: %#v", color)
	}
	if got := color.ToFilter(nil, false); got != "progid:DXImageTransform.Microsoft.gradient(GradientType = 1, startColorstr=#66000000,endColorstr=#66000000)" {
		t.Fatalf("Spin(NaN) gradient metadata = %s", got)
	}

	alphaColor, _ := FromCompat("rgba(255, 0, 0, .4)", false)
	alphaColor.Greyscale()
	if alphaColor.ToHexString() != "#808080" || alphaColor.Alpha() != .4 {
		t.Fatalf("Greyscale() = %#v", alphaColor)
	}
}

func TestMix(t *testing.T) {
	black, _ := FromCompat("#000", false)
	white, _ := FromCompat("#fff", false)
	if got := Mix(black, white, 50).ToHexString(); got != "#808080" {
		t.Fatalf("Mix(50) = %s", got)
	}
	if got := Mix(black, white, 0).ToHexString(); got != "#000000" {
		t.Fatalf("Mix(0) = %s", got)
	}
	if got := Mix(white, black, 90).ToHexString(); got != "#1a1a1a" {
		t.Fatalf("Mix(90) = %s", got)
	}

	transparent, _ := FromCompat("transparent", false)
	mixed := Mix(transparent, black, 25)
	if mixed.ToHexString() != "#000000" || mixed.Alpha() != .25 {
		t.Fatalf("transparent Mix = %#v", mixed)
	}
	if transparent.Alpha() != 0 || black.ToHexString() != "#000000" {
		t.Fatal("Mix mutated an input")
	}

	if got := Mix(black, white, 200).ToHexString(); got != "#ffffff" {
		t.Fatalf("out-of-range Mix = %s", got)
	}
}
