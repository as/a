package main

import (
	"image"

	"github.com/as/ui/tag"
	"github.com/as/ui/win"
	"golang.org/x/mobile/event/mouse"
)

var (
	winid  = make(map[int]*win.Win)
	wincol = make(map[int]*Col)
	xx     Cursor
	down   uint
	pt     image.Point
)

func readmouse(e mouse.Event) mouse.Event {
	switch e.Direction {
	case 1:
		down |= 1 << uint(e.Button)
	case 2:
		down &^= 1 << uint(e.Button)
	}
	activate(p(e), g)
	pt = p(e).Add(act.Loc().Min)
	e.X -= float32(act.Sp.X)
	e.Y -= float32(act.Sp.Y)
	return e
}

func (g *Grid) dragCol(c *Col, e mouse.Event, mousein <-chan mouse.Event) {
	xx.srcCol = actCol
	xx.sweepCol = true
	xx.src = nil
	g.detach(g.ID(xx.srcCol))
	g.fill()

	for e = range mousein {
		e = readmouse(e)
		if down == 0 {
			break
		}
	}
	e.X += float32(act.Sp.X)
	e.Y += float32(act.Sp.Y)
	pt.X = int(e.X)
	pt.Y = int(e.Y)
	activate(p(e), g)

	xx.sweepCol = false
	g.fill()
	g.Attach(xx.srcCol, pt.X)
	moveMouse(xx.srcCol.Loc().Min)
}

func (g *Grid) dragTag(c *Col, t *tag.Tag, e mouse.Event, mousein <-chan mouse.Event) {
	detachtag(g, c, t)
	for e = range mousein {
		e = readmouse(e)
		if down == 0 {
			break
		}
	}
	e.X += float32(act.Sp.X)
	e.Y += float32(act.Sp.Y)
	pt.X = int(e.X)
	pt.Y = int(e.Y)
	activate(p(e), g)
	xx.sweep = false
	xx.srcCol.fill()
	if xx.src == nil {
		return
	}
	c = actCol
	c.Attach(xx.src, pt.Y)
	moveMouse(xx.src.Loc().Min)
}

func detachtag(g *Grid, c *Col, t *tag.Tag) {
	xx.srcCol = c
	xx.src = t
	c.detach(c.ID(t))
}
func detachcol(g *Grid, c *Col) {
	xx.sweepCol = true
	xx.srcCol = c
	xx.sweepCol = true
	xx.src = nil
	g.detach(g.ID(c))
	g.fill()
}

func detacxh(g *Grid, it interface{}) {
	switch item := it.(type) {
	case nil:
		g.aerr("cant detach nil")
	case *Col:
		xx.srcCol = item
		xx.src = nil
		xx.sweepCol = true
		g.detach(g.ID(item))
		g.fill()
	case *tag.Tag:
		g.aerr("%#v\n", item)
		panic("constraint violation")
	case *win.Win:
		xx.srcCol = actCol
		xx.src = actTag
		xx.srcCol.detach(xx.srcCol.ID(xx.src))
	case interface{}:
		g.aerr("%#v\n", item)
		panic("constraint violation")
	}
}

type Cursor struct {
	sweep    bool
	sweepCol bool
	srcCol   *Col
	src      Plane
}
