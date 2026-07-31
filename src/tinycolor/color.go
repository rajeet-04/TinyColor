// Package tinycolor exposes the compatibility facade over normalized parsing.
package tinycolor

import (
	"fmt"
	"math"
	"strconv"

	"github.com/rajeet-04/tinycolor-go/internal/color"
	"github.com/rajeet-04/tinycolor-go/internal/parser"
)

type Color struct{ model color.Model }

func FromCompat(input any, fromRatio bool) (Color, error) {
	if fromRatio {
		return Color{model: parser.ParseFromRatio(input)}, nil
	}
	return Color{model: parser.Parse(input)}, nil
}

func (c Color) Valid() bool    { return c.model.Valid }
func (c Color) Format() string { return string(c.model.Format) }
func (c Color) Alpha() float64 { return c.model.A }
func (c Color) Original() any  { return c.model.Original }
func (c Color) RGB() map[string]any {
	return map[string]any{"r": math.Round(c.model.R), "g": math.Round(c.model.G), "b": math.Round(c.model.B), "a": c.model.A}
}

// String supplies the minimal source-compatible string snapshot used by the
// JSONL oracle. Full public output APIs remain Phase 3 work.
func (c Color) String() string {
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
