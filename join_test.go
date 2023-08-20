package main

import (
	"image"
	"os"
	"testing"
)

func TestJoin(t *testing.T) {
	if !imagesEq(getImgFile(basediffpath), getImgFile("combined.png")) {
		t.Error("Combined image doesn't match original")
	} /*else {
		fmt.Println("combined image matches!")
	}*/
}

func getImgFile(path string) image.Image {
	f, err := os.Open(path)
	check(err)
	return getImg(f)
}

func imagesEq(img0, img1 image.Image) bool {
	b0 := img0.Bounds()
	b1 := img1.Bounds()
	if !b0.Eq(b1) {
		return false
	}
	for y := b0.Min.Y; y < b0.Max.Y; y++ {
		for x := b0.Min.X; x < b0.Max.X; x++ {
			if !colorEq(img0.At(x, y), img1.At(x, y)) {
				return false
			}
		}
	}
	return true
}
