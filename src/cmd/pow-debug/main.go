package main

import (
	"fmt"
	"math"
)

func main() {
	for channel := 11; channel < 256; channel++ {
		normalized := float64(channel) / 255
		fmt.Printf("%d %x\n", channel, math.Float64bits(math.Pow((normalized+.055)/1.055, 2.4)))
	}
}
