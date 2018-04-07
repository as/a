package main // import "github.com/as/a"

import (
	"flag"
	"fmt"
	"image"
	"os"

	"github.com/as/event"
	"github.com/as/shiny/screen"
	mus "github.com/as/text/mouse"
	"golang.org/x/mobile/event/lifecycle"

	"github.com/as/edit"
	"github.com/as/ui/tag"
)

var (
	Version = "0.6.1"
	eprint  = fmt.Println
	timefmt = "15.04.05"
)

var focused = false

var (
	utf8       = flag.Bool("u", false, "enable utf8 experiment")
	elastic    = flag.Bool("elastic", false, "enable elastic tabstops")
	oled       = flag.Bool("b", false, "OLED display mode (black)")
	ftsize     = flag.Int("ftsize", 11, "font size")
	srvaddr    = flag.String("l", "", "listen (extermely dangerous) announce and serve file system clients on given endpoint")
	clientaddr = flag.String("d", "", "dial to a remote file system on given endpoint")
	quiet      = flag.Bool("q", false, "dont interact with the user or graphical subsystem (use with -l)")
)

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

func banner() {
	logf("ver=%s", Version)
	logf("pid=%d", os.Getpid())
	logf("args=%q", os.Args)
	repaint()
}

var D *screen.Device

func main() {
	defer trypprof()()
	list := argparse()
	if *quiet {
		banner()
		createnetworks()
		<-done
		os.Exit(0)
	}

	dev, wind, d, ft := frameinstall()
	D = d
	g = NewGrid(dev, image.ZP, winSize, ft, list...)
	setLogFunc(g.aerr)
	banner()
	createnetworks()
	actinit(g)

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
				if borderHit(e) {
					procBorderHit(e)
				} else {
					procButton(e)
				}
				repaint()
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
			if !unthrottled() {
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
				g.Look(e)
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
		repaint()
	}

}
