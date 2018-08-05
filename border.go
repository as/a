package main

import (
	"github.com/as/shiny/event/mouse"
	"github.com/as/ui/col"
	"github.com/as/ui/win"
)

func borderHit(e mouse.Event) bool {
	pt := p(e)
	return inSizer(pt) || inScroll(pt)
}

func procBorderHit(e mouse.Event) {
	apt := p(e)
	e = rel(e, act)
	pt := p(e)

	switch {
	case inSizer(pt):
		if HasButton(1, down) {
			if !apt.In(g.Area()) {
				for down != 0 {
					g.Move(apt)
					g.Refresh()
					col.Fill(g)
					apt = p(readmouse(<-D.Mouse))
				}
			} else if canopy(apt) {
				dragCol(g, actCol, e, D.Mouse)
			} else {
				dragTag(actCol, actTag, e, D.Mouse)
			}
			break
		}
		switch down {
		case Button(2):
		case Button(3):
			//			actCol.PrintList()
			actCol.RollUp(actCol.ID(actTag), act.Bounds().Min.Y)
			//			actCol.PrintList()
			moveMouse(act.Bounds().Min)
		}
		for down != 0 {
			readmouse(<-D.Mouse)
		}
	case inScroll(pt):
		switch down {
		case Button(1):
			scroll(act, ScrollEvent{Dy: 10, Event: e})
		case Button(2):
			w, _ := act.(*win.Win)
			if w == nil {
				break
			}
			w.Clicksb(pt, 0)
			repaint()
			for HasButton(2, down) {
				w.Clicksb(p(rel(readmouse(<-D.Mouse), w)), 0)
				repaint()
			}
		case Button(3):
			scroll(act, ScrollEvent{Dy: -10, Event: e})
		}
	default:
		logf("unknown border action at pt: %s", pt)
	}
}
