package main // import "github.com/as/a"

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/as/event"
	"github.com/as/shiny/screen"
	"github.com/as/text/find"
	mus "github.com/as/text/mouse"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/time/rate"

	"github.com/as/edit"
	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/ui"
	"github.com/as/ui/tag"
)

var (
	Version = "0.6.0"
	eprint  = fmt.Println
	timefmt = "2006.01.02 15.04.05"
)

var focused = false

func argparse() (list []string) {
	if len(flag.Args()) > 0 {
		list = append(list, flag.Args()...)
	} else {
		list = append(list, "guide")
		list = append(list, ".")
	}
	return
}

var (
	utf8       = flag.Bool("u", false, "enable utf8 experiment")
	elastic    = flag.Bool("elastic", false, "enable elastic tabstops")
	oled       = flag.Bool("b", false, "OLED display mode (black)")
	ftsize     = flag.Int("ftsize", 11, "font size")
	srvaddr    = flag.String("l", "", "listen (extermely dangerous) announce and serve file system clients on given endpoint")
	clientaddr = flag.String("d", "", "dial to a remote file system on given endpoint")
	quiet      = flag.Bool("q", false, "dont interact with the user or graphical subsystem (use with -l)")
)

var dirty bool
var ck = repaint

func repaint() {
	if !dirty || act == nil {
		return
	}
	select {
	case act.Window().Device().Paint <- paint.Event{}:
	default:
	}
}

var (
	g         *Grid
	events    = make(chan interface{}, 5)
	done      = make(chan bool)
	moribound = make(chan bool, 1)
)

func init() {
	// this grants the capability to shut down the program
	// it happens exactly once
	moribound <- true
}

func init() {
	flag.Parse()
}

