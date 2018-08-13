package main

import (
	"github.com/as/edit"
	"github.com/as/rgba"
	"github.com/as/text"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

var (
	diffIns = edit.MustCompile(`,x,.+\n,x,^\\+.+\n,h`)
	diffDel = edit.MustCompile(`,x,.+\n,x,^-.+\n,h`)
)

func Diff(t *tag.Tag) {
	w, _ := t.Window.(*win.Win)
	if w == nil {
		return
	}
	w.FuncInstall(func(w *win.Win) {
		fr := w.Frame
		buf := text.BufferFrom(w.Bytes()[w.Origin() : w.Origin()+fr.Len()])
		ed, _ := text.Open(buf)
		diffIns.Run(ed)
		diffDel.Run(ed)
		pal := fr.Color.Palette
		pal.Text = rgba.Black
		pal.Back = rgba.Green
		for _, dot := range diffIns.Emit.Dot {
			w.Frame.Recolor(fr.PointOf(dot.Q0), dot.Q0, dot.Q1, pal)
		}
		pal.Back = rgba.Red
		for _, dot := range diffDel.Emit.Dot {
			w.Frame.Recolor(fr.PointOf(dot.Q0), dot.Q0, dot.Q1, pal)
		}
	})
}
