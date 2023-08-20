package main

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

var labPalette []colorful.Color

func init() {
	labPalette = make([]colorful.Color, len(palette))
	for i, pc := range palette {
		lc, ok := colorful.MakeColor(pc)
		if !ok {
			panic("Couldn't convert color")
		}
		labPalette[i] = lc
	}
}

func toPaletteLab(c color.Color) color.Color {
	ez := palette.Convert(c)
	if colorEq(c, ez) {
		return ez
	}
	cc, ok := colorful.MakeColor(c)
	if !ok {
		return color.Transparent
	}
	var min float64 = 1000.0
	var best colorful.Color
	for _, pc := range labPalette {
		colorful.MakeColor(pc)
		d := cc.DistanceCIE94(pc)
		if d < min {
			min = d
			best = pc
		}
		if min == 0 {
			break
		}
	}
	if min == 1000.0 {
		panic("?")
	}
	return best
}
