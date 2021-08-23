package win

func (w *Win) fixEnd() {
	p0, p1 := w.Frame.Dot()
	p0 = clamp(p0, 0, w.Frame.Len())
	p1 = clamp(p1, 0, w.Frame.Len())
	if p1 > p0 {
		w.Redraw(w.PointOf(p1-1), p1-1, p1, true)
	}
	w.Redraw(w.PointOf(p1), p1, p1+w.Frame.Len(), false)
}

func (w *Win) Fill() {
	if w.Frame.Full() {
		return
	}
	for !w.Frame.Full() {
		qep := w.org + w.Nchars
		n := max(0, min(w.Len()-qep, 2000))
		if n == 0 {
			break
		}
		rp := w.Bytes()[qep : qep+n]

		nl := w.MaxLine() - w.Line()
		m := 0
		i := int64(0)
		for i < n {
			if rp[i] == '\n' {
				m++
				if m >= nl {
					i++
					break
				}
			}
			i++
		}
		w.Frame.Insert(rp[:i], w.Nchars)
		w.dirty = true
	}
}
