package main // import "github.com/as/a"

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/as/shiny/screen"
	mus "github.com/as/text/mouse"
	"golang.org/x/mobile/event/lifecycle"

	"github.com/as/edit"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

var (
	Version = "0.6.9"
	eprint  = fmt.Println
	timefmt = "15.04.05"
)

var (
	g         *Grid
	D         *screen.Device
	events    = make(chan interface{}, 301)
	done      = make(chan bool)
	moribound = make(chan bool, 1)
	sigterm   = make(chan bool)
	focused   = false
)

func init() {
	// this grants the capability to shut down the program
	// it happens exactly once
	moribound <- true

	// error.go:/logFunc/
	log.SetFlags(log.Llongfile)
	log.SetPrefix("a: ")
}

func banner() {
	logf("ver=%s", Version)
	logf("pid=%d", os.Getpid())
	logf("args=%q", os.Args)
	repaint()
}

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

	g = NewGrid(dev, GridConfig)
	sp, size := image.Pt(0, 0), image.Pt(900, 900)
	g.Move(sp)
	g.Resize(size)

	for _, v := range list {
		col.Attach(g, NewCol(dev, ft, image.ZP, image.ZP, v), sp)
		sp.X += size.X / len(list)
	}
	col.Fill(g)
	g.Refresh()

	setLogFunc(g.aerr)
	banner()
	createnetworks()
	actinit(g)
	assert("actinit", g)
	go func() {
		for {
			select {
			case e := <-D.Scroll:
				activate(p(e), g)
				scroll(act, mus.ScrollEvent{Dy: 5, Event: e})
			case e := <-D.Mouse:
				activate(p(e), g)
				e = readmouse(e)
				if down == 0 {
					continue
				}
				if borderHit(rel(e, act)) {
					procBorderHit(e)
				} else {
					// assert("procButton", g) //
					procButton(rel(e, act))
				}
				repaint()
			case <-sigterm:
				logf("mouse: sigterm")
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-D.Key:
				kbdin(e, actTag, act)
				repaint()
			case <-sigterm:
				logf("kbd: sigterm")
				return
			}
		}
	}()

Loop:
	for {
		select {
		case <-sigterm:
			logf("mainselect: sigterm")
			break Loop
		case e := <-D.Size:
			winSize = image.Pt(e.WidthPx, e.HeightPx)
			g.Resize(winSize)
			repaint()
		case e := <-D.Paint:
			if throttled() {
				continue
			}
			if e.External {
				g.Resize(winSize)
			}
			g.Upload()
			wind.Publish()
		case e := <-D.Lifecycle:
			procLifeCycle(e)
			repaint()
		case e := <-events:
			switch e := e.(type) {
			case tag.GetEvent:
				t := New(actCol, e.Basedir, e.Name)
				if e.Addr != "" {
					actTag = t.(*tag.Tag)
					act = actTag.Body
					//actTag.Handle(actTag.Body, edit.MustCompile(e.Addr))
					MoveMouse(act)
				} else {
					moveMouse(t.Loc().Min)
				}
			case edit.File:
				g.acolor(e)
			case edit.Print:
				g.aout(string(e))
			case error:
				logf("unspecified error: %s", e)
			case interface{}:
				logf("missing event: %#v\n", e)
				continue
			}
			repaint()
		}
	}
}

func teardown() {
	select {
	case clean := <-moribound:
		if clean {
			setLogFunc(log.Printf)
			logf("TODO: polite shutdown")
			close(sigterm)
			close(moribound)
		}
	default:
	}
}

func procLifeCycle(e lifecycle.Event) {
	if e.To == lifecycle.StageDead {
		teardown()
		return
	}
	if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
		focused = false
	} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
		focused = true
	}
}
