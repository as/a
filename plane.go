package main

import (
	"image"
)

type Plane interface {
	Loc() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
	Refresh()
}

func eq(a, b Plane) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Loc() == b.Loc()
}

func sizeof(r image.Rectangle) image.Point {
	return r.Max.Sub(r.Min)
}
 