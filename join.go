package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
)

func join(args []string) {
	dir := splitpath
	if len(args) > 0 {
		dir = os.Args[0]
	}

	combined := image.NewNRGBA(image.Rect(0, 0, 3000, 2000))

	for major := 'A'; major <= 'X'; major++ {
		for minor := 1; minor <= 16; minor++ {
			f, err := os.Open(filepath.Join(dir, fmt.Sprintf("%c%d.png", major, minor)))
			check(err)
			simg := getImg(f)
			sb := simg.Bounds()
			draw.Draw(combined, minorRect(major, minor), simg, sb.Min, draw.Src)
		}
	}

	savePng(combined, "combined.png")

	//entries, err := os.ReadDir(args[0])
	//check(err)
	//for _, e := range entries {
	//	e.Name()
	//}
}
