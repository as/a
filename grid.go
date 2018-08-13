package main

import (
	"image"

	"github.com/as/edit"
	"github.com/as/frame"
	"github.com/as/rgba"
	"github.com/as/text"
	"github.com/as/ui"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

type Grid struct {
	col.Table2
}

var (
	GridLabel = "Newcol Killall Dump Load Exit  guru^(callees callers callstack definition describe freevars implements peers pointsto referrers what whicherrs)"
)

func NewGrid(dev ui.Dev, conf *tag.Config, files ...string) *Grid {
	g := &Grid{col.NewTable2(dev, conf)}
	g.Tag.Label.InsertString(GridLabel, 0)
	return g
}

func (g *Grid) Move(sp image.Point) {
	g.Table2.Move(sp)
}

func (g *Grid) Resize(size image.Point) {
	g.ForceSize(size)
	g.Tag.Resize(image.Pt(size.X, g.Config.TagHeight()))
	col.Fill(g)
}

var InstallPalette = frame.Palette{
	Back: rgba.Seagreen,
	Text: frame.A.Text,
}

// Install places the given edit script in between
// calls to the target windows SetOrigin method. This
// is an experiment to test out highlighting with
// structural regular expressions.
//
// The current implementation will change and it
// has unfavorable performance characteristics (i.e., compiling
// the script every time), however, this isn't usually noticable
// unless the command is long
//
// Conventionally, the command should be in the form
//	,x,string,h
// Any other use is undefined and untested for now
func (g *Grid) Install(t *tag.Tag, srcprog string) {
	w, _ := t.Window.(*win.Win)
	if w == nil {
		return
	}

	prog, err := edit.Compile(srcprog)
	if err != nil {
		g.aerr(err.Error())
		return
	}

	w.FuncInstall(func(w *win.Win) {
		fr := w.Frame
		buf := text.BufferFrom(w.Bytes()[w.Origin() : w.Origin()+fr.Len()])
		ed, _ := text.Open(buf)
		prog.Run(ed)
		for _, dot := range prog.Emit.Dot {
			w.Frame.Recolor(fr.PointOf(dot.Q0), dot.Q0, dot.Q1, InstallPalette)
		}
		//prog.Emit = &edit.Emitted{}
	})
}
