package tag

import (
	"fmt"
	"image"

	"github.com/as/frame"
	"github.com/as/text"
)

type painter interface {
	Recolor(image.Point, int64, int64, frame.Palette)
	text.Projector
}

func OpenGraphic(b text.Buffer, col frame.Color) (w text.Editor, err error) {
	if b == nil {
		return nil, fmt.Errorf("bad buffer")
	}
	return &client{b, 0, 0, 0, &col}, nil
}

func framepaint(t painter, col *frame.Color, s, r0, r1, q0, q1 int64) {
	switch direction(s, r0, r1, q0, q1) {
	case -1:
		if r0 < q0 {
			t.Recolor(t.PointOf(r0), r0, q0, col.Hi)
		} else {
			t.Recolor(t.PointOf(q0), q0, r0, col.Palette)
		}
	case 1:
		if q1 < r1 {
			t.Recolor(t.PointOf(q1), q1, r1, col.Hi)
		} else {
			t.Recolor(t.PointOf(r1), r1, q1, col.Palette)
		}
	case 0:
		t.Recolor(t.PointOf(q0), q0, q1, col.Palette)
		t.Recolor(t.PointOf(r0), r0, r1, col.Hi)
	}
}

func direction(s, r0, r1, q0, q1 int64) (dir int) {
	switch {
	case r0 == s && q0 == s:
		return 1
	case r1 == s && q1 == s:
		return -1
	case r0 == s && q1 == s:
		return 0
	case r1 == s && q0 == s:
		return 0
	default:
		return 0
	}
}
