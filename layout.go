package main

import (
	//	"image"

	"image"
	"time"

	"github.com/as/ui/col"
	"github.com/as/ui/tag"

	"golang.org/x/mobile/event/mouse"
)

var (
	DragArea    = image.Rect(-50, -50, 50, 50)
	DragTimeout = time.Second * 1
	down        uint
)

func readmouse(e mouse.Event) mouse.Event {
	switch e.Direction {
	case 1:
		down |= 1 << uint(e.Button)
	case 2:
		down &^= 1 << uint(e.Button)
	}
	return e
}

func (g *Grid) dragCol(c *Col, e mouse.Event, mousein <-chan mouse.Event) {
	c0 := actCol
	col.Detach(g, g.ID(c0))
	col.Fill(c0)
	for e = range mousein {
		e = readmouse(e)
		if down == 0 {
			break
		}
	}
	activate(p(e), g)
	col.Attach(g, c0, p(e))
	//g.Attach(c0, p(e).X)
	moveMouse(c0.Loc().Min)
}

func (g *Grid) dragTag(c *Col, t *tag.Tag, e mouse.Event, mousein <-chan mouse.Event) {
	c.Detach(c.ID(t))
	t0 := time.Now()
	r0 := DragArea.Add(p(e).Add(t.Bounds().Min))
	for e = range mousein {
		e = readmouse(e)
		if down == 0 {
			break
		}
	}
	if time.Since(t0) < DragTimeout && p(e).In(r0) {
		actCol.Attach(t, p(e).Y-100)
	} else {
		activate(p(e), g)
		col.Fill(c)
		if t == nil {
			return
		}
		actCol.Attach(t, p(e).Y)
	}
	moveMouse(t.Loc().Min)
}
