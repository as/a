package font

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
)

type Cliche interface {
	Cache
	LoadBox(s []byte, fg, bg color.Color) image.Image
}

type boxsig struct {
	b  string
	fg color.RGBA
	bg color.RGBA
}

func NewCliche(f font.Face) Cliche {
	if f, ok := f.(Cliche); ok {
		return f
	}
	return &cliche{
		Cache: NewCache(f),
		cache: make(map[boxsig]image.Image),
	}
}

type cliche struct {
	Cache
	cache map[boxsig]image.Image
}

func (c *cliche) LoadBox(b []byte, fg, bg color.Color) image.Image {
	if len(b) == 0 {
		return nil
	}
	sig := boxsig{string(b), convert(fg), convert(bg)}
	if img, ok := c.cache[sig]; ok {
		return img
	}
	var list []image.Image
	dx := 0
	for _, v := range b {
		img := c.LoadGlyph(rune(v), fg, bg)
		dx += img.Bounds().Dx()
		list = append(list, img)
	}
	r := list[0].Bounds()
	r.Max.X += dx
	boximg := image.NewRGBA(r)
	for _, img := range list {
		dx := img.Bounds().Dx()
		r.Max.X += dx
		draw.Draw(boximg, r, img, image.ZP, draw.Src)
		r.Min.X += dx
	}
	c.cache[sig] = boximg
	return boximg
}
