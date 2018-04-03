package main

import "golang.org/x/mobile/event/paint"

var (
	dirty bool
)

func repaint() {
	if !dirty || act == nil {
		return
	}
	select {
	case act.Window().Device().Paint <- paint.Event{}:
	default:
	}
}
