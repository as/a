package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/as/edit"
	"github.com/as/event"
	"github.com/as/path"
	"github.com/as/text"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

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
			break
			prog, err := edit.Compile(s, &edit.Options{Sender: nil, Origin: abs})
			if err != nil {
				logf(err.Error())
				return
			}
			ed := text.Editor(e.To[0])
			if e.To[0] == actTag.Win {
				ed = actTag.Body
			}
			prog.Run(ed)
			ajump2(ed, false)
		} else if strings.HasPrefix(s, "Install ") {
			s = s[8:]
			g.Install(actTag, s)
		} else {
			x := strings.Fields(s)
			if len(x) < 1 {
				logf("empty command")
				return
			}
			tagname := fmt.Sprintf("%s%c-%s", path.DirOf(abs), filepath.Separator, x[0])
			to := g.afinderr(path.DirOf(abs), tagname)
			cmdexec(to.Body, path.DirOf(abs), s)
			//			setdirty()
		}
	}
}

func cmdexec(f text.Editor, dir string, argv string) {
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
	q0, q1 := f.Dot()
	f.Delete(q0, q1)
	q1 = q0
	var fd0 io.WriteCloser
	fd1, _ := cmd.StdoutPipe()
	fd2, _ := cmd.StderrPipe()
	fd0, _ = cmd.StdinPipe()

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

	go func() {
		defer wg.Done()
		b := make([]byte, 65536)
		for {
			select {
			case <-donec:
				return
			default:
				n, err := fd2.Read(b)
				if n > 0 {
					errc <- append([]byte{}, b[:n]...)
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
	cmd.Start()
	go func() {
		_, err := io.Copy(fd0, bytes.NewReader(append([]byte{}, f.Bytes()[q0:q1]...)))
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
			repaint()
		}
	}()

}
