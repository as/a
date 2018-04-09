package main

import (
	"github.com/as/edit"
	"github.com/as/frame"
	"github.com/as/ui/win"
)

func usedarkcolors() {
	// TODO(as): shouldn't have to edit the frame package here
	frame.A.Text = frame.MTextW
	frame.A.Back = frame.MBodyW
	frame.ATag0.Text = frame.MTextW
	frame.ATag0.Back = frame.MTagG
	frame.ATag1.Text = frame.MTextW
	frame.ATag1.Back = frame.MTagC
}

func (g *Grid) acolor(e edit.File) {
	// TODO(as): O(n*m) -> O(1)
	if t := g.FindName(e.Name); t != nil {
		if t.Body == nil {
			return
		}
		win := t.Body.(*win.Win)
		if win == nil {
			return
		}
		fr := win.Frame
		p0, p1 := e.Q0, e.Q1
		p0 -= clamp(p0-t.Body.Origin(), 0, fr.Len())
		p1 -= clamp(p1-t.Body.Origin(), 0, fr.Len())
		fr.Recolor(fr.PointOf(p0), p0, p1, frame.Mono.Palette)
		fr.Mark()
	}
}
