package win

import (
	"image"

	"github.com/as/ui/scroll"
)

const minSbWidth = 10

func (w *Win) Scroll(dl int) {
	if !w.Config.Scroll.Enable {
		return
	}
	org := w.org
	if dl < 0 {
		org = w.backNL(org, -dl)
		w.SetOrigin(org, true)
	} else {
		// TODO(as): Forward scrolling will be broken in non-graphical mode
		// Needs to be fixed here
		if !w.graphical() {
			return
		}
		if org+w.Frame.Nchars == w.Len() {
			return
		}
		r := w.Frame.Bounds()
		nline := w.Frame.Line()
		if nline == 0 {
			nline = 1
		}
		mul := int64(dl / nline)
		if mul == 0 {
			mul++
		}
		dx := w.IndexOf(image.Pt(r.Min.X, r.Min.Y+dl*w.Face().Dy())) * mul
		org += dx
		w.SetOrigin(org, true)
	}
	w.updatesb()
	w.drawsb()
	w.dirty = true
}

func (w *Win) backNL(p int64, n int) int64 {
	if n == 0 && p > 0 && w.Bytes()[p-1] != '\n' {
		n = 1
	}
	for i := n; i > 0 && p > 0; {
		i--
		p--
		if p == 0 {
			break
		}
		for j := 512; j-1 > 0 && p > 0; p-- {
			j--
			if p-1 < 0 || p-1 > w.Len() || w.Bytes()[p-1] == '\n' {
				break
			}
		}
	}
	return p
}
func (w *Win) Clicksb(pt image.Point, dir int) {
	w.clicksb(pt, dir)
	w.drawsb()
	w.dirty = true
}
func (w *Win) clicksb(pt image.Point, dir int) {
	fl := float64(w.Frame.Len())
	n := w.org
	switch dir {
	case 0:
		//rat := w.Bar.Delta(pt) // float64(pt.Y) / float64(w.ScrollBar.r.Dy())
		n = int64(float64(w.Len()) * w.sb.Delta(pt))
	case 1, -1:
		//rat = (barY1 / ptY)
		delta := int64(fl * w.sb.Delta(pt))
		n += delta * int64(dir)
	}
	w.SetOrigin(n, false)
	w.updatesb()
}

func (w *Win) scrollinit(pad image.Point) {
	if w.Frame == nil {
		return
	}
	sr := w.Frame.RGBA().Bounds()
	if pad.X > minSbWidth+2 {
		sr.Max.X = minSbWidth
	}
	w.sb = scroll.New(sr, w.Config.Scroll)
	w.updatesb()
	w.refreshsb()
}

func (w *Win) updatesb() {
	if !w.Config.Scroll.Enable || w.Len() == 0 {
		return
	}
	rat0 := float64(w.org) / float64(w.Len())
	rat1 := float64(w.org+w.Frame.Len()) / float64(w.Len())
	w.dirty = w.sb.Put(rat0, rat1) || w.dirty
}
func (w *Win) drawsb() {
	if !w.Graphical() {
		return
	}
	w.dirty = w.sb.Update(w.Frame.RGBA(), w.Frame) || w.dirty
}
func (w *Win) refreshsb() {
	if !w.Graphical() {
		return
	}
	w.sb.Refresh(w.Frame.RGBA(), w.Frame)
	w.dirty = true
}
