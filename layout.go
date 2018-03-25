package main

import (
	//	"image"

	"github.com/as/ui/tag"
	"github.com/as/ui/win"
	"golang.org/x/mobile/event/mouse"
)

var (
	winid  = make(map[int]*win.Win)
	wincol = make(map[int]*Col)
	down   uint
)

func readmouse0(e mouse.Event) mouse.Event {
	switch e.Direction {
	case 1:
		down |= 1 << uint(e.Button)
	case 2:
		down &^= 1 << uint(e.Button)
	}
	activate(p(e), g)
	return e
}

func (g *Grid) dragCol(c *Col, e mouse.Event, mousein <-chan mouse.Event) {
	c0 := actCol
	g.detach(g.ID(c0))
	g.fill()
	for e = range mousein {
		e = readmouse0(e)
		if down == 0 {
			break
		}
	}
	activate(p(e), g)
	g.fill()
	g.Attach(c0, p(e).X)
	moveMouse(c0.Loc().Min)
}

func (g *Grid) dragTag(c *Col, t *tag.Tag, e mouse.Event, mousein <-chan mouse.Event) {
	c.detach(c.ID(t))
	for e = range mousein {
		e = readmouse0(e)
		if down == 0 {
			break
		}
	}
	activate(p(e), g)
	c.fill()
	if t == nil {
		return
	}
	actCol.Attach(t, p(e).Y)
	moveMouse(t.Loc().Min)
}
