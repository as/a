package main

import (
	"image"
)

type Tile interface {
	Delta(int) image.Point
	Kid(n int) Plane
	Len() int
}

func fill(t Tile) {
	if t.Len() == 0 {
		return
	}
	for n := 0; n != t.Len(); n++ {
		pt := t.Delta(n)
		if pt == image.ZP {
			panic("zp")
		}
		k := t.Kid(n)
		k.Resize(pt)
	}
}

func identity(x, y int) image.Point {
	if x == 0 || y == 0 {
		return image.ZP
	}
	return image.Pt(x, y)
}
