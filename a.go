package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/as/ui/tag"
)

func addrfmt(label string, q0, q1 int64) string {
	if q0 == q1 {
		return fmt.Sprintf("%s:#%d", label, q1)
	}
	return fmt.Sprintf("%s:#%d,#%d", label, q0, q1)
}

func (g *Grid) aerr(fm string, i ...interface{}) {
	t := g.afinderr(".", "")
	if t == nil || t.Window == nil {
		return
	}
	fmt.Fprintf(t, "%s: %s\n", time.Now().Format(timefmt), fmt.Sprintf(fm, i...))
	ajump(t.Window, cursorNop)
}
func (g *Grid) aout(fm string, i ...interface{}) {
	t := g.afinderr(".", "")
	fmt.Fprintf(t, fm+"\n", i...)
	ajump(t.Window, cursorNop)
}
func (g *Grid) aguru(fm string, i ...interface{}) {
	t := g.afindguru(".", "")
	fmt.Fprintf(t, "%s: %s", time.Now().Format(timefmt), fmt.Sprintf(fm, i...))
	ajump(t.Window, cursorNop)
}

func (g *Grid) afindguru(wd string, name string) *tag.Tag {
	if !strings.HasSuffix(name, "+Guru") {
		name += "+Guru"
	}
	t := g.FindName(name)
	if t == nil {
		c := g.List[len(g.List)-1].(*Col)
		t = New(c, "", name).(*tag.Tag)
		if t == nil {
			panic("cant create tag")
		}
	}
	return t
}

func (g *Grid) afinderr(wd string, name string) *tag.Tag {
	name = strings.TrimSpace(name)
	if !strings.HasSuffix(name, "+Errors") {
		name += "+Errors"
	}
	t := g.FindName(name)
	if t == nil {
		c := g.List[len(g.List)-1].(*Col)
		t = New(c, "", name).(*tag.Tag)
		//		r := t.Bounds()
		//		r1 := underText(t)
		//moveMouse(t.Bounds().Min)
	}
	return t
}
