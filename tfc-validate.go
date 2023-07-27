package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

const oradir string = "ora-unzipped"

var palette = loadPalette("r-slash-place-2023.gpl")

var layersToIgnore = []string{
	"Sector Map",
	"Noise Removal",
	"BASE LAYER",
}

func main() {
	od := parseOraStack()
	for i := len(od.Layers) - 1; i >= 0; i-- {
		checkLayer(od.Layers[i])
	}
}

func checkLayer(l Layer) {
	if contains(layersToIgnore, l.Name) {
		return
	}
	if l.Opacity != 1.0 {
		fmt.Printf(`[WARNING] "%v" has opacity %v`+"\n", l.Name, l.Opacity)
	}
	img := getImg(filepath.Join(oradir, l.Src))
	b := img.Bounds()
	var tc int  // translucent pixel counter
	var npc int // non-palette pixel counter
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			_, _, _, a := c.RGBA()
			if a > 0 && a < 0xffff {
				tc++
				//fmt.Println(c)
			}
			if !rgbEq(c, palette.Convert(c)) {
				npc++
			}
		}
	}
	if tc > 0 {
		fmt.Printf(`[WARNING] "%v" has %v translucent pixels`+"\n", l.Name, tc)
	}
	if npc > 0 {
		fmt.Printf(`[WARNING] "%v" has %v pixels not matching the palette`+"\n", l.Name, npc)
	}
}

func rgbEq(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB)
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
		p = append(p, color.RGBA{r, g, b, 255})
	}
	return p
}

func getImg(path string) image.Image {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	img, _, err := image.Decode(f)
	check(err)
	return img
}

func savePng(img image.Image, path string) {
	f, err := os.Create(path)
	check(err)
	png.Encode(f, img)
	f.Close()
}

func parseOraStack() *oraImage {
	stackxml, err := os.ReadFile(filepath.Join(oradir, "stack.xml"))
	check(err)
	var oi oraImage
	xml.Unmarshal(stackxml, &oi)
	return &oi
}

type oraImage struct {
	W      int     `xml:"w,attr"`
	H      int     `xml:"h,attr"`
	Layers []Layer `xml:"stack>layer"`
}

type Layer struct {
	Src     string  `xml:"src,attr"`
	Name    string  `xml:"name,attr"`
	Opacity float32 `xml:"opacity,attr"`
	X       int     `xml:"x,attr"`
	Y       int     `xml:"y,attr"`
}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
