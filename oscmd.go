package main

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/as/text"
)

func chansend(ctx context.Context, dst chan []byte, src ...io.ReadCloser) {
	done := ctx.Done()
	var wg sync.WaitGroup

	for _, fd := range src {
		wg.Add(1)
		go func(fd io.ReadCloser) {
			defer wg.Done()
			defer fd.Close()
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

}

func chanrecv(ctx context.Context, src chan []byte, dst ...io.Writer) {
	fd := io.MultiWriter(dst...)
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case b, ok := <-src:
			if !ok {
				return
			}
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

func cmdexec(ctx context.Context, dst, src text.Editor, dir string, args ...string) (fin chan error) {
	fin = make(chan error, 1)
	if len(args) == 0 {
		close(fin)
		return fin
	}
	name := args[0]
	args = args[1:]

	done := ctx.Done()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	fd1, _ := cmd.StdoutPipe()
	fd2, _ := cmd.StderrPipe()
	fd0, _ := cmd.StdinPipe()

	cdst := make(chan []byte)
	csrc := make(chan []byte)
	go chansend(ctx, cdst, fd1, fd2)
	go chanrecv(ctx, csrc, fd0)

	err := cmd.Start()
	if err != nil {
		logf("exec: %s: %s", args, err)
		fin <- err
		close(fin)
		return fin
	}

	go func() {
		if src != nil {
			q0, q1 := src.Dot()
			csrc <- append([]byte{}, src.Bytes()[q0:q1]...)
		}
		close(csrc)
		fin <- cmd.Wait()
		close(fin)
	}()

	// paces with ndel, the global delete counter
	delctr := atomic.LoadUint32(&ndel)

	go func() {
		for {
			select {
			case p := <-cdst:
				// Check the state of delctr before each write, if we're out of phase
				// we reload the error window to avoid a nil pointer dereference
				//
				// TODO(as): Do this for all writes
				// TODO(as): Do this for all reads
				if ctr := atomic.LoadUint32(&ndel); dst == nil || ctr != delctr {
					dst = g.afinderr(dir, cmdlabel(name, dir))
					delctr = ctr
				}
				dst.Write(p)
				repaint()
			case <-done:
				return
			}
		}
	}()
	return fin
}

func newOSCmd(dir, argv string) (name string, c Cmd) {
	x := strings.Fields(argv)
	if len(x) == 0 {
		logf("|: nothing on rhs")
		return "", nil
	}

	cmd := exec.Command(x[0], x[1:]...)
	cmd.Dir = dir
	return x[0], &OSCmd{
		Cmd: cmd,
	}
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
