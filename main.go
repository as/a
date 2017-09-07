package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/as/event"
	"github.com/as/text/action"
	"github.com/as/text/find"
	mus "github.com/as/text/mouse"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"github.com/as/cursor"
	"github.com/as/edit"
	"github.com/as/frame"
	"github.com/as/frame/font"
	"github.com/as/frame/tag"
	window "github.com/as/ms/win"
	"github.com/as/text"
)

var xx Cursor
var eprint = fmt.Println

// Put
var (
	winSize = image.Pt(1900, 1000)
	pad     = image.Pt(15, 15)
	fsize   = 11
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

var cols = frame.A

func Tagtext(s string, w Plane) {
	switch w := w.(type) {
	case *Grid:
		Tagtext(s, w.List[0])
	case *Col:
		Tagtext(s, w.List[0])
	case *tag.Tag:
		t := w.Win
		q0, q1 := t.Dot()
		t.Delete(q0, q1)
		q1 = q0
		t.InsertString(s, q1)
		t.Select(q1, q1+int64(len(s)))
	case Plane:
	}
}

type CmdEvent struct {
	grid *Grid
	col  *Col
	tag  *tag.Tag
	act  Plane
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func p(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

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
func moveMouse(pt image.Point) {
	cursor.MoveTo(window.ClientAbs().Min.Add(pt))
}

// Put
func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	list := argparse()
	driver.Main(func(src screen.Screen) {
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "A"})
		wind.Send(paint.Event{})
		ft := font.NewGoMono(fsize)

		g := NewGrid(src, wind, ft, image.ZP, image.Pt(winSize.X, winSize.Y), list...)
		actCol = g.List[1].(*Col)
		actTag = actCol.List[1].(*tag.Tag)
		act = actTag.Body

		go func() {
			sc := bufio.NewScanner(bufio.NewReader(os.Stdin))
			for sc.Scan() {
				if x := sc.Text(); x == "u" || x == "r" {
					act.SendFirst(x)
					continue
				}
				act.SendFirst(edit.MustCompile(sc.Text()))
			}
		}()

		var pt image.Point
		r := act.Bounds()
		mousein := mus.NewMouse(time.Second/3, wind)
		mousein.Machine.SetRect(image.Rect(r.Min.X, r.Min.Y+pad.Y, r.Max.X, r.Max.Y-pad.Y))
		var dirty bool
		ck := func() {
			if dirty || (act != nil && act.Dirty()) {
				wind.Send(paint.Event{})
			}
		}
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
		sizerOf := func(p Plane) image.Rectangle {
			r := p.Loc()
			r.Max = r.Min.Add(image.Pt(11, 11))
			return r
		}
		sizerHit := func(p Plane, pt image.Point) bool {
			in := pt.In(sizerOf(p))
			fmt.Printf("win=%s sizer=%s pt=%s in=%v", p.Loc(), sizerOf(p), pt, in)
			return in
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
		}
		tophit := func() bool {
			return pt.Y > g.sp.Y+g.tdy && pt.Y < g.sp.Y+g.tdy*2
		}
		timefmt := "2006.01.02 15.04.05"
		aerr := func(fm string, i ...interface{}) {
			t := g.FindName("+Errors")
			if t == nil {
				t = New(actCol, "+Errors").(*tag.Tag)
				if t == nil {
					panic("cant create aerr window")
				}
				moveMouse(t.Loc().Min)
			}
			q1 := t.Body.Len()
			t.Body.Select(q1, q1)
			n := int64(t.Body.Insert([]byte(time.Now().Format(timefmt)+": "+fmt.Sprintf(fm, i...)+"\n"), q1))
			t.Body.Select(q1+n, q1+n)
		}
		ajump := func(ed text.Editor, cursor bool) {
			proj, ok := ed.(text.Projector)
			if !ok {
				return
			}
			nchars := proj.IndexOf(image.Pt(9999, 9999))
			q0, q1 := ed.Dot()
			sc, ok := ed.(text.Scroller)
			if !ok {
				return
			}
			if text.Region5(q0, q1, sc.Origin(), sc.Origin()+nchars) != 0 {
				sc.SetOrigin(q0, true)
				sc.Scroll(-3)
				ck()
			}
			if cursor {
				jmp := proj.PointOf(q0 - sc.Origin())
				moveMouse(sc.(text.Plane).Bounds().Min.Add(jmp))
			}
		}
		alook := func(e event.Look) {
			// First we find out if its coming from a tag or
			// a window body. Then we find out if its an address
			// and lastly we look
			var t *tag.Tag
			fromdir := ""
			istag := false
			_, istag = e.From.(*tag.Tag)
			str := string(e.P)
			name, addr := action.SplitPath(str)
			if !filepath.IsAbs(name) {
				fromdir = e.FromFile
				if !action.IsDir(fromdir) {
					fromdir = action.Dirof(t.FileName())
				}
			}

			t2 := g.FindName(name)
			if name == "" && addr != "" {
				// Just an address with no name:
				// jump to it in the current file
				prog := edit.MustCompile(addr)
				for _, ed := range e.To {
					prog.Run(ed)
					ajump(ed, !istag)
				}
			} else if name != "" && t2 != nil {
				if addr != "" {
					prog := edit.MustCompile(addr)
					prog.Run(t2.Body)
					ajump(t2.Body, true)
				} else {
					moveMouse(t2.Bounds().Min)
				}
			} else if name := filepath.Join(fromdir, name); action.IsFile(name) || action.IsDir(name) {
				// If the path is relative, it's combined with the tag's cwd
				// jump to the window by the same name if it's already open
				t2 := g.FindName(name)
				if t2 == nil {
					t2 = New(actCol, name).(*tag.Tag)
					moveMouse(t2.Loc().Min)
				}
				t2.Body.Select(0, 0)
				if addr != "" {
					prog := edit.MustCompile(addr)
					prog.Run(t2.Body)
					ajump(t2.Body, true)
				} else {
					moveMouse(t2.Bounds().Min)
				}
				println("open file")
			} else {
				// Find the literal string in the caller
				// this is useful because Go's regexp doesn't support
				// non-utf8 runes in the regexp compiler
				for _, ed := range e.To {
					q0, q1 := find.FindNext(ed, e.P)
					ed.Select(q0, q1)
					ajump(ed, !istag)
				}
			}
		}
		var (
			scrollbar = 1
			sizer     = 2
			window    = 4
			context   = 0
		)

		aerr("pid=%d", os.Getpid())
		aerr("args=%q", os.Args)

		for {
			e := wind.NextEvent()
			switch e := e.(type) {
			case mus.Drain:
			DrainLoop:
				for {
					switch wind.NextEvent().(type) {
					case mus.DrainStop:
						break DrainLoop
					}
				}
			case tag.GetEvent:
				t := New(actCol, e.Path)
				if e.Addr != "" {
					actTag = t.(*tag.Tag)
					act = actTag.Body
					actTag.Handle(actTag.Body, edit.MustCompile(e.Addr))
					p0, _ := act.Frame.Dot()
					moveMouse(act.Loc().Min.Add(act.PointOf(p0)))
				} else {
					moveMouse(t.Loc().Min)
				}
			case mouse.Event:
				pt = p(e).Add(act.Loc().Min)
				if context == 0 {
					activate(p(e), g)
				}
				e.X -= float32(act.Sp.X)
				e.Y -= float32(act.Sp.Y)
				mousein.Sink <- e
			case mus.MarkEvent:
				context = 0
				pt = p(e.Event).Add(act.Loc().Min)
				if sizerHit(actTag, pt) {
					if e.Button == 2 {
						if tophit() {
							detachcol()
						} else {
							detachwin() //markwin()
						}
						context = sizer
					} else {
						growshrink(e.Event)
					}

				} else if x := int(e.X); x >= 0 && x < 15 {
					context = scrollbar
					act.Clicksb(p(e.Event), int(e.Button))
				} else {
					//TODO: start here
					markwin = markwin
					actTag.Handle(act, e)
					context = window
				}
				ck()
			case mus.ClickEvent:
				println("ClickEvent")
				switch context {
				case scrollbar:
					act.Clicksb(p(e.Event), 0)
				case sizer:
				case window:
					actTag.Handle(act, e)
				}
				context = 0
				ck()
			case mus.ScrollEvent:
				if e.Button == -1 {
					e.Dy = -e.Dy
				}
				actTag.Body.Scroll(e.Dy)
				ck()
			case mus.SweepEvent:
				switch context {
				case scrollbar:
					act.Clicksb(p(e.Event), 0)
				case sizer:
				case window:
					actTag.Handle(act, e)
				}
				ck()
			case mus.CommitEvent:
				//context = 0
			case mus.SnarfEvent, mus.InsertEvent:
				actTag.Handle(act, e)
				ck()
			case mus.SelectEvent:
				println("SelectEvent")
				switch context {
				case scrollbar:
					act.Clicksb(p(e.Event), 0)
					wind.SendFirst(mus.Drain{})
					wind.Send(mus.DrainStop{})
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
				context = 0
				ck()
			case key.Event:
				actTag.Handle(act, e)
				ck()
			case event.Look:
				alook(e)
				ck()
			case event.Cmd:
				s := string(e.P)
				log.Printf("string: %q\n", s)
				switch s {
				case "Put", "Get":
					actTag.Handle(act, s)
					aerr(s)
				case "New":
					moveMouse(New(actCol, "").Loc().Min)
				case "Newcol":
					moveMouse(NewCol2(g, "").Loc().Min)
				case "Del":
					Del(actCol, actCol.ID(actTag))
				default:
					if len(e.To) == 0 {
						aerr("cmd has no destination: %q", s)
					}
					if strings.HasPrefix(s, "Edit ") {
						if len(s) > 5 {
							prog := edit.MustCompile(s[5:])
							prog.Run(e.To[0])
						}
					} else {
						cmd(e.To[0], s)
					}
				}
				ck()
			case size.Event:
				winSize = image.Pt(e.WidthPx, e.HeightPx)
				g.Resize(winSize)
				ck()
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
			case interface{}:
				log.Printf("missing event: %#v\n", e)
			}
		}
	})

}

func cmd(f text.Editor, argv string) {
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
	_, err = io.Copy(fd0, bytes.NewReader(append([]byte{}, f.Bytes()[q0:q1]...)))
	if err != nil {
		eprint(err)
		return
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
	go func() {
		cmd.Start()
		cmd.Wait()
		close(donec)
	}()
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

}
