package main

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/rajeet-04/tinycolor-go/internal/compat"
	"github.com/rajeet-04/tinycolor-go/tinycolor"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		request, err := compat.Decode(line)
		if err != nil {
			response, _ := compat.Failure("", err.Error())
			write(response)
			continue
		}
		write(handle(request))
	}
}

func handle(request compat.Request) compat.Response {
	var (
		color tinycolor.Color
		err   error
	)
	args, _ := request.Args.(map[string]any)
	switch request.Operation {
	case "inspect", "string":
		color, err = tinycolor.FromCompat(request.Input, false)
	case "fromRatio":
		color, err = tinycolor.FromCompat(request.Input, true)
	case "output", "analysis", "clone", "modify", "mix", "readability", "isReadable", "mostReadable", "palette":
		color, err = tinycolor.FromCompatWithOptions(request.Input, false, options(args))
	case "equals", "randomInvariant":
	default:
		response, _ := compat.Failure(request.ID, "unsupported operation")
		return response
	}
	if err != nil {
		response, _ := compat.Failure(request.ID, err.Error())
		return response
	}
	if request.Operation == "string" {
		response, _ := compat.Success(request.ID, color.String())
		return response
	}
	switch request.Operation {
	case "output":
		return output(request.ID, color, args)
	case "analysis":
		return analysis(request.ID, color, args["method"])
	case "clone":
		response, _ := compat.Success(request.ID, color.Clone().Inspect())
		return response
	case "modify":
		return modify(request.ID, &color, args)
	case "mix":
		return mix(request.ID, color, args)
	case "readability":
		return readability(request.ID, color, args)
	case "isReadable":
		return isReadable(request.ID, color, args)
	case "mostReadable":
		return mostReadable(request.ID, color, args)
	case "palette":
		return palette(request.ID, color, args)
	case "equals":
		response, _ := compat.Success(request.ID, tinycolor.Equals(request.Input, args["other"]))
		return response
	case "randomInvariant":
		randomColor := tinycolor.Random()
		rgb := randomColor.ToRGB()
		response, _ := compat.Success(request.ID, map[string]any{
			"valid":      randomColor.Valid(),
			"alpha":      float64(1),
			"rgbInRange": rgb.R >= 0 && rgb.R <= 255 && rgb.G >= 0 && rgb.G <= 255 && rgb.B >= 0 && rgb.B <= 255,
		})
		return response
	}
	response, _ := compat.Success(request.ID, color.Inspect())
	return response
}

func modify(id string, color *tinycolor.Color, args map[string]any) compat.Response {
	before := color.Inspect()
	method, _ := args["method"].(string)
	var returned *tinycolor.Color
	switch method {
	case "lighten":
		returned = color.Lighten(defaultAmount(args, 10))
	case "brighten":
		returned = color.Brighten(defaultAmount(args, 10))
	case "darken":
		returned = color.Darken(defaultAmount(args, 10))
	case "saturate":
		returned = color.Saturate(defaultAmount(args, 10))
	case "desaturate":
		returned = color.Desaturate(defaultAmount(args, 10))
	case "greyscale":
		returned = color.Greyscale()
	case "spin":
		if _, ok := args["amount"]; ok {
			returned = color.Spin(number(args["amount"]))
		} else {
			returned = color.Spin(math.NaN())
		}
	default:
		response, _ := compat.Failure(id, "unsupported method")
		return response
	}
	response, _ := compat.Success(id, map[string]any{
		"before":       before,
		"after":        color.Inspect(),
		"sameReceiver": returned == color,
	})
	return response
}

func mix(id string, first tinycolor.Color, args map[string]any) compat.Response {
	second, err := tinycolor.FromCompat(args["other"], false)
	if err != nil {
		response, _ := compat.Failure(id, err.Error())
		return response
	}
	response, _ := compat.Success(id, tinycolor.Mix(first, second, defaultAmount(args, 50)).Inspect())
	return response
}

func readability(id string, first tinycolor.Color, args map[string]any) compat.Response {
	second, _ := tinycolor.FromCompat(args["other"], false)
	response, _ := compat.Success(id, tinycolor.Readability(first, second))
	return response
}

func isReadable(id string, first tinycolor.Color, args map[string]any) compat.Response {
	second, _ := tinycolor.FromCompat(args["other"], false)
	response, _ := compat.Success(id, tinycolor.IsReadable(first, second, wcagOptions(args)))
	return response
}

