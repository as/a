package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"sync"

	window "github.com/as/ms/win"
	"github.com/as/cursor"

	"github.com/as/frame"
	"github.com/as/frame/tag"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func moveMouse(pt image.Point) {
	cursor.MoveTo(window.ClientAbs().Min.Add(pt))
}

func mkfont(size int) frame.Font {
	return frame.NewTTF(gomono.TTF, size)
}

// Put
var (
	winSize = image.Pt(1900, 1000)
	pad     = image.Pt(15, 5)
	fsize   = 11
)

type Plane interface {
	Loc() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
}

// Put
func active(pt image.Point, act Plane, list ...Plane) (x Plane) {
	if tag.Buttonsdown != 0 {
		return act
	}
	if act != nil {
		list = append([]Plane{act}, list...)
	}
	for i, w := range list {
		r := w.Loc()
		if pt.In(r) {
			return list[i]
		}
	}
	return act
}

type Col struct {
	sp   image.Point
	size image.Point
	src  screen.Screen
	wind screen.Window
	Tag  *tag.Tag
	tdy  int
	List []Plane
}

var cols = frame.Acme

func sizeof(r image.Rectangle) image.Point{
	return r.Max.Sub(r.Min)
}


func New(co *Col, filename string) (w Plane){
	last := co.List[len(co.List)-1]
	last.Loc()
	tw := co.Tag.Wtag
	t := tag.NewTag(co.src, co.wind, tw.Font, co.sp, image.Pt(co.size.X, co.tdy*2), pad, tw.Color)
	t.Open(filename)
	lsize := sizeof(last.Loc())
	lsize.Y -= lsize.Y/3
	last.Resize(lsize)
	co.attach(t, len(co.List))
	co.fill()
	return t
}

func Del(co *Col, id int){
	type Releaser interface{
		Release()
	}
	w := co.detach(id)
	if t, ok := w.(Releaser); ok{
		t.Release()
	}
	co.fill()
}

func NewCol(src screen.Screen, wind screen.Window, ft frame.Font, sp, size image.Point, files ...string) *Col {
	N := len(files)
	tdy := ft.Dy() * 2
	T := tag.NewTag(src, wind, ft, image.Pt(sp.X, sp.Y), image.Pt(size.X, tdy), pad, cols)
	//T.Open("tag")
	T.Wtag.InsertString("New Delcol Sort", 0)
	T.Wtag.Scroll = nil
	size.Y -= tdy
	sp.Y += tdy
	dy := image.Pt(size.X, size.Y/N)
	col := &Col{sp: sp, src: src, size: size, wind: wind, Tag: T, tdy: tdy, List: make([]Plane, len(files))}
	for i, v := range files {
		t := tag.NewTag(src, wind, ft, sp, dy, pad, cols)
		t.Open(v)
		col.List[i] = t
		sp.Y += dy.Y
	}
	col.List = append([]Plane{T}, col.List...)
	return col
}
// Put
func (co *Col) Upload(wind screen.Window) {
	for _, t := range co.List {
		t.(*tag.Tag).Upload(wind)
	}
}
func (co *Col) Move(sp image.Point) {
	co.sp = sp
	y := 0
	for _, t := range co.List {
		t.Move(sp.Add(image.Pt(0, y)))
		y += t.Loc().Dy()
	}
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
		t.Move(sp)
		t.Resize(dy)
		sp.Y += dy.Y
	}
}

func (co *Col) Loc() image.Rectangle {
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
	if id < 1 {
		return
	}
	// 0 1 2
	// a b c
	// a w b c
	co.List=append(co.List[:id], append([]Plane{w}, co.List[id:]...)...)
	r := co.List[id-1].Loc()
	w.Move(image.Pt(r.Min.X, r.Max.Y))
}
func (co *Col) fill() {
	x := co.size.X
	y1 := co.Loc().Max.Y
	for n := len(co.List) - 1; n > 0; n-- {
		y0 := co.List[n].Loc().Min.Y
		co.List[n].Resize(image.Pt(x, y1-y0))
		y1 = y0
		if false{
		t := co.List[n].(*tag.Tag).Wtag
		q0, q1 := t.Dot()
		t.Delete(q0, q1)
		q1 = q0
		s := fmt.Sprintf("id=%d r=%s", n, co.List[n].Loc())
		t.InsertString(s, q1)
		t.Select(q1,q1+int64(len(s)))
		}
	}
}


func eq(a, b Plane) bool{
	if a == nil || b == nil{
		return false
	}
	return a.Loc() == b.Loc()
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
	fmt.Printf("id of plane %#v is %d\n", w, id)
	return id
}

func (co *Col) MoveWin(id int, y int) {
	if id == 0 || id >= len(co.List) {
		return
	}
	s := co.detach(id)
	co.fill()
	co.Attach(s, y)
}

