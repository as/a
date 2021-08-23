package win

import (
	"github.com/as/text"
)

// Insert inserts the bytes in p at position q0. When q0
// is zero, Insert prepends the bytes in p to the underlying
// buffer.
func (w *Win) Insert(p []byte, q0 int64) int {
	if EnableUndoExperiment {
		w.ops.Insert(p, q0)
	}
	return w.insert(p, q0)
}

func (w *Win) insert(p []byte, q0 int64) int {
	if w.Editor == nil {
		panic("nil editor")
	}
	if len(p) == 0 {
		return 0
	}
	n := w.Editor.Insert(p, q0)
	if !w.graphical() {
		return n
	}

	// If at least one point in the region overlaps the
	// frame's visible area then we alter the frame. Otherwise
	// there's no point in moving text down, it's just annoying.

	switch q1 := q0 + int64(len(p)); text.Region5(q0, q1, w.org-1, w.org+w.Frame.Len()+1) {
	case -2:
		w.org += q1 - q0
	case -1:
		// Insertion to the left
		w.Frame.Insert(p[q1-w.org:], 0)
		w.org += w.org - q0
		w.dirty = true
	case 1:
		w.Frame.Insert(p, q0-w.org)
		w.dirty = true
	case 0:
		if q0 < w.org {
			p0 := w.org - q0
			w.Frame.Insert(p[p0:], 0)
			w.org += w.org - q0
		} else {
			w.Frame.Insert(p, q0-w.org)
		}
		w.dirty = true
	}
	return n
}

func (w *Win) Write(p []byte) (n int, err error) {
	return w.Insert(p, w.Len()), nil
}