func mostReadable(id string, base tinycolor.Color, args map[string]any) compat.Response {
	inputs, _ := args["candidates"].([]any)
	candidates := make([]tinycolor.Color, 0, len(inputs))
	for _, input := range inputs {
		candidate, _ := tinycolor.FromCompat(input, false)
		candidates = append(candidates, candidate)
	}
	result, ok := tinycolor.MostReadable(base, candidates, wcagOptions(args))
	if !ok {
		response, _ := compat.Success(id, nil)
		return response
	}
	response, _ := compat.Success(id, result.Inspect())
	return response
}

func palette(id string, color tinycolor.Color, args map[string]any) compat.Response {
	method, _ := args["method"].(string)
	var colors []tinycolor.Color
	switch method {
	case "complement":
		colors = []tinycolor.Color{color.Complement()}
	case "splitcomplement":
		colors = color.SplitComplement()
	case "triad":
		colors = color.Triad()
	case "tetrad":
		colors = color.Tetrad()
	case "analogous":
		colors = color.Analogous(defaultCount(args["results"], 6), defaultCount(args["slices"], 30))
	case "monochromatic":
		colors = color.Monochromatic(defaultCount(args["results"], 6))
	default:
		response, _ := compat.Failure(id, "unsupported method")
		return response
	}
	inspections := make([]map[string]any, len(colors))
	for index, paletteColor := range colors {
		inspections[index] = paletteColor.Inspect()
	}
	response, _ := compat.Success(id, inspections)
	return response
}

func defaultAmount(args map[string]any, fallback float64) float64 {
	if value, ok := args["amount"]; ok {
		amount := number(value)
		if amount == 0 {
			return 0
		}
		return amount
	}
	return fallback
}

func defaultCount(value any, fallback int) int {
	if !truthy(value) {
		return fallback
	}
	return int(number(value))
}

func number(value any) float64 {
	amount, _ := value.(float64)
	return amount
}

func options(args map[string]any) tinycolor.CompatOptions {
	options, _ := args["options"].(map[string]any)
	format, _ := options["format"].(string)
	gradientType, _ := options["gradientType"].(bool)
	return tinycolor.CompatOptions{Format: format, GradientType: gradientType}
}

func wcagOptions(args map[string]any) tinycolor.WCAG2Options {
	raw, _ := args["options"].(map[string]any)
	level, _ := raw["level"].(string)
	size, _ := raw["size"].(string)
	return tinycolor.WCAG2Options{
		Level:                 level,
		Size:                  size,
		IncludeFallbackColors: truthy(raw["includeFallbackColors"]),
	}
}

func analysis(id string, color tinycolor.Color, method any) compat.Response {
	var result any
	switch method {
	case "brightness":
		result = color.Brightness()
	case "luminance":
		result = color.Luminance()
	case "isDark":
		result = color.IsDark()
	case "isLight":
		result = color.IsLight()
	default:
		response, _ := compat.Failure(id, "unsupported method")
		return response
	}
	response, _ := compat.Success(id, result)
	return response
}

func output(id string, color tinycolor.Color, args map[string]any) compat.Response {
	var result any
	switch args["method"] {
	case "toHex":
		result = color.ToHex()
	case "toHex8":
		result = color.ToHex8()
	case "toHexString":
		result = color.ToHexString()
	case "toHex8String":
		result = color.ToHex8String()
	case "toRgbString":
		result = color.ToRGBString()
	case "toPercentageRgbString":
		result = color.ToPercentageRGBString()
	case "toHslString":
		result = color.ToHSLString()
	case "toHsvString":
		result = color.ToHSVString()
	case "toString":
		format, _ := args["format"].(string)
		result = color.ToString(format)
	case "toName":
		name, ok := color.ToName()
		if ok {
			result = name
		} else {
			result = false
		}
	case "toFilter":
		var second *tinycolor.Color
		if truthy(args["secondColor"]) {
			parsed, _ := tinycolor.FromCompat(args["secondColor"], false)
			second = &parsed
		}
		result = color.ToFilter(second, false)
	default:
		response, _ := compat.Failure(id, "unsupported method")
		return response
	}
	response, _ := compat.Success(id, result)
	return response
}

func truthy(value any) bool {
	switch value := value.(type) {
	case nil:
		return false
	case bool:
		return value
	case float64:
		return value != 0
	case string:
		return value != ""
	default:
		return true
	}
}

func write(response compat.Response) {
	encoded, err := compat.Encode(response)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(encoded))
}
