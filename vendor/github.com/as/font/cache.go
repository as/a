package font

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Cache interface {
	Face
	LoadGlyph(r rune, fg, bg color.Color) image.Image
}

func NewCache(f font.Face) Cache {
	var f0 Face
	switch f := f.(type) {
	case Cache:
		return f
	case Rune:
		return newRuneCache(f)
	case Face:
		f0 = f
	case font.Face:
		f0 = Open(f)
	}
	cf := &cachedFace{
		Face:  f0,
		cache: make(map[signature]*image.RGBA),
	}
	for i := range cf.cachewidth {
		cf.cachewidth[i] = f0.Dx([]byte{byte(i)})
	}
	return cf
}

type cachedFace struct {
	Face
	cache      map[signature]*image.RGBA
	cachewidth [256]int
}

type signature struct {
	r  rune
	fg color.RGBA
	bg color.RGBA
}

func (f *cachedFace) LoadGlyph(r rune, fg, bg color.Color) image.Image {
	sig := signature{r, convert(fg), convert(bg)}
	if img, ok := f.cache[sig]; ok {
		return img
	}
	mask, r0 := f.genChar(r)
	img := image.NewRGBA(r0)
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.ZP, draw.Src)
	draw.DrawMask(img, img.Bounds(), image.NewUniform(fg), image.ZP, mask, r0.Min, draw.Over)
	f.cache[sig] = img
	if int(r) < len(f.cache) {
		f.cachewidth[byte(r)] = f.Dx([]byte{byte(r)})
	}
	return img
}

func (f *cachedFace) Fits(p []byte, limitDx int) (n int) {
	var c byte
	for n, c = range p {
		limitDx -= f.cachewidth[c]
		if limitDx < 0 {
			return n
		}
	}
	return n
}

func (f *cachedFace) Dx(p []byte) (dx int) {
	for _, c := range p {
		dx += f.cachewidth[c]
	}
	return dx
}

func (f *cachedFace) genChar(r rune) (*image.Alpha, image.Rectangle) {
	dr, mask, maskp, adv, _ := f.Face.Glyph(fixed.P(0, f.Height()), r)
	r0 := image.Rect(0, 0, Fix(adv), f.Dy())
	m := image.NewAlpha(r0)
	r0 = r0.Add(image.Pt(dr.Min.X, dr.Min.Y))
	draw.Draw(m, r0, mask, maskp, draw.Src)
	return m, m.Bounds()
}

func convert(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
}
