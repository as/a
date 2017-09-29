package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/as/edit"
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

func (g *Grid) Look(e event.Look) {
	name, addr := action.SplitPath(string(e.P))
	if g.meta(e.To[0]) {
		return
	}
	if name == "" && addr == "" {
		return
	}

	// Existing window label?
	if label := g.Lookup(name); label != nil {
		fn := moveMouse
		if name == "" {
			fn = nil
		}
		if addr != "" {
			//TODO(as): danger, edit needs a way to ensure it will only jump to an address
			// we can expose an address parsing function from edit
			prog, err := edit.Compile(addr)
			if err != nil {
				g.aerr(err.Error())
				return
			}

			if t := label.(*tag.Tag); t.Body != nil {
				prog.Run(t.Body)
				ajump(t.Body, fn)
			} else {
				prog.Run(t)
				ajump(t, fn)
			}
		}
		return
	}

	isdir := false
	abspath := ""
	visible := ""
	exists := false
	switch {
	case filepath.IsAbs(name):
		if !path.Exists(name) {
			break
		}
		abspath = filepath.Clean(name)
		visible = abspath
	case filepath.IsAbs(e.Name):
		if !path.Exists(e.Name) {
			// The tag might point to a non-existent file but its parent directory
			// might be valid. In this case we should look inside the directory
			// even though the file doesn't exist. This makes +Error windows work
			// as intended
			e.Name = filepath.Dir(e.Name)
		}
		if !path.Exists(e.Name) {
			break
		}
		abspath = filepath.Join(path.DirOf(e.Name), name)
		visible = abspath
	case filepath.IsAbs(e.Basedir):
		if !path.Exists(e.Basedir) {
			break
		}
		abspath = path.DirOf(e.Basedir)
		visible = filepath.Join(path.DirOf(e.Name), name)
	default:
	}

	stat := func(name string) os.FileInfo {
		fi, _ := os.Stat(name)
		return fi
	}
	VisitAll(g, func(p Named) {
		if abspath == p.FileName() {
			exists = true
		} else if os.SameFile(stat(abspath), stat(p.FileName())) {
			exists = true
		}
	})

	var t *tag.Tag
	isdir = isdir
	if exists {
		q := g.Lookup(abspath)
		if q, ok := q.(*tag.Tag); ok {
			t = q
		}
	} else if abspath == visible && path.Exists(visible) {
		isdir = path.IsDir(abspath)
		t = New(actCol, path.DirOf(abspath), visible).(*tag.Tag)
	} else if realpath := filepath.Join(abspath, visible); path.Exists(realpath) {
		isdir = path.IsDir(realpath)
		t = New(actCol, path.DirOf(abspath), visible).(*tag.Tag)
	}
	if t != nil {
		if addr != "" {
			//TODO(as): danger, edit needs a way to ensure it will only jump to an address
			edit.MustCompile(addr).Run(t.Body)
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
			lookliteral(p.(*tag.Tag).Body, e.P, cursorNop)
		})
	} else {
		lookliteral(e.To[0], e.P, moveMouse)
	}

}
func (g *Grid) afinderr(wd string, name string) *tag.Tag {
	if !strings.HasSuffix(name, "+Errors") {
		name += "+Errors"
	}
	t := g.FindName(name)
	if t == nil {
		t = New(actCol, "", name).(*tag.Tag)
		if t == nil {
			panic("cant create tag")
		}
		moveMouse(t.Loc().Min)
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
func lookliteral(ed text.Editor, p []byte, mouseFunc func(image.Point)) {
	// String literal
	q0, q1 := find.FindNext(ed, p)
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
	return
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
		cursor(p.Loc().Min)
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
