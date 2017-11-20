package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	//	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/as/event"
	//	"github.com/as/font/vga"
	mus "github.com/as/text/mouse"
	"golang.org/x/exp/shiny/screen"
	//	"golang.org/x/image/font/plan9font"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"github.com/as/cursor"
	"github.com/as/edit"
	"github.com/as/frame"
	"github.com/as/frame/font"
	window "github.com/as/ms/win"
	"github.com/as/path"
	"github.com/as/text"
	"github.com/as/ui"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

var (
	Version = "0.4.3"
	xx      Cursor
	eprint  = fmt.Println
	timefmt = "2006.01.02 15.04.05"
)

var (
	winSize   = image.Pt(1024, 768)
	fsize     = 12 // Put
	pad       = image.Pt(15, 15)
	tagHeight = fsize*2 + fsize/2 - 2
	scrollX   = 10
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

var cols = frame.Mono

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

var (
	utf8    = flag.Bool("u", false, "enable utf8 experiment")
	elastic = flag.Bool("elastic", false, "enable elastic tabstops")
	oled    = flag.Bool("b", false, "OLED display mode (black)")
	ftsize  = flag.Int("ftsize", 16, "font size")
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
	frame.A.Text = image.NewUniform(color.RGBA{192, 192, 232, 255})
	frame.ATag1.Back, frame.ATag1.Text = frame.ATag1.Text, frame.ATag1.Back
	frame.ATag1.Text = frame.A.Text
	frame.ATag0.Back, frame.ATag0.Text = frame.ATag0.Text, frame.ATag0.Back
	frame.ATag0.Text = frame.A.Text
	frame.ATag0.Back = image.Black
	tag.Gray = image.NewUniform(color.RGBA{192, 192, 232, 255})
	tag.LtGray = image.NewUniform(color.RGBA{192, 192, 232, 255})
	tag.X = image.NewUniform(color.RGBA{192, 192, 232, 255})
	frame.A.Back = image.Black
}

var dirty bool

func ck() {
	if dirty || (act != nil && act.Dirty()) {
		act.Window().Send(paint.Event{})
	}
}

var g *Grid

// Put
func main() {
	flag.Parse()
	defer trypprof()()
	frame.ForceUTF8 = *utf8
	frame.ForceElastic = *elastic

	if *oled {
		black()
	}

	list := argparse()
	dev, err := ui.Init(&screen.NewWindowOptions{Width: winSize.X, Height: winSize.Y, Title: "A"})
	if err != nil {
		log.Fatalln(err)
	}
	wind := dev.Window()

	// Linux will segfault here if X is not present
	wind.Send(paint.Event{})
	ft := font.NewGoMedium(fsize)
	g = NewGrid(dev, image.ZP, winSize, ft, list...)

	// This in particular needs to go
	actCol = g.List[1].(*Col)
	actTag = actCol.List[1].(*tag.Tag)
	act = actTag.Body

	var pt image.Point
	r := act.Bounds()
	mousein := mus.NewMouse(time.Second/3, wind)
	mousein.Machine.SetRect(image.Rect(r.Min.X, r.Min.Y+pad.Y, r.Max.X, r.Max.Y-pad.Y))

	//
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
	// Returns the bounding box of the invisible sizer you can use
	// to draw windows and columns around.
	sizerOf := func(p Plane) image.Rectangle {
		r := p.Loc()
		r.Max = r.Min.Add(image.Pt(scrollX, tagHeight))
		return r
	}
	sizerHit := func(p Plane, pt image.Point) bool {
		in := pt.In(sizerOf(p))
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
		switch e.Button {
		case 3:
			actCol.RollUp(id, dy)
			//actCol.MoveWin(id, dy)
		case 2:
			dy -= fsize * 2
			actCol.MoveWin(id, dy)
		case 1:
			actCol.Grow(id, actCol.bestGrowth(id, tagHeight))
		}
		moveMouse(actTag.Loc().Min)
	}
	tophit := func() bool {
		return pt.Y > g.sp.Y+g.tdy && pt.Y < g.sp.Y+g.tdy*2
	}

	ajump := func(ed text.Editor, cursor bool) {
		fn := moveMouse
		if cursor == false {
			fn = nil
		}
		if ed, ok := ed.(text.Jumper); ok {
			ed.Jump(fn)
		}
	}
	ismeta := func(ed Plane) bool {
		return ed == g.List[0].(*tag.Tag).Body
	}
	ismeta = ismeta
	alook := func(e event.Look) {
		g.Look(e)
	}
	aerr := g.aerr
	var (
		scrollbar = 1
		sizer     = 2
		window    = 4
		context   = 0
	)

	aerr("ver=%s", Version)
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
						detachwin()
					}
					context = sizer
				} else {
					growshrink(e.Event)
				}

			} else if x := int(e.X); x >= 0 && x < 10 {
				context = scrollbar
				act.Clicksb(p(e.Event), int(e.Button))
			} else {
				actTag.Handle(act, e)
				context = window
			}
			ck()
		case mus.ScrollEvent:
			doScrollEvent(act, e)
		case mus.SweepEvent:
			switch context {
			case scrollbar:
				act.Clicksb(p(e.Event), 0)
			case sizer:
			case window:
				// A very effective optimization eliminates several function
				// calls completely if the cursor isn't moving and the use is still
				// in the window's bounds.
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
			context = 0
			ck()
		case mus.SnarfEvent, mus.InsertEvent:
			actTag.Handle(act, e)
			ck()
		case mus.SelectEvent:
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
			dirty = true
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
				aerr("Delcol: TODO")
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
					prog, err := edit.Compile(s, &edit.Options{Sender: wind, Origin: abs})
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
		case size.Event:
			winSize = image.Pt(e.WidthPx, e.HeightPx)
			g.Resize(winSize)
			ck()
		case paint.Event:
			if !focused {
				g.Resize(winSize)
			}
			g.Upload(wind)
			wind.Publish()
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				return
			}
			// NT doesn't repaint the window if another window covers it
			if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
				focused = false
				wind.Send(paint.Event{})
			} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
				focused = true
			}
		case error:
			aerr(e.Error())
		case interface{}:
			log.Printf("missing event: %#v\n", e)
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
