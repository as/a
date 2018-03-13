package main

import (
	"image"
	"time"

	"github.com/as/ui/win"
	"golang.org/x/mobile/event/mouse"
)

const (
	scrollX = 10
	sbWidth = 10
)

var (
	tagHeight   = *ftsize*2 + *ftsize/2 - 2
	winSize     = image.Pt(1024, 768)
	pad         = image.Pt(15, 15)
	sizerR      = image.Rect(0, 0, scrollX, tagHeight)
	dcPerimeter = image.Rect(-4, -4, 4, 4)
)

func doubleclick(pt0, pt1 image.Point, deadline time.Time) bool {
	return !time.Now().After(deadline) && pt1.In(dcPerimeter.Add(pt0))
}

func absP(e mouse.Event, sp image.Point) image.Point {
	return p(e).Add(sp)
}

func relP(e mouse.Event, sp image.Point) image.Point {
	return p(e).Sub(sp)
}

func canopy(pt image.Point) bool {
	return pt.Y > g.sp.Y+g.tdy && pt.Y < g.sp.Y+g.tdy*2
}

func p(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func inSizer(pt image.Point) bool {
	return pt.In(sizerR)
}

func inScroll(pt image.Point) bool {
	return pt.X >= 0 && pt.X < sbWidth
}

func yRegion(y, ymin, ymax int) int {
	if y < ymin {
		return 1
	}
	if y > ymax {
		return -1
	}
	return 0
}

func sweep(w *win.Win, e mouse.Event, s, q0, q1 int64) (int64, int64, int64) {
	r := image.Rectangle{image.ZP, w.Size()}
	y := int(e.Y)
	padY := 15
	lo := r.Min.Y + padY
	hi := r.Dy() - padY
	units := w.Bounds().Dy()
	reg := yRegion(y, lo, hi)

	if reg != 0 {
		if reg == 1 {
			w.Scroll(-((lo-y)%units + 1) * 3)
		} else {
			w.Scroll(+((y-hi)%units + 1) * 3)
		}
	}
	q := w.IndexOf(image.Pt(int(e.X), int(e.Y))) + w.Origin()
	if q0 == s {
		if q < q0 {
			return q0, q, q0
		}
		return q0, q0, q
	}
	if q > q1 {
		return q1, q1, q
	}
	return q1, q, q1
}
