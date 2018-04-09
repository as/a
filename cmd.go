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
)

func acmd(e event.Cmd) {
	g.aerr("cmd: %#v\n", e)
	s := string(e.P)
	switch s {
	case "Put", "Get":
		actTag.Handle(act, s)
		repaint()
	case "New":
		newtag := New(actCol, "", "")
		logf("%v\n", newtag.Loc())
		moveMouse(newtag.Loc().Min)
	case "Newcol":
		moveMouse(NewCol2(g, "").Loc().Min)
	case "Del":
		logf("Del -> %#v\n", e)
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
			dirty = true
		}
	}

}

func cmdexec(f text.Editor, dir string, argv string) {
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
			repaint()
		}
	}()

}
