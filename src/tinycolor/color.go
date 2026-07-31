// Package tinycolor exposes the compatibility facade over normalized parsing.
package tinycolor

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/rajeet-04/tinycolor-go/internal/color"
	"github.com/rajeet-04/tinycolor-go/internal/parser"
)

type Color struct {
	model        color.Model
	gradientType bool
}
type CompatOptions struct {
	Format       string
	GradientType bool
}
type RGB struct {
	R, G, B int
	A       float64
}
type HSL struct{ H, S, L, A float64 }
type HSV struct{ H, S, V, A float64 }

func FromCompat(input any, fromRatio bool) (Color, error) {
	if fromRatio {
		return Color{model: parser.ParseFromRatio(input)}, nil
	}
	return Color{model: parser.Parse(input)}, nil
}
func FromCompatWithOptions(input any, fromRatio bool, options CompatOptions) (Color, error) {
	c, e := FromCompat(input, fromRatio)
	if options.Format != "" {
		c.model.Format = color.Format(options.Format)
	}
	c.gradientType = options.GradientType
	return c, e
}

func (c Color) Valid() bool    { return c.model.Valid }
func (c Color) Format() string { return string(c.model.Format) }
func (c Color) Alpha() float64 { return c.model.A }
func (c Color) Original() any  { return c.model.Original }
func (c Color) RGB() map[string]any {
	return map[string]any{"r": math.Round(c.model.R), "g": math.Round(c.model.G), "b": math.Round(c.model.B), "a": c.model.A}
}
func (c Color) ToRGB() RGB {
	return RGB{int(math.Round(c.model.R)), int(math.Round(c.model.G)), int(math.Round(c.model.B)), c.model.A}
}
func (c Color) ToHSL() HSL {
	h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
	return HSL{h * 360, s, l, c.model.A}
}
func (c Color) ToHSV() HSV {
	h, s, v := rgbToHSV(c.model.R, c.model.G, c.model.B)
	return HSV{h * 360, s, v, c.model.A}
}
func (c Color) ToRGBString() string {
	x := c.ToRGB()
	if x.A == 1 {
		return fmt.Sprintf("rgb(%d, %d, %d)", x.R, x.G, x.B)
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %s)", x.R, x.G, x.B, roundedAlpha(x.A))
}
func (c Color) ToPercentageRGB() RGB {
	x := c.ToRGB()
	return RGB{percent(x.R), percent(x.G), percent(x.B), x.A}
}
func (c Color) ToPercentageRGBString() string {
	x := c.ToPercentageRGB()
	if x.A == 1 {
		return fmt.Sprintf("rgb(%d%%, %d%%, %d%%)", x.R, x.G, x.B)
	}
	return fmt.Sprintf("rgba(%d%%, %d%%, %d%%, %s)", x.R, x.G, x.B, roundedAlpha(x.A))
}
func (c Color) ToHSLString() string {
	x := c.ToHSL()
	p := "hsl"
	if x.A < 1 {
		p = "hsla"
	}
	s := fmt.Sprintf("%s(%d, %d%%, %d%%", p, mathRound(x.H), mathRound(x.S*100), mathRound(x.L*100))
	if x.A < 1 {
		s += ", " + roundedAlpha(x.A)
	}
	return s + ")"
}
func (c Color) ToHSVString() string {
	x := c.ToHSV()
	p := "hsv"
	if x.A < 1 {
		p = "hsva"
	}
	s := fmt.Sprintf("%s(%d, %d%%, %d%%", p, mathRound(x.H), mathRound(x.S*100), mathRound(x.V*100))
	if x.A < 1 {
		s += ", " + roundedAlpha(x.A)
	}
	return s + ")"
}
func (c Color) ToHex() string {
	return fmt.Sprintf("%02x%02x%02x", mathRound(c.model.R), mathRound(c.model.G), mathRound(c.model.B))
}
func (c Color) ToHexString() string           { return "#" + c.ToHex() }
func (c Color) ToHex8() string                { return c.ToHex() + fmt.Sprintf("%02x", mathRound(c.model.A*255)) }
func (c Color) ToHex8String() string          { return "#" + c.ToHex8() }
func (c Color) ToString(format string) string { return c.toString(format, format != "") }
func (c Color) toString(format string, explicit bool) string {
	if format == "" {
		format = string(c.model.Format)
	}
	if !explicit && c.model.A < 1 && (format == "hex" || format == "hex6" || format == "hex3" || format == "hex4" || format == "hex8" || format == "name") {
		if format == "name" && c.model.A == 0 {
			return "transparent"
		}
		return c.ToRGBString()
	}
	switch format {
	case "rgb":
		return c.ToRGBString()
	case "prgb":
		return c.ToPercentageRGBString()
	case "hsl":
		return c.ToHSLString()
	case "hsv":
		return c.ToHSVString()
	case "hex", "hex6":
		return c.ToHexString()
	case "hex8":
		return c.ToHex8String()
	case "name":
		if n, ok := c.ToName(); ok {
			return n
		}
		return c.ToHexString()
	}
	return c.ToHexString()
}
func (c Color) ToName() (string, bool) {
	if c.model.A == 0 {
		return "transparent", true
	}
	if c.model.A < 1 {
		return "", false
	}
	n := parser.NameForRGB(mathRound(c.model.R), mathRound(c.model.G), mathRound(c.model.B))
	return n, n != ""
}
func (c Color) ToFilter(second *Color, gradient bool) string {
	start := "#" + fmt.Sprintf("%02x", mathRound(c.model.A*255)) + c.ToHex()
	end := start
	if second != nil {
		end = "#" + fmt.Sprintf("%02x", mathRound(second.model.A*255)) + second.ToHex()
	}
	prefix := ""
	if gradient || c.gradientType {
		prefix = "GradientType = 1, "
	}
	return "progid:DXImageTransform.Microsoft.gradient(" + prefix + "startColorstr=" + start + ",endColorstr=" + end + ")"
}
func (c Color) Brightness() float64 { x := c.ToRGB(); return float64(x.R*299+x.G*587+x.B*114) / 1000 }
func (c Color) Luminance() float64 {
	x := c.ToRGB()
	f := func(v int) float64 {
		z := float64(v) / 255
		if z <= .03928 {
			return z / 12.92
		}
		return math.Pow((z+.055)/1.055, 2.4)
	}
	return .2126*f(x.R) + .7152*f(x.G) + .0722*f(x.B)
}
func (c Color) IsDark() bool  { return c.Brightness() < 128 }
func (c Color) IsLight() bool { return !c.IsDark() }
func (c Color) Clone() Color  { n, _ := FromCompat(c.String(), false); return n }
func Equals(a, b any) bool {
	if a == nil || b == nil {
		return false
	}
	x, _ := FromCompat(a, false)
	y, _ := FromCompat(b, false)
	return x.ToRGBString() == y.ToRGBString()
}
func Random() Color {
	return Color{model: color.Model{R: rand.Float64() * 255, G: rand.Float64() * 255, B: rand.Float64() * 255, A: 1, Valid: true, Format: color.FormatRGB}}
}
func mathRound(v float64) int { return int(math.Floor(v + .5)) }

