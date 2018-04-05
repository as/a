package main

import (
	"image"

	"github.com/as/text"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
	"golang.org/x/mobile/event/mouse"
)

func Button(n uint) uint {
	return 1 << n
}
func HasButton(n, mask uint) bool {
	return Button(n)&mask != 0
}

// mouseMove(pt image.Point) // defined in mouse_other.go and mouse_linux.go

func sweepFunc(w *win.Win, e mouse.Event, mc <-chan mouse.Event) (q0, q1 int64, e1 mouse.Event) {
	start := down
	q0, q1 = w.Dot()
	act.Sq = q0
	for down == start {
		w.Sq, q0, q1 = sweep(w, e, w.Sq, q0, q1)
		w.Select(q0, q1)
		repaint()
		e = rel(readmouse(<-mc), w)
	}
	return q0, q1, e
}

func cursorNop(p image.Point) {}

func shouldCursor(p Plane) (fn func(image.Point)) {
	switch p.(type) {
	case Named:
		return cursorNop
	default:
		return moveMouse
	}
}
func ajump2(ed text.Editor, cursor bool) {
	fn := moveMouse
	if !cursor {
		fn = nil
	}
	if ed, ok := ed.(text.Jumper); ok {
		ed.Jump(fn)
	}
}

func ajump(p interface{}, cursor func(image.Point)) {
	switch p := p.(type) {
	case nil:
		return //TODO(as): error message without a recursive call
	case *tag.Tag:
		if p != nil {
			cursor(p.Loc().Min)
		}
	case text.Jumper:
		p.Jump(cursor)
	case Plane:
		if cursor == nil {
			cursor = shouldCursor(p)
		}
		cursor(p.Loc().Min)
	}
}
