package main

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/as/text"
)

func chansend(ctx context.Context, dst chan []byte, src ...io.ReadCloser) {
	done := ctx.Done()
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, fd := range src {
		wg.Add(1)
		go func(fd io.ReadCloser) {
			defer wg.Done()
			var b [65536]byte
			for {
				select {
				case <-done:
					return
				default:
					n, err := fd.Read(b[:])
					if n > 0 {
						dst <- append([]byte{}, b[:n]...)
					}
					if err != nil {
						if err != io.EOF {
							eprint(err)
						}
						return
					}
				}
			}
		}(fd)
	}

	wg.Wait()
}

func chanrecv(ctx context.Context, src chan []byte, dst ...io.Writer) {
	fd := io.MultiWriter(dst...)
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case b := <-src:
			if _, err := fd.Write(b); err != nil {
				return
			}
		}
	}
}

func parsecmd(s string) (name string, args []string) {
	if s == "" {
		return "", nil
	}
	a := strings.Fields(s)
	return a[0], a[len(a):]
}

func cmdexec(ctx context.Context, input text.Editor, dir string, args ...string) {
	if len(args) == 0{
		return
	}
	name := args[0]
	args = args[1:]

	done := ctx.Done()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	fd1, _ := cmd.StdoutPipe()
	fd2, _ := cmd.StderrPipe()
	fd0, _ := cmd.StdinPipe()

	dst := make(chan []byte)
	src := make(chan []byte)
	go chansend(ctx, dst, fd1, fd2)
	go chanrecv(ctx, src, fd0)

	err := cmd.Start()
	if err != nil {
		logf("exec: %s: %s", args, err)
		return
	}

	var f text.Editor
	lazyinit := func() {
		f = g.afinderr(dir, cmdlabel(name, dir))
	}

	go func() {
		if input != nil {
			q0, q1 := input.Dot()
			src <- append([]byte{}, input.Bytes()[q0:q1]...)
		}
		close(src)
		cmd.Wait()
	}()

	go func() {
		for {
			select {
			case p := <-dst:
				lazyinit()
				f.Write(p)
				repaint()
			case <-done:
				return
			}
		}
	}()
}

func newOSCmd(dir, argv string) (name string, c Cmd) {
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

type OSCmd struct {
	*exec.Cmd
}

func (o *OSCmd) Arg() []string { return o.Cmd.Args }
func (o *OSCmd) Fd(n int) (io.ReadWriter, error) {
	if n < 0 {
		return nil, ErrBadFD
	}
	switch n {
	case 0:
		return ro{o.Stdin}, nil
	case 1:
		return wo{o.Stdout}, nil
	case 2:
		return wo{o.Stderr}, nil
	}
	return nil, ErrNoFD
}
func (o *OSCmd) Env() []string {
	return o.Cmd.Env
}

func (o *OSCmd) Redir(n int, rw io.ReadWriter) error {
	if n < 0 {
		panic("negative")
	}
	if n-3 >= len(o.ExtraFiles) {
		panic("too big")
	}
	switch n {
	case 0:
		o.Stdin = rw
	case 1:
		o.Stdout = rw
	case 2:
		o.Stderr = rw
	}
	return nil
}

type wo struct {
	io.Writer
}
type ro struct {
	io.Reader
}

func (w wo) Read(p []byte) (int, error)  { return 0, nil }
func (r ro) Write(p []byte) (int, error) { return 0, nil }
