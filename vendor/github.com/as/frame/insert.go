package frame

import (
	"image"
	"image/draw"
)

// Insert inserts the contents of s at index p0 in
// the frame and returns the number of characters
// written.
func (f *Frame) Insert(s []byte, p0 int64) (wrote int) {
	if p0 > f.Nchars || len(s) == 0 || f.b == nil {
		return
	}

	// find p0, it's box, and its point in the box its in
	b0 := f.Find(0, 0, p0)
	//	eb := f.StartCell(b0)
	//ob0 := b0
	cb0 := p0
	b1 := b0
	pt0 := f.pointOfBox(p0, b0)

	// find p1
	ppt0, pt1 := f.boxscan(s, pt0)
	opt0 := pt0
	ppt1 := pt1
	// Line wrap
	if b0 < f.Nbox {
		b := &f.Box[b0]
		pt0 = f.wrapMax(pt0, b)
		ppt1 = f.wrapMin(ppt1, b)
	}
	f.modified = true

	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}

	cb0, b0, pt0, pt1 = f.boxalign(cb0, b0, pt0, pt1)
	f.boxpush(p0, b0, b1, pt0, pt1, ppt1)
	f.bitblt(cb0, b0, pt0, pt1, opt0)
	text, back := f.pick(p0, f.p0+1, f.p1+1)
	f.Paint(ppt0, ppt1, back)
	f.redrawRun0(f.ir, ppt0, text, back)
	f.Run.Combine(f.ir, b1)
	if b1 > 0 && f.Box[b1-1].Nrune >= 0 && ppt0.X-f.Box[b1-1].Width >= f.r.Min.X {
		b1--
		ppt0.X -= f.Box[b1].Width
	}
	b0 += f.ir.Nbox
	if b0 < f.Nbox-1 {
		b0++
	}
	f.clean(ppt0, b1, b0)
	f.Nchars += f.ir.Nchars
	if p0 <= f.p0 {
		f.p0 += f.ir.Nchars
		f.p1 += f.ir.Nchars
	} else if p0 < f.p1 {
		f.p1 += f.ir.Nchars
	}
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	f.badElasticAlg()
	return int(f.ir.Nchars)
}

func (f *Frame) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	if f == nil {
		panic("nil frame")
	}
	if f.Drawer == nil {
		panic("nil drawer")
	}
	f.Drawer.Draw(dst, r, src, sp, op)
	f.Flush(r)
}
func (f *Frame) badElasticAlg() {
	if f.elastic() {
		if f.Nbox <= 1 {
			return
		}
		b := 0
		b1 := 0
		for b != f.Nbox {
			b1 = b
			b = f.NextCell(b)
			if b == 0 || b == b1 {
				break
			}
		}
		f.Stretch(f.Nbox)
		b = b1
		for ; b1 > 1; b1 = f.Stretch(b1) {
			if b == b1 {
				break
			}
			b = b1
		}
		f.Stretch(b1)
		f.Refresh()
	}
}

// Mark marks the frame as dirty
func (f *Frame) Mark() {
	f.modified = true
}

// boxalign collects a list of pts of each box
// on the frame before and after an insertion occurs
func (f *Frame) boxalign(cb0 int64, b0 int, pt0, pt1 image.Point) (int64, int, image.Point, image.Point) {
	type Pts [2]image.Point
	f.pts = f.pts[:0]

	// collect the start pts for each box on the plane
	// before and after the insertion
	for {
		if pt0.X == pt1.X || pt1.Y == f.r.Max.Y || b0 == f.Nbox {
			break
		}
		b := &f.Box[b0]
		pt0 = f.wrapMax(pt0, b)
		pt1 = f.wrapMin(pt1, b)
		if b.Nrune > 0 {
			if n := f.fits(pt1, b); n != b.Nrune {
				f.Split(b0, n)
				b = &f.Box[b0]
			}
		}

		// early exit - point went off the frame
		if pt1.Y == f.r.Max.Y {
			break
		}

		f.pts = append(f.pts, Pts{pt0, pt1})

		pt0 = f.advance(pt0, b)
		pt1.X += f.plot(pt1, b)

		cb0 += int64(b.Len())
		b0++
	}
	return cb0, b0, pt0, pt1
}

// boxpush moves boxes down the frame to make room for an insertion
// from pt0 to pt1
func (f *Frame) boxpush(p0 int64, b0, b1 int, pt0, pt1, ppt1 image.Point) {
	h := f.Face.Dy()
	// delete boxes that ran off the frame
	// and update the char count
	if pt1.Y == f.r.Max.Y && b0 < f.Nbox {
		f.Nchars -= f.Count(b0)
		f.Run.Delete(b0, f.Nbox-1)
	}

	// update the line count
	if b0 == f.Nbox {
		f.Nlines = (pt1.Y - f.r.Min.Y) / h
		if pt1.X > f.r.Min.X {
			f.Nlines++
		}
		return
	}

	if pt1.Y == pt0.Y {
		// insertion won't propagate down
		return
	}

	qt0 := pt0.Y + h
	qt1 := pt1.Y + h
	f.Nlines += (qt1 - qt0) / h
	if f.Nlines > f.maxlines {
		f.trim(ppt1, p0, b1)
	}

	// shift down the existing boxes
	// on the bitmap
	if r := f.r; pt1.Y < r.Max.Y {
		r.Min.Y = qt1

		// rectangular group of boxes
		if qt1 < f.r.Max.Y {
			f.Draw(f.b, r, f.b, image.Pt(f.r.Min.X, qt0), f.op)
		}

		// partial line
		r.Min = pt1
		r.Max.X = pt1.X + (f.r.Max.X - pt0.X)
		r.Max.Y = qt1
		f.Draw(f.b, r, f.b, pt0, f.op)
	}
}

