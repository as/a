package win

import (
	"io"

	"github.com/as/text"
)

func (w *Win) WriteAt(p []byte, at int64) (n int, err error) {
	n, err = w.Editor.(io.WriterAt).WriteAt(p, at)
	q0, q1 := at, at+int64(len(p))

	switch text.Region5(q0, q1, w.org-1, w.org+w.Frame.Len()+1) {
	case -2:
		// Logically adjust origin to the left (up)
		w.org -= q1 - q0
	case -1:
		// Remove the visible text and adjust left
		w.Frame.Delete(0, q1-w.org)
		w.Frame.Insert(p, 0)
		w.org = q0
		w.Fill()
		w.dirty = true
	case 0:
		p0 := clamp(q0-w.org, 0, w.Frame.Len())
		p1 := clamp(q1-w.org, 0, w.Frame.Len())
		w.Frame.Delete(p0, p1)
		w.Frame.Insert(p, p0)
		w.Fill()
		w.dirty = true
	case 1:
		w.Frame.Delete(q0-w.org, w.Frame.Len())
		w.Frame.Insert(p, q0-w.org)
		w.Fill()
		w.dirty = true
	case 2:
	}
	return
}
