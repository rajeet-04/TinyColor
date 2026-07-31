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
	mixedHalf := Mix(black, white, 50)
	if got := mixedHalf.ToHexString(); got != "#808080" {
		t.Fatalf("Mix(50) = %s", got)
	}
	if got, ok := mixedHalf.Original().(map[string]float64); !ok || got["r"] != 127.5 || got["g"] != 127.5 || got["b"] != 127.5 || got["a"] != 1 {
		t.Fatalf("Mix(50) original = %#v", mixedHalf.Original())
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

func TestReadability(t *testing.T) {
	black, _ := FromCompat("#000", false)
	white, _ := FromCompat("#fff", false)
	if got := Readability(black, black); got != 1 {
		t.Fatalf("same-color readability = %v", got)
	}
	if got := Readability(black, white); got != 21 {
		t.Fatalf("black-on-white readability = %v", got)
	}
}

func TestIsReadable(t *testing.T) {
	for _, test := range []struct {
		name     string
		ratio    float64
		options  WCAG2Options
		expected bool
	}{
		{"AA small threshold", 4.5, WCAG2Options{Level: "AA", Size: "small"}, true},
		{"AA large threshold", 3, WCAG2Options{Level: "AA", Size: "large"}, true},
		{"AAA small threshold", 7, WCAG2Options{Level: "AAA", Size: "small"}, true},
		{"AAA large threshold", 4.5, WCAG2Options{Level: "AAA", Size: "large"}, true},
		{"below threshold", 4.499999999999999, WCAG2Options{Level: "AA", Size: "small"}, false},
		{"mixed case", 7, WCAG2Options{Level: "aaa", Size: "LARGE"}, true},
		{"invalid values default", 4.5, WCAG2Options{Level: "invalid", Size: "invalid"}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := isReadableRatio(test.ratio, test.options); got != test.expected {
				t.Fatalf("isReadableRatio(%v, %#v) = %t", test.ratio, test.options, got)
			}
		})
	}
}

func TestMostReadable(t *testing.T) {
	base, _ := FromCompat("#000", false)
	first, _ := FromCompat("white", false)
	second, _ := FromCompat("#fff", false)
	if got, ok := MostReadable(base, []Color{first, second}, WCAG2Options{}); !ok || got.Original() != "white" {
		t.Fatalf("first tied candidate = %#v, %t", got, ok)
	}

	gray, _ := FromCompat("#777", false)
	matchingGray, _ := FromCompat("#777", false)
	if got, ok := MostReadable(gray, []Color{matchingGray}, WCAG2Options{IncludeFallbackColors: true}); !ok || got.ToHexString() != "#000000" {
		t.Fatalf("fallback candidate = %#v, %t", got, ok)
	}
	if got, ok := MostReadable(gray, []Color{matchingGray}, WCAG2Options{}); !ok || got.ToHexString() != "#777777" {
		t.Fatalf("fallback-disabled candidate = %#v, %t", got, ok)
	}
	if _, ok := MostReadable(base, nil, WCAG2Options{}); ok {
		t.Fatal("empty candidates without fallback must have no result")
	}
	if _, ok := MostReadable(base, nil, WCAG2Options{IncludeFallbackColors: true}); ok {
		t.Fatal("empty candidates with fallback must have no result")
	}
}
