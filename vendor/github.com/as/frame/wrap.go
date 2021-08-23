package frame

import (
	"image"

	"github.com/as/frame/box"
)

// wrapMax returns the point where b should go on a plane
// if b doesn't fit entirely on the plane at pt, wrapMax
// returns a pt on the next line
func (f *Frame) wrapMax(pt image.Point, b *box.Box) image.Point {
	width := b.Width
	if b.Nrune < 0 {
		width = b.Minwidth
	}
	if width > f.r.Max.X-pt.X {
		return f.wrap(pt)
	}
	return pt
}

// wrapMin is like wrapMax, except it lazily wraps lines if
// no chars in the box fit on the plane at pt.
func (f *Frame) wrapMin(pt image.Point, b *box.Box) image.Point {
	if f.fits(pt, b) == 0 {
		return f.wrap(pt)
	}
	return pt
}

func (f *Frame) wrap(pt image.Point) image.Point {
	pt.X = f.r.Min.X
	pt.Y += f.Face.Dy()
	return pt
}

func (f *Frame) advance(pt image.Point, b *box.Box) (x image.Point) {
	if b.Nrune < 0 && b.Break() == '\n' {
		pt = f.wrap(pt)
	} else {
		pt.X += b.Width
	}
	return pt
}

// fits returns the number of runes that can fit on the line at pt. A newline yields 1.
func (f *Frame) fits(pt image.Point, b *box.Box) (nr int) {
	left := f.r.Max.X - pt.X
	if b.Nrune < 0 {
		if b.Minwidth <= left {
			return 1
		}
		return 0
	}
	if left >= b.Width {
		return b.Nrune
	}
	return f.Face.Fits(b.Ptr, left)
}
func (f *Frame) plot(pt image.Point, b *box.Box) int {
	b.Width = f.project(pt, b)
	return b.Width
}
func (f *Frame) project(pt image.Point, b *box.Box) int {
	c := f.r.Max.X
	x := pt.X
	if b.Nrune >= 0 || b.Break() != '\t' { //
		return b.Width
	}
	if f.elastic() && b.Break() == '\t' {
		return b.Minwidth
	}
	if x+b.Minwidth > c {
		pt.X = f.r.Min.X
		x = pt.X
	}
	x += f.maxtab
	x -= (x - f.r.Min.X) % f.maxtab
	if x-pt.X < b.Minwidth || x > c {
		x = pt.X + b.Minwidth
	}
	return x - pt.X
}
