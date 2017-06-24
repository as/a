package main

import (
	"bufio"
	"image"
	"os"

	"github.com/as/cursor"
	window "github.com/as/ms/win"

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

var cols = frame.Acme

func Tagtext(s string, w Plane) {
	switch w := w.(type) {
	case *Grid:
		Tagtext(s, w.List[0])
	case *Col:
		Tagtext(s, w.List[0])
	case *tag.Tag:
		t := w.Wtag
		q0, q1 := t.Dot()
		t.Delete(q0, q1)
		q1 = q0
		t.InsertString(s, q1)
		t.Select(q1, q1+int64(len(s)))
	case Plane:
	}
}

// Put
func main() {
	driver.Main(func(src screen.Screen) {
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "A"})
		wind.Send(paint.Event{})
		focused := false
		focused = focused
		ft := mkfont(fsize)
		var list = []string{}
		if len(os.Args) > 1 {
			list = append(list, os.Args[1:]...)
		} else {
			list = append(list, "guide")
			list = append(list, ".")
		}
		g := NewGrid(src, wind, ft, image.ZP, image.Pt(winSize.X, winSize.Y), list...)
		actCol = g.List[1].(*Col)
		actTag = actCol.List[1].(*tag.Tag)
		act = actTag.W

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

		var xx struct {
			sweep    bool
			sweepCol bool
			sr       image.Rectangle
			srcCol   *Col
			src      Plane
			dst      Plane
			detach   func()
		}
		for {
			// Put
			switch e := act.NextEvent().(type) {
			case tag.GetEvent:
				t := New(actCol, e.Path)
				if e.Addr != ""{
					actTag = t.(*tag.Tag)
					act = actTag.W
					actTag.Handle(actTag.W, tag.Cmdparse(e.Addr))
					p0, _ := act.Frame.Dot()
					moveMouse(act.Loc().Min.Add(act.PointOf(p0)))
				} else {
					moveMouse(t.Loc().Min)
				}
			case mouse.Event:
				rpt := image.Pt(int(e.X), int(e.Y))
				pt := rpt.Add(act.Loc().Min)

				activate(pt, g)

				if e.Button == 2 && e.Direction == 2 && (xx.sweep || xx.sweepCol) {
					if xx.sweepCol {
						xx.sweepCol = false
						g.fill()
						g.Attach(xx.srcCol, pt.X)
						moveMouse(xx.srcCol.Loc().Min)
						act.SendFirst(paint.Event{})
					} else {
						xx.sweep = false
						xx.srcCol.fill()
						actCol.Attach(xx.src, pt.Y)
						moveMouse(xx.src.Loc().Min)
					}
					act.SendFirst(paint.Event{})
					continue
				}
				if (xx.sweep || xx.sweepCol) && e.Button == 0 {
					continue
				}

				{
					r := actTag.Loc()
					dy := r.Min.Y
					r.Max = r.Min.Add(image.Pt(20, 20))
					if e.Direction == 1 && pt.In(r) {
						if e.Button == 2 {
							xx.srcCol = actCol
							if pt.Y > g.sp.Y+g.tdy && pt.Y < g.sp.Y+g.tdy*2 {
								xx.sweepCol = true
								xx.src = nil
								g.detach(g.ID(xx.srcCol))
								g.fill()
							} else {
								xx.sweep = true
								xx.src = actTag
								xx.srcCol.detach(xx.srcCol.ID(xx.src))
							}
						} else {
							id := actCol.ID(actTag)
							if e.Button == 3 && id > 1 {
								a := actCol.List[id-1].Loc()
								by := actCol.List[id].Loc().Min.Y
								dy = by - (a.Max.Y - a.Min.Y) + fsize*2
								dy += fsize * 2
							} else {
								dy -= fsize * 2
							}
							actCol.MoveWin(id, dy)
							moveMouse(actTag.Loc().Min)
							act.SendFirst(paint.Event{})
						}
						continue
					}
				}
				actTag.Handle(act, e)
			case string, *tag.Command, tag.ScrollEvent, key.Event:
				if s, ok := e.(string); ok {
					if s == "New" {
						moveMouse(New(actCol, "mink").Loc().Min)
						act.SendFirst(paint.Event{})
						continue
					} else if s == "Newcol"{
						moveMouse(NewCol2(g, "mink").Loc().Min)
						act.SendFirst(paint.Event{})
					} else if s == "Del" {
						Del(actCol, actCol.ID(actTag))
						act.SendFirst(paint.Event{})
						continue
					}
				}
				actTag.Handle(act, e)
			case size.Event:
				winSize = image.Pt(e.WidthPx, e.HeightPx)
				g.Resize(winSize)
				act.SendFirst(paint.Event{})
			case paint.Event:
				g.Upload(wind)
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
