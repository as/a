package main

import (
	"image"
	"path/filepath"
	"github.com/as/edit"
	"github.com/as/event"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"
	"github.com/as/path"
	"github.com/as/text"
	"github.com/as/text/action"
	"github.com/as/text/find"
)

func AbsOf(basedir, path string) string{
	if filepath.IsAbs(path){
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
	if len(e.To) > 0 && g.meta(e.To[0]) {

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
			if label, ok := label.(*tag.Tag); ok {
				//TODO(as): danger, edit needs a way to ensure it will only jump to an address
				edit.MustCompile(addr).Run(label.Body)
				ajump(label.Body, fn)
			}
		}
		return
	}
	isdir := false
	abspath := ""
	visible := ""
	switch {
	case filepath.IsAbs(name):
		if !path.Exists(name){
			break
		}
		abspath = filepath.Clean(name)
		visible = abspath
	case filepath.IsAbs(e.Name):
		if !path.Exists(e.Name){
			break
		}
		abspath = filepath.Join(path.DirOf(e.Name), name)
		visible = abspath
	case filepath.IsAbs(e.Basedir):
		if !path.Exists(e.Basedir){
			break
		}
		abspath = path.DirOf(e.Basedir)
		visible = filepath.Join(path.DirOf(e.Name), name)
	default:
	}
	
	var t *tag.Tag
	isdir=isdir
	if abspath == visible && path.Exists(visible){
		isdir = path.IsDir(abspath)
		t = New(actCol, path.DirOf(abspath), visible).(*tag.Tag)
	} else if realpath := filepath.Join(abspath, visible); path.Exists(realpath){
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
	
	// String literal
	for _, ed := range e.To {
		q0, q1 := find.FindNext(ed, e.P)
		ed.Select(q0, q1)
		fn := cursorNop
		if e.From == ed {
			fn = moveMouse
		}
		ajump(ed, fn)
	}
}

func (g *Grid) meta(p interface{}) bool {
	if w, ok := p.(*win.Win); ok {
		return w == g.List[0].(*tag.Tag).Win
	}
	return false
}

func VisitAll(root Plane, fn func(p Plane)) {
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
		if root.(*tag.Tag) != nil {
			VisitAll(root, fn)
		}
	case Plane:
		fn(root)
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
