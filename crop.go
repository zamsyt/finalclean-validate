package main

import (
	"image"
	"image/color"
	"image/draw"
)

func isTransparent(c color.Color) bool {
	return alpha(c) == 0
}

func crop(img image.Image) image.Image {
	b := img.Bounds()
	minX := b.Max.X
	maxX := b.Min.X
	minY := b.Max.Y
	maxY := b.Min.Y

minY:
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if !isTransparent(img.At(x, y)) {
				minY = y
				break minY
			}
		}
	}

minX:
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := minY; y < b.Max.Y; y++ {
			if !isTransparent(img.At(x, y)) {
				minX = x
				break minX
			}
		}
	}

maxX:
	for x := b.Max.X; x > minX; x-- {
		for y := minY; y < b.Max.Y; y++ {
			if !isTransparent(img.At(x, y)) {
				maxX = x + 1
				break maxX
			}
		}
	}

maxY:
	for y := b.Max.Y; y > minY; y-- {
		for x := minX; x < maxX; x++ {
			if !isTransparent(img.At(x, y)) {
				maxY = y + 1
				break maxY
			}
		}
	}

	cropped := image.NewNRGBA(image.Rect(minX, minY, maxX, maxY))
	draw.Draw(cropped, cropped.Bounds(), img, image.Point{minX, minY}, draw.Src)

	return cropped
}
