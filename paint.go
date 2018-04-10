package main

import (
	"golang.org/x/mobile/event/paint"
)

func repaint() {
	select {
	case D.Paint <- paint.Event{}:
	default:
	}
}
