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
			ed = actTag.Body.(*win.Win)
		}
		prog.Run(ed)
	case *tag.Tag:
		prog.Run(ed.Body)
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

func acmd(e event.Cmd) {
	s := string(e.P)
	switch s {
	case "Load":
		Load(g, "a.dump")
	case "Dump":
		Dump(g, g.cwd(), "gomono", "goregular")
	case "Elastic":
		t := actTag
		w, _ := t.Body.(*win.Win)
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
		actTag.Get(actTag.FileName())
		repaint()
	case "New":
		newtag := New(actCol, "", "")
		moveMouse(newtag.Loc().Min)
	case "Newcol":
		moveMouse(NewColParams(g, "").Loc().Min)
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

func cmdexec(input text.Editor, dir string, argv string) {
	x := strings.Fields(argv)
	if len(x) == 0 {
		logf("|: nothing on rhs")
		return
	}
	n := x[0]
	var a []string
	if len(x) > 1 {
		a = x[1:]
	}

	cmd := exec.Command(n, a...)
	cmd.Dir = dir
	var fd0 io.WriteCloser
	fd1, _ := cmd.StdoutPipe()
	fd2, _ := cmd.StderrPipe()
	fd0, _ = cmd.StdinPipe()

	fd0.Close()
	var wg sync.WaitGroup
	donec := make(chan bool)
	outc := make(chan []byte)
	errc := make(chan []byte)
	for _, fd := range []io.ReadCloser{fd1, fd2} {
		fd := fd
		wg.Add(1)
		go func() {
			defer wg.Done()
			var b [65536]byte
			for {
				select {
				case <-donec:
					return
				default:
					n, err := fd.Read(b[:])
					if n > 0 {
						outc <- append([]byte{}, b[:n]...)
					}
					if err != nil {
						if err != io.EOF {
							eprint(err)
						}
						return
					}
				}
			}

		}()
	}

	err := cmd.Start()
	if err != nil {
		logf("exec: %s: %s", argv, err)
		close(donec)
		return
	}

	var (
		q0, q1 int64
		f      text.Editor
	)
	lazyinit := func() {
		to := g.afinderr(dir, cmdlabel(n, dir))
		f = to.Body
		q0, q1 := f.Dot()
		f.Delete(q0, q1)
		q1 = q0
	}

	go func() {
		stdin := io.Reader(null)
		if input != nil {
			stdin = bytes.NewReader(append([]byte{}, input.Bytes()[q0:q1]...))
		}
		if _, err := io.Copy(fd0, stdin); err != nil {
			eprint(err)
			return
		}
		cmd.Wait()
		close(donec)
	}()
	go func() {
		select {
		case p := <-outc:
			lazyinit()
			f.Insert(p, q1)
			q1 += int64(len(p))
		case p := <-errc:
			lazyinit()
			f.Insert(p, q1)
			q1 += int64(len(p))
		case <-donec:
			return
		}
		repaint()
		for {
			select {
			case p := <-outc:
				f.Insert(p, q1)
				q1 += int64(len(p))
			case p := <-errc:
				f.Insert(p, q1)
				q1 += int64(len(p))
			case <-donec:
				return
			}
			repaint()
		}
	}()
}

func cmdlabel(name, dir string) (label string) {
	return fmt.Sprintf("%s%c-%s", dir, filepath.Separator, name)

}
