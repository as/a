package frame

import (
	"image"
	"image/draw"

	. "github.com/as/font"
	"golang.org/x/image/font"
)

// Drawer implements the set of methods a frame needs to draw on a draw.Image. The frame's default behavior is to use
// the native image/draw package and x/exp/font packages to satisfy this interface.
type Drawer interface {
	Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op)
	//DrawMask(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op)

	// StringBG draws a string to dst at point p
	StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int

	// Flush requests that prior calls to the draw and string methods are flushed from an underlying soft-screen. The list of rectangles provide
	// optional residency information. Implementations may refresh a superset of r, or ignore it entirely, as long as the entire region is
	// refreshed
	Flush(r ...image.Rectangle) error
}

func NewDefaultDrawer() Drawer {
	return &defaultDrawer{}
}

type defaultDrawer struct{}

func (d *defaultDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	draw.Draw(dst, r, src, sp, op)
}

func (d *defaultDrawer) StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int {
	return StringBG(dst, p, src, sp, ft, s, bg, bgp)
}

func (d *defaultDrawer) Flush(r ...image.Rectangle) error {
	return nil
}

func negotiateFace(f font.Face, flags int) Face {
	if flags&FrUTF8 != 0 {
		return NewCache(NewRune(f))
	}
	switch f := f.(type) {
	case Face:
		return f
	case font.Face:
		return Open(f)
	}
	return Open(f)
}
