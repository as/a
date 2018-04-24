package main

import (
	"image"

	"github.com/as/edit"
	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/text"
	"github.com/as/ui"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

type Grid struct {
	*Col
}

var (
	GridLabel = "Newcol Killall Exit    guru^(callees callers callstack definition describe freevars implements peers pointsto referrers what whicherrs)"
)

func NewGrid(dev ui.Dev, sp, size image.Point, ft font.Face, files ...string) *Grid {
	conf := GridConfig
	g := &Grid{col.NewGridHack(dev, sp, size, tagHeight, ft)}
	T := tag.New(dev, sp, image.Pt(size.X, tagHeight), conf)
	T.Win.InsertString(GridLabel, 0)
	g.Tag = T
	N := len(files)
	g.List = make([]Plane, N)

	size.Y -= tagHeight
	sp.Y += tagHeight // Put
	d := image.Pt(size.X/N, size.Y)
	for i, v := range files {
		g.List[i] = NewCol(dev, ft, sp, d, v)
		sp.X += d.X
	}

	g.Refresh()
	return g
}

func (g *Grid) Attach(src Plane, x int) {
	r := g.Tag.Loc()
	r.Min.Y = r.Max.Y
	pt := r.Min
	if len(g.List) == 0 {
		src.Move(pt)
		g.AttachFill(src, 0)
		return
	}
	pt.X = x
	src.Move(pt)
	did := g.IDPoint(pt)
	g.AttachFill(src, did+1)
}

func (g *Grid) Delta(n int) image.Point {
	x0 := g.List[n].Loc().Min.X
	x1 := g.Loc().Max.X

	if n+1 != len(g.List) {
		x1 = g.List[n+1].Loc().Min.X
	}
	return identity(x1-x0, g.Loc().Dy()-g.Tag.Loc().Dy())
}

func (g *Grid) Label() *win.Win { return g.Tag.Win }

func (g *Grid) Move(sp image.Point) {
	panic("never call this")
}

func (g *Grid) Resize(size image.Point) {
	g.ForceSize(size)
	g.Tag.Resize(image.Pt(size.X, g.Tag.Loc().Dy()))
	g.fill()
}

func (g *Grid) fill() {
	fill(g)
	g.Tag.Resize(g.Tag.Loc().Size())
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
	w, _ := t.Body.(*win.Win)
	if w == nil {
		return
	}
	var green = frame.Palette{
		Back: frame.Green,
		Text: frame.A.Text,
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
			w.Frame.Recolor(fr.PointOf(dot.Q0), dot.Q0, dot.Q1, green)
		}
		//prog.Emit = &edit.Emitted{}
	})
}
