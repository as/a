package frame

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame/box"
)

const (
	// strict enables panic on the condition that the frame is too small to fix
	// any characters
	strict = false
)

// bxscan resets the measuring function and calls Bxscan in the embedded run
func (f *Frame) boxscan(s []byte, pt image.Point) (image.Point, image.Point) {
	switch f.Face.(type) {
	case font.Rune:
		f.ir.Runescan(s, f.maxlines)
	case interface{}:
		f.ir.Boxscan(s, f.maxlines)
	}
	if f.elastic() {
		// TODO(as): remove this after adding tests since its redundant
		//
		// Just to see if the algorithm works not ideal to sift through all of
		// the boxes per insertion, although surprisingly faster than expected
		// to the point of where its almost unnoticable without the print
		// statements
		bn := f.ir.Nbox
		for bn > 0 {
			bn = f.ir.Stretch(bn)
		}
		f.ir.Stretch(bn)
	}
	pt = f.wrapMin(pt, &f.ir.Box[0])
	return pt, f.boxscan2D(f.ir, pt)
}

func (f *Frame) boxscan2D(r *box.Run, pt image.Point) image.Point {
	n := 0
	for nb := 0; nb < r.Nbox; nb++ {
		b := &r.Box[nb]
		pt = f.wrapMin(pt, b)
		if pt.Y == f.r.Max.Y {
			r.Nchars -= r.Count(nb)
			r.Delete(nb, r.Nbox-1)
			break
		}
		if b.Nrune > 0 {
			if n = f.fits(pt, b); n == 0 {
				if strict {
					panic("boxscan2D: fits 0")
				}
				return pt
			}
			if n != b.Nrune {
				r.Split(nb, n)
				b = &r.Box[nb]
			}
			pt.X += b.Width
		} else {
			if b.Break() == '\n' {
				pt = f.wrap(pt)
			} else {
				pt.X += f.plot(pt, b)
			}
		}
	}
	return pt
}
