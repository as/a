package font

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"

	"golang.org/x/image/math/fixed"
)

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int {
	if bg == nil {
		return StringNBG(dst, p, src, sp, ft, s)
	}
	if fg, bg, ok := canCache(src, bg); ok {
		switch ft := ft.(type) {
		case Cliche:
			img := ft.LoadBox(s, fg, bg)
			dr := img.Bounds().Add(p)
			draw.Draw(dst, dr, img, img.Bounds().Min, draw.Src)
			return dr.Dx()
		case Cache:
			switch ft := ft.(type) {
			case Rune:
				return staticRuneBG(dst, p, ft.(Cache), s, fg, bg)
			}
			return staticStringBG(dst, p, ft, s, fg, bg)

		}
	}
	switch ft := ft.(type) {
	case *runeface:
		return runeBG(dst, p, src, sp, ft, s, bg, bgp)
	case Face:
		return stringBG(dst, p, src, sp, ft, s, bg, bgp)
	}
	return stringBG(dst, p, src, sp, Open(ft), s, bg, bgp)
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte) int {
	var (
		f  Face
		ok bool
	)
	if f, ok = ft.(Face); !ok {
		f = Open(ft)
	}
	p.Y += f.Height()
	for _, b := range s {
		dr, mask, maskp, advance, _ := f.Glyph(fixed.P(p.X, p.Y), rune(b))
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
	}
	return p.X
}

func canCache(f image.Image, b image.Image) (fg, bg color.Color, ok bool) {
	if f, ok := f.(*image.Uniform); ok {
		if b, ok := b.(*image.Uniform); ok {
			return f.C, b.C, true
		}
	}
	return fg, bg, false
}

func runeBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft Face, s []byte, bg image.Image, bgp image.Point) int {
	p.Y += ft.Height()
	for _, b := range string(s) {
		dr, mask, maskp, advance, _ := ft.Glyph(fixed.P(p.X, p.Y), b)
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
	}
	return p.X
}
func staticRuneBG(dst draw.Image, p image.Point, ft Cache, s []byte, fg, bg color.Color) int {
	r := image.Rectangle{p, p}
	r.Max.Y += ft.Dy()

	for _, b := range string(s) {
		img := ft.LoadGlyph(b, fg, bg)
		dx := img.Bounds().Dx()
		r.Max.X += dx
		draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)
		r.Min.X += dx
	}
	return r.Min.X - p.X
}

func stringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft Face, s []byte, bg image.Image, bgp image.Point) int {
	p.Y += ft.Height()
	for _, b := range s {
		dr, mask, maskp, advance, _ := ft.Glyph(fixed.P(p.X, p.Y), rune(b))
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
	}
	return p.X
}

func staticStringBG(dst draw.Image, p image.Point, ft Cache, s []byte, fg, bg color.Color) int {
	r := image.Rectangle{p, p}
	r.Max.Y += ft.Dy()

	for _, b := range s {
		img := ft.LoadGlyph(rune(b), fg, bg)
		dx := img.Bounds().Dx()
		r.Max.X += dx
		draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)
		r.Min.X += dx
	}
	return r.Min.X - p.X
}

/*
func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		//draw.Draw(dst, r.Add(p), bg, bgp, draw.Src)
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx()
	}
	return p.X
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx()
	}
	return p.X
}

func RuneBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			panic("RuneBG")
		}
		//draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
		if len(s)-size == 0 {
			break
		}
		s = s[size:]
	}
	return p.X
}

func RuneNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			panic("RuneBG")
		}
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
		if len(s)-size == 0 {
			break
		}
		s = s[size:]
	}
	return p.X
}
*/
