package main

import (
	"fmt"
	"image/color"
	"os"
	"testing"
)

func TestCount(t *testing.T) {
	f, err := os.Open(outpath)
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
	fmt.Println("Count:", count)
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
