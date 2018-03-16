package main // import "github.com/as/a"

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/as/event"
	"github.com/as/shiny/screen"
	mus "github.com/as/text/mouse"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/time/rate"

	"github.com/as/edit"
	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/path"
	"github.com/as/text"
	"github.com/as/ui"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

var (
	Version = "0.5.3"
	xx      Cursor
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

/*
func vgaface() font.Face {
	data, err := ioutil.ReadFile("u_vga16.font")
	if err != nil {
		panic(err)
	}
	face, err := plan9font.ParseFont(data, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}
	return face
}
*/
// TODO(as): refactor frame so this stuff doesn't have to exist here
func black() {
	frame.A.Text = frame.MTextW
	frame.A.Back = frame.MBodyW

	frame.ATag0.Text = frame.MTextW
	frame.ATag0.Back = frame.MTagG

	frame.ATag1.Text = frame.MTextW
	frame.ATag1.Back = frame.MTagC
}

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
	frame.ForceUTF8 = *utf8
	frame.ForceElastic = *elastic

	// Startup
	err := createnetworks()
	if err != nil {
		log.Fatalln(err)
	}

	lim := rate.NewLimiter(rate.Every(time.Second/120), 2)

	if *oled {
		black()
	}

	list := argparse()
	if *quiet {
		<-done
		os.Exit(0)
	}

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

	// This in particular needs to go
	actCol = g.List[1].(*Col)
	actTag = actCol.List[1].(*tag.Tag)
	act = actTag.Body

	var pt image.Point
	r := act.Bounds()

	// Temporary just until col and tag can be seperated into their own packages. The
	// majority of these closures will disappear as the program becomes more stable
	sweepend := func() {
		if xx.sweepCol {
			xx.sweepCol = false
			g.fill()
			g.Attach(xx.srcCol, pt.X)
			moveMouse(xx.srcCol.Loc().Min)
		} else {
			xx.sweep = false
			xx.srcCol.fill()
			if xx.src == nil {
				return
			}
			actCol.Attach(xx.src, pt.Y)
			moveMouse(xx.src.Loc().Min)
		}
	}
	markwin := func() {
		xx.srcCol = actCol
		xx.src = actTag
	}
	detachcol := func() {
		xx.srcCol = actCol
		xx.sweepCol = true
		xx.src = nil
		g.detach(g.ID(xx.srcCol))
		g.fill()
	}
	detachwin := func() {
		markwin()
		xx.srcCol.detach(xx.srcCol.ID(xx.src))
	}
	growshrink := func(e mouse.Event) {
		dy := r.Min.Y
		id := actCol.ID(actTag)
		switch e.Button {
		case 3:
			actCol.RollUp(id, dy)
			//actCol.MoveWin(id, dy)
		case 2:
			dy -= *ftsize * 2
			actCol.MoveWin(id, dy)
		case 1:
			actCol.Grow(id, actCol.bestGrowth(id, tagHeight))
		}
		moveMouse(actTag.Loc().Min)
	}
	ajump := func(ed text.Editor, cursor bool) {
		fn := moveMouse
		if !cursor {
			fn = nil
		}
		if ed, ok := ed.(text.Jumper); ok {
			ed.Jump(fn)
		}
	}
	alook := func(e event.Look) {
		g.Look(e)
	}
	aerr := g.aerr
	var (
		scrollbar = 1
		sizer     = 2
		window    = 4
		cont      = 0
	)

	aerr("ver=%s", Version)
	aerr("pid=%d", os.Getpid())
	aerr("args=%q", os.Args)
	if srv != nil {
		aerr("listening for remote connections")
	}
	if client != nil {
		aerr("connected to remote filesystem")
	}

	var down uint

	readmouse := func(e mouse.Event) mouse.Event {
		if e.Button != 0 {
			switch e.Direction {
			case 1:
				down |= 1 << uint(e.Button)
			case 2:
				down ^= 1 << uint(e.Button)
			}
		}
		activate(p(e), g)
		pt = p(e).Add(act.Loc().Min)
		e.X -= float32(act.Sp.X)
		e.Y -= float32(act.Sp.Y)
		return e
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
				for down == start {
					act.Sq, q0, q1 = sweep(act, e, act.Sq, q0, q1) //TODO (nil was act)
					act.Select(q0, q1)
					ck()

					select {
					case e = <-D.Mouse:
						e = readmouse(e)
					}
				}

				switch start {
				case 1 << 1:
					if inSizer(p(e)) {
						aerr("InSizer: %s", p(e))
					} else if inScroll(p(e)) {
						aerr("inScroll: %s", p(e))
					} else {
						aerr("sweepOrClock: %s", p(e))
						// sweep or click
						//					c0 = time.Now().Add(time.Second/3)
						act.Select(q0, q1)
						ck()
					}
				case 1 << 2:

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
						//						Basedir: t.basedir,
						Name: t.FileName(),
					}
				case 1 << 3:
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
						//						Basedir: t.basedir,
						Name: t.FileName(),
					}
				}

				//Drel <- e
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-D.Key:
				actTag.Handle(act, e)
				dirty = true
				repaint()
			}
		}
	}()

	for {
		select {
		case e := <-D.Lifecycle:
			if e.To == lifecycle.StageDead {
				return
			}
			// NT doesn't repaint the window if another window covers it
			if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
				focused = false
				ck()
			} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
				focused = true
			}
		case e := <-D.Paint:
			if !lim.Allow() {
				continue
			}
			if e.External {
				println("fffffffff")
				g.Resize(winSize)
			}
			g.Upload(wind)
			wind.Publish()
		case e := <-D.Size:
			winSize = image.Pt(e.WidthPx, e.HeightPx)
			g.Resize(winSize)
			ck()
		case e := <-events:
			//		e := wind.NextEvent()
			switch e := e.(type) {
			case mus.DrainStop:
			case mus.Drain:
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
			case mus.MarkEvent:
				cont = 0
				pt = p(e.Event).Add(act.Loc().Min)
				if inSizer(p(e.Event)) {
					switch e.Button {
					case 1:
						if canopy(pt) {
							detachcol()
						} else {
							detachwin()
						}
						cont = sizer
					default:
						growshrink(e.Event)

					}
				} else if p := p(e.Event); inScroll(p) {
					cont = scrollbar
					act.Clicksb(p, int(e.Button))
				} else {
					actTag.Handle(act, p)
					cont = window
				}
				ck()
			case mus.ScrollEvent:
				//			doScrollEvent(act, e)
			case mus.SweepEvent:
				switch cont {
				case scrollbar:
					act.Clicksb(p(e.Event), 0)
				case sizer:
				case window:
					if !e.Motion() && act != nil {
						r := act.Frame.Bounds()
						if p(e.Event).In(r) {
							continue
						}
					}
					actTag.Handle(act, e)
				}
				ck()
			case mus.CommitEvent:
				cont = 0
				ck()
			case mus.SnarfEvent, mus.InsertEvent:
				actTag.Handle(act, e)
				ck()
			case mus.SelectEvent:
				switch cont {
				case scrollbar:
					act.Clicksb(p(e.Event), 0)
				case sizer:
					e.X += float32(act.Sp.X)
					e.Y += float32(act.Sp.Y)
					pt.X = int(e.X)
					pt.Y = int(e.Y)
					activate(p(e.Event), g)
					sweepend()
				case window:
					actTag.Handle(act, e)
					//q0, q1 := act.Dot()
					if e.Button == 2 {
						//	s := string(act.Bytes()[q0:q1])
						//	actTag.Handle(act, s)
						//	wind.Send(s)
					}
				}
				cont = 0
				ck()
			case event.Look:
				alook(e)
				ck()
			case event.Cmd:
				s := string(e.P)
				switch s {
				case "Put", "Get":
					actTag.Handle(act, s)
					ck()
				case "New":
					moveMouse(New(actCol, "", "").Loc().Min)
				case "Newcol":
					moveMouse(NewCol2(g, "").Loc().Min)
				case "Del":
					Del(actCol, actCol.ID(actTag))
				case "Sort":
					aerr("Sort: TODO")
				case "Delcol":
					Delcol(g, g.ID(actCol))
				case "Exit":
					aerr("Exit: TODO")
				default:
					if len(e.To) == 0 {
						aerr("cmd has no destination: %q", s)
					}
					abs := AbsOf(e.Basedir, e.Name)
					if strings.HasPrefix(s, "Edit ") {
						s = s[5:]
						// The event sink shouldn't be specified during
						// compile time, but its the easiest way to
						// see it works correctly with the editor
						prog, err := edit.Compile(s, &edit.Options{Sender: nil, Origin: abs})
						if err != nil {
							aerr(err.Error())
							continue
						}
						prog.Run(e.To[0])
						w := e.To[0].(*win.Win)
						w.Resize(w.Size())
						//e.To[0].(*win.Win).Refresh()
						ajump(e.To[0], false)
					} else if strings.HasPrefix(s, "Install ") {
						s = s[8:]
						g.Install(actTag, s)
					} else {
						x := strings.Fields(s)
						if len(x) < 1 {
							aerr("empty command")
							continue
						}
						tagname := fmt.Sprintf("%s%c-%s", path.DirOf(abs), filepath.Separator, x[0])
						to := g.afinderr(path.DirOf(abs), tagname)
						cmd(to.Body, path.DirOf(abs), s)
						dirty = true
					}
				}
				ck()
			case edit.File:
				g.acolor(e)
			case edit.Print:
				g.aout(string(e))
			case error:
				aerr(e.Error())
			case interface{}:
				log.Printf("missing event: %#v\n", e)
			}
		}
	}

}

