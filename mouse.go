package main

import (
	"image"
	"time"

	"github.com/as/event"
	"github.com/as/text"
	"github.com/as/text/find"
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

var (
	last   uint
	lastpt image.Point
	t0     = time.Now()
)

func procButton(e mouse.Event) {
	double := false
	if last == down {
		if time.Since(t0) < time.Second/2 && lastpt.In(image.Rect(-3, -3, 3, 3).Add(p(e))) {
			double = true
		}
	}
	t0 = time.Now()
	last = down
	lastpt = p(e)
	t := actTag
	w, _ := act.(*win.Win)
	if w == nil {
		return
	}

	s0, s1 := w.Dot()
	q0 := w.IndexOf(p(e)) + w.Origin()
	q1 := q0
	act.Select(q0, q1)
	repaint()
	switch down {
	case Button(1):
		if double {
			q0, q1 = find.FreeExpand(w, q0)
			double = false
			w.Select(q0, q1)
		} else {
			// In Acme and Sam, the double click action doesn't maintain
			// a hold on the selection if the mouse is moved out of a rectangular
			// region. I don't do the same thing here because it's sometimes
			// advantageous to make the selection hold for a scrolling select
			// operation.
			q0, q1, e = sweepFunc(w, e, D.Mouse)
		}
		for down != 0 {
			w.Select(q0, q1)
			if HasButton(2, down) {
				tag.Snarf(w, e)
				q1 = q0
			} else if HasButton(3, down) {
				q0, q1 = tag.Paste(w, e)
			}
			repaint()
			e = rel(readmouse(<-D.Mouse), t)
		}
		t0 = time.Now()
		w.Select(q0, q1)
	case Button(2):
		q0, q1, _ := sweepFunc(w, e, D.Mouse)
		if q0 == q1 {
			if text.Region3(q0, s0, s1) == 0 {
				q0, q1 = s0, s1
			} else {
				q0, q1 = find.ExpandFile(w.Bytes(), q0)
			}
		}
		w.Select(s0, s1)
		acmd(event.Cmd{
			Name: t.FileName(),
			From: t, To: []event.Editor{w},
			Rec: event.Rec{Q0: q0, Q1: q0, P: w.Bytes()[q0:q1]},
		})
	case Button(3):
		q0, q1, _ := sweepFunc(w, e, D.Mouse)
		if q0 == q1 {
			q0, q1 = find.ExpandFile(w.Bytes(), q0)
		}
		w.Select(s0, s1)
		g.Look(event.Look{
			Name: t.FileName(),
			From: t, To: []event.Editor{w},
			Rec: event.Rec{Q0: q0, Q1: q1, P: w.Bytes()[q0:q1]},
		})
	}
}

// moveMouse(pt image.Point) // defined in mouse_other.go and mouse_linux.go
func MoveMouse(address interface{}) {
	switch a := address.(type) {
	case *win.Win:
		p0, _ := a.Frame.Dot()
		moveMouse(a.Loc().Min.Add(a.PointOf(p0)))
	case image.Point:
		moveMouse(a)
		return
	case int64:
		w, _ := act.(*win.Win)
		if w == nil {
			return
		}
		p0, _ := w.Frame.Dot()
		moveMouse(w.PointOf(p0))
		return
	}
	logf("MoveMove: error %T act=%#v", address, act)
}

func sweepFunc(w *win.Win, e mouse.Event, mc <-chan mouse.Event) (q0, q1 int64, e1 mouse.Event) {
	start := down
	q0, q1 = w.Dot()
	w.Sq = q0
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
