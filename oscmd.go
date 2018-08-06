package main

import (
	"io"
	"os/exec"
	"strings"
)

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
