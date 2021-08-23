package frame

import (
	"errors"
	"image"
	"image/draw"

	"github.com/as/font"
	"github.com/as/frame/box"
)

const (
	FrElastic = 1 << iota
	FrUTF8
)

var (
	ForceElastic bool
	ForceUTF8    bool
)

var (
	ErrBadDst = errors.New("bad dst")
)

// Frame is a write-only container for editable text
type Frame struct {
	box.Run
	p0 int64
	p1 int64
	b  draw.Image
	r  image.Rectangle
	ir *box.Run

	Face font.Face
	Color
	Ticked bool
	Scroll func(int)
	Drawer
	op draw.Op

	mintab int
	maxtab int
	full   int

	tick     draw.Image
	tickback draw.Image
	tickoff  bool
	maxlines int
	modified bool

	pts [][2]image.Point

	flags int
}

func New(dst draw.Image, r image.Rectangle, conf *Config) *Frame {
	if dst == nil {
		return nil
	}
	if conf == nil {
		conf = &Config{}
	}
	conf.check()
	fl := conf.Flag
	face := negotiateFace(conf.Face, fl)
	mintab, maxtab := tabMinMax(face, fl&FrElastic != 0)
	f := &Frame{
		Face:   face,
		Color:  conf.Color,
		Drawer: conf.Drawer,
		Run:    box.NewRun(mintab, 5000, face),
		op:     draw.Src,
		mintab: mintab,
		maxtab: maxtab,
		flags:  fl,
	}
	f.setrects(r, dst)
	f.inittick()
	run := box.NewRun(mintab, 5000, face)
	f.ir = &run
	return f
}

func (f *Frame) Config() *Config {
	return &Config{
		Flag:   f.flags,
		Color:  f.Color,
		Face:   f.Face,
		Drawer: f.Drawer,
	}
}

var zc Color

// Flags returns the flags currently set for the frame
func (f *Frame) Flags() int {
	return f.flags
}

// Flag sets the flags for the frame. At this time
// only FrElastic is supported.
func (f *Frame) SetFlags(flags int) {
	fl := getflag(flags)
	f.flags = fl
	f.mintab, f.maxtab = tabMinMax(f.Face, f.elastic())
	//	f.Reset( f.r, f.RGBA(),f.Font)
	//	f.mintab, f.maxtab = tabMinMax(f.Font, f.elastic())
}

func (f *Frame) elastic() bool {
	return f.flags&FrElastic != 0
}

func tabMinMax(ft font.Face, elastic bool) (min, max int) {
	mintab := ft.Dx([]byte{' '})
	maxtab := mintab * 4
	if elastic {
		mintab = maxtab
	}
	return mintab, maxtab
}

func (f *Frame) RGBA() *image.RGBA {
	rgba, _ := f.b.(*image.RGBA)
	return rgba
}
func (f *Frame) Size() image.Point {
	return f.r.Size()
}

// Dirty returns true if the contents of the frame have changes since the last redraw
func (f *Frame) Dirty() bool {
	return f.modified
}

// SetDirty alters the frame's internal state
func (f *Frame) SetDirty(dirty bool) {
	f.modified = dirty
}

func (f *Frame) SetFont(ft font.Face) {
	f.Face = font.Open(ft)
	f.Run.Reset(f.Face)
	f.Refresh()
}

func (f *Frame) SetOp(op draw.Op) {
	f.op = op
}

// Close closes the frame
func (f *Frame) Close() error {
	return nil
}

// Reset resets the frame to display on image b with bounds r and font ft.
func (f *Frame) Reset(r image.Rectangle, b *image.RGBA, ft font.Face) {
	f.r = r
	f.b = b
	f.SetFont(ft)
}

// Bounds returns the frame's clipping rectangle
func (f *Frame) Bounds() image.Rectangle {
	return f.r.Bounds()
}

// Full returns true if the last line in the frame is full.
func (f *Frame) Full() bool {
	if f == nil {
		return true
	}
	return f.full == 1
}

// Maxline returns the max number of wrapped lines fitting on the frame
func (f *Frame) MaxLine() int {
	if f == nil {
		return 0
	}
	return f.maxlines
}

// Line returns the number of wrapped lines currently in the frame
func (f *Frame) Line() int {
	if f == nil {
		return 0
	}
	return f.Nlines
}

// Len returns the number of bytes currently in the frame
func (f *Frame) Len() int64 {
	if f == nil {
		return 0
	}
	return f.Nchars
}

// Dot returns the range of the selected text
func (f *Frame) Dot() (p0, p1 int64) {
	return f.p0, f.p1
}

func (f *Frame) setrects(r image.Rectangle, b draw.Image) {
	f.b = b
	f.r = r
	h := f.Face.Dy()
	f.r.Max.Y -= f.r.Dy() % h
	f.maxlines = f.r.Dy() / h
}
