package main

import (
	"github.com/as/ui/tag"
	"golang.org/x/mobile/event/mouse"
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
