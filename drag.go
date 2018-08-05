package main

import (
	"image"
	"time"

	"github.com/as/ui/col"
	"github.com/as/ui/tag"

	"github.com/as/shiny/event/mouse"
)

var (
	DragArea    = image.Rect(-50, -50, 50, 50)
	DragTimeout = time.Second * 1
	down        uint
)

func readmouse(e mouse.Event) mouse.Event {
	// km := KMod()
	if e.Button == 1 {
		if km := e.Modifiers; km&KCtrl != 0 {
			e.Button = 3
		} else if km&KAlt != 0 {
			e.Button = 2
		}
	}
	if dir := e.Direction; dir == 1 {
		down |= 1 << uint(e.Button)
	} else if dir == 2 {
		down &^= 1 << uint(e.Button)
	}
	return e
}

func dragCol(g *Grid, c *Col, e mouse.Event, mousein <-chan mouse.Event) {
	c0 := actCol
	for e = range mousein {
		e = readmouse(e)
		if down == 0 {
			break
		}
		// uncomment for really stupid stuff
		//col.Detach(g, g.ID(c0))
		//col.Fill(g)
		//col.Attach(g, c0, p(e))
		//g.Upload()
	}
	col.Detach(g, g.ID(c0))
	col.Fill(g)
	col.Attach(g, c0, p(e))
	g.Upload()
	moveMouse(c0.Bounds().Min)
}

func dragTag(c *Col, t *tag.Tag, e mouse.Event, mousein <-chan mouse.Event) {
	c.Detach(c.ID(t))
	t0 := time.Now()
	r0 := DragArea.Add(p(e).Add(t.Bounds().Min))
	for e = range mousein {
		e = readmouse(e)
		if down == 0 {
			break
		}
	}
	pt := p(e)
	if time.Since(t0) < DragTimeout && p(e).In(r0) {
		pt.Y -= 100
		col.Attach(actCol, t, pt)
	} else {
		activate(p(e), g)
		col.Fill(c)
		if t == nil {
			return
		}
		col.Attach(actCol, t, pt)
	}
	moveMouse(t.Bounds().Min)
}
