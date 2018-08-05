package main

import (
	"image"

	"github.com/as/shiny/event/mouse"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

const (
	scrollX = 10
	sbWidth = 10
)

var (
	tagHeight = tag.Height(*ftsize)
	sizerR    = image.Rect(0, 0, scrollX, tagHeight)
)

func rel(e mouse.Event, p Plane) mouse.Event {
	pt := p.Bounds().Min
	e.X -= float32(pt.X)
	e.Y -= float32(pt.Y)
	return e
}

func canopy(pt image.Point) bool {
	r := g.Bounds()
	r.Max.Y = r.Min.Y + g.Tag.Bounds().Dy()*2
	return pt.In(r)
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

func clamp(v, l, h int64) int64 {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}

func sweep(w *win.Win, e mouse.Event, s, q0, q1 int64) (int64, int64, int64) {
	r := image.Rectangle{image.ZP, w.Size()}
	y := int(e.Y)
	padY := tagHeight
	lo := r.Min.Y + padY
	hi := r.Dy() - padY
	units := w.Bounds().Dy()
	if units == 0 {
		units++
	}
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
