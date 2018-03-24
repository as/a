package main

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"time"

	"github.com/as/event"
	"github.com/as/path"
	"github.com/as/text"
	"github.com/as/text/action"
	"github.com/as/text/find"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

func AbsOf(basedir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(basedir, path)
}

// Looks are done in the following order
// 1). Tag name - if it exists already jump to the tag, if there's an address jump
//	to that address in the tag
//
// 2). Readable absolute file - if the name matches a readable file in the namespaces
//	file system
//
// 3). Readable relative file - if the name matches the name of a readable file in
//	the tags base directory
//
// Address lookups follow for each of the three above
//
// 4). An address in the destination

var resolver = &fileresolver{ // from fs.go:/resolver/
	Fs: newfsclient(), // called in :/Grid..Look/
}

func (g *Grid) Look(e event.Look) {
	if g.meta(e.To[0]) {
		return
	}

	ed := e.To[0]
	t, _ := ed.(*tag.Tag)
	body, _ := ed.(*win.Win)
	if t != nil {
		body = t.Body
	}

	e.Q0, e.Q1 = expand3(ed, e.Q0, e.Q1)
	e.P = ed.Bytes()[e.Q0:e.Q1]

	name, addr := action.SplitPath(string(e.P))
	if name == "" && addr == "" {
		return
	}
	if name == "" {
		if g.EditRun(addr, ed) {
			ajump(ed, cursorNop)
		}
		return
	}

	// Existing window label?
	if label := g.Lookup(name); label != nil && body != nil {
		if g.EditRun(addr, body) {
			ajump(body, moveMouse)
		}
		return
	}

	g.aerr("res: %#v\n", resolver)
	info, exists := resolver.look(pathinfo{tag: e.Name, name: name})
	t, exists = g.Lookup(info.abspath).(*tag.Tag)

	if exists {
	} else if info.abspath == info.visible && path.Exists(info.visible) {
		t, _ = New(actCol, path.DirOf(info.abspath), info.visible).(*tag.Tag)
	} else if realpath := filepath.Join(info.abspath, info.visible); path.Exists(realpath) {
		t, _ = New(actCol, path.DirOf(info.abspath), info.visible).(*tag.Tag)
	}

	g.guru(e.Name, e.Q0, e.Q1)

	if t != nil {
		if g.EditRun(addr, t.Body) || g.EditRun(addr, t) {
			ajump(t.Body, moveMouse)
		} else {
			ajump(t, moveMouse)
		}
		return
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
		lookliteral(e.To[0], e, moveMouse)
	}

}
func (g *Grid) afinderr(wd string, name string) *tag.Tag {
	if !strings.HasSuffix(name, "+Errors") {
		name += "+Errors"
	}
	t := g.FindName(name)
	if t == nil {
		c := g.List[len(g.List)-1].(*Col)
		t = New(c, "", name, SizeThirdOf).(*tag.Tag)
		if t == nil {
			panic("cant create tag")
		}
		//moveMouse(t.Loc().Min)
	}
	return t
}
func (g *Grid) aerr(fm string, i ...interface{}) {
	t := g.afinderr(".", "")
	q1 := t.Body.Len()
	t.Body.Select(q1, q1)
	n := int64(t.Body.Insert([]byte(time.Now().Format(timefmt)+": "+fmt.Sprintf(fm, i...)+"\n"), q1))
	t.Body.Select(q1, q1+n)
	t.Body.Jump(cursorNop)
}
func (g *Grid) aout(fm string, i ...interface{}) {
	t := g.afinderr(".", "")
	q1 := t.Body.Len()
	t.Body.Select(q1, q1)
	n := int64(t.Body.Insert([]byte(fmt.Sprintf(fm, i...)+"\n"), q1))
	t.Body.Select(q1, q1+n)
	t.Body.Jump(cursorNop)
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

func lookliteral(ed text.Editor, e event.Look, mouseFunc func(image.Point)) {
	// The behavior of look:
	//
	// Independent of the dot range, mark the given range as the starting point.
	// Advance to the end of the starting point
	// Search for the value repsesenting the range under the original starting point.
	// If the found range is identical to the starting point, no result has been found

	t0, t1 := ed.Dot()
	g.aerr("lookliteral:  dot(%d:%d)", t0, t1)
	g.aerr("lookliteral: find(%d:%d) [%q]", e.Q0, e.Q1, e.P)
	q0, q1 := find.FindNext(ed, e.Q0, e.Q1, e.P)
	g.aerr("lookliteral: next(%d:%d)", q0, q1)
	if q0 == e.Q0 && q1 == e.Q1 {
		g.aerr("lookliteral: not found, same output(%d:%d)", q0, q1)
		ed.Select(t0, t1)
		return
	}
	g.aerr("lookliteral: found, diff output(%d:%d) != input(%d:%d)", q0, q1, e.Q0, e.Q1)
	ed.Select(q0, q1)
	ajump(ed, mouseFunc)
}

func (g *Grid) meta(p interface{}) bool {
	if w, ok := p.(*win.Win); ok {
		return w == g.List[0].(*tag.Tag).Win
	}
	return false
}

func VisitAll(root Plane, fn func(p Named)) {
	switch root := root.(type) {
	case *Grid:
		for _, k := range root.List[1:] {
			VisitAll(k, fn)
		}
	case *Col:
		for _, k := range root.List[1:] {
			VisitAll(k, fn)
		}
	case Named:
		fn(root)
	case Plane:
	case interface{}:
		panic("bad visitor")
	}
}

func (col *Col) Kids() []Plane {
	return col.List
}

func (col *Col) Dirty() bool {
	return true
}

func (grid *Grid) Lookup(pid interface{}) Plane {
	for _, k := range grid.Kids() {
		if k, ok := k.(Indexer); ok {
			tag := k.Lookup(pid)
			if tag != nil {
				return tag
			}
		}
	}
	return nil
}

func (col *Col) Lookup(pid interface{}) Plane {
	kids := col.Kids()
	if len(kids) == 0 {
		return nil
	}
	switch pid := pid.(type) {
	case int:
		if pid >= len(kids) {
			pid = len(kids) - 1
		}
		return col.Kids()[pid]
	case string:
		for i, v := range col.Kids() {
			if v, ok := v.(Named); ok {
				if v.FileName() == pid {
					return col.Kids()[i]
				}
			}
		}
	case image.Point:
		return ptInAny(pid, col.Kids()...)
	case interface{}:
		panic("")
	}
	return nil
}

type Named interface {
	Plane
	FileName() string
}
type Indexer interface {
	Lookup(interface{}) Plane
}

func ptInPlane(pt image.Point, p Plane) bool {
	if p == nil {
		return false
	}
	return pt.In(p.Loc())
}

func ptInAny(pt image.Point, list ...Plane) (x Plane) {
	for i, w := range list {
		if ptInPlane(pt, w) {
			return list[i]
		}
	}
	return nil
}

func ajump(p interface{}, cursor func(image.Point)) {
	switch p := p.(type) {
	case *tag.Tag:
		if p != nil {
			cursor(p.Loc().Min)
		}
	case text.Jumper:
		p.Jump(cursor)
	case Plane:
		if cursor == nil {
			cursor = shouldCursor(p)
		}
		cursor(p.Loc().Min)
	}
}

func cursorNop(p image.Point) {}
func shouldCursor(p Plane) (fn func(image.Point)) {
	switch p.(type) {
	case Named:
		return cursorNop
	default:
		return moveMouse
	}
}
