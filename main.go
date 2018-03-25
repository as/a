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
	Version = "0.5.3"
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

// Put
func main() {
	defer trypprof()()

	// Startup
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
	aerr = g.aerr

	// This in particular needs to go
	actCol = g.List[1].(*Col)
	actTag = actCol.List[1].(*tag.Tag)
	act = actTag.Body

	alook := func(e event.Look) {
		g.Look(e)
	}

	aerr("ver=%s", Version)
	aerr("pid=%d", os.Getpid())
	aerr("args=%q", os.Args)
	if srv != nil {
		aerr("listening for remote connections")
	}
	if client != nil {
		aerr("connected to remote filesystem")
	}

	go func() {
		//		var c0 time.Time
		for {
			select {
			case e := <-D.Scroll:
				activate(p(e), g)
				doScrollEvent(act, mus.ScrollEvent{Dy: 5, Event: e})
			case e := <-D.Mouse:

				e = readmouse(e)
				if down == 0 {
					continue
				}
				s0, s1 := act.Dot()

				org := act.Origin()
				q0 := org + act.IndexOf(p(e))
				q1 := q0
				act.Sq = q0

				act.Select(q0, q1)
				ck()
				start := down
				sweepFunc := func() {
					for down == start {
						act.Sq, q0, q1 = sweep(act, e, act.Sq, q0, q1)
						act.Select(q0, q1)
						ck()
						e = readmouse(<-D.Mouse)
					}
				}
				switch start {
				case 1 << 1:
					if inSizer(p(e)) {
						aerr("InSizer: %s", p(e))
						if canopy(absP(e, act.Bounds().Min)) {
							g.dragCol(actCol, e, D.Mouse)
						} else {
							g.dragTag(actCol, actTag, e, D.Mouse)
						}
					} else if inScroll(p(e)) {
						sweepFunc()
						aerr("inScroll: %s", p(e))
					} else {
						sweepFunc()
						for down != 0 {
							if down&(1<<2) != 0 {
								act.Select(q0, q1)
								tag.Snarf(act, e)
								aerr("should snarf")

								ck()
							} else if down&(1<<3) != 0 {
								act.Select(q0, q1)
								tag.Paste(act, e)
								//act.Select(q0,q1)
								aerr("should paste")
								ck()
							}
							e = readmouse(<-D.Mouse)
						}
						act.Select(q0, q1)
						ck()
					}
				case 1 << 2:
					sweepFunc()
					act.Select(s0, s1)
					aerr("execute: %s", p(e))
					t := actTag
					act.Ctl() <- event.Cmd{
						Rec: event.Rec{
							Q0: q0, Q1: q0,
							P: act.Bytes()[q0:q1],
						},
						From: t,
						To:   []event.Editor{t.Body},
						Name: t.FileName(),
					}
				case 1 << 3:
					sweepFunc()
					act.Select(s0, s1)
					aerr("look: %s", p(e))
					t := actTag
					act.Ctl() <- event.Look{
						Rec: event.Rec{
							Q0: q0, Q1: q1,
							P: act.Bytes()[q0:q1],
						},
						From: t,
						To:   []event.Editor{t.Body},
						Name: t.FileName(),
					}
				}
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
				aerr(e.Error())
			case interface{}:
				log.Printf("missing event: %#v\n", e)
				continue Main
			}
		}
		ck()
	}

}
