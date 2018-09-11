package main

import (
	"image"

	"github.com/as/edit"
	"github.com/as/text"
)

func reload(ed text.Editor) {
	if ref, ok := ed.(interface {
		Size() image.Point
		Resize(image.Point)
	}); ok {
		ref.Resize(ref.Size())
	}
	repaint()
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
	// TODO(as): should this be in the function? Probably not
	reload(ed)
	return err == nil
}
