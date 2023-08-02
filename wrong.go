package main

import (
	"image"
	"image/color"
)

func genDiff(img image.Image) (diff image.Image, wrong, total int) {
	var dp = []color.Color{
		color.Transparent,
		rgb{0x22, 0x00, 0x80},
		rgb{0x19, 0x51, 0xBA},
		rgb{0xFF, 0xF1, 0x52},
		rgb{0xFF, 0x79, 0x29},
		rgb{0xFF, 0x00, 0x00},
	}
	b := img.Bounds()
	diff = image.NewPaletted(b, dp)
	//diff = image.NewNRGBA(b)
	//var overflow int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			if alpha(c) == 0 {
				continue
			}
			total++
			if !colorEq(c, palette.Convert(c)) {
				//diff.(draw.Image).Set(x, y, c)
				index := simpleColorDiff(c, palette.Convert(c))/12 + 1
				if int(index) >= len(dp) {
					//overflow++
					index = uint8(len(dp) - 1)
				}
				diff.(*image.Paletted).SetColorIndex(x, y, index)
				wrong++
			}
		}
	}
	//fmt.Println("overflow:", overflow)
	return
}

func alpha(c color.Color) uint8 {
	switch v := c.(type) {
	case color.NRGBA:
		return v.A
	case rgb:
		return 255
	default:
		return color.RGBAModel.Convert(c).(color.RGBA).A
		//panic(fmt.Sprintf("Unexpected color type %T\n", v))
	}
}

func colorEq(a, b color.Color) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB &&
		aA == bA)
}

func simpleColorDiff(a, b color.Color) uint8 {
	na := color.NRGBAModel.Convert(a).(color.NRGBA)
	nb := color.NRGBAModel.Convert(b).(color.NRGBA)
	return max(diff(na.R, nb.R), diff(na.G, nb.G), diff(na.B, nb.B), diff(na.A, nb.A))
}

type ordered interface {
	uint8 | uint32 | int
}

func diff[T ordered](a, b T) T {
	if a > b {
		return a - b
	}
	return b - a
}
func max[T ordered](nums ...T) T {
	var largest T = 0
	for _, v := range nums {
		if v > largest {
			largest = v
		}
	}
	return largest
}
