package win

import "github.com/as/text"
import "image"

const (
	// Extra lines to scroll down to comfortably display the result of a look operation
	JumpScrollMargin = -3
)

// Select selects the range [q0:q1] inclusive and returns
// the old selection (equivalent to dot before select was called)
func (w *Win) Select(q0, q1 int64) {
	if q0 > q1 {
		q0, q1 = q1, q0
	}
	w.Editor.Select(q0, q1)
	if !w.graphical() {
		return
	}
	reg := text.Region3(q0, w.org-1, w.org+w.Frame.Len())
	w.dirty = true
	p0, p1 := q0-w.org, q1-w.org
	w.Frame.Select(p0, p1)
	if q0 == q1 && reg != 0 {
		w.Untick()
	}
}

// Jump scrolls the active selection into view. An optional mouseFunc
// is given the transfer coordinates to move the mouse cursor under
// the selection.
func (w *Win) Jump(mouseFunc func(image.Point)) {
	q0, q1 := w.Dot()
	if text.Region5(q0, q1, w.Origin(), w.Origin()+w.Frame.Len()) != 0 {
		w.SetOrigin(q0, true)
		w.Scroll(JumpScrollMargin)
	}
	if mouseFunc != nil {
		jmp := w.PointOf(q0 - w.org)
		mouseFunc(w.Bounds().Min.Add(jmp))
	}
}
