package main

import (
	"image"

	"github.com/as/font"
	"github.com/as/shiny/screen"
	"github.com/as/text/find"
	"github.com/as/text/kbd"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
	"golang.org/x/mobile/event/key"
)

type Window interface {
	Bytes() []byte

	Select(q0, q1 int64)
	SetOrigin(org int64, exact bool)
	Dot() (int64, int64)
	Len() int64
	Origin() int64

	Scroll(dl int)
	Insert(p []byte, q0 int64) (n int)
	//      WriteAt(p []byte, at int64) (n int, err error)
	Delete(q0, q1 int64) (n int)

	Fill()
	Blank()
	Refresh()
	Dirty() bool
	Loc() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
	Close() error
	Upload()
	Window() screen.Window
}

func kbdin(e key.Event, t *tag.Tag, act Window) {
	if e.Direction == 2 {
		return
	}
	if e.Code == key.CodeI && e.Modifiers == key.ModControl {
		runGoImports(t, e)
		return
	}
	switch e.Code {
	case key.CodeEqualSign, key.CodeHyphenMinus:
		if e.Modifiers == key.ModControl {
			win, _ := t.Body.(*win.Win)
			if win == nil {
				return
			}
			size := win.Frame.Face.Height()
			if key.CodeHyphenMinus == e.Code {
				size -= 1
			} else {
				size += 1
			}
			if size < 3 {
				size = 6
			}
			t.SetFont(font.NewFace(size))
			return
		}
	}
	ntab := int64(-1)
	if (e.Rune == '\n' || e.Rune == '\r') && act == t.Body {
		q0, q1 := act.Dot()
		if q0 == q1 {
			p := act.Bytes()
			l0, _ := find.Findlinerev(p, q0, 0)
			ntab = find.Accept(p, l0, []byte{'\t'})
			ntab -= l0 + 1
		}
	}
	kbd.SendClient(act, e)
	for ntab >= 0 {
		e.Rune = '\t'
		kbd.SendClient(act, e)
		ntab--
	}
	t.Mark()
}
