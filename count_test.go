package main

import (
	"fmt"
	"image/color"
	"os"
	"testing"
)

func TestBaseCount(t *testing.T) {
	f, err := os.Open(basediffpath)
	check(err)
	img := getImg(f)
	b := img.Bounds()
	var count int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			if alpha(c) != 0 {
				count++
			}
		}
	}
	fmt.Println("base diff count:", count)
}

func TestPaletteCount(t *testing.T) {
	f, err := os.Open(palettediffpath)
	check(err)
	img := getImg(f)
	b := img.Bounds()
	var count int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			if alpha(c) != 0 {
				count++
				if colorEq(c, palette.Convert(c)) {
					t.Error("Expected palette diff to only have non-palette pixels")
					return
				}
			}
		}
	}
	fmt.Println("palette diff count:", count)
}

func alpha(c color.Color) uint8 {
	switch v := c.(type) {
	case color.NRGBA:
		return v.A
	case rgb:
		return 255
	default:
		return color.RGBAModel.Convert(c).(color.RGBA).A
	}
}
