package main // import "github.com/as/a"

import (
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"

	"github.com/as/edit"
	"github.com/as/shiny/event/lifecycle"
	"github.com/as/shiny/screen"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

var (
	eprint  = fmt.Println
	timefmt = "15.04.05"
)

var (
	g         *Grid
	D         *screen.Device
	events    = make(chan interface{}, 25)
	done      = make(chan bool)
	moribound = make(chan bool, 1)
	sigterm   = make(chan bool)
	focused   = false
)

var (
	utf8     = flag.Bool("u", false, "enable utf8 experiment")
	elastic  = flag.Bool("elastic", false, "enable elastic tabstops")
	images   = flag.Bool("img", os.Getenv("img") == "auto", "render images in editor (experimental/unstable)")
	oled     = flag.Bool("b", false, "OLED display mode (black)")
	ftsize   = flag.Int("ftsize", defaultFaceSize(), "font size")
	srvaddr  = flag.String("srv", "", "(dangerous) announce and serve file system clients on given endpoint")
	load     = flag.String("l", "", "load state from a dump file in acme format")
	dialaddr = flag.String("dial", "", "dial to a remote file system on endpoint")
	quiet    = flag.Bool("q", false, "dont interact with the graphical subsystem (use with -l)")
)

func init() {
	// this grants the capability to shut down the program
	// it happens exactly once
	moribound <- true

	// error.go:/logFunc/
	log.SetFlags(log.Llongfile)
	log.SetPrefix("a: ")

	flag.Parse()
}

// showbanner is set to true only if the user ran the program without
// arguments. In any other case, it's just annoying to have this on
// see args.go:/showbanner/
var showbanner = false

func banner() {
	if !showbanner {
		return
	}
	logf("ver=%s pid=%d args=%q", Version, os.Getpid(), os.Args)
	repaint()
}

func trap() {
	err := recover()
	if err != nil {
		teardown(true)
		panic(err)
	}
}

func main() {
	defer trypprof()()
	defer trap()

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
	{
		sp, size := image.Pt(0, 0), image.Pt(900, 900)
		g.Move(sp)
		g.Resize(size)

		if *load != "" {
			Load(g, *load)
		} else {
			for _, v := range list {
				c := NewCol(dev, ft, image.ZP, image.ZP, v)
				if v == "-" {
					c = NewCol(dev, ft, image.ZP, image.ZP)
					io.Copy(New(c, "", "").(*tag.Tag), os.Stdin)
				}
				col.Attach(g, c, sp)
				sp.X += size.X / len(list)
			}
			col.Fill(g)
		}
		g.Refresh()
	}

	setLogFunc(g.aerr)
	banner()

	createnetworks()
	actinit(g)
	assert("actinit", g)
	go func() {
		defer trap()
		for {
			select {
			case e := <-D.Scroll:
				activate(p(e), g)
				scroll(act, ScrollEvent{Dy: 7, Event: e})
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
		defer trap()
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
				getcmd(t.(*tag.Tag))
				if e.Addr != "" {
					actTag = t.(*tag.Tag)
					act = actTag.Window
					//actTag.Handle(actTag.Window, edit.MustCompile(e.Addr))
					MoveMouse(act)
				} else {
					moveMouse(t.Bounds().Min)
				}
				repaint()
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

func teardown(save bool) {
	select {
	case clean := <-moribound:
		if clean {
			setLogFunc(log.Printf)
			if save {
				Dump(g, g.cwd(), "gomono", "goregular")
				println("crash: saved: use 'a -l a.dump' to restore")
			}
			close(sigterm)
			close(moribound)
		}
	default:
	}
}

func procLifeCycle(e lifecycle.Event) {
	if e.To == lifecycle.StageDead {
		teardown(false)
		return
	}
	if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
		focused = false
	} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
		focused = true
	}
}
