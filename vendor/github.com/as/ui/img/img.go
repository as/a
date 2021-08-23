package img

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/png"

	"github.com/as/shiny/screen"
	"github.com/as/text"
	"github.com/as/ui"
)

type Node struct {
	sp     image.Point
	size   image.Point
	margin image.Point
	dirty  bool
}

//func (n Node) Size() image.Point {
//	return n.size
//}
//func (n Node) Pad() image.Point {
//	return n.Sp.Add(n.pad)
//}

func (w *Img) Area() image.Rectangle {
	return w.Bounds().Add(w.margin)
}
func (w *Img) area() image.Rectangle {
	return image.Rectangle{w.margin, w.size}
}

var DefaultConfig = Config{
	Margin: image.Pt(13, 3),
}

type Config struct {
	Name   string
	Margin image.Point
	Editor text.Editor

	// Ctl is a channel provided by the window owner. It carries window messages
	// back to the creator. Valid types are event.Look and event.Cmd
	Ctl chan interface{}
}

type Img struct {
	Node
	ui.Dev
	margin image.Point
	img    image.Image
	b      screen.Buffer
	ScrollBar
	org int64
	text.Editor
}

type ScrollBar struct {
	Scrollr image.Rectangle
}

func New(dev ui.Dev, conf *Config) *Img {
	if conf == nil {
		c := DefaultConfig
		conf = &c
	}
	ed, _ := text.Open(text.NewBuffer())

	var img image.Image
	if ed.Len() != 0 {
		img, _, _ = image.Decode(bytes.NewReader(ed.Bytes()))
	}

	w := &Img{
		Dev:    dev,
		margin: conf.Margin,
		Editor: ed,
		img:    img,
	}
	w.init()
	return w
}

var (
	MinRect = image.Rect(0, 0, 10, 10)
)

func (w *Img) reallocimage(size image.Point) bool {
	if w.b != nil {
		w.b.Release()
		w.b = nil
	}
	if small(size) {
		return false
	}
	b, err := w.NewBuffer(size)
	if err != nil {
		panic(size.String())
	}
	w.b = b
	return true
}
func small(size image.Point) bool {
	return size.X == 0 || size.Y == 0 || size.In(MinRect)
}
func (w *Img) init() {
	w.Blank()
	w.Fill()
	q0, q1 := w.Dot()
	w.Select(q0, q1)
	w.Mark()
}

func (w *Img) Blank() {
	if w.b == nil {
		return
	}
	r := w.b.RGBA().Bounds()
	draw.Draw(w.b.RGBA(), r, image.Black, image.ZP, draw.Src)
	if w.sp.Y > 0 {
		r.Min.Y--
	}
	w.Mark()
	//	w.drawsb()
}

func (w *Img) Graphical() bool {
	return w != nil && w.img != nil && w.Dev != nil && w.b != nil && !w.size.In(MinRect) && w.size != image.ZP
}

func (w *Img) Mark()                           { w.dirty = true }
func (w *Img) Bounds() image.Rectangle         { return image.Rectangle{w.sp, w.sp.Add(w.size)} }
func (w *Img) Buffer() screen.Buffer           { return w.b }
func (w *Img) Dirty() bool                     { return w.dirty }
func (w *Img) Bytes() []byte                   { return w.Editor.Bytes() }
func (w *Img) Len() int64                      { return w.Editor.Len() }
func (w *Img) Move(sp image.Point)             { w.sp = sp }
func (w *Img) Origin() int64                   { return w.org }
func (w *Img) Fill()                           {}
func (w *Img) Clicksb(pt image.Point, dir int) {}
func (w *Img) Scroll(dl int)                   {}
func (w *Img) SetOrigin(org int64, exact bool) {}
func (w *Img) Refresh() {
	w.dirty = true
	w.Upload()
}
func (w *Img) Upload() {
	if !w.dirty || !w.Graphical() {
		return
	}
	b, r := w.b, w.img.Bounds()
	draw.Draw(b.RGBA(), b.RGBA().Bounds().Add(w.margin), w.img, r.Min, draw.Src)
	w.Window().Upload(w.sp, w.b, w.minbounds())
	w.dirty = false
}
func (w *Img) Resize(size image.Point) {
	w.size = size
	if !w.reallocimage(w.size) {
		if w == nil {
			return
		}
		w.img = nil
		return
	}
	w.dirty = true
	if w.Editor.Len() != 0 {
		w.img, _, _ = image.Decode(bytes.NewReader(w.Editor.Bytes()))
	}
	w.init()
	w.Refresh()
}

func (w Img) minbounds() image.Rectangle {
	return image.Rectangle{image.ZP, w.Bounds().Size()}.Union(w.b.Bounds())
}
