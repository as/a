package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/as/edit"
	"github.com/as/event"
	"github.com/as/frame"
	"github.com/as/path"
	"github.com/as/shiny/event/lifecycle"
	"github.com/as/text"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

var (
	ErrBadFD = errors.New("bad file descriptor")
	ErrNoFD  = errors.New("no fd")
)

type Cmd interface {
	Arg() []string
	Fd(int) (io.ReadWriter, error)
	Env() []string

	Start() error
	Wait() error
	Redir(fd int, src io.ReadWriter) error
}

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
		if ed == actTag.Label {
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

var (
	// ndel is a delete counter. writers monitor its status to avoid
	// writing to deleted windows
	ndel uint32
)

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
	case "Undo", "Redo":
		do(s == "Undo")
	case "Del":
		w := Del(actCol, actCol.ID(actTag))
		atomic.AddUint32(&ndel, +1)
		w.Close()
	case "Sort":
		logf("Sort: TODO")
	case "Delcol":
		ws := Delcol(g, g.ID(actCol))
		atomic.AddUint32(&ndel, 1)
		for _, w := range ws {
			if w, _ := w.(io.Closer); w != nil {
				w.Close()
			}
		}
	case "Exit":
		D.Lifecycle <- lifecycle.Event{To: lifecycle.StageDead}
		return
	case "Diff":
		Diff(actTag)
	default:
		if len(e.To) == 0 {
			logf("cmd has no destination: %q", s)
		}
		abs := AbsOf(e.Basedir, e.Name)
		p := prefixer{string: strings.TrimSpace(s)}
		switch {
		case p.Prefix("|"):
			editcmd(e.To[0], abs, s)
		case p.Prefix("<"):
			cmdexec(context.Background(), actTag, actTag, path.DirOf(abs), strings.Fields(p.Chop())...)
		case p.Prefix(">"):
			cmdexec(context.Background(), nil, actTag, path.DirOf(abs), strings.Fields(p.Chop())...)
		case p.Prefix("Edit"):
			editcmd(e.To[0], abs, p.Chop())
		case p.Prefix("Look"):
			// Determine where this command is executed. If it's at the grid or column
			// labels, we run Look on all the text.Editors below that point. A column
			// runs Look on all planes under the column. A Grid does the same for
			// all columns.
			//
			// Use the track data structure to determine what Editor was last selected
			// and use that as the search parameter. Seperate the two cases so the
			// new feature is less likely to blow up the process if it has bugs--it's
			// less likely to be used anyway
			if k, from := KindOf(e.From); k.List() {
				data := track.win.Rdsel()
				VisitAll(from, func(p Named) {
					switch ed := p.(type) {
					case nil:
						return
					case text.Editor:
						q0, q1 := ed.Dot()
						lookliteraltag(ed, q0, q1, data)
					}
				})
			} else {
				switch ed := actTag.Window.(type) {
				case *win.Win:
					q0, q1 := ed.Dot()
					lookliteraltag(ed, q0, q1, ed.Rdsel())
				}
			}
		case p.Prefix("Install"):
			g.Install(actTag, p.Chop())
		default:
			cmdexec(context.Background(), nil, nil, path.DirOf(abs), strings.Fields(s)...)
		}
		reload(e.To[0])
	}
}

type prefixer struct {
	string
	n int
}

func (p *prefixer) Prefix(pre string) bool {
	p.n = 0
	if !strings.HasPrefix(p.string, pre) {
		return false
	}
	p.n = len(pre)
	return true
}
func (p *prefixer) Chop() string {
	if p.n != 0 {
		return p.string[p.n:]
	}
	return p.string
}

func cmdlabel(name, dir string) (label string) {
	return fmt.Sprintf("%s%c-%s", dir, filepath.Separator, name)
}
