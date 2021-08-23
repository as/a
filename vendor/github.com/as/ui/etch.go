package ui

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/as/shiny/screen"
	"github.com/as/shiny/math/f64"
)

var screenRect = image.Rect(0, 0, 2048, 2048)

func NewEtch() *Etch {
	return &Etch{}
}

type Etch struct {
	dots *image.RGBA
	d    *screen.Device
}

func (e *Etch) WritePNG(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return png.Encode(fd, e.dots)
}
func (e *Etch) Screenshot(r image.Rectangle) *image.RGBA {
	if e.dots == nil {
		return image.NewRGBA(image.ZR)
	}
	dst := image.NewRGBA(r)
	draw.Draw(dst, r, e.dots, r.Min, draw.Src)
	return dst
}

func (e *Etch) Blank() {
	if e.dots != nil {
		draw.Draw(e.dots, e.dots.Bounds(), image.Black, image.ZP, draw.Src)
	}
}
func (e *Etch) Window() screen.Window {
	return e // cheap
}
func (e *Etch) Device() *screen.Device { return e.d }
func (e *Etch) Release()               {}
func (e *Etch) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	r := image.Rectangle{dp, dp.Add(src.RGBA().Bounds().Size())}
	if e.dots == nil {
		e.dots = image.NewRGBA(screenRect)
	}
	draw.Draw(e.dots, r, src.RGBA(), sr.Min, draw.Src)
}
func (e *Etch) Fill(dr image.Rectangle, src color.Color, op draw.Op) {}
func (e *Etch) Publish() (s screen.PublishResult)                    { return s }
func (e *Etch) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}
func (e *Etch) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}
func (e *Etch) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}
func (e *Etch) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}

func (e *Etch) Screen() screen.Screen {
	return e
}

func (e *Etch) NewBuffer(size image.Point) (screen.Buffer, error) {
	r := image.Rectangle{image.ZP, size}
	return &EtchBuffer{image.NewRGBA(r)}, nil
}
func (e *Etch) NewTexture(size image.Point) (screen.Texture, error) {
	panic("NewTexture is not implemented")
}
func (e *Etch) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("NewWindow is not implemented")
}

type EtchBuffer struct {
	img *image.RGBA
}

func (eb *EtchBuffer) Release()          {}
func (eb *EtchBuffer) Size() image.Point { return eb.img.Bounds().Size() }
func (eb *EtchBuffer) Bounds() image.Rectangle {
	return image.Rectangle{image.ZP, eb.Size()}
}
func (eb *EtchBuffer) RGBA() *image.RGBA {
	return eb.img
}
