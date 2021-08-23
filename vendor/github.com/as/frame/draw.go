package frame

import (
	"image"

	"github.com/as/frame/box"
)

// Refresh renders the entire frame, including the underlying
// bitmap. Refresh should not be called after insertion and deletion
// unless the frame's RGBA bitmap was painted over by another
// draw operation.
func (f *Frame) Refresh() {
	cols := f.Color
	if f.p0 == f.p1 {
		ticked := f.Ticked
		if ticked {
			f.tickat(f.PointOf(f.p0), false)
		}
		f.drawsel(f.PointOf(0), 0, f.Nchars, cols.Back, cols.Text)
		if ticked {
			f.tickat(f.PointOf(f.p0), true)
		}
		return
	}
	pt := f.PointOf(0)
	pt = f.drawsel(pt, 0, f.p0, cols.Back, cols.Text)
	pt = f.drawsel(pt, f.p0, f.p1, cols.Hi.Back, cols.Hi.Text)
	f.drawsel(pt, f.p1, f.Nchars, cols.Back, cols.Text)
}

// RedrawAt renders the frame's bitmap starting at pt and working downwards.
func (f *Frame) RedrawAt(pt image.Point, text, back image.Image) {
	f.redrawRun0(&(f.Run), pt, text, back)
}

// Redraw draws the range [p0:p1] at the given pt.
func (f *Frame) Redraw(pt image.Point, p0, p1 int64, issel bool) {
	if f.Ticked {
		f.tickat(f.PointOf(f.p0), false)
	}

	if p0 == p1 {
		f.tickat(pt, issel)
		return
	}

	pal := f.Color.Palette
	if issel {
		pal = f.Color.Hi
	}
	f.drawsel(pt, p0, p1, pal.Back, pal.Text)
}

// Recolor redraws the range p0:p1 with the given palette
func (f *Frame) Recolor(pt image.Point, p0, p1 int64, cols Palette) {
	f.drawsel(pt, p0, p1, cols.Back, cols.Text)
	f.modified = true
}

func (f *Frame) redrawRun0(r *box.Run, pt image.Point, text, back image.Image) image.Point {
	nb := 0
	for ; nb < r.Nbox; nb++ {
		b := &r.Box[nb]
		pt = f.wrapMax(pt, b)
		//if !f.noredraw && b.nrune >= 0 {
		if b.Nrune >= 0 {
			f.StringBG(f.b, pt, text, image.ZP, f.Face, b.Ptr, back, image.ZP)
		}
		pt.X += b.Width
	}
	return pt
}

// widthBox returns the width of box n. If the length of
// alt is different than the box, alt is measured and
// returned instead.
func (f *Frame) widthBox(b *box.Box, alt []byte) int {
	if b.Nrune < 0 || len(alt) == b.Len() {
		return b.Width
	}
	return f.Face.Dx(alt)
}

func (f *Frame) drawsel(pt image.Point, p0, p1 int64, back, text image.Image) image.Point {
	{
		// doubled
		p0, p1 := int(p0), int(p1)
		q0 := 0
		trim := false

		var (
			nb int
			b  *box.Box
		)

		for nb = 0; nb < f.Nbox; nb++ {
			b = &f.Box[nb]
			l := q0 + b.Len()
			if l > p0 {
				break
			}
			q0 = l
		}

		// Step into box, start coloring it
		// How much does this lambda slow things down?
		stepFill := func() {
			qt := pt
			pt = f.wrapMax(pt, b)
			if pt.Y > qt.Y {
				r := image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y)
				f.Draw(f.b, r, back, qt, f.op)
			}
		}
		for ; nb < f.Nbox && q0 < p1; nb++ {
			b = &f.Box[nb]
			if q0 >= p0 { // region 0 or 1 or 2
				stepFill()
			}
			ptr := b.Ptr[:b.Len()]
			if q0 < p0 {
				// region -1: shift p right inside the selection
				ptr = ptr[p0-q0:]
				q0 = p0
			}

			trim = false
			if q1 := q0 + len(ptr); q1 > p1 {
				// region 1: would draw too much, retract the selection
				lim := len(ptr) - (q1 - p1)
				ptr = ptr[:lim]
				trim = true
			}
			w := f.widthBox(b, ptr)
			if b.Nrune > 0 {
				f.Draw(f.b, image.Rect(pt.X, pt.Y, min(pt.X+w, f.r.Max.X), pt.Y+f.Face.Dy()), back, pt, f.op)
				f.StringBG(f.b, pt, text, image.ZP, f.Face, ptr, back, image.ZP)
			} else {
				f.Draw(f.b, image.Rect(pt.X, pt.Y, min(pt.X+w, f.r.Max.X), pt.Y+f.Face.Dy()), back, pt, f.op)
			}
			pt.X += w

			if q0 += len(ptr); q0 > p1 {
				break
			}
		}
		if p1 > p0 && nb != 0 && nb < f.Nbox && f.Box[nb-1].Len() > 0 && !trim {
			b = &f.Box[nb]
			stepFill()
		}
		return pt
	}
}
