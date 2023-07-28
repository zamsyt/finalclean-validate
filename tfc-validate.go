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
	//"Noise Removal",
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
		fmt.Printf("[WARN] %q has opacity %v\n", l.Name, l.Opacity)
	}
	img := getImg(filepath.Join(oradir, l.Src))
	total, wrong, _, correctTranslucent, diff := checkColors(img)
	if wrong > 0 {
		fmt.Printf("[WARN] %q has %v/%v pixels not matching the palette", l.Name, wrong, total)
		if correctTranslucent > 0 {
			fmt.Printf(". %v are translucent but otherwise correct", correctTranslucent)
		}
		fmt.Println()
		os.MkdirAll("diff", 0755)
		savePng(diff, filepath.Join("diff", l.Name+"_"+filepath.Base(l.Src)))
	}
	if total == 0 {
		fmt.Printf("[INFO] %q is empty\n", l.Name)
	}
}

// Opaque color
type rgb struct {
	r, g, b uint8
}

// Implement color.Color
func (c rgb) RGBA() (r, g, b, a uint32) {
	return color.RGBA{c.r, c.g, c.b, 255}.RGBA()
}

func checkColors(img image.Image) (total, wrong, translucent, correctTranslucent int, diff image.Image) {
	b := img.Bounds()
	diff = image.NewPaletted(b, []color.Color{color.Transparent, rgb{255, 0, 0}})
	var transparent int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			a := alpha(c)
			if a == 0 {
				transparent++
				continue
			}
			palettized := palette.Convert(c)
			if colorEq(c, palettized) {
				continue
			}
			wrong++
			if a < 255 {
				translucent++
				if rgbEq(c, palettized) {
					correctTranslucent++
				}
			}
			diff.(*image.Paletted).SetColorIndex(x, y, 1)
		}
	}
	area := b.Dx() * b.Dy()
	total = area - transparent
	return
}

func alpha(c color.Color) uint8 {
	switch v := c.(type) {
	case color.NRGBA:
		return v.A
	case rgb:
		return 255
	default:
		panic(fmt.Sprintf("Unexpected color type %T\n", v))
	}
}

func rgbEq(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB)
}

func colorEq(a, b color.Color) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB &&
		aA == bA)
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