func (f *Frame) bitblt(cb0 int64, b0 int, pt0, pt1, opt0 image.Point) (res image.Rectangle) {
	h := f.Face.Dy()
	y := 0
	if pt1.Y == f.r.Max.Y {
		y = pt1.Y
	}
	x := len(f.pts)
	if x != 0 {
		res = image.Rectangle{pt0, pt1}
		res.Canon().Max.Add(f.pts[x-1][0])
	}
	run := f.Box[b0-x:]
	x--
	_, back := f.pick(cb0, f.p0, f.p1)
	for ; x >= 0; x-- {
		b := &run[x]
		br := image.Rect(0, 0, b.Width, h)
		pt := f.pts[x]
		if b.Nrune > 0 {
			f.Draw(f.b, br.Add(pt[1]), f.b, pt[0], f.op)
			// clear bit hanging off right
			if x == 0 && pt[1].Y > pt0.Y {

				// new char was wider than the old
				// one so the line wrapped anyway

				_, back = f.pick(cb0, f.p0, f.p1)
				r := br.Add(opt0)
				r.Max.X = f.r.Max.X
				f.Draw(f.b, r, back, r.Min, f.op)

			} else if pt[1].Y < y {

				// copy from left to right, bottom to top

				_, back = f.pick(cb0, f.p0, f.p1)
				r := image.ZR.Add(pt[1])
				r.Min.X += b.Width
				r.Max.X += f.r.Max.X
				r.Max.Y += h
				f.Draw(f.b, r, back, r.Min, f.op)

			}
			y = pt[1].Y
			cb0 -= int64(b.Nrune)
		} else {
			r := br.Add(pt[1])
			if r.Max.X >= f.r.Max.X {
				r.Max.X = f.r.Max.X
			}
			cb0--
			_, back = f.pick(cb0, f.p0, f.p1)
			f.Draw(f.b, r, back, r.Min, f.op)
			y = 0
			if pt[1].X == f.r.Min.X {
				y = pt[1].Y
			}
		}
	}
	return res
}

func (f *Frame) pick(c, p0, p1 int64) (text, back image.Image) {
	if p0 <= c && c < p1 {
		return f.Color.Hi.Text, f.Color.Hi.Back
	}
	return f.Color.Text, f.Color.Back
}

/*
//
// Below ideas

type offset struct {
	p0   int64
	b0   int
	cb0  int64
	b1   int64
	pt0  image.Point
	pt1  image.Point
	opt0 image.Point
}

func (f *Frame) zInsertElastic(s []byte, p0 int64) (wrote int) {
	type Pts [2]image.Point
	if p0 > f.Nchars || len(s) == 0 || f.b == nil {
		return
	}

	// find p0, it's box, and its point in the box its in
	b0 := f.Find(0, 0, p0)
	cb0 := p0

	eb0 := f.StartCell(b0)
	pt0 := f.pointOfBox(p0, b0)
	ppt0, pt1 := f.boxscan(s, pt0)
	//	ept0 := f.pointOfBox(p0, eb0)
	f.Box = append(f.Box[:b0], append(f.ir.Box[:f.ir.Nbox], f.Box[b0:]...)...)
	f.Nbox += f.ir.Nbox

	b1 := b0

	for bn := f.NextCell(b1 + f.ir.Nbox); bn > eb0; bn = f.Stretch(bn) {
	}
	f.Stretch(eb0)

	//		pt0 = f.pointOfBox(p0, b0)
	//		pt1 = f.pointOfBox(65535, b1)

	opt0 := pt0
	ppt1 := pt1
	// Line wrap
	if b0 < f.Nbox {
		b := &f.Box[b0]
		pt0 = f.wrapMax(pt0, b)
		ppt1 = f.wrapMin(ppt1, b)
	}
	f.modified = true

	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}

	cb0, b0, pt0, pt1 = f.boxalign(cb0, b0+f.ir.Nbox, pt0, pt1)
	f.boxpush(p0, b0+f.ir.Nbox, b1+f.ir.Nbox, pt0, pt1, ppt1)
	f.bitblt(cb0, b0, pt0, pt1, opt0)
	text, back := f.pick(p0, f.p0+1, f.p1+1)
	f.Paint(ppt0, ppt1, back)
	f.redrawRun0(f.ir, ppt0, text, back)

	b1 = b0
	if b1 > 0 && f.Box[b1-1].Nrune >= 0 && ppt0.X-f.Box[b1-1].Width >= f.r.Min.X {
		b1--
		ppt0.X -= f.Box[b1].Width
	}

	b0 += f.ir.Nbox
	if b0 < f.Nbox-1 {
		b0++
	}
	f.clean(f.pointOfBox(p0, b1), b0, b1)
	f.Nchars += f.ir.Nchars

	f.p0, f.p1 = coInsert(p0, p0+f.Nchars, f.p0, f.p1)
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	return int(f.ir.Nchars)

}
*/
