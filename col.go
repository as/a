package main

import (
	"fmt"
	"image"
	"io"

	"github.com/as/edit"
	"github.com/as/frame"
	"github.com/as/frame/font"
	"github.com/as/text"
	"github.com/as/ui"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
	"golang.org/x/exp/shiny/screen"
)

type Col struct {
	dev  *ui.Dev
	ft   *font.Font
	sp   image.Point
	size image.Point
	Tag  *tag.Tag
	tdy  int
	List []Plane
}

func NewCol(dev *ui.Dev, ft *font.Font, sp, size image.Point, files ...string) *Col {
	N := len(files)
	tdy := ft.Dy() + ft.Dy()/2
	tagpad := image.Pt(pad.X, 3)
	T := tag.New(dev, sp, image.Pt(size.X, tdy), tagpad, ft, frame.ATag1)
	//T.Open(path.NewPath(""))
	T.Win.InsertString("New Delcol Sort", 0)
	col := &Col{dev: dev, sp: sp, size: size, ft: ft, Tag: T, tdy: tdy, List: make([]Plane, len(files))}
	size.Y -= tdy
	sp.Y += tdy
	dy := image.Pt(size.X, size.Y/N)
	for i, v := range files {
		t := tag.New(dev, sp, dy, pad, ft, frame.ATag1)
		t.Get(v)
		t.Insert([]byte(" [Edit  ,x]"), t.Len())
		col.List[i] = t
		sp.Y += dy.Y
	}
	col.List = append([]Plane{T}, col.List...)
	return col
}

func (col *Col) PrintList() {
	for i, v := range col.List {
		fmt.Printf("%d: %#v\n", i, v)
	}
}

func (col *Col) Refresh() {
	for _, v := range col.List {
		v.Refresh()
	}
}

