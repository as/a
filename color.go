package main

import (
	"github.com/as/edit"
	"github.com/as/frame"
	. "github.com/as/rgba"
	"github.com/as/ui/reg"
	"github.com/as/ui/win"
)

func uselightcolors() {
	Palette = Lightpalette
	updatecolors()
}

func usedarkcolors() {
	Palette = Darkpalette
	updatecolors()
}

func updatecolors() {
	GridConfig.Color[0] = Palette["grid"]
	ColConfig.Color[0] = Palette["col"]
	TagConfig.Color[0] = Palette["tag"]
	TagConfig.Color[1] = Palette["win"]
	reg.PutImage(Shade, "scroll/fg")
	reg.PutImage(Gray, "scroll/bg")
}

var (
	Palette      = Lightpalette
	Lightpalette = map[string]frame.Color{
		"grid": frame.Theme(Gray, Storm, White, Mauve),
		"col":  frame.Theme(Gray, Strata, White, Mauve),
		"tag":  frame.Theme(Gray, Strata, White, Mauve),
		"win":  frame.Theme(Gray, Peach, White, Mauve),
	}
	Darkpalette = map[string]frame.Color{
		"grid": frame.Theme(LightGray, Darkbluegray, White, Darkbluegray),
		"col":  frame.Theme(LightGray, Gray, White, Darkbluegray),
		"tag":  frame.Theme(LightGray, Gray, White, Darkbluegray),
		"win":  frame.Theme(LightGray, Bluegray, White, Darkbluegray),
	}
)

func (g *Grid) acolor(e edit.File) {
	if t := g.FindName(e.Name); t != nil {
		if t.Window == nil {
			return
		}
		win := t.Window.(*win.Win)
		if win == nil {
			return
		}
		fr := win.Frame
		p0, p1 := e.Q0, e.Q1
		p0 -= clamp(p0-t.Window.Origin(), 0, fr.Len())
		p1 -= clamp(p1-t.Window.Origin(), 0, fr.Len())
		fr.Recolor(fr.PointOf(p0), p0, p1, frame.Mono.Palette)
		fr.Mark()
	}
}
