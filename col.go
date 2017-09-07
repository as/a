package main

import (
	"fmt"
	"github.com/as/frame"
	"github.com/as/frame/font"
	"github.com/as/frame/tag"
	"golang.org/x/exp/shiny/screen"
	"image"
)

type Col struct {
	ft   *font.Font
	sp   image.Point
	size image.Point
	src  screen.Screen
	wind screen.Window
	Tag  *tag.Tag
	tdy  int
	List []Plane
}

func NewCol(src screen.Screen, wind screen.Window, ft *font.Font, sp, size image.Point, files ...string) *Col {
	N := len(files)
	tdy := ft.Dy() + ft.Dy()/2
	tagpad := image.Pt(pad.X, 3)
	T := tag.NewTag(src, wind, ft, image.Pt(sp.X, sp.Y), image.Pt(size.X, tdy), tagpad, frame.ATag1)
	//T.Open("tag")
	T.Win.InsertString("New Delcol Sort", 0)
	col := &Col{sp: sp, src: src, size: size, wind: wind, ft: ft, Tag: T, tdy: tdy, List: make([]Plane, len(files))}
	size.Y -= tdy
	sp.Y += tdy
	dy := image.Pt(size.X, size.Y/N)
	for i, v := range files {
		t := tag.NewTag(src, wind, ft, sp, dy, pad, frame.ATag1)
		t.Get(v)
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
	fmt.Printf("newcol sp=%s size=%s\n", sp, size)
	col := NewCol(g.src, g.wind, g.ft, sp, size, filenames...)
	g.attach(col, len(g.List))
	g.fill()
	return col
}

func New(co *Col, filename string) (w Plane) {
	last := co.List[len(co.List)-1]
	last.Loc()
	tw := co.Tag.Win
	co.PrintList()
	t := tag.NewTag(co.src, co.wind, tw.Font, co.sp, image.Pt(co.size.X, co.tdy*2), pad, tw.Color)
	t.Get(filename)
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
	if t, ok := w.(Releaser); ok {
		t.Release()
	}
	co.fill()
}

func (co *Col) Move(sp image.Point) {
	co.sp = sp
	dy := 0
	for _, t := range co.List {
		sp0 := image.Pt(sp.X, sp.Y+dy)
		fmt.Printf("movewin -> %s\n", sp0)
		t.Move(sp0)
		dy = t.Loc().Dy()
	}
	co.fill()
}

func (co *Col) Resize(size image.Point) {
	co.size = size
	co.fill()
	return

	size.Y = co.size.Y - co.tdy
	sp := co.sp
	sp.Y += co.tdy
	N := len(co.List) - 1
	dy := image.Pt(size.X, size.Y/N)
	for _, t := range co.List[1:] {
		//		t.Move(sp)
		t.Resize(dy)
		sp.Y += dy.Y
	}
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
