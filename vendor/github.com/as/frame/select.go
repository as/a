package frame

import (
	"image"
	"image/draw"
)

// Paint paints the color col on the frame at points pt0-pt1. The result is a Z shaped fill
// consisting of at-most 3 rectangles. No text is redrawn.
func (f *Frame) Paint(p0, p1 image.Point, col image.Image) {
	if f.b == nil {
		panic("selectpaint: b == 0")
	}
	if f.r.Max.Y == p0.Y {
		return
	}
	h := f.Face.Dy()
	q0, q1 := p0, p1
	q0.Y += h
	q1.Y += h
	n := (p1.Y - p0.Y) / h

	if n == 0 { // one line
		f.Draw(f.b, image.Rectangle{p0, q1}, col, image.ZP, draw.Over)
	} else {
		if p0.X >= f.r.Max.X {
			p0.X = f.r.Max.X // - 1
		}
		f.Draw(f.b, image.Rect(p0.X, p0.Y, f.r.Max.X, q0.Y), col, image.ZP, draw.Over)
		if n > 1 {
			f.Draw(f.b, image.Rect(f.r.Min.X, q0.Y, f.r.Max.X, p1.Y), col, image.ZP, draw.Over)
		}
		f.Draw(f.b, image.Rect(f.r.Min.X, p1.Y, q1.X, q1.Y), col, image.ZP, draw.Over)
	}
}

// Select selects the region [p0:p1). The operation highlights
// the range of text under that region. If p0 = p1, a tick is
// drawn to indicate a null selection.
func (f *Frame) Select(p0, p1 int64) {
	pp0, pp1 := f.Dot()
	if pp1 <= p0 || p1 <= pp0 || p0 == p1 || pp1 == pp0 {
		f.Redraw(f.PointOf(pp0), pp0, pp1, false)
		f.Redraw(f.PointOf(p0), p0, p1, true)
	} else {
		if p0 < pp0 {
			f.Redraw(f.PointOf(p0), p0, pp0, true)
		} else if p0 > pp0 {
			f.Redraw(f.PointOf(pp0), pp0, p0, false)
		}
		if p1 > pp1 {
			f.Redraw(f.PointOf(pp1), pp1, p1, true)
		} else if p1 < pp1 {
			f.Redraw(f.PointOf(p1), p1, pp1, false)
		}
	}
	f.modified = true
	f.p0, f.p1 = p0, p1
}
