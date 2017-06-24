package main

import (
	"image"

	"github.com/as/frame"
	"github.com/as/frame/tag"
	"golang.org/x/exp/shiny/screen"
)

type Grid struct {
	*Col
}

func NewGrid(src screen.Screen, wind screen.Window, ft frame.Font, sp, size image.Point, files ...string) *Grid {
	N := len(files)
	tdy := ft.Dy() * 2
	T := tag.NewTag(src, wind, ft, image.Pt(sp.X, sp.Y), image.Pt(size.X, tdy), pad, cols)
	T.Wtag.InsertString("Newcol Killall Exit", 0)
	T.Wtag.Scroll = nil
	g := &Grid{&Col{sp: sp, src: src, size: size, wind: wind, ft: ft, Tag: T, tdy: tdy, List: make([]Plane, len(files))}}
	size.Y -= tdy
	sp.Y += tdy
	d := image.Pt(size.X/N, size.Y)
	for i, v := range files {
		g.List[i] = NewCol(src, wind, ft, sp, size, v)
		sp.X += d.X
	}
	g.List = append([]Plane{T}, g.List...)
	return g
}

func (g *Grid) Move(sp image.Point) {
	panic("never call this")
	g.sp = sp
	x := 0
	g.List[0].Move(sp)
	sp.Y += g.List[0].Loc().Dy()
	for _, co := range g.List[1:] {
		co.Move(sp.Add(image.Pt(x, 0)))
		x += co.Loc().Dx()
	}
}

// attach inserts w in position id, shifting the original right
func (g *Grid) attach(w Plane, id int) {
	if id < 1 {
		return
	}
	g.List = append(g.List[:id], append([]Plane{w}, g.List[id:]...)...)
	r := g.List[id-1].Loc()
	if id-1 == 0{
		r = image.Rect(g.sp.X, g.sp.Y+g.tdy, g.sp.X, g.sp.Y+g.size.Y)
	}
	w.Move(image.Pt(r.Max.X, g.sp.Y+g.tdy))
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

func (g *Grid) fill() {
	tdy := g.tdy
	g.List[0].Resize(image.Pt(g.size.X, tdy))
	y := g.size.Y - tdy
	x1 := g.Loc().Max.X
	//	Tagtext(fmt.Sprintf("id=maintag r=%s", g.List[0].Loc()), g.List[0])
	for n := len(g.List) - 1; n > 0; n-- {
		x0 := g.List[n].Loc().Min.X
		g.List[n].Resize(image.Pt(x1-x0, y))
		x1 = x0
	}
}

func (g *Grid) Resize(size image.Point) {
	g.size = size
	g.fill()
}
