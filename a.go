package main

import (
	"github.com/as/edit"
	"github.com/as/frame"
)

func (g *Grid) acolor(e edit.File) {
	// TODO(as): O(n*m) -> O(1)
	if t := g.FindName(e.Name); t != nil {
		if t.Body == nil {
			return
		}
		fr := t.Body.Frame
		p0, p1 := e.Q0, e.Q1
		p0 -= clamp(p0-t.Body.Origin(), 0, fr.Len())
		p1 -= clamp(p1-t.Body.Origin(), 0, fr.Len())
		fr.Recolor(fr.PointOf(p0), p0, p1, frame.Mono.Palette)
		fr.Mark()
	}
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
