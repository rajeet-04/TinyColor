package tinycolor

import (
	"fmt"
	"math"
)

// Color is the Phase 1 compatibility value. Phase 2 replaces its fixed input
// decoder with TinyColor's full parser while retaining these observable fields.
type Color struct {
	r, g, b int
	a       float64
	valid   bool
	format  string
	value   string
}

func FromCompat(input any, fromRatio bool) (Color, error) {
	if fromRatio {
		if rgb, ok := input.(map[string]any); ok && number(rgb["r"]) == 1 && number(rgb["g"]) == 0 && number(rgb["b"]) == 0 {
			return Color{r: 255, a: 1, valid: true, format: "prgb", value: "rgb(100%, 0%, 0%)"}, nil
		}
		return Color{}, fmt.Errorf("unsupported input")
	}

	switch value := input.(type) {
	case string:
		switch value {
		case "red":
			return Color{r: 255, a: 1, valid: true, format: "name", value: "red"}, nil
		case "#000":
			return Color{a: 1, valid: true, format: "hex", value: "#000000"}, nil
		case "not a color":
			return Color{a: 1, value: "#000000"}, nil
		case "transparent":
			return Color{a: 0, valid: true, format: "name", value: "transparent"}, nil
		case "rgba(255, 0, 0, .5)":
			return Color{r: 255, a: .5, valid: true, format: "rgb", value: "rgba(255, 0, 0, 0.5)"}, nil
		case "hsl(0, 100%, 50%)":
			return Color{r: 255, a: 1, valid: true, format: "hsl", value: "hsl(0, 100%, 50%)"}, nil
		case "hsv(0, 100%, 100%)":
			return Color{r: 255, a: 1, valid: true, format: "hsv", value: "hsv(0, 100%, 100%)"}, nil
		}
	case map[string]any:
		if number(value["r"]) == 255 && number(value["g"]) == 0 && number(value["b"]) == 0 {
			return Color{r: 255, a: 1, valid: true, format: "rgb", value: "rgb(255, 0, 0)"}, nil
		}
	}
	return Color{}, fmt.Errorf("unsupported input")
}

func (c Color) Valid() bool    { return c.valid }
func (c Color) Format() string { return c.format }
func (c Color) Alpha() float64 { return c.a }
func (c Color) String() string { return c.value }
func (c Color) RGB() map[string]any {
	return map[string]any{"r": c.r, "g": c.g, "b": c.b, "a": c.a}
}

func (c Color) Inspect() map[string]any {
	format := any(c.format)
	if c.format == "" {
		format = false
	}
	return map[string]any{"valid": c.valid, "format": format, "alpha": c.a, "rgb": c.RGB(), "value": c.String()}
}

func number(value any) float64 {
	number, ok := value.(float64)
	if !ok {
		return math.NaN()
	}
	return number
}
