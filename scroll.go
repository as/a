package main

import (
	"github.com/as/shiny/event/mouse"
	"github.com/as/ui/tag"
)

type ScrollEvent struct {
	Dy int
	mouse.Event
}

func scroll(act tag.Window, e ScrollEvent) {
	if e.Button == -1 {
		e.Dy = -e.Dy
	}
	actTag.Scroll(e.Dy)
	repaint()
}
