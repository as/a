package main

import (
	"github.com/as/edit"
	"github.com/as/text"
)

func editRefresh(ed text.Editor) {
	type r interface {
		Refresh()
	}
	type u interface {
		Upload()
	}

	if ed, ok := ed.(r); ok {
		ed.Refresh()
	}
	if ed, ok := ed.(u); ok {
		ed.Upload()
	}
}

func (g *Grid) EditRun(prog string, ed text.Editor) (ok bool) {
	//TODO(as): danger, edit needs a way to ensure it will only jump to an address
	if prog == "" {
		return false
	}
	if ed == nil {
		g.aerr("edit: ed == nil")
		return false
	}
	cmd, err := edit.Compile(prog)
	if err != nil {
		g.aerr("edit: compile: %q: %s", prog, err)
		return false
	}
	err = cmd.Run(ed)
	if err != nil {
		g.aerr("edit: run: %q: %s", prog, err)
	}

	return err == nil
}
