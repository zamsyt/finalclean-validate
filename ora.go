package main

import (
	"archive/zip"
	"encoding/xml"
	"image"
	"image/draw"
	"io"
	"log"
)

// Flat ora file (no nested stacks)
type oraStack struct {
	W      int        `xml:"w,attr"`
	H      int        `xml:"h,attr"`
	Layers []oraLayer `xml:"stack>layer"`
}

type oraLayer struct {
	Src     string  `xml:"src,attr"`
	Name    string  `xml:"name,attr"`
	Opacity float32 `xml:"opacity,attr"`
	X       int     `xml:"x,attr"`
	Y       int     `xml:"y,attr"`
}

type ora struct {
	zip *zip.ReadCloser
	oraStack
}

// func (o *ora) Close() error { return o.zip.Close() }
func OpenORA(path string) *ora {
	r, err := zip.OpenReader(path)
	check(err)
	stackxml, err := r.Open("stack.xml")
	check(err)
	data, err := io.ReadAll(stackxml)
	check(err)
	check(stackxml.Close())
	var stack oraStack
	check(xml.Unmarshal(data, &stack))
	return &ora{r, stack}
}

func (o *ora) Layer(name string) image.Image {
	for _, l := range o.Layers {
		if l.Name == name {
			f, err := o.zip.Open(l.Src)
			check(err)
			return resizeCanvas(getImg(f), image.Rect(0, 0, o.W, o.H), image.Pt(l.X, l.Y))
		}
	}
	log.Printf("ora: image %q not found\n", name)
	return nil
}

func resizeCanvas(img image.Image, r image.Rectangle, offset image.Point) image.Image {
	resized := image.NewNRGBA(r)
	b := img.Bounds()
	draw.Draw(resized, b.Add(offset), img, b.Min, draw.Src)
	return resized
}
