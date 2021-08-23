package col

import (
	"image"
)

type Tile interface {
	Plane
	Axis
	Kid(n int) Plane
	Len() int
	IDPoint(image.Point) int
	attach(Plane, int)
	detach(int) Plane
}

func Attach(t Tile, src Plane, pt image.Point) {
	if t.Len() == 0 {
		pt = t.Minor(image.ZP)
		src.Move(pt)
		t.attach(src, 0)
		Fill(t)
		return
	}
	did := t.IDPoint(pt)
	pt = t.Minor(pt)
	src.Move(pt)
	t.attach(src, did+1)
	Fill(t)
}

func Fill(t Tile) {
	if t.Len() == 0 {
		return
	}
	for n := 0; n != t.Len(); n++ {
		pt := delta(t, n)
		if pt == image.ZP {
			//panic("zp")
		}
		//		defer t.Kid(n).Resize(pt)
		t.Kid(n).Resize(pt)
	}
}

func delta(c Tile, n int) image.Point {
	y0 := c.Minor(c.Area().Min)
	y1 := c.Major(c.Area().Max)
	if n != c.Len() {
		y0 = c.Minor(c.Kid(n).Bounds().Min)
	}
	if n+1 != c.Len() {
		y1 = c.Major(c.Kid(n + 1).Bounds().Min)
	}
	return y1.Sub(y0)
}

func Detach(t Tile, id int) Plane {
	return t.detach(id)
}

//func (co *Col) Attach(src Plane, y int) {
//	Attach(co, src, image.Pt(0, y))
//}

func (co *Col) Detach(id int) Plane {
	return Detach(co, id)
}

func identity(x, y int) image.Point {
	if x == 0 || y == 0 {
		return image.ZP
	}
	return image.Pt(x, y)
}
