package frame

import (
	"image"
	"image/draw"
)

const (
	TickOff = 0
	TickOn  = 1
)

func (f *Frame) Untick() {
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}
}
func (f *Frame) Tick() {
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), true)
	}
}

func (f *Frame) SetTick(style int) {
	f.tickoff = style == TickOff
}
func mktick(fontY int) (boxw int, linew int) {
	const magic = 12
	boxw = 3 + 1*(fontY/magic)
	for boxw%3 != 0 {
		boxw--
	}
	if boxw < 3 {
		boxw = 3
	}

	linew = boxw / 3
	for boxw%linew != 0 {
		boxw--
	}
	if linew < 1 {
		linew = 1
	}
	return
}

func (f *Frame) tickbg() image.Image {
	return f.Color.Text
	/*
		r, g, b, a := f.Color.Hi.Back.At(0,0).RGBA()
		a=a
		return image.NewUniform(color.RGBA{
			uint8(r),
			uint8(g),
			uint8(b),
			uint8(0),
		})
	*/

}

func (f *Frame) inittick() {

	he := f.Face.Height()
	as := f.Face.Ascent()
	de := f.Face.Descent()
	boxw, linew := mktick(he)
	linew2 := linew / 2
	if linew < 1 {
		linew = 1
	}
	z0 := de - 2
	r := image.Rect(0, z0, boxw, he-(he-as)/2+f.Face.Letting()/2)
	r = r.Sub(image.Pt(r.Dx()/2, 0))
	f.tick = image.NewRGBA(r)
	f.tickback = image.NewRGBA(r)
	draw.Draw(f.tick, f.tick.Bounds(), f.Color.Back, image.ZP, draw.Src)
	tbg := f.tickbg()
	drawtick := func(x0, y0, x1, y1 int) {
		draw.Draw(f.tick, image.Rect(x0, y0, x1, y1), tbg, image.ZP, draw.Src)
	}
	drawtick(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+boxw)
	drawtick(r.Min.X, r.Max.Y-(boxw), r.Max.X, r.Max.Y)
	if boxw%2 != 0 {
		drawtick(-linew2, 0, linew2+1, r.Max.Y)
	} else {
		drawtick(-linew2, 0, linew2, r.Max.Y)
	}

}

// Put
func (f *Frame) tickat(pt image.Point, ticked bool) {
	if f.Ticked == ticked || f.tick == nil || !pt.In(f.Bounds().Inset(-1)) {
		return
	}
	pt.X -= 1
	//pt.Y -= f.Font.Letting() / 4
	r := f.tick.Bounds().Add(pt)
	if r.Max.X > f.r.Max.X {
		r.Max.X = f.r.Max.X
	}
	if ticked {
		f.Draw(f.tickback, f.tickback.Bounds(), f.b, pt.Add(f.tickback.Bounds().Min), draw.Src)
		f.Draw(f.b, r, f.tick, f.tick.Bounds().Min, draw.Src)
	} else {
		f.Draw(f.b, r, f.tickback, f.tickback.Bounds().Min, draw.Src)
	}
	//f.Flush(r.Inset(-1))
	f.Ticked = ticked
}
