package main

import (
	"testing"
)

func TestSector(t *testing.T) {
	for i := 'A'; i <= 'X'; i++ {
		//fmt.Println(string(i))
		br := majorRect(i)
		for j := 1; j <= 16; j++ {
			//fmt.Println(string(i), j)
			//fmt.Println(j)
			r := minorRect(i, j)
			if !r.In(br) {
				t.Error()
				return
			}
			xx, yy := sectorXY(i, j)
			if xx*125 != r.Min.X ||
				yy*125 != r.Min.Y {
				t.Error()
				return
				//fmt.Println(xx, yy, x, y, xx*125+x, x-r.Min.X, yy*125+y, y-r.Min.Y)
			}
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					major, minor := sector(x, y)
					if major != i || minor != j {
						t.Error()
						return
					}
				}
			}
		}
	}
}
