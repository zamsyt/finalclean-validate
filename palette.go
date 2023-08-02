package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

const diffpath = "diff"
const fixedpath = "fix"

func checkPalette(args []string) {
	if len(args) == 0 {
		fmt.Println("not enough arguments")
		fmt.Println("Usage:\n\tcheckpalette <my-drawpile.ora>")
		os.Exit(1)
	}
	o := OpenORA(args[0])
	f, err := o.zip.Open("mergedimage.png")
	check(err)
	diff, wrong, total := genDiff(getImg(f))
	fmt.Printf("%v/%v pixels differ from the palette in the merged image\n", wrong, total)
	savePng(diff, "merged-diff.png")

	os.MkdirAll(diffpath, 0755)
	for _, l := range o.Layers {
		checkLayer(o, l)
	}
}

var layersToIgnore = []string{
	"Sector Map",
	//"Noise Removal",
	//"BASE LAYER",
}

var cropFlag bool
var diffFlag bool
var alphaTresholdFlag uint8

func init() {
	pflag.BoolVar(&cropFlag, "crop", false, "")
	pflag.BoolVar(&diffFlag, "diff", false, "")
	pflag.Uint8Var(&alphaTresholdFlag, "alpha-treshold", 0, "")
}

func checkLayer(o *ora, l oraLayer) {
	if contains(layersToIgnore, l.Name) {
		return
	}
	if l.Opacity != 1.0 {
		fmt.Printf("[WARN] %q has opacity %v\n", l.Name, l.Opacity)
	}
	f, err := o.zip.Open(l.Src)
	check(err)
	diff, wrong, total := genDiff(getImg(f))
	if total == 0 {
		fmt.Printf("[INFO] %q is empty\n", l.Name)
		return
	}
	if wrong == 0 {
		return
	}
	fmt.Printf("[WARN] %q has %v/%v pixels not matching the palette\n", l.Name, wrong, total)
	if cropFlag {
		diff = crop(diff)
	}
	savePng(diff, filepath.Join(diffpath, l.Name+"_"+filepath.Base(l.Src)))
}

func fixPalette(args []string) {
	if len(args) == 0 {
		fmt.Println("not enough arguments")
		fmt.Println("Usage:\n\tfixpalette <my-drawpile.ora>")
		os.Exit(1)
	}
	o := OpenORA(args[0])
	os.MkdirAll(fixedpath, 0755)
	for _, l := range o.Layers {
		if contains(layersToIgnore, l.Name) {
			continue
		}
		f, err := o.zip.Open(l.Src)
		check(err)
		img := getImg(f)
		var fixed image.Image
		fixed, changed, toTransparent, toOpaque := posterize(img, palette, alphaTresholdFlag)
		if changed == 0 {
			//if !imagesEq(img, fixed) { panic("nothing shuould've changed") }
			continue
		}
		fmt.Printf("%v: %v pixels changed", l.Name, changed)
		if toTransparent > 0 || toOpaque > 0 {
			fmt.Printf(". %v translucent pixels converted to opaque, %v to transparent", toOpaque, toTransparent)
		}
		fmt.Println()
		if cropFlag {
			fixed = crop(fixed)
		}
		savePng(fixed, filepath.Join(fixedpath, l.Name+"_"+filepath.Base(l.Src)))
	}
}

func posterize(img image.Image, p color.Palette, alphaTreshold uint8) (paletted *image.Paletted, changed, toTransparent, toOpaque int) {
	tp := append([]color.Color{color.Transparent}, p...)
	b := img.Bounds()
	paletted = image.NewPaletted(b, tp)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			if alpha(c) == 0 {
				continue
			}
			if alpha(c) < alphaTreshold {
				toTransparent++
				changed++
				continue
			}
			if alpha(c) < 255 {
				toOpaque++
			}
			converted := p.Convert(c)
			if !colorEq(c, converted) {
				changed++
			} else if diffFlag {
				continue
			}
			paletted.Set(x, y, converted)
		}
	}
	//draw.Draw(paletted, b, img, b.Min, draw.Src)
	return

}

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

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
