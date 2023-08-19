package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"
)

func TestPalette(t *testing.T) {
	tp := append(color.Palette{color.Transparent}, palette...)
	f, err := os.Open(outpath)
	check(err)
	img := getImg(f)
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			if !colorEq(c, tp.Convert(c)) {
				fmt.Println("Output doesn't match the palette!")
				fmt.Printf("%v: %+v converts to %+v\n", image.Pt(x, y), c, tp.Convert(c))
				os.Exit(1)
			}
		}
	}
}
