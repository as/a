package kbd

import (
	"github.com/as/font"
	"github.com/as/shiny/event/key"
	"github.com/as/text"
	"github.com/as/text/find"
	"github.com/as/ui/win"
)

// markDirt calls Mark if the editor implements
// the find.Dirt interface
func markDirt(ed text.Editor) {
	if ed, ok := ed.(text.Dirt); ok {
		ed.Mark()
	}
}

// Send process a keyboard event with the editor
func SendClient(hc text.Editor, e key.Event) {
	if e.Direction == key.DirRelease {
		return
	}
	e = preProcess(e)
	defer markDirt(hc)
	q0, q1 := hc.Dot()
	switch e.Code {
	case key.CodeEqualSign, key.CodeHyphenMinus:
		if e.Direction == key.DirRelease {
			return
		}

		if e.Modifiers == key.ModControl {
			return
			switch hc := hc.(type) {
			case *win.Win:
				df := 2
				if key.CodeHyphenMinus == e.Code {
					df = -2
				}
				if ft, ok := hc.Face().(*font.Resizer); ok {
					hc.SetFont(ft.New(ft.Dy() + df))
				}
				return
			}
		}
	case key.CodeUpArrow, key.CodePageUp, key.CodeDownArrow, key.CodePageDown:
		n := 1
		if e.Code == key.CodePageUp || e.Code == key.CodePageDown {
			n *= 10
		}
		if e.Code == key.CodeUpArrow || e.Code == key.CodePageUp {
			n = -n
		}
		if hc, ok := hc.(text.Scroller); ok {
			hc.Scroll(n)
		}
		//		hc.Mark()
		return
	case key.CodeLeftArrow, key.CodeRightArrow:
		if e.Code == key.CodeLeftArrow {
			if e.Modifiers&key.ModShift == 0 {
				q1--
			}
			q0--
		} else {
			if e.Modifiers&key.ModShift == 0 {
				q0++
			}
			q1++
		}
		hc.Select(q0, q1)
		//		hc.Mark()
		return
	}
	switch e.Rune {
	case -1:
		return
	case '\x01', '\x05', '\x08', '\x15', '\x17':
		if q0 == 0 && q1 == 0 {
			return
		}
		if q0 == q1 && q0 != 0 {
			q0--
		}
		switch e.Rune {
		case '\x15', '\x01': // ^U, ^A
			p := hc.Bytes()
			if q0 < int64(len(p))-1 {
				q0++
			}
			n0, n1 := find.Findlinerev(hc.Bytes(), q0, 0)
			if e.Rune == '\x15' {
				hc.Delete(n0, n1)
			}
			hc.Select(n0, n0)
		case '\x05': // ^E
			_, n1 := find.Findline3(hc.Bytes(), q1, 1)
			if n1 > 0 {
				n1--
			}
			hc.Select(n1, n1)
		case '\x17':
			if find.Isany(hc.Bytes()[q0], find.AlphaNum) {
				q0 = find.Acceptback(hc.Bytes(), q0, find.AlphaNum)
			}
			hc.Delete(q0, q1)
			hc.Select(q0, q0)
		case '\x08':
			fallthrough
		default:
			if q0 > q1 {
				q0, q1 = q1, q0
			}
			hc.Delete(q0, q1)
			hc.Select(q0, q0)
		}
		//		hc.Mark()
		return
	}
	ch := []byte(string(e.Rune))
	if q1 != q0 {
		hc.Delete(q0, q1)
		//		hc.Mark()
		q1 = q0
	}
	q1 += int64(hc.Insert(ch, q0))
	q0 = q1
	hc.Select(q0, q1)

}

func preProcess(e key.Event) key.Event {
	if e.Rune == '\t' {
		return e
	}
	r, ok := code2rune[e.Code]
	if ok {
		e.Rune = r
	}
	if e.Rune == '\r' {
		e.Rune = '\n'
	}
	return e
}

var code2rune = map[key.Code]rune{
	key.CodeReturnEnter:     '\n',
	key.CodeDeleteBackspace: '\x08',
	key.CodeTab:             '\x09',
	key.CodeSpacebar:        '\x20',
}