func cmd(f text.Editor, dir string, argv string) {
	x := strings.Fields(argv)
	if len(x) == 0 {
		eprint("|: nothing on rhs")
		return
	}
	n := x[0]
	var a []string
	if len(x) > 1 {
		a = x[1:]
	}

	cmd := exec.Command(n, a...)
	cmd.Dir = dir
	q0, q1 := f.Dot()
	f.Delete(q0, q1)
	q1 = q0
	var fd0 io.WriteCloser
	fd1, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	fd2, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	fd0, err = cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	fd0.Close()
	var wg sync.WaitGroup
	donec := make(chan bool)
	outc := make(chan []byte)
	errc := make(chan []byte)
	wg.Add(2)
	go func() {
		defer wg.Done()
		b := make([]byte, 65536)
		for {
			select {
			case <-donec:
				return
			default:
				n, err := fd1.Read(b)
				if err != nil {
					if err == io.EOF {
						break
					}
					eprint(err)
				}
				outc <- append([]byte{}, b[:n]...)
			}
		}
	}()

	go func() {
		defer wg.Done()
		b := make([]byte, 65536)
		for {
			select {
			case <-donec:
				return
			default:
				n, err := fd2.Read(b)
				if err != nil {
					if err == io.EOF {
						break
					}
				}
				errc <- append([]byte{}, b[:n]...)
			}
		}
	}()
	cmd.Start()
	go func() {
		_, err = io.Copy(fd0, bytes.NewReader(append([]byte{}, f.Bytes()[q0:q1]...)))
		if err != nil {
			eprint(err)
			return
		}
		cmd.Wait()
		close(donec)
	}()
	go func() {
	Loop:
		for {
			select {
			case p := <-outc:
				f.Insert(p, q1)
				q1 += int64(len(p))
			case p := <-errc:
				f.Insert(p, q1)
				q1 += int64(len(p))
			case <-donec:
				break Loop
			}
		}
	}()

}
