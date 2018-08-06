package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/as/edit"
	"github.com/as/event"
	"github.com/as/frame"
	"github.com/as/path"
	"github.com/as/text"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

var null, _ = os.Open(os.DevNull)

func editcmd(ed interface{}, origin, cmd string) {
	prog, err := edit.Compile(cmd, &edit.Options{Sender: nil, Origin: origin})
	if err != nil {
		logf("editcmd: %s", err)
		return
	}
	runeditcmd(prog, ed)
	{
		ed, _ := ed.(text.Editor)
		if ed != nil {
			ajump2(ed, false)
		}
	}
}

func runeditcmd(prog *edit.Command, ed interface{}) {
	switch ed := ed.(type) {
	case *win.Win:
		if ed == actTag.Win {
			ed = actTag.Window.(*win.Win)
		}
		prog.Run(ed)
	case *tag.Tag:
		prog.Run(ed.Window)
	case *Grid:
		for _, ed := range ed.List {
			runeditcmd(prog, ed)
		}
	case *Col:
		for _, ed := range ed.List {
			runeditcmd(prog, ed)
		}
	case text.Editor:
		prog.Run(ed)
	case interface{}:
		logf("dont know what %T is", ed)
	}
}

func getcmd(t *tag.Tag) {
	// Add rendering here if image?
	if *images && tryImage(t.FileName()) {
		render(t)
	} else {
		t.Get(t.FileName())
	}
}

func acmd(e event.Cmd) {
	s := string(e.P)
	switch s {
	case "Img":
		renderimage(actTag)
		repaint()
	case "Load":
		Load(g, "a.dump")
	case "Dump":
		Dump(g, g.cwd(), "gomono", "goregular")
	case "Elastic":
		t := actTag
		w, _ := t.Window.(*win.Win)
		if w != nil && w.Frame != nil {
			cf := &t.Config.Body.Frame
			if cf.Flag&frame.FrElastic == 0 {
				cf.Flag |= frame.FrElastic
			} else {
				cf.Flag &^= frame.FrElastic
			}
			cf.Flag |= frame.FrElastic
			w.Frame.SetFlags(cf.Flag)
			w.Resize(w.Size())
		}
		repaint()
	case "Font":
		if actTag == g.Tag {
			nextFace(g)
		} else if actTag == actCol.Tag {
			nextFace(actCol)
		} else {
			nextFace(actTag)
		}
	case "Put":
		actTag.Put()
		repaint()
	case "Get":
		getcmd(actTag)

		//repaint()
	case "New":
		newtag := New(actCol, "", "")
		moveMouse(newtag.Bounds().Min)
	case "Newcol":
		moveMouse(NewColParams(g, "").Bounds().Min)
	case "Del":
		Del(actCol, actCol.ID(actTag))
	case "Sort":
		logf("Sort: TODO")
	case "Delcol":
		Delcol(g, g.ID(actCol))
	case "Exit":
		logf("Exit: TODO")
	default:
		if len(e.To) == 0 {
			logf("cmd has no destination: %q", s)
		}
		abs := AbsOf(e.Basedir, e.Name)
		if strings.HasPrefix(s, "Edit ") {
			s = s[5:]
			editcmd(e.To[0], abs, s)
			editRefresh(e.To[0])
		} else if strings.HasPrefix(s, "Install ") {
			s = s[8:]
			g.Install(actTag, s)
		} else {
			x := strings.Fields(s)
			if len(x) < 1 {
				logf("empty command")
				return
			}
			cmdexec(nil, path.DirOf(abs), s)
		}
	}
}

type Funnel struct {
	sync.Mutex
	io.Writer
}

func (f *Funnel) Read(p []byte) (n int, err error) {
	return
}
func (f *Funnel) Write(p []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	return f.Writer.Write(p)
}
func (f *Funnel) Unlock() {
	f.Mutex.Unlock()
	repaint()
}

func oscmd(dir, argv string) (name string, c Cmd) {
	x := strings.Fields(argv)
	if len(x) == 0 {
		logf("|: nothing on rhs")
		return "", nil
	}
	n := x[0]
	var a []string
	if len(x) > 1 {
		a = x[1:]
	}
	oc := &OSCmd{
		Cmd: exec.Command(n, a...),
	}
	oc.Dir = dir
	return n, oc
}

func cmdexec(src text.Editor, dir string, argv string) {
	input := []byte{}
	if src != nil {
		q0, q1 := src.Dot()
		input = append([]byte{}, src.Bytes()[q0:q1]...)
	}

	n, cmd := oscmd(dir, argv)

	dst := g.afinderr(dir, cmdlabel(n, dir))
	dst.Delete(dst.Dot())

	fun := &Funnel{Writer: dst}
	cmd.Redir(0, bytes.NewBuffer(input))
	cmd.Redir(1, fun)
	cmd.Redir(2, fun)

	err := cmd.Start()
	if err != nil {
		logf("exec: %s: %s", argv, err)
		return
	}
}

func cmdlabel(name, dir string) (label string) {
	return fmt.Sprintf("%s%c-%s", dir, filepath.Separator, name)
}
