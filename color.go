package main

import (
	"github.com/as/edit"
	"github.com/as/frame"
	"github.com/as/ui/win"
)

var (
	Tag0  = frame.ATag0
	Tag1  = frame.ATag1
	Tag2  = frame.ATag1
	Body2 = frame.A
)

func usedarkcolors() {
	Body2.Text = frame.MTextW
	Body2.Back = frame.MBodyW

	Tag0.Text = frame.MTextW
	Tag0.Back = frame.MTagG

	Tag1.Text = frame.MTextW
	Tag1.Back = frame.MTagC

	Tag2.Text = frame.MTextW
	Tag2.Back = frame.MTagC

	GridConfig.Color[0] = Tag0
	ColConfig.Color[0] = Tag1
	TagConfig.Color[0] = Tag2
	TagConfig.Color[1] = Body2

	SB := frame.Color{
		Palette: frame.Palette{
			Text: frame.MTagC,
			Back: frame.MTagG,
		},
	}
	GridConfig.Color[2] = SB
	ColConfig.Color[2] = SB
	TagConfig.Color[2] = SB
}

func (g *Grid) acolor(e edit.File) {
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
