package main

import (
	"github.com/as/edit"
	"github.com/as/frame"
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
}

type color32 uint32

func (c color32) RGBA() (r, g, b, a uint32) {
	u := uint32(c)
	r, g, b, a = u>>24, u<<8>>24, u<<16>>24, u<<24>>24
	return r << 8, g << 8, b << 8, a << 8
}

var (
	Alpha     color32 = 0xff
	White             = 0xffffff00 + Alpha
	Black             = 0x00000000 + Alpha
	LightGray         = 0xb0b0c000 + Alpha
	Gray              = 0x1c1f2600 + Alpha
	DarkGray          = 0x120a1400 + Alpha
	Yellow            = 0xfffffd00 + Alpha
	Red               = 0xffe8ef00 + Alpha
	Green             = 0xefffe800 + Alpha
	Blue              = 0xe8efff00 + Alpha
	Peach             = 0xfff8e800 + Alpha
	Strata            = 0xf8f2f800 + Alpha
	Storm             = 0xd8d8e800 + Alpha
	Mauve             = 0x9090C000 + Alpha

	Paleblue  = 0xf3f8fe00 + Alpha
	Palegreen = 0xe2ebe800 + Alpha
	Palegray  = 0xe2e1e800 + Alpha
	Palepink  = 0xfce8fc00 + Alpha

	Blueviolet   = 0x66558800 + Alpha
	Bluegray     = 0x2b323b00 + Alpha
	Darkbluegray = 0x1c1f2600 + Alpha

	Seagreen = 0x99cc9900 + Alpha
)

var (
	Palette      = Lightpalette
	Lightpalette = map[string]frame.Color{
		"grid": frame.NewUniform(Gray, Storm, White, Mauve),
		"col":  frame.NewUniform(Gray, Strata, White, Mauve),
		"tag":  frame.NewUniform(Gray, Strata, White, Mauve),
		"win":  frame.NewUniform(Gray, Peach, White, Mauve),
	}
	Darkpalette = map[string]frame.Color{
		"grid": frame.NewUniform(LightGray, Darkbluegray, Mauve, White),
		"col":  frame.NewUniform(LightGray, Gray, Mauve, White),
		"tag":  frame.NewUniform(LightGray, Gray, Mauve, White),
		"win":  frame.NewUniform(LightGray, Bluegray, Mauve, White),
	}
)

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
