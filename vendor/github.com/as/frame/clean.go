package frame

import (
	"image"
)

func (f *Frame) clean(pt image.Point, n0, n1 int) {
	c := f.r.Max.X
	for ; n0 < n1-1; n0++ {
		b0 := &f.Box[n0]
		b1 := &f.Box[n0+1]
		pt = f.wrapMax(pt, b0)
		for b0.Nrune >= 0 && n0 < n1-1 && b1.Nrune >= 0 && pt.X+b0.Width+b1.Width < c {
			f.Merge(n0)
			n1--
		}

		pt = f.advance(pt, &f.Box[n0]) // dont simplify this
	}

	for ; n0 < f.Nbox; n0++ {
		b0 := &f.Box[n0]
		pt = f.wrapMax(pt, b0)
		pt = f.advance(pt, b0)
	}

	f.full = 0
	if pt.Y >= f.r.Max.Y {
		f.full = 1
	}
}