// String supplies the minimal source-compatible string snapshot used by the
// JSONL oracle. Full public output APIs remain Phase 3 work.
func (c Color) String() string {
	return c.toString("", false)
	/*
		r, g, b := int(math.Round(c.model.R)), int(math.Round(c.model.G)), int(math.Round(c.model.B))
		if c.model.Format == color.FormatName {
			if c.model.A == 0 {
				return "transparent"
			}
			if name := parser.NameForRGB(r, g, b); name != "" {
				return name
			}
		}
		if c.model.Format == color.FormatPRGB {
			if c.model.A < 1 {
				return fmt.Sprintf("rgba(%d%%, %d%%, %d%%, %s)", percent(r), percent(g), percent(b), roundedAlpha(c.model.A))
			}
			return fmt.Sprintf("rgb(%d%%, %d%%, %d%%)", percent(r), percent(g), percent(b))
		}
		if c.model.Format == color.FormatHSL {
			h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
			if c.model.A < 1 {
				return fmt.Sprintf("hsla(%d, %d%%, %d%%, %s)", int(math.Round(h*360)), int(math.Round(s*100)), int(math.Round(l*100)), roundedAlpha(c.model.A))
			}
			return fmt.Sprintf("hsl(%d, %d%%, %d%%)", int(math.Round(h*360)), int(math.Round(s*100)), int(math.Round(l*100)))
		}
		if c.model.Format == color.FormatHSV {
			h, s, v := rgbToHSV(c.model.R, c.model.G, c.model.B)
			if c.model.A < 1 {
				return fmt.Sprintf("hsva(%d, %d%%, %d%%, %s)", int(math.Round(h*360)), int(math.Round(s*100)), int(math.Round(v*100)), roundedAlpha(c.model.A))
			}
			return fmt.Sprintf("hsv(%d, %d%%, %d%%)", int(math.Round(h*360)), int(math.Round(s*100)), int(math.Round(v*100)))
		}
		if c.model.Format == color.FormatRGB || c.model.A < 1 {
			if c.model.A == 1 {
				return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
			}
			return fmt.Sprintf("rgba(%d, %d, %d, %s)", r, g, b, roundedAlpha(c.model.A))
		}
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	*/
}

