package main

import "fmt"

func main() {
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#fh9921",
		"white": "#ffffff",
	}

	// var colors map[string]string

	// colors := make(map[string]string)
	// colors["red"] = "#ff0001"
	// delete(colors, "red")

	for color, hex := range colors {
		fmt.Printf("Color %v has the %v hex\n", color, hex)
	}
}
