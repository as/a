package main

import (
	"image"
	"image/color"
	"runtime"

	"github.com/as/edit"
	"github.com/as/frame"
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

var (
	// Transparent and other colors
	// These might not reflect any logical convention or standard
	Transparent = uniform(0x00000000)
	White       = uniform(0xffffffff)
	Black       = uniform(0x000000ff)

	Shade     = uniform(0x121217ff)
	Gray      = uniform(0x1c1f26ff)
	LightGray = uniform(0xb0b0c0ff)
	DarkGray  = uniform(0x120a14ff)

	Yellow = uniform(0xfffffd0f)
	Red    = uniform(0xffd0d8ff)
	Green  = uniform(0xd8ffd0ff)
	Blue   = uniform(0xd8d0ffff)

	Mauve  = uniform(0x9090C0ff)
	Peach  = uniform(0xfff8e8ff)
	Strata = uniform(0xf8f2f8ff)
	Storm  = uniform(0xd8d8e8ff)
	Scroll = uniform(0x9d9da7ff)

	Paleblue  = uniform(0xf3f8feff)
	Palegreen = uniform(0xe2ebe8ff)
	Palegray  = uniform(0xe2e1e8ff)
	Palepink  = uniform(0xfce8fcff)

	Blueviolet   = uniform(0x665588ff)
	Bluegray     = uniform(0x2b323bff)
	Darkbluegray = uniform(0x1c1f26ff)
	Seagreen     = uniform(0x99cc99ff)
)

// Uniform is short for image.NewUniform(Hex(rgba)). On linux,
// we are doing something nasty by pre-swizzling the uniform
// colors. This is until I can fix the swizzle in as/shiny for linux
var uniform = func() func(rgba uint32) *image.Uniform {
	if runtime.GOOS == "linux" {
		return linuxuniform
	}
	return plan9uniform
}()

// plan9uniform is short for image.NewUniform(Hex(rgba))
func plan9uniform(rgba uint32) *image.Uniform {
	return image.NewUniform(hex(rgba))
}

// linuxuniform is short for image.NewUniform(Hex(rgba))\
// this function exists because I haven't fixed swizzle on
// linux yet.
func linuxuniform(rgba uint32) *image.Uniform {
	c := hex(rgba)
	c.R, c.B = c.B, c.R
	return image.NewUniform(c)
}

// hex converts a 32-bit RGBA quad to a color.RGBA
func hex(rgba uint32) color.RGBA {
	return color.RGBA{
		R: uint8(rgba >> 24),
		G: uint8(rgba << 8 >> 24),
		B: uint8(rgba << 16 >> 24),
		A: uint8(rgba << 24 >> 24),
	}
}