func (c Color) Inspect() map[string]any {
	format := any(c.Format())
	if c.Format() == "" {
		format = false
	}
	return map[string]any{"valid": c.Valid(), "format": format, "alpha": c.Alpha(), "rgb": c.RGB(), "value": c.String(), "original": c.Original()}
}

func percent(channel int) int { return int(math.Round(float64(channel) / 255 * 100)) }
func roundedAlpha(alpha float64) string {
	return strconv.FormatFloat(math.Round(alpha*100)/100, 'f', -1, 64)
}

func rgbToHSL(r, g, b float64) (float64, float64, float64) {
	r, g, b = r/255, g/255, b/255
	max, min := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	l := (max + min) / 2
	if max == min {
		return 0, 0, l
	}
	d := max - min
	s := d / (2 - max - min)
	if l < 0.5 {
		s = d / (max + min)
	}
	var h float64
	switch max {
	case r:
		h = (g-b)/d + map[bool]float64{true: 6, false: 0}[g < b]
	case g:
		h = (b-r)/d + 2
	default:
		h = (r-g)/d + 4
	}
	return h / 6, s, l
}

func rgbToHSV(r, g, b float64) (float64, float64, float64) {
	r, g, b = r/255, g/255, b/255
	max, min := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	d := max - min
	if max == 0 {
		return 0, 0, 0
	}
	if d == 0 {
		return 0, 0, max
	}
	var h float64
	switch max {
	case r:
		h = (g-b)/d + map[bool]float64{true: 6, false: 0}[g < b]
	case g:
		h = (b-r)/d + 2
	default:
		h = (r-g)/d + 4
	}
	return h / 6, d / max, max
}

func (c *Color) Lighten(amount float64) *Color {
	h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
	c.setHSL(h, s, clamp01(l+amount/100))
	return c
}

func (c *Color) Darken(amount float64) *Color {
	h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
	c.setHSL(h, s, clamp01(l-amount/100))
	return c
}

func (c *Color) Saturate(amount float64) *Color {
	h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
	c.setHSL(h, clamp01(s+amount/100), l)
	return c
}

func (c *Color) Desaturate(amount float64) *Color {
	h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
	c.setHSL(h, clamp01(s-amount/100), l)
	return c
}

func (c *Color) Greyscale() *Color { return c.Desaturate(100) }

func (c *Color) Brighten(amount float64) *Color {
	delta := float64(mathRound(255 * amount / 100))
	c.model.R = clampChannel(c.model.R + delta)
	c.model.G = clampChannel(c.model.G + delta)
	c.model.B = clampChannel(c.model.B + delta)
	return c
}

func (c *Color) Spin(amount float64) *Color {
	if math.IsNaN(amount) {
		c.model.R, c.model.G, c.model.B = 0, 0, 0
		return c
	}
	h, s, l := rgbToHSL(c.model.R, c.model.G, c.model.B)
	h = math.Mod(h*360+amount, 360)
	if h < 0 {
		h += 360
	}
	c.setHSL(h/360, s, l)
	return c
}

func Mix(first, second Color, amount float64) Color {
	p := amount / 100
	firstRGB, secondRGB := first.ToRGB(), second.ToRGB()
	return Color{model: color.Model{
		R:      clampChannel((float64(secondRGB.R)-float64(firstRGB.R))*p + float64(firstRGB.R)),
		G:      clampChannel((float64(secondRGB.G)-float64(firstRGB.G))*p + float64(firstRGB.G)),
		B:      clampChannel((float64(secondRGB.B)-float64(firstRGB.B))*p + float64(firstRGB.B)),
		A:      clampAlpha((second.model.A-first.model.A)*p + first.model.A),
		Valid:  true,
		Format: color.FormatRGB,
	}}
}

func (c *Color) setHSL(h, s, l float64) {
	c.model.R, c.model.G, c.model.B = hslToRGB(h, s, l)
}

func hslToRGB(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		channel := l * 255
		return channel, channel, channel
	}
	var q float64
	if l < .5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	channel := func(t float64) float64 {
		if t < 0 {
			t++
		}
		if t > 1 {
			t--
		}
		switch {
		case t*6 < 1:
			return p + (q-p)*6*t
		case t*2 < 1:
			return q
		case t*3 < 2:
			return p + (q-p)*(2.0/3.0-t)*6
		default:
			return p
		}
	}
	return channel(h+1.0/3.0) * 255, channel(h) * 255, channel(h-1.0/3.0) * 255
}

func clamp01(value float64) float64      { return math.Min(1, math.Max(0, value)) }
func clampChannel(value float64) float64 { return math.Min(255, math.Max(0, value)) }
func clampAlpha(value float64) float64 {
	if value < 0 || value > 1 || math.IsNaN(value) {
		return 1
	}
	return value
}
