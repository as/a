package font

import (
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func NewHex(dy int) Face {
	f := &hex{}
	f.genHexChars(dy)
	return f
}

type hex struct {
	glyphs                [256]*image.Alpha
	h, a, d, l, dx, dy, s int
}

func (f *hex) Stride() int  { return f.s }
func (f *hex) Letting() int { return f.l }
func (f *hex) Height() int  { return f.h }
func (f *hex) Ascent() int  { return f.a }
func (f *hex) Descent() int { return f.d }
func (f *hex) Dy() int      { return f.dy }

func (f *hex) Dx(p []byte) (dx int) {
	return (f.dx * f.s) + (f.dx * len(p))
}

func (f *hex) Fits(p []byte, limitDx int) (n int) {
	n = limitDx / f.dx
	if n > len(p) {
		n = len(p)
	}
	return n
}
func (f *hex) Close() error {
	return nil
}
func (f *hex) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if r > 255 {
		return image.ZR, nil, image.ZP, 0, false
	}
	dot0 := image.Pt(dot.X.Ceil(), dot.Y.Ceil())
	dot0 = dot0.Add(image.Pt(0, -10))
	img := f.glyphs[byte(r)]
	r0 := img.Bounds()
	dr.Max.X = r0.Dx()
	dr.Max.Y = r0.Dy()
	dr = dr.Add(dot0)
	return dr, img, image.ZP, fixed.I(f.dx), true
}
func (f *hex) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if r > 255 {
		return
	}
	r0 := f.glyphs[byte(r)].Bounds()
	bounds.Max = bounds.Max.Add(fixed.P(r0.Dx(), r0.Dy()))
	return bounds, fixed.I(f.dx), true
}
func (f *hex) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if r > 255 {
		return
	}
	//r0 := f.glyphs[byte(r)].Bounds()
	return fixed.I(f.dx), true //fixed.I(r0.Dx()), true
}
func (f *hex) Kern(r0, r1 rune) fixed.Int26_6 { return 0 }
func (f *hex) Metrics() (m font.Metrics)      { return }
func (f *hex) genHexChars(dy int) {

	var helper [16]*image.Alpha

	{
		ft := NewGoMedium(dy/5 + dy/3 + 3)

		for i, c := range []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'} {
			dr, mask, maskp, adv, _ := ft.Glyph(fixed.P(0, ft.Height()), c)
			r := image.Rect(0, 0, Fix(adv), ft.Dy())
			m := image.NewAlpha(r)
			r = r.Add(image.Pt(dr.Min.X, dr.Min.Y))
			draw.Draw(m, r, mask, maskp, draw.Src)
			helper[i] = m
		}
	}

	d0 := f.Descent()
	for i := 0; i != 256; i++ {
		g0 := helper[i/16]
		g1 := helper[i%16]
		r := image.Rect(2, d0, g0.Bounds().Dx()+g1.Bounds().Dx()+7, dy-3)
		m := image.NewAlpha(r)
		draw.Draw(m, r.Add(image.Pt(2, 0)), g0, image.ZP, draw.Over)
		r.Min.X += g0.Bounds().Dx()
		draw.Draw(m, r.Add(image.Pt(-d0/4+2, d0-d0*2)), g1, image.ZP, draw.Over)
		f.glyphs[i] = m
	}

	ft := NewGoMedium(dy)
	m := ft.Metrics()
	f.a = m.Ascent.Ceil()
	f.h = m.Height.Ceil()
	f.d = m.Descent.Ceil()
	f.dy = f.h + f.h/2
	f.l = dy / 2
	f.dx = f.glyphs[0].Bounds().Dx()
}
