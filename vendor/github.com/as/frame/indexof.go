package frame

import (
	"image"
)

// IndexOf returns the chracter index under the
// point pt.
func (f *Frame) IndexOf(pt image.Point) (p int64) {
	pt.X += 1
	return f.indexOf(pt)
}

func (f *Frame) indexOf(pt image.Point) (p int64) {
	pt = f.grid(pt)
	qt := f.r.Min
	bn := 0
	for ; bn < f.Nbox && qt.Y < pt.Y; bn++ {
		b := &f.Box[bn]
		qt = f.wrapMax(qt, b)
		if qt.Y >= pt.Y {
			break
		}
		qt = f.advance(qt, b)
		p += int64(b.Len())
	}

	for ; bn < f.Nbox && qt.X <= pt.X; bn++ {
		b := &f.Box[bn]
		qt = f.wrapMax(qt, b)
		if qt.Y > pt.Y {
			break
		}
		if qt.X+b.Width > pt.X {
			if b.Nrune < 0 {
				qt = f.advance(qt, b)
			} else {
				left := pt.X - qt.X
				p += int64(f.Face.Fits(b.Ptr, left))
				qt.X += left
			}
		} else {
			p += int64(b.Len())
			qt = f.advance(qt, b)
		}
	}
	return p
}