func (co *Col) Attach(src Plane, y int){
	did := co.IDPoint(image.Pt(co.sp.X, y))
	if did != 0 && did < len(co.List){
		d := co.List[did]
		y := y-d.Loc().Min.Y
		x := sizeof(d.Loc()).X
		d.Resize(image.Pt(x, y))
	}
	co.attach(src, did+1)
	co.fill()
}

func (co *Col) Handle(act *tag.Invertable, e interface{}) {
	for i := range co.List {
		t := co.List[i].(*tag.Tag)
		t.Handle(t.W, e)
	}
}

// Put
func main() {
	driver.Main(func(src screen.Screen) {
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y})
		wind.Send(paint.Event{})
		focused := false
		focused = focused
		ft := mkfont(fsize)

		list := []string{"/dev/stdin"}
		if len(os.Args) > 1 {
			list = os.Args[1:]
		}
		if len(list) < 2 {
			list = append(list, "guide")
		}
		sp := image.ZP
		dx := winSize.X / 2
		x := dx
		co0 := NewCol(src, wind, ft, sp, image.Pt(sp.X+x, winSize.Y), list[len(list)/2:]...)
		sp.X += dx
		co1 := NewCol(src, wind, ft, sp, image.Pt(sp.X+x, winSize.Y), list[:len(list)/2]...)
		actCol := co0
		actTag := actCol.List[1]
		act := actTag.(*tag.Tag).W

		go func() {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() {
				if x := sc.Text(); x == "u" || x == "r" {
					act.SendFirst(x)
					continue
				}
				act.SendFirst(tag.Cmdparse(sc.Text()))
			}
		}()

		var xx struct{
			sweep bool
			sweepCol bool
			sr image.Rectangle
			srcCol *Col
			src Plane
			dst Plane
			detach func()
		}
		for {
			// Put
			switch e := act.NextEvent().(type) {
			case mouse.Event:
				rpt := image.Pt(int(e.X), int(e.Y))
				pt := rpt.Add(act.Loc().Min)
				actCol = active(pt, actCol, co0, co1).(*Col)
				actTag = active(pt, actTag, actCol.List...).(*tag.Tag)
				act = active(pt, act, actTag.(*tag.Tag).W, actTag.(*tag.Tag).Wtag).(*tag.Invertable)	// Put

				if  (e.Button == 2 && e.Direction == 2 && (xx.sweep || xx.sweepCol)){
					if xx.sweepCol {
						xx.sweepCol = false
					}
					xx.sweep = false
					xx.srcCol.fill()
					actCol.Attach(xx.src, pt.Y)
					moveMouse(xx.src.Loc().Min)
					act.SendFirst(paint.Event{})
					continue
				}
				if xx.sweep && e.Button == 0{
				}

				{
					r := actTag.Loc()
					dy := r.Min.Y
					r.Max = r.Min.Add(image.Pt(20, 20))
					if e.Direction == 1 && pt.In(r){
						if e.Button == 2 {
							xx.sweep = true
							xx.srcCol = actCol
							if xx.srcCol.ID(xx.src) == 0{
								xx.sweepCol = true
							} else {
								xx.src = actTag
								xx.srcCol.detach(xx.srcCol.ID(xx.src))
							}
						} else {
							id := actCol.ID(actTag)
							if e.Button == 3 && id > 1{
								a := actCol.List[id-1].Loc()
								by := actCol.List[id].Loc().Min.Y
								dy = by-(a.Max.Y-a.Min.Y)+fsize*2
								dy += fsize*2
							} else {
								dy -= fsize*2
							}
							actCol.MoveWin(id, dy)
							moveMouse(actTag.Loc().Min)
							act.SendFirst(paint.Event{})
						}
						continue
					}
				}
				actTag.(*tag.Tag).Handle(act, e)
			case string, *tag.Command, tag.ScrollEvent, key.Event:
				if s, ok := e.(string); ok {
					if s == "New"{
						moveMouse(New(actCol, "mink").Loc().Min)
						act.SendFirst(paint.Event{})
						continue
					} else if s == "Del"{
						Del(actCol, actCol.ID(actTag))
						act.SendFirst(paint.Event{})
						continue
					}
				}
				actTag.(*tag.Tag).Handle(act, e)
			case size.Event:
				winSize = image.Pt(e.WidthPx, e.HeightPx)
				x := 0; dx := co0.size.X
				co0.Resize(image.Pt(dx, winSize.Y)); x += dx; dx = winSize.X-dx;
				co1.Resize(image.Pt(dx, winSize.Y))
				act.SendFirst(paint.Event{})
			case paint.Event:
				var wg sync.WaitGroup
				wg.Add(2)
				go func() { co0.Upload(wind); wg.Done() }()
				go func() { co1.Upload(wind); wg.Done() }()
				wg.Wait()
				wind.Publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
				// NT doesn't repaint the window if another window covers it
				if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
					focused = false
				} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
					focused = true
				}
			}
		}
	})

}
