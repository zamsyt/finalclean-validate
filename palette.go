package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
)

// Opaque color
type rgb struct {
	r, g, b uint8
}

// Implement color.Color
func (c rgb) RGBA() (r, g, b, a uint32) {
	return color.RGBA{c.r, c.g, c.b, 255}.RGBA()
}

var palette = loadPalette("r-slash-place-2023.gpl")

// Load palette from GIMP palette file (.gpl)
func loadPalette(path string) color.Palette {
	var p []color.Color
	f, err := os.Open(path)
	check(err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		ln := s.Text()
		var r, g, b uint8
		n, err := fmt.Sscanf(ln, "%d %d %d", &r, &g, &b) // ignoring color name
		if err != nil {
			continue
		}
		if n != 3 {
			panic(fmt.Sprintf("Expected 3 numbers, have: %v", ln))
		}
		p = append(p, rgb{r, g, b})
	}
	return p
}
