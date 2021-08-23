package tag

import "image"

type Vis int

const (
	VisNone Vis = 0
	VisTag  Vis = 1
	VisBody Vis = 1 << 1
	VisFull Vis = VisTag | VisBody
)

func (t *Tag) Move(pt image.Point) {
	t.sp = pt
	t.Label.Move(pt)
	pt.Y += t.Label.Bounds().Dy()
	t.Window.Move(pt)
}

func (t *Tag) Resize(pt image.Point) {
	ts := t.Config.TagHeight()
	if ts > pt.Y {
		pt.Y = 0
		t.size = pt
		t.Label.Resize(pt)
		t.Window.Resize(pt)
		t.Vis = VisNone
		return
	}
	t.dirty = true
	if ts*2 > pt.Y {
		// Theres enough room for the label but the body wouldn't
		// have enough room.
		pt.Y = ts
		t.size = pt
		t.Label.Resize(pt)

		// Coherence: window always under tag
		t.align()

		pt.Y = 0
		t.Window.Resize(pt)
		t.Vis = VisTag
		return
	}
	t.size = pt
	t.Label.Resize(image.Pt(pt.X, ts))
	t.align()
	t.Window.Resize(image.Pt(pt.X, pt.Y-ts))
	t.Vis = VisFull
}

func (t *Tag) align() {
	// Coherence: window always under tag
	r := t.Label.Bounds()
	r.Min.Y = r.Max.Y
	t.Window.Move(r.Min)
}

func (t *Tag) Bounds() image.Rectangle {
	return image.Rectangle{t.sp, t.sp.Add(t.size)}
}
