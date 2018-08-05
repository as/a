package main

import (
	"github.com/as/shiny/event/paint"
)

func repaint() {
	select {
	case D.Paint <- paint.Event{}:
	default:
	}
}
