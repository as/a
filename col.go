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
	pt := p.Bounds().Min
	t, _ := p.(*tag.Tag)
	if t == nil {
		return pt
	}
	w, _ := t.Window.(*win.Win)
	if w == nil || !w.Graphical() {
		return pt
	}

	ar := w.Area()
	qt := ar.Min.Add(w.PointOf(w.Frame.Len()))

	if w.Frame.Full() || qt.Y >= ar.Max.Y {
		qt.Y = ar.Min.Y + ar.Dy()/2
		return qt
	}

	return qt
}

// New creates opens a names resource as a tagged window in column c
func New(c *Col, basedir, name string) (w Plane) {
	t := tag.New(c.Dev(), TagConfig)
	t.Open(basedir, name)
	t.Label.Write([]byte(" [Edit ] "))
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
	c.Tag.Label.InsertString("New Delcol Sort	|", 0)
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
		r = g.List[len(g.List)-1].Bounds()
	}
	c := NewCol(g.Dev(), g.Face(), r.Min, r.Size(), filenames...)
	col.Attach(g, c, phi(r))
	return c
}

func Delcol(g *Grid, id int) []Plane {
	co := col.Detach(g, id)
	x := co.Bounds().Min.X
	y := co.Bounds().Min.Y
	for ; id < len(g.List); id++ {
		x2 := g.List[id].Bounds().Min.X
		g.List[id].Move(image.Pt(x, y))
		x = x2
	}
	col.Fill(g)

	if co, _ := co.(interface{ Kids() []Plane }); co != nil {
		return co.Kids()
	}
	return nil
}

func Del(co *Col, id int) io.Closer {
	w := co.Detach(id)
	y := w.Bounds().Min.Y
	x := co.Bounds().Min.X
	for ; id < len(co.List); id++ {
		y2 := co.List[id].Bounds().Min.Y
		co.List[id].Move(image.Pt(x, y))
		y = y2
	}
	col.Fill(co)
	return w.(io.Closer)
}
