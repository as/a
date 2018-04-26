package main

// Put
import (
	"image"
	"io"

	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/ui"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

var GridConfig = &tag.Config{
	Margin:     image.Pt(15, 0),
	Filesystem: newfsclient(),
	Facer:      font.NewFace,
	FaceHeight: *ftsize,
	Color: [3]frame.Color{
		0: frame.ATag0,
	},
	Image: true,
	Ctl:   events,
}

var TagConfig = &tag.Config{
	Margin:     image.Pt(15, 0),
	Filesystem: newfsclient(),
	Facer:      font.NewFace,
	FaceHeight: *ftsize,
	Color: [3]frame.Color{
		0: frame.ATag1,
	},
	Image: true,
	Ctl:   events,
}

type Col = col.Col

func NewTag(dev ui.Dev, basedir, name string) *tag.Tag {
	t := tag.New(dev, image.ZP, image.ZP, TagConfig)
	t.Open(basedir, name)
	t.Insert([]byte(" [Edit  ,x]	|"), t.Len())
	return t
}

func NewCol(dev ui.Dev, ft font.Face, sp, size image.Point, files ...string) *Col {
	col := col.New(dev, sp, size, TagConfig)
	for _, name := range files {
		New(col, ".", name)
	}
	return col
}

func NewColParams(g *Grid, filenames ...string) *Col {
	r := g.Loc()
	r.Min.X += g.Tag.Loc().Dx()
	if len(g.List) == 0 {
		r = g.List[len(g.List)-1].Loc()
		r.Min.X += r.Size().X / 2
	}
	return NewCol(g.Dev(), g.Face(), r.Min, r.Size(), filenames...)
}

func New(c *Col, basedir, name string, sizerFunc ...func(int) int) (w Plane) {
	t := NewTag(c.Dev(), basedir, name)
	if len(c.List) == 0 {
		c.Attach(t, c.Tag.Loc().Max.Y)
		return t
	}
	r := c.List[len(c.List)-1].Loc()
	c.Attach(t, r.Min.Y+r.Dy()/2)
	return t
}

func Delcol(g *Grid, id int) {
	co := col.Detach(g, id)
	x := co.Loc().Min.X
	y := co.Loc().Min.Y
	for ; id < len(g.List); id++ {
		x2 := g.List[id].Loc().Min.X
		g.List[id].Move(image.Pt(x, y))
		x = x2
	}
	col.Fill(g)
}

func Del(co *Col, id int) {
	w := co.Detach(id)
	y := w.Loc().Min.Y
	x := co.Loc().Min.X
	w.(io.Closer).Close()
	for ; id < len(co.List); id++ {
		y2 := co.List[id].Loc().Min.Y
		co.List[id].Move(image.Pt(x, y))
		y = y2
	}
	col.Fill(co)
}
