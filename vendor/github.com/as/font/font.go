package font

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func Fix(i fixed.Int26_6) int {
	return i.Ceil()
}

type Face interface {
	font.Face
	Ruler
}
type Ruler interface {
	Ascent() int
	Descent() int
	Height() int
	Letting() int
	Stride() int
	Dy() int
	Dx(p []byte) (dx int)
	Fits(p []byte, limitDx int) (n int)
}

func Open(f font.Face) Face {
	if f == nil {
		panic("open: nil face")
	}
	switch f := f.(type) {
	case Face:
		return f
	}
	m := f.Metrics()
	a := m.Ascent.Ceil()
	h := m.Height.Ceil()
	d := m.Descent.Ceil()
	dy := h + h/2
	l := dy / 2
	return &face{
		Face: f,
		s:    0,
		a:    a,
		d:    d,
		h:    h,
		l:    l,
		dy:   dy,
	}
}

type face struct {
	h, a, d, l, dy, s int
	font.Face
}

func (f *face) Stride() int  { return f.s }
func (f *face) Letting() int { return f.l }
func (f *face) Height() int  { return f.h }
func (f *face) Ascent() int  { return f.a }
func (f *face) Descent() int { return f.d }
func (f *face) Dy() int      { return f.dy }
func (f *face) Dx(p []byte) (dx int) {
	for _, c := range p {
		w, _ := f.Face.GlyphAdvance(rune(c))
		dx += Fix(w)
	}
	return dx + f.Stride()*len(p)
}
func (f *face) Fits(p []byte, limitDx int) (n int) {
	var c byte
	stride := f.Stride()
	for n, c = range p {
		w, _ := f.Face.GlyphAdvance(rune(c))
		limitDx -= Fix(w) + stride
		if limitDx < 0 {
			return n
		}
	}
	return n
}

/*
// NewBasic always returns a 7x13 basic font
func NewRaster(f Face, size int) *Font {
	ft := &Font{
		Face:    f,
		size:    size,
		ascent:  2,
		descent: 1,
		letting: 0,
		stride:  0,
	}
	ft.dy = ft.ascent + ft.descent + ft.size
	hexFt := makefont(gomono.TTF, ft.Dy()/4+3)
	ft.hexDx = ft.genChar('_').Bounds().Dx()
	for i := 0; i != 256; i++ {
		ft.cache[i] = ft.genChar(byte(i))
		if ft.cache[i] == nil {
			ft.cache[i] = hexFt.genHexChar(ft.Dy(), byte(i))
		}
	}
	return ft
}

func makefont(data []byte, size int) *Font {
	reply := make(chan interface{})
	fontIRQ <- fontPKT{
		id:    string(crc32.NewIEEE().Sum(data)),
		reply: reply,
		data:  data,
	}
	rx := <-reply
	switch rx := rx.(type) {
	case error:
		println(rx)
		return nil
	case *truetype.Font:
		ft := FromFace(truetype.NewFace(rx,
			&truetype.Options{
				Size:              float64(size),
				GlyphCacheEntries: 512,
				SubPixelsX:        1,
			}), size)
		ft.data = data
		ft.dy = ft.ascent + ft.descent + ft.size
		return ft
	}
	panic("makefont")
}

func (f *Font) genChar(b byte) *Glyph {
	dr, mask, maskp, adv, _ := f.Glyph(fixed.P(0, f.size), rune(b))
	if !f.Printable(b) {
		return nil
	}
	r := image.Rect(0, 0, Fix(adv), f.Dy())
	m := image.NewAlpha(r)
	r = r.Add(image.Pt(dr.Min.X, dr.Min.Y))
	draw.Draw(m, r, mask, maskp, draw.Src)
	return &Glyph{mask: m, Rectangle: m.Bounds()}
}

func (f *Font) genHexChar(dy int, b byte) *Glyph {
	s := fmt.Sprintf("%02x", b)
	g0 := f.genChar(s[0])
	g1 := f.genChar(s[1])
	r := image.Rect(2, f.descent+f.ascent, g0.Bounds().Dx()+g1.Bounds().Dx()+6, dy)
	m := image.NewAlpha(r)
	draw.Draw(m, r, g0.Mask(), image.ZP, draw.Over)
	r.Min.X += g0.Mask().Bounds().Dx()
	draw.Draw(m, r.Add(image.Pt(-f.descent/4, f.descent*2)), g1.Mask(), image.ZP, draw.Over)
	return &Glyph{mask: m, Rectangle: m.Bounds()}
}
*/