func main() {
	defer trypprof()()

	err := createnetworks()
	if err != nil {
		log.Fatalln(err)
	}

	lim := rate.NewLimiter(rate.Every(time.Second/120), 2)

	if *oled {
		usedarkcolors()
	}

	list := argparse()
	if *quiet {
		<-done
		os.Exit(0)
	}

	frame.ForceUTF8 = *utf8
	frame.ForceElastic = *elastic
	dev, err := ui.Init(&screen.NewWindowOptions{Width: winSize.X, Height: winSize.Y, Title: "A"})
	if err != nil {
		log.Fatalln(err)
	}

	wind := dev.Window()
	D := wind.Device()

	// Linux will segfault here if X is not present
	repaint()
	ft := font.NewFace(*ftsize)
	g = NewGrid(dev, image.ZP, winSize, ft, list...)
	setLogFunc(g.aerr)

	// This in particular needs to go
	actCol = g.List[1].(*Col)
	actTag = actCol.List[1].(*tag.Tag)
	act = actTag.Body

	alook := func(e event.Look) {
		g.Look(e)
	}

	logf("ver=%s", Version)
	logf("pid=%d", os.Getpid())
	logf("args=%q", os.Args)
	if srv != nil {
		logf("listening for remote connections")
	}
	if client != nil {
		logf("connected to remote filesystem")
	}

	var (
		double bool
		last   = down
		lastpt image.Point
	)
	last = 0
	var t0 = time.Now()
	go func() {
		for {
			select {
			case e := <-D.Scroll:
				activate(p(e), g)
				scroll(act, mus.ScrollEvent{Dy: 5, Event: e})
			case e := <-D.Mouse:
				activate(p(e), g)
				e = rel(readmouse(e), act)
				if down == 0 {
					continue
				}
				if last == down {
					if time.Since(t0) < time.Second/2 && lastpt.In(image.Rect(-3, -3, 3, 3).Add(p(e))) {
						double = true
					}
				}
				t0 = time.Now()
				last = down
				lastpt = p(e)
				if pt := p(e); inSizer(pt) || inScroll(pt) {
					if inSizer(pt) {
						if HasButton(1, down) {
							if canopy(absP(e, act.Bounds().Min)) {
								g.dragCol(actCol, e, D.Mouse)
							} else {
								g.dragTag(actCol, actTag, e, D.Mouse)
							}
						} else {
							switch down {
							case Button(2):
							case Button(3):
								actCol.RollUp(actCol.ID(actTag), act.Loc().Min.Y)
								moveMouse(act.Loc().Min)
							}
							for down != 0 {
								readmouse(<-D.Mouse)
							}
						}
					} else if inScroll(pt) {
						switch down {
						case Button(1):
							scroll(act, mus.ScrollEvent{Dy: 5, Event: e})
						case Button(2):
							w := act
							e = rel(readmouse(e), w)
							scroll(w, mus.ScrollEvent{Dy: 5, Event: e})
							ck()
						case Button(3):
							scroll(act, mus.ScrollEvent{Dy: -5, Event: e})
						}
						logf("inScroll: %s", p(e))
					}
					continue
				}

				t, w := actTag, act
				s0, s1 := w.Dot()
				q0 := w.IndexOf(p(e)) + w.Origin()
				q1 := q0
				act.Select(q0, q1)
				ck()

				switch down {
				case Button(1):
					if double {
						q0, q1 = find.FreeExpand(w, q0)
						double = false
					} else {
						q0, q1, e = sweepFunc(w, e, D.Mouse)
						for down != 0 {
							t.Select(q0, q1)
							if HasButton(2, down) {
								tag.Snarf(w, e)
							} else if HasButton(3, down) {
								tag.Paste(w, e)
							}
							ck()
							e = rel(readmouse(<-D.Mouse), t)
						}
						t0 = time.Now()
					}
					w.Select(q0, q1)
					ck()
				case Button(2):
					q0, q1, _ := sweepFunc(w, e, D.Mouse)
					if q0 == q1 {
						q0, q1 = find.ExpandFile(act.Bytes(), q0)
					}
					w.Select(s0, s1)
					w.Ctl() <- event.Cmd{
						Name: t.FileName(),
						From: t, To: []event.Editor{lookTarget(act, t)},
						Rec: event.Rec{Q0: q0, Q1: q0, P: w.Bytes()[q0:q1]},
					}
				case Button(3):
					q0, q1, _ := sweepFunc(w, e, D.Mouse)
					if q0 == q1 {
						q0, q1 = find.ExpandFile(act.Bytes(), q0)
					}
					w.Select(s0, s1)
					w.Ctl() <- event.Look{
						Name: t.FileName(),
						From: t, To: []event.Editor{lookTarget(act, t)},
						Rec: event.Rec{Q0: q0, Q1: q1, P: w.Bytes()[q0:q1]},
					}
				}
				ck()
			}
		}
	}()

	go func() {
		for e := range D.Key {
			actTag.Handle(act, e)
			dirty = true
			repaint()
		}
	}()

Main:
	for {
		select {
		case e := <-D.Size:
			winSize = image.Pt(e.WidthPx, e.HeightPx)
			g.Resize(winSize)
		case e := <-D.Paint:
			if !lim.Allow() {
				continue Main
			}
			if e.External {
				g.Resize(winSize)
			}
			g.Upload(wind)
			wind.Publish()
			continue Main
		case e := <-events:
			switch e := e.(type) {
			case tag.GetEvent:
				t := New(actCol, e.Basedir, e.Name)
				if e.Addr != "" {
					actTag = t.(*tag.Tag)
					act = actTag.Body
					actTag.Handle(actTag.Body, edit.MustCompile(e.Addr))
					p0, _ := act.Frame.Dot()
					moveMouse(act.Loc().Min.Add(act.PointOf(p0)))
				} else {
					moveMouse(t.Loc().Min)
				}
			case mus.SnarfEvent, mus.InsertEvent:
				actTag.Handle(act, e)
			case event.Look:
				alook(e)
			case event.Cmd:
				acmd(e)
			case edit.File:
				g.acolor(e)
			case edit.Print:
				g.aout(string(e))
			case error:
				logf(e.Error())
			case interface{}:
				logf("missing event: %#v\n", e)
				continue Main
			}
		case e := <-D.Lifecycle:
			if e.To == lifecycle.StageDead {
				return
			}
			// NT doesn't repaint the window if another window covers it
			if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
				focused = false
			} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
				focused = true
				continue Main
			}
		}
		ck()
	}

}
