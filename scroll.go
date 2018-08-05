package main

import (
	mus "github.com/as/text/mouse"
	"github.com/as/ui/tag"
)

func scroll(act tag.Window, e mus.ScrollEvent) {
	if e.Button == -1 {
		e.Dy = -e.Dy
	}
	actTag.Scroll(e.Dy)
	repaint()
}
