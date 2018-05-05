package main

import (
	"image"
	"io"

	"github.com/as/font"
	"github.com/as/ui"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

type Col = col.Col

func phi(r image.Rectangle) image.Point {
	size := r.Size()
	size = size.Sub(size.Div(3))
	return r.Min.Add(size)
}

func underText(p Plane) image.Point {
	pt := phi(p.Loc())
	t, _ := p.(*tag.Tag)
	if t == nil {
		return pt
	}
	w, _ := t.Body.(*win.Win)
	if w == nil || !w.Graphical() {
		return pt
	}
	if w.Frame.Full() {
		return pt
	}
	return w.Area().Min.Add(w.PointOf(w.Frame.Len()))
}

// New creates opens a names resource as a tagged window in column c
func New(c *Col, basedir, name string) (w Plane) {
	t := tag.New(c.Dev(), TagConfig)
	t.Open(basedir, name)
	t.Insert([]byte(" [Edit  ,x]"), t.Len())
	r := c.Area()
	if c.Len() > 0 {
		r.Min = underText(c.List[len(c.List)-1])
	} else {
		r.Min = phi(r)
	}
	col.Attach(c, t, r.Min)
	return t
}

func NewCol(dev ui.Dev, ft font.Face, sp, size image.Point, files ...string) *Col {
	c := col.New(dev, ColConfig)
	c.Tag.InsertString("New Delcol Sort	|", 0)
	c.Move(sp)
	c.Resize(size)
	for _, name := range files {
		New(c, "", name)
	}
	return c
}

func NewColParams(g *Grid, filenames ...string) *Col {
	r := g.Area()
	if len(g.List) != 0 {
		r = g.List[len(g.List)-1].Loc()
	}
	c := NewCol(g.Dev(), g.Face(), r.Min, r.Size(), filenames...)
	col.Attach(g, c, phi(r))
	return c
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
