package tinycolor

import "testing"

func TestPhaseOneInputsMatchTinyColor(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		ratio  bool
		valid  bool
		format string
		alpha  float64
		value  string
	}{
		{"name", "red", false, true, "name", 1, "red"},
		{"hex", "#000", false, true, "hex", 1, "#000000"},
		{"invalid", "not a color", false, false, "", 1, "#000000"},
		{"transparent", "transparent", false, true, "name", 0, "transparent"},
		{"rgba", "rgba(255, 0, 0, .5)", false, true, "rgb", .5, "rgba(255, 0, 0, 0.5)"},
		{"hsl", "hsl(0, 100%, 50%)", false, true, "hsl", 1, "hsl(0, 100%, 50%)"},
		{"hsv", "hsv(0, 100%, 100%)", false, true, "hsv", 1, "hsv(0, 100%, 100%)"},
		{"rgb object", map[string]any{"r": float64(255), "g": float64(0), "b": float64(0)}, false, true, "rgb", 1, "rgb(255, 0, 0)"},
		{"from ratio", map[string]any{"r": float64(1), "g": float64(0), "b": float64(0)}, true, true, "prgb", 1, "rgb(100%, 0%, 0%)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := FromCompat(tt.input, tt.ratio)
			if err != nil {
				t.Fatal(err)
			}
			if color.Valid() != tt.valid || color.Format() != tt.format || color.Alpha() != tt.alpha || color.String() != tt.value {
				t.Fatalf("got valid=%v format=%q alpha=%v string=%q", color.Valid(), color.Format(), color.Alpha(), color.String())
			}
		})
	}
}

func TestPhaseOneRejectsUnsupportedInput(t *testing.T) {
	if _, err := FromCompat("blue", false); err == nil {
		t.Fatal("expected unsupported input error")
	}
}
