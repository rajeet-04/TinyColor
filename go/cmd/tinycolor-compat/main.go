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
			write(compat.Response{ID: "", Error: err.Error()})
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
	switch request.Operation {
	case "inspect", "string":
		color, err = tinycolor.FromCompat(request.Input, false)
	case "fromRatio":
		color, err = tinycolor.FromCompat(request.Input, true)
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
	response, _ := compat.Success(request.ID, color.Inspect())
	return response
}

func write(response compat.Response) {
	encoded, err := compat.Encode(response)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(encoded))
}
