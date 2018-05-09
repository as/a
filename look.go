package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/as/event"
	"github.com/as/srv/fs"
	"github.com/as/text"
	"github.com/as/text/action"
	"github.com/as/text/find"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
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

var resolver = &fs.Resolver{ // from fs.go:/resolver/
	Fs: newfsclient(), // called in :/Grid..Look/
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
		if g.EditRun(addr, e.Tag.Body) {
			ajump(e.Tag.Body, cursorNop)
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
		if g.EditRun(addr, t.Body) {
			ajump(t.Body, moveMouse)
		}
		return
	}

	// A file on the filesystem
	info, exists := resolver.look(pathinfo{tag: e.Name, name: name})
	t, exists = g.Lookup(info.abspath).(*tag.Tag)


	panic("unfinished")
}
*/

func (g *Grid) cwd() string {
	s, _ := os.Getwd()
	return s
}
func (g *Grid) Look(e event.Look) {
	if g.meta(g.Tag) {
		return
	}

	ed := e.To[0]
	t, _ := ed.(*tag.Tag)

	e.Q0, e.Q1 = expand3(ed, e.Q0, e.Q1)
	name, addr := action.SplitPath(string(e.P))
	//	e.P = ed.Bytes()[e.Q0:e.Q1]
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
		if t.Body == nil {
			return
		}
		if g.EditRun(addr, t.Body) {
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
		if t.Body == nil || !g.EditRun(addr, t.Body) {
			ajump(t, cursorNop)
			return
		} else if t.Body != nil {
			ajump(t.Body, moveMouse)
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
	}
	if t != nil {
		if g.EditRun(addr, t.Body) {
			ajump(t.Body, moveMouse)
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
	if e.To[0].(*win.Win) == nil || e.To[0].(Plane).Loc().Max.Y < 48 {
		VisitAll(g, func(p Named) {
			if p == nil {
				return
			}
			lookliteral(p.(*tag.Tag).Body, e, cursorNop)
		})
	} else {
		if e.To[0] != e.From {
			lookliteraltag(e.To[0], e.Q0, e.Q1, e.P)
		} else {
			lookliteral(e.To[0], e, moveMouse)
		}
	}
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
		if t == nil {
			panic("cant create tag")
		}
		//moveMouse(t.Loc().Min)
	}
	return t
}
func (g *Grid) aerr(fm string, i ...interface{}) {
	t := g.afinderr(".", "")
	if t == nil || t.Body == nil {
		return
	}
	q1 := t.Body.Len()
	t.Body.Select(q1, q1)
	n := int64(t.Body.Insert([]byte(time.Now().Format(timefmt)+": "+fmt.Sprintf(fm, i...)+"\n"), q1))
	t.Body.Select(q1+n, q1+n)
	ajump(t.Body, cursorNop)
}
func (g *Grid) aout(fm string, i ...interface{}) {
	t := g.afinderr(".", "")
	q1 := t.Body.Len()
	t.Body.Select(q1, q1)
	n := int64(t.Body.Insert([]byte(fmt.Sprintf(fm, i...)+"\n"), q1))
	t.Body.Select(q1, q1+n)
	ajump(t.Body, cursorNop)
}

// expand3 return (r0:r1) if and only if that range is wide and
// not inside ed's dot, otherwise it returns dot
func expand3(ed text.Editor, r0, r1 int64) (int64, int64) {
	q0, q1 := ed.Dot()
	if r0 == r1 && text.Region3(r0, q0, q1) == 0 {
		return q0, q1
	}
	return r0, r1
}

func lookliteraltag(ed text.Editor, q0, q1 int64, what []byte) {
	q0, q1 = ed.Dot()
	s0, s1 := find.FindNext(ed, q0, q1, what)
	ed.Select(s0, s1)
	ajump(ed, nil)
}

func lookliteral(ed text.Editor, e event.Look, mouseFunc func(image.Point)) {
	// The behavior of look:
	//
	// Independent of the dot range, mark the given range as the starting point.
	// Advance to the end of the starting point
	// Search for the value repsesenting the range under the original starting point.
	// If the found range is identical to the starting point, no result has been found

	t0, t1 := ed.Dot()
	//	g.aerr("lookliteral:  dot(%d:%d)", t0, t1)
	//	g.aerr("lookliteral: find(%d:%d) [%q]", e.Q0, e.Q1, e.P)
	q0, q1 := find.FindNext(ed, e.Q0, e.Q1, e.P)
	//	g.aerr("lookliteral: next(%d:%d)", q0, q1)
	if q0 == e.Q0 && q1 == e.Q1 {
		//		g.aerr("lookliteral: not found, same output(%d:%d)", q0, q1)
		ed.Select(t0, t1)
		return
	}
	//	g.aerr("lookliteral: found, diff output(%d:%d) != input(%d:%d)", q0, q1, e.Q0, e.Q1)
	ed.Select(q0, q1)
	ajump(ed, mouseFunc)
}

func (g *Grid) meta(p interface{}) bool {
	return p == g.Tag.Win
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