func NewCol2(g *Grid, filenames ...string) (w Plane) {
	x0 := g.List[0].Loc().Min.X
	y0 := g.List[0].Loc().Dy()
	x1 := g.sp.X + g.size.X
	y1 := g.sp.X + g.size.Y - y0
	if len(g.List) > 1 {
		last := g.List[len(g.List)-1]
		last.Resize(image.Pt(last.Loc().Dx()/2, last.Loc().Dy()))
		x0 = last.Loc().Max.X
		x1 = x0 + last.Loc().Dx()/2
	}
	sp := image.Pt(x0, y0)
	size := image.Pt(x1-x0, y1-y0)
	col := NewCol(g.dev, g.ft, sp, size, filenames...)
	g.attach(col, len(g.List))
	g.fill()
	return col
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
//    ,x,string,h
// Any other use is undefined and untested for now
func (g *Grid) Install(t *tag.Tag, srcprog string) {

	var green = frame.Palette{
		Back: frame.Green,
		Text: frame.A.Text,
	}

	prog, err := edit.Compile(srcprog)
	if err != nil {
		g.aerr(err.Error())
		return
	}

	if t.Body != nil {
		t.Body.FuncInstall(func(w *win.Win) {
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
}

func New(co *Col, basedir, name string) (w Plane) {
	last := co.List[len(co.List)-1]
	last.Loc()
	tw := co.Tag.Win
	t := tag.New(co.dev, co.sp, image.Pt(co.size.X, co.tdy*2), pad, tw.Font, tw.Color)
	t.Open(basedir, name)
	t.Insert([]byte(" [Edit  ,x]"), t.Len())
	lsize := sizeof(last.Loc())
	lsize.Y -= lsize.Y / 3
	last.Resize(lsize)
	co.attach(t, len(co.List))
	co.fill()
	return t
}

func Del(co *Col, id int) {
	type Releaser interface {
		Release()
	}
	w := co.detach(id)
	y := w.Loc().Min.Y
	x := co.Loc().Min.X
	w.(io.Closer).Close()
	for ; id < len(co.List); id++ {
		y2 := co.List[id].Loc().Min.Y
		co.List[id].Move(image.Pt(x, y))
		y = y2
	}
	co.fill()
}

func (co *Col) Close() error {
	for _, t := range co.List {
		if t == nil {
			continue
		}
		if t, ok := t.(io.Closer); ok {
			t.Close()
		}
	}
	co.List = nil
	return nil
}

func (co *Col) RollUp(id int, dy int) {
	if id <= 0 || id >= len(co.List) {
		return
	}
	for x := 2; x <= id; x++ {
		a := co.List[x-1].Loc()
		dy = a.Min.Y + tagHeight
		co.MoveWin(x, dy)
	}
	co.MoveWin(id, dy)
}

func (co *Col) RollDown(id int, dy int) {
	/*
		if id >= len(co.List) {
			return
		}
		x:=id
		a := co.List[x].Loc()
		for x+1 < len(co.List){
			b := co.List[x+1].Loc()
			if extra := a.Min.Y+tagHeight+dy-b.Min.Y; extra > 0{
				co.MoveWin(x, a.Min.Y+extra)
			}
		}
		if a.Min.Y+dy > co.Loc().Max.Y{
			dy = co.Loc().Max.Y - tagHeight
		}
		co.MoveWin(id, a.Min.Y+dy)
	*/
}

func (co *Col) Move(sp image.Point) {
	co.sp.X = sp.X
	for _, t := range co.List {
		sp := image.Pt(sp.X, t.Loc().Min.Y)
		t.Move(sp)
	}
	co.fill()
}

func (co *Col) Resize(size image.Point) {
	co.size = size
	co.fill()
}

func (co *Col) FindName(name string) *tag.Tag {
	for _, v := range co.List[1:] {
		switch v := v.(type) {
		case *Col:
			t := v.FindName(name)
			if t != nil {
				return t
			}
		case *tag.Tag:
			if v.FileName() == name {
				return v
			}

		}
	}
	return nil
}

func (co *Col) Upload(wind screen.Window) {
	type Uploader interface {
		Upload(screen.Window)
		//		Dirty() bool
	}
	for _, t := range co.List {
		if t, ok := t.(Uploader); ok {
			//if co.Dirty(){
			t.Upload(wind)
			//}
		}
	}
}
func (co *Col) Loc() image.Rectangle {
	if co == nil {
		return image.ZR
	}
	return image.Rectangle{co.sp, co.sp.Add(co.size)}
}

func (co *Col) detach(id int) Plane {
	if id < 1 || id > len(co.List)-1 {
		return nil
	}
	w := co.List[id]
	copy(co.List[id:], co.List[id+1:])
	co.List = co.List[:len(co.List)-1]
	return w
}

// attach inserts w in position id, shifting the original forwards
func (co *Col) attach(w Plane, id int) {
	if w == nil || w == co.List[0] || id < 1 {
		return
	}
	co.List = append(co.List[:id], append([]Plane{w}, co.List[id:]...)...)
	r := co.List[id-1].Loc()
	if len(co.List) > 2 {
		w.Move(image.Pt(r.Min.X, r.Max.Y))
	}
}

func (co *Col) fill() {
	if co == nil || co.List[0] == nil {
		return
	}
	ty := co.List[0].Loc().Dy()
	co.List[0].Resize(image.Pt(co.size.X, ty))
	//		Tagtext(fmt.Sprintf("id=tagtag r=%s", co.List[0].Loc()), co.List[0])

	x := co.size.X
	y1 := co.Loc().Max.Y
	for n := len(co.List) - 1; n > 0; n-- {
		y0 := co.List[n].Loc().Min.Y
		co.List[n].Resize(image.Pt(x, y1-y0))
		y1 = y0
		//		Tagtext(fmt.Sprintf("id=%d r=%s", n, co.List[n].Loc()), co.List[n])
	}
}

func (co *Col) MoveWin(id int, y int) {
	if id == 0 || id >= len(co.List) {
		return
	}
	s := co.detach(id)
	co.fill()
	co.Attach(s, y)
}

func (co *Col) Attach(src Plane, y int) {
	did := co.IDPoint(image.Pt(co.sp.X, y))
	if did == 0 || did >= len(co.List) {
		return
	}
	d := co.List[did]
	y -= d.Loc().Min.Y
	x := sizeof(d.Loc()).X
	d.Resize(image.Pt(x, y))
	co.attach(src, did+1)
	co.fill()
}

func (co *Col) Handle(e interface{}) {
	for i := range co.List {
		t := co.List[i]
		switch t := t.(type) {
		case (*tag.Tag):
			t.Handle(t.Body, e)
		}
	}
}

func (co *Col) IDPoint(pt image.Point) (id int) {
	for id = 0; id < len(co.List); id++ {
		if pt.In(co.List[id].Loc()) {
			break
		}
	}
	return id
}
func (co *Col) ID(w Plane) (id int) {
	for id = 0; id < len(co.List); id++ {
		if eq(w, co.List[id]) {
			break
		}
	}
	return id
}
