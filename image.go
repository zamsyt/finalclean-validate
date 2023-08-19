package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
)

var palette = loadPalette("r-slash-place-2023.gpl")

func posterize(img image.Image, p color.Palette) *image.Paletted {
	b := img.Bounds()
	paletted := image.NewPaletted(b, p)
	draw.Draw(paletted, b, img, b.Min, draw.Src)
	return paletted
}

// Opaque color
type rgb struct {
	r, g, b uint8
}

// Implement color.Color
func (c rgb) RGBA() (r, g, b, a uint32) {
	return color.RGBA{c.r, c.g, c.b, 255}.RGBA()
}

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

func genDiff(img0, img1 image.Image) (diff image.Image, count int) {
	b := img0.Bounds()
	if !b.Eq(img1.Bounds()) {
		log.Fatal("Images have different sizes:", b, img1.Bounds())
	}
	diff = image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c0 := img0.At(x, y)
			c1 := img1.At(x, y)
			if !colorEq(c0, c1) {
				diff.(*image.NRGBA).Set(x, y, c1)
				count++
			}
		}
	}
	return
}

func colorEq(a, b color.Color) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB &&
		aA == bA)
}

func getImg(r io.ReadCloser) image.Image {
	img, _, err := image.Decode(r)
	check(err)
	check(r.Close())
	return img
}

func savePng(img image.Image, path string) {
	f, err := os.Create(path)
	check(err)
	png.Encode(f, img)
	f.Close()
}
