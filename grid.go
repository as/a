package main

import (
	"image"

	"github.com/as/edit"
	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/text"
	"github.com/as/ui"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

type Grid struct {
	*Col
}

func NewGrid(dev *ui.Dev, sp, size image.Point, ft font.Face, files ...string) *Grid {
	N := len(files)
	tdy := tag.TagSize(ft)
	tagpad := tag.TagPad(pad)
	conf := &tag.Config{
		Filesystem: newfsclient(),
		Margin:     tagpad,
		Facer:      font.NewFace,
		FaceHeight: ft.Height(),
		Color: [3]frame.Color{
			0: frame.ATag0,
		},
		Ctl: events,
	}
	T := tag.New(dev, sp, image.Pt(size.X, tdy), conf)
	T.Win.InsertString("Newcol Killall Exit    guru^(callees callers callstack definition describe freevars implements peers pointsto referrers what whicherrs)", 0)
	g := &Grid{&Col{dev: dev, sp: sp, size: size, ft: ft, Tag: T, tdy: tdy, List: make([]Plane, len(files))}}
	size.Y -= tdy
	sp.Y += tdy
	d := image.Pt(size.X/N, size.Y)
	for i, v := range files {
		g.List[i] = NewCol(dev, ft, sp, size, v)
		sp.X += d.X
	}
	g.List = append([]Plane{T}, g.List...)
	g.Refresh()
	return g
}

func (g *Grid) Attach(src Plane, x int) {
	did := g.IDPoint(image.Pt(x, g.sp.Y+g.tdy))
	if did != 0 && did < len(g.List) {
		d := g.List[did]
		x := x - d.Loc().Min.X
		y := sizeof(d.Loc()).Y
		d.Resize(image.Pt(x, y))
	}
	g.attach(src, did+1)
	g.fill()
}

func (g *Grid) Move(sp image.Point) {
	panic("never call this")
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

func (g *Grid) Resize(size image.Point) {
	g.size = size
	g.fill()
}

// attach inserts w in position id, shifting the original right
func (g *Grid) attach(w Plane, id int) {
	if id < 1 {
		return
	}
	g.List = append(g.List[:id], append([]Plane{w}, g.List[min(id, len(g.List)):]...)...)
	r := g.List[id-1].Loc()
	if id-1 == 0 {
		r = image.Rect(g.sp.X, g.sp.Y+g.tdy, g.sp.X, g.sp.Y+g.size.Y)
	}
	w.Move(image.Pt(r.Max.X, g.sp.Y+g.tdy))
}

func (g *Grid) fill() {
	tdy := g.tdy
	g.List[0].Resize(image.Pt(g.size.X, tdy))
	y := g.size.Y - tdy
	x1 := g.Loc().Max.X
	for n := len(g.List) - 1; n > 0; n-- {
		x0 := g.List[n].Loc().Min.X
		g.List[n].Resize(image.Pt(x1-x0, y))
		x1 = x0
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
