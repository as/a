package main

import (
	"image"
	"os"
	"path/filepath"

	"github.com/as/event"
	"github.com/as/srv/fs"
	"github.com/as/text"
	"github.com/as/text/action"
	"github.com/as/text/find"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

type Named interface {
	Plane
	FileName() string
}
type Indexer interface {
	Lookup(interface{}) Plane
}

func AbsOf(basedir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(basedir, path)
}

var resolver = &fs.Resolver{
	Fs: newfsclient(),
}

type vis struct {
}

func (v *vis) Look(e event.Look) {
}

func (g *Grid) cwd() string {
	s, _ := os.Getwd()
	return s
}

// hiclick returns true if the nil address (q0,q1) intersects
// the highlighted selection (r0:r1)
func hiclick(r0, r1, q0, q1 int64) bool {
	x := r0 == r1 && text.Region3(r0, q0-1, q1) == 0
	return x
}

func (g *Grid) Look(e event.Look) {
	if g.meta(g.Tag) {
		return
	}

	ed := e.To[0]
	t, _ := ed.(*tag.Tag)

	p0, p1 := e.From.Dot() // pre-sweep
	// e.Q0 and e.Q1 is post sweep

	if e.Q0 == e.Q1 && !hiclick(e.Q0, e.Q1, p0, p1) {
		// one click outside old selection
		// expand the address
		a1 := expandFile(d2a(e.Q0, e.Q1), e.From)
		e.Q0, e.Q1 = a1.Dot()
	} else if e.Q0 == e.Q1 {
		// click inside old selection
		// use the old selection
		e.Q0, e.Q1 = p0, p1
	} else {
		// selection overlaps old selection
		// we don't care about that
	}

	//fmt.Printf("name and dot %q %d %d\n", e.P, e.Q0, e.Q1)
	e.P = e.P[e.Q0:e.Q1] //ed.Bytes()

	name, addr := action.SplitPath(string(e.P))
	if name == "" && addr == "" {
		return
	}

	if matches(httpLink.Plumb(&Plumbmsg{Data: e.P})) {
		return
	}
	if name == "" {
		if t == nil {
			return
		}
		if t.Window == nil {
			return
		}
		if g.EditRun(addr, t.Window) {
			ajump(ed, cursorNop)
		}
		return
	}

	if label := g.Lookup(name); label != nil {
		// TODO(as): This fails to find labels like +Error
		// and is overall useless
		t, _ := label.(*tag.Tag)
		if t == nil {
			logf("nil tag")
			return
		}
		if t.Window == nil || !g.EditRun(addr, t.Window) {
			ajump(t, cursorNop)
			return
		} else if t.Window != nil {
			ajump(t.Window, moveMouse)
			return
		}
		ajump(t, moveMouse)
		return
	}

	exists := false
	info, existsRemote := resolver.Look(fs.Path{Root: g.cwd(), Tag: e.Name, Pred: name})

	t, exists = g.Lookup(info.Path).(*tag.Tag)
	if !exists && existsRemote {
		t, _ = New(actCol, info.Dir, info.Path).(*tag.Tag)
		getcmd(t)
	}
	if t != nil {
		if g.EditRun(addr, t.Window) {
			ajump(t.Window, moveMouse)
		} else {
			ajump(t, moveMouse)
		}
		return
	}

	if !exists && !existsRemote {
		advance, _ := g.guru(e.Name, e.Q0, e.Q1)
		if !advance {
			return
		}
	}

	//TODO(as): fix this so it doesn't compare hard coded coordinates
	if e.To[0] == nil {
		VisitAll(g, func(p Named) {
			if p == nil {
				return
			}
			lookliteral(p.(*tag.Tag).Window, e, cursorNop)
		})
	} else {

		if e.To[0] != e.From {
			lookliteraltag(e.To[0], e.Q0, e.Q1, e.P)
		} else {
			lookliteral(e.To[0], e, moveMouse)
		}
	}
}

func stub(g *Grid, p Plane) bool {
	if p == nil {
		return false
	}
	if p == g.Tag {
		return true
	}
	l := g.Kids()
	for i := range l {
		c, _ := l[i].(*col.Col)
		if c != nil && p == c.Tag {
			return true
		}
	}
	return false
}

func lookliteraltag(ed text.Editor, q0, q1 int64, what []byte) {
	q0, q1 = ed.Dot()
	s0, s1 := find.FindNext(ed, q0, q1, what)
	ed.Select(s0, s1)
	ajump(ed, cursorNop)
}

func lookliteral(ed text.Editor, e event.Look, mouseFunc func(image.Point)) {
	// The behavior of look:
	//
	// Independent of the dot range, mark the given range as the starting point.
	// Advance to the end of the starting point
	// Search for the value repsesenting the range under the original starting point.
	// If the found range is identical to the starting point, no result has been found

	t0, t1 := ed.Dot()

	q0, q1 := find.FindNext(ed, e.Q1, e.Q1, e.P)
	//fmt.Printf("after q0,q1: %d,%d\n\n\t%q\n", e.Q0, e.Q1, e.P)

	if q0 == e.Q0 && q1 == e.Q1 {
		ed.Select(t0, t1)
		return
	}
	ed.Select(q0, q1)
	ajump(ed, mouseFunc)
}

func (g *Grid) meta(p interface{}) bool {
	return p == g.Tag.Label
}

func VisitAll(root Plane, fn func(p Named)) {
	type List interface {
		Kids() []Plane
	}

	switch root := root.(type) {
	case List:
		for _, k := range root.Kids() {
			VisitAll(k, fn)
		}
	case Named:
		fn(root)
	case Plane:
	case interface{}:
		panic("bad visitor")
	}
}

func (g *Grid) FindName(name string) *tag.Tag {
	for _, p := range g.List {
		c, _ := p.(*Col)
		if c == nil {
			continue
		}
		t := c.FindName(name)
		if t != nil {
			return t
		}
	}
	return nil
}

func (g *Grid) Lookup(pid interface{}) Plane {
	for _, k := range g.Kids() {
		if k, ok := k.(Indexer); ok {
			tag := k.Lookup(pid)
			if tag != nil {
				return tag
			}
		}
	}
	return nil
}

/*

type Looker struct {
	*tag.Tag // owning tag
	*win.Win // source of the address below
	Q0, Q1   int64
	P        []byte
	err      error
}

func (e *Looker) Err() error {
	if e.err == nil {
		if e.Win == nil {
			e.err = ErrNoWin
		}
		if e.Tag == nil {
			e.err = ErrNoTag
		}
	}
	return e.err
}

func (l *Looker) FromTag() bool {
	return l.Tag.Win == l.Win
}

func (e *Looker) SplitAddr() (name, addr string) {
	e.Q0, e.Q1 = expand3(e.Win, e.Q0, e.Q1)
	return action.SplitPath(string(e.Win.Bytes()[e.Q0:e.Q1]))
}

func (e *Looker) LookGrid(g *Grid) (error) {
	panic("unfinished")
	if e.Err() != nil {
		return e.Err()
	}

	name, addr := e.SplitAddr()
	if name == "" && addr == "" {
		return nil
	}

	if name == "" {
		if g.EditRun(addr, e.Tag.Window) {
			ajump(e.Tag.Window, cursorNop)
		}
		return nil
	}

	// Existing window label?
	if label := g.Lookup(name); label != nil  {
		t, _ := label.(*tag.Tag)
		if t == nil{
			logf("look d: tag is nil")
			return
		}
		if g.EditRun(addr, t.Window) {
			ajump(t.Window, moveMouse)
		}
		return
	}

	// A file on the filesystem
	info, exists := resolver.look(pathinfo{tag: e.Name, name: name})
	t, exists = g.Lookup(info.abspath).(*tag.Tag)


	panic("unfinished")
}
*/
