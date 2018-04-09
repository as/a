package main

import (
	mus "github.com/as/text/mouse"
	"golang.org/x/mobile/event/mouse"
)

func borderHit(e mouse.Event) bool {
	pt := p(e)
	return inSizer(pt) || inScroll(pt)
}

func procBorderHit(e mouse.Event) {
	pt := p(e)
	switch {
	case inSizer(pt):
		if HasButton(1, down) {
			if canopy(absP(e, act.Loc().Min)) {
				g.dragCol(actCol, e, D.Mouse)
			} else {
				g.dragTag(actCol, actTag, e, D.Mouse)
			}
			break
		}
		switch down {
		case Button(2):
		case Button(3):
			actCol.RollUp(actCol.ID(actTag), act.Loc().Min.Y)
			moveMouse(act.Loc().Min)
		}
		for down != 0 {
			readmouse(<-D.Mouse)
		}
	case inScroll(pt):
		switch down {
		case Button(1):
			scroll(act, mus.ScrollEvent{Dy: 5, Event: e})
		case Button(2):
			w := act
			e = rel(readmouse(e), w)
			scroll(w, mus.ScrollEvent{Dy: 5, Event: e})
			repaint()
		case Button(3):
			scroll(act, mus.ScrollEvent{Dy: -5, Event: e})
		}
		logf("inScroll: %s", p(e))
	default:
		logf("unknown border action at pt: %s", pt)

	}
}
