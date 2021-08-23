package win

import (
	"image/png"
	"os"
)

func (w *Win) SetOrigin(org int64, exact bool) {
	org = clamp(org, 0, w.Len())
	if org == w.org {
		return
	}
	if org > 0 && !exact {
		for i := 0; i < 2048 && org < w.Len(); i++ {
			if w.Bytes()[org] == '\n' {
				org++
				break
			}
			org++
		}
	}
	if w.graphical() {
		w.setOrigin(clamp(org, 0, w.Len()))
		w.Mark()
		w.UserFunc(w)
	} else {
		w.org = org
		w.dirty = true
	}
}

func (w *Win) setOrigin(org int64) {
	if org == w.org {
		return
	}

	f := w.Frame
	q0, q1 := w.Dot()
	delta := org - w.org
	fix := false
	switch {
	case abs(delta) >= f.Len():
		f.Delete(0, f.Len())
	case delta > 0:
		func() {
			end := w.org + f.Len()
			if q0 < end && q1 >= end {
				w.Swap()
				defer w.Swap()
			}
			f.Delete(0, delta)
		}()
		fix = true
	default:
		f.Insert(w.Bytes()[org:org-delta], 0)
	}
	w.org = org
	w.Fill()
	w.drawsb()
	w.Select(q0, q1)

	if fix {
		w.fixEnd()
	}
	w.dirty = true
}

func (w *Win) pngwrite(name string) {
	fd, _ := os.Create(name)
	png.Encode(fd, w.Frame.RGBA())
	fd.Close()
}
func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}
