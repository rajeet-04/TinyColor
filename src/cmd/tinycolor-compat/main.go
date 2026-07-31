package main

import (
	"bufio"
	"fmt"
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
	case "output", "analysis", "clone":
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

func options(args map[string]any) tinycolor.CompatOptions {
	options, _ := args["options"].(map[string]any)
	format, _ := options["format"].(string)
	gradientType, _ := options["gradientType"].(bool)
	return tinycolor.CompatOptions{Format: format, GradientType: gradientType}
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
