package main

import (
	"image"

	"github.com/as/shiny/event/key"
	"github.com/as/shiny/screen"
	"github.com/as/text/find"
	"github.com/as/ui/kbd"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

// TODO(as): This interface is too big and ugly
// get rid of it
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
	Bounds() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
	Close() error
	Upload()
	Window() screen.Window
	Write(p []byte) (n int, err error)
}

func kbdin(e key.Event, t *tag.Tag, act Window) {
	//	SetKMod(KFlag(e.Modifiers))
	if e.Direction == 2 {
		return
	}

	switch e.Code {
	case key.CodeEscape:
		track.esc()
		return
	case key.CodeI:
		if e.Modifiers == key.ModControl {
			runGoImports(t, e)
			reload(t)
			return
		}
	case key.CodeEqualSign, key.CodeHyphenMinus:
		if e.Modifiers == key.ModControl {
			win, _ := t.Window.(*win.Win)
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
			t.SetFont(t.Config.Facer(size))
			return
		}
	}

	track.set(false)

	ntab := int64(-1)
	if (e.Rune == '\n' || e.Rune == '\r') && act == t.Window {
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
