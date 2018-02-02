package main

import (
	"fmt"
	"image"
	"io"
	"sync"

	"github.com/as/frame"
	"github.com/as/frame/font"
	"github.com/as/ui"
	"github.com/as/ui/tag"
	"github.com/as/shiny/screen"
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

func New(co *Col, basedir, name string, sizerFunc ...func(int) int) (w Plane) {
	last := co.List[len(co.List)-1]
	last.Loc()
	tw := co.Tag.Win
	t := tag.New(co.dev, co.sp, image.Pt(co.size.X, co.tdy*2), pad, tw.Font, tw.Color)
	if name != "+Errors" {
		//g.aerr(fmt.Sprintf("open basedir=%q name=%q", basedir, name))
	}
	t.Open(basedir, name)
	t.Insert([]byte(" [Edit  ,x]"), t.Len())
	lsize := sizeof(last.Loc())

	fn := SizeThirdOf
	if len(sizerFunc) != 0 {
		fn = sizerFunc[0]
	}
	lsize.Y = fn(lsize.Y)
	last.Resize(lsize)
	co.attach(t, len(co.List))
	co.fill()
	return t
}

func Delcol(g *Grid, id int) {
	co := g.detach(id)
	x := co.Loc().Min.X
	y := co.Loc().Min.Y
	for ; id < len(g.List); id++ {
		x2 := g.List[id].Loc().Min.X
		g.List[id].Move(image.Pt(x, y))
		x = x2
	}
	g.fill()
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

func (col *Col) PrintList() {
	for i, v := range col.List {
		fmt.Printf("%d: %#v\n", i, v)
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

func (co *Col) Loc() image.Rectangle {
	if co == nil {
		return image.ZR
	}
	return image.Rectangle{co.sp, co.sp.Add(co.size)}
}

func (co *Col) Move(sp image.Point) {
	co.sp.X = sp.X
	for _, t := range co.List {
		sp := image.Pt(sp.X, t.Loc().Min.Y)
		t.Move(sp)
	}
	co.fill()
}

func (col *Col) Refresh() {
	for _, v := range col.List {
		v.Refresh()
	}
}

func (co *Col) Resize(size image.Point) {
	co.size = size
	co.fill()
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
func (co *Col) badID(id int) bool {
	return id <= 0
}
func (co *Col) bestGrowth(id int, dy int) int {
	if co.badID(id) || co.badID(id-1) {
		return dy
	}
	nclicks := co.List[id-1].Loc().Dy() / dy
	if nclicks < 3 {
		return dy
	}
	if nclicks < 5 {
		return dy * 2
	}
	if nclicks < 8 {
		return dy * 3
	}
	return dy * 4
}
func (co *Col) Grow(id int, dy int) {
	a, b := id-1, id
	if co.badID(a) || co.badID(b) {
		return
	}
	ra, rb := co.List[a].Loc(), co.List[b].Loc()
	ra.Max.Y -= dy
	if dy := ra.Dy() - tagHeight; dy < 0 {
		co.Grow(a, -dy)
	}
	fmt.Printf("min: %d, dy: %d, min-dy: %d\n", rb.Min.Y, dy, rb.Min.Y-dy)
	co.MoveWin(b, rb.Min.Y-dy)
}

/*
func (co *Col) Show(id int, dy int){
	if co.badID(id){
		return
	}
	r := co.List[id].Loc()
	if r.Dy() > dy{
		return
	}
}
*/

func (co *Col) RollDown(id int, dy int) {
	/*
		if id >= len(co.List) {
			return
		}
		r := co.List[x].Loc()
		if r.Min.Y+dy
		a := r.Min.Y
		b := a+co.List[x].Loc().Min.Y

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

func (co *Col) MoveWin(id int, y int) {
	if id == 0 || id >= len(co.List) {
		return
	}
	maxy := co.List[len(co.List)-1].Loc().Max.Y - tagHeight
	if y >= maxy {
		return
	}
	s := co.detach(id)
	co.fill()
	co.Attach(s, y)
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

func (co *Col) detach(id int) Plane {
	if id < 1 || id > len(co.List)-1 {
		return nil
	}
	w := co.List[id]
	copy(co.List[id:], co.List[id+1:])
	co.List = co.List[:len(co.List)-1]
	return w
}

func (co *Col) fill() {
	if co == nil || co.List[0] == nil {
		return
	}
	var wg sync.WaitGroup
	defer wg.Wait()

	ty := co.List[0].Loc().Dy()

	co.List[0].Resize(image.Pt(co.size.X, ty))
	//		Tagtext(fmt.Sprintf("id=tagtag r=%s", co.List[0].Loc()), co.List[0])

	x := co.size.X
	y1 := co.Loc().Max.Y
	for n := len(co.List) - 1; n > 0; n-- {
		n := n
		y0 := co.List[n].Loc().Min.Y
		wg.Add(1)
		pt := image.Pt(x, y1-y0)
		go func() {
			co.List[n].Resize(pt)
			defer wg.Done()
		}()
		y1 = y0
		//		Tagtext(fmt.Sprintf("id=%d r=%s", n, co.List[n].Loc()), co.List[n])
	}
}
