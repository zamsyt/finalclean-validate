package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
)

const splitpath = "split"

type croppable interface {
	SubImage(r image.Rectangle) image.Image
}

func split(args []string) {
	if len(args) == 0 {
		fmt.Println("not enough arguments")
		fmt.Println("Usage:\n\tsplit <my-image.png>")
		os.Exit(1)
	}
	f, err := os.Open(args[0])
	check(err)
	img, ok := getImg(f).(croppable)
	if !ok {
		log.Fatalf("Image (%T) not croppable", img)
	}
	os.MkdirAll(splitpath, 0755)
	for major := 'A'; major <= 'X'; major++ {
		for minor := 1; minor <= 16; minor++ {
			//fmt.Printf("%c", major)
			r := minorRect(major, minor)
			sectorImg := img.SubImage(r)
			savePng(sectorImg, filepath.Join(splitpath, fmt.Sprintf("%c%d.png", major, minor)))
		}
	}
}

func sector(x, y int) (major rune, minor int) {
	major = rune('A' + y/500*6 + x/500)
	remX, remY := x%500, y%500
	minor = 1 + remY/125*4 + remX/125
	return
}

func sectorXY(major rune, minor int) (x, y int) {
	bx := int(major-'A') % 6
	by := int(major-'A') / 6

	sx := (minor - 1) % 4
	sy := (minor - 1) / 4

	return bx*4 + sx, by*4 + sy
}

func minorRect(major rune, minor int) image.Rectangle {
	bx := int(major-'A') % 6
	by := int(major-'A') / 6

	sx := (minor - 1) % 4
	sy := (minor - 1) / 4

	x0 := bx*500 + sx*125
	y0 := by*500 + sy*125
	x1 := bx*500 + sx*125 + 125
	y1 := by*500 + sy*125 + 125

	return image.Rect(x0, y0, x1, y1)
}

func majorRect(major rune) image.Rectangle {
	x := int(major-'A') % 6
	y := int(major-'A') / 6
	return image.Rect(x*500, y*500, x*500+500, y*500+500)
}
