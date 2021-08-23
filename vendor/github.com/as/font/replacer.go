package font

import (
	"image"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Replacer returns a face that conditionally replaces a glyph's
// mask and dimensions based on a decision function. For best
// results, the returned Face should be cached with NewCache; the functions
// that compute Fits and Dx are particularly inefficient when used
// without caching.
//
// Note: The current implementation currently returns a cached
// face, but callers shouldn't assume this will always be true
// and use NewCache(NewReplacer(a,b,cond)) anyway.
func Replacer(a, b Face, cond func(r rune) bool) Face {
	if cond == nil {
		cond = func(r rune) bool {
			return r > 127 || r <= 0 || !unicode.IsGraphic(r)
		}
	}
	return NewCache(&repl{
		Face: a,
		b:    b,
		fn:   cond,
	})
}

type repl struct {
	Face
	b  Face
	fn func(r rune) bool
}

func (f *repl) Dx(p []byte) (dx int) {
	for n, c := range p {
		if f.fn(rune(c)) {
			dx += f.b.Dx(p[n : n+1])
		} else {
			dx += f.Face.Dx(p[n : n+1])
		}
	}
	return dx
}

func (f *repl) Fits(p []byte, limitDx int) (n int) {
	var c byte
	for n, c = range p {
		if f.fn(rune(c)) {
			limitDx -= f.b.Dx(p[n : n+1])
		} else {
			limitDx -= f.Face.Dx(p[n : n+1])
		}
		if limitDx < 0 {
			return n
		}
	}
	return n
}
func (f *repl) Close() error {
	return nil
}
func (f *repl) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if f.fn(r) {
		return f.b.Glyph(dot, r)
	}
	return f.Face.Glyph(dot, r)
}
func (f *repl) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if f.fn(r) {
		return f.b.GlyphBounds(r)
	}
	return f.Face.GlyphBounds(r)
}
func (f *repl) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if f.fn(r) {
		return f.b.GlyphAdvance(r)
	}
	return f.Face.GlyphAdvance(r)
}
func (f *repl) Kern(r0, r1 rune) fixed.Int26_6 {
	if f.fn(r0) {
		return f.b.Kern(r0, r1)
	}
	return f.Face.Kern(r0, r1)
}

func (f *repl) Metrics() (m font.Metrics) {
	return f.Face.Metrics()
}
