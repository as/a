package frame

import (
	"image"

	_ "github.com/as/etch"
)

// Delete deletes the range [p0:p1) and
// returns the number of characters deleted
func (f *Frame) Delete(p0, p1 int64) int {

	if p0 >= f.Nchars || p0 == p1 || f.b == nil {
		return 0
	}

	if p1 > f.Nchars {
		p1 = f.Nchars
	}
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}
	n0 := f.Find(0, 0, p0)
	//	eb := f.StartCell(n0)
	//	eb = eb
	nn0 := n0
	n1 := f.Find(n0, p0, p1)
	pt0 := f.pointOfBox(p0, n0)
	ppt0 := pt0
	pt1 := f.PointOf(p1)
	f.Free(n0, n1-1)
	f.modified = true

	// Advance forward, copying the first un-deleted box
	// on the right all the way to the left, splitting them
	// when necessary to fit on a wrapped line. A bit of draw
	// computation is saved by keeping track of the selection
	// and interpolating its drawing routine into the same
	// loop.
	//
	// Might have to rethink this when adding support for
	// multiple selections.
	//
	// pt0/pt1: deletion start/stop
	// n0/n1: deleted box/first surviving box
	// int64(p1): char index of the surviving box

	pt0, pt1, n0, n1 = f.delete(pt0, pt1, n0, n1, int64(p1))

	if n1 == f.Nbox && pt0.X != pt1.X {
		f.Paint(pt0, pt1, f.Color.Back)
	}

	// Delete more than a line. All the boxes have been shifted
	// but the bitmap might still have a copy of them down below
	if pt1.Y != pt0.Y {
		pt0, pt1, n1 = f.fixTrailer(pt0, pt1, n1)
	}

	f.Run.Close(n0, n1-1)
	if nn0 > 0 && f.Box[nn0-1].Nrune >= 0 && ppt0.X-f.Box[nn0-1].Width >= f.r.Min.X {
		nn0--
		ppt0.X -= f.Box[nn0].Width
	}

	if n0 < f.Nbox-1 {
		n0++
	}
	f.clean(ppt0, nn0, n0)

	if f.p1 > p1 {
		f.p1 -= p1 - p0
	} else if f.p1 > p0 {
		f.p1 = p0
	}
	if f.p0 > p1 {
		f.p0 -= p1 - p0
	} else if f.p0 > p0 {
		f.p0 = p0
	}

	f.Nchars -= p1 - p0
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	pt0 = f.PointOf(f.Nchars)
	extra := 0
	if pt0.X > f.r.Min.X {
		extra = 1
	}
	h := f.Face.Dy()
	f.Nlines = (pt0.Y-f.r.Min.Y)/h + extra
	f.badElasticAlg()
	return int(p1 - p0) //n - f.Nlines
}
func (f *Frame) delete(pt0, pt1 image.Point, n0, n1 int, cn1 int64) (image.Point, image.Point, int, int) {
	h := f.Face.Dy()
	for pt1.X != pt0.X && n1 < f.Nbox {
		b := &f.Box[n1]
		pt0 = f.wrapMin(pt0, b)
		pt1 = f.wrapMax(pt1, b)
		r := image.Rectangle{pt0, pt0}
		r.Max.Y += h

		if b.Nrune > 0 { // non-newline
			n := f.fits(pt0, b)
			if n != b.Nrune {
				f.Split(n1, n)
				b = &f.Box[n1]
			}
			r.Max.X += b.Width
			f.Draw(f.b, r, f.b, pt1, f.op)
			//drawBorder(f.b, r.Inset(-4), Green, image.ZP, 8)
			cn1 += int64(b.Nrune)
		} else {
			r.Max.X = min(r.Max.X+f.project(pt0, b), f.r.Max.X)
			_, col := f.pick(cn1, f.p0, f.p1)
			f.Draw(f.b, r, col, pt0, f.op)
			cn1++
		}
		pt1 = f.advance(pt1, b)
		pt0.X += f.plot(pt0, b)
		f.Box[n0] = f.Box[n1]
		n0++
		n1++
	}
	return pt0, pt1, n0, n1
}
func (f *Frame) fixTrailer(pt0, pt1 image.Point, n1 int) (image.Point, image.Point, int) {
	if n1 == f.Nbox && pt0.X != pt1.X {
		f.Paint(pt0, pt1, f.Color.Back)
	}
	h := f.Face.Dy()
	pt2 := f.pointOf(65536, pt1, n1)
	if pt2.Y > f.r.Max.Y {
		pt2.Y = f.r.Max.Y - h
	}
	if n1 < f.Nbox {
		q0 := pt0.Y + h
		q1 := pt1.Y + h
		q2 := pt2.Y + h
		if q2 > f.r.Max.Y {
			q2 = f.r.Max.Y
		}
		f.Draw(f.b, image.Rect(pt0.X, pt0.Y, pt0.X+(f.r.Max.X-pt1.X), q0), f.b, pt1, f.op)
		f.Draw(f.b, image.Rect(f.r.Min.X, q0, f.r.Max.X, q0+(q2-q1)), f.b, image.Pt(f.r.Min.X, q1), f.op)
		f.Paint(image.Pt(pt2.X, pt2.Y-(pt1.Y-pt0.Y)), pt2, f.Color.Back)
	} else {
		f.Paint(pt0, pt2, f.Color.Back)
	}
	return pt0, pt1, n1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
