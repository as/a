package col

import (
	"fmt"
	"image"
	"io"

	"github.com/as/ui/tag"
	"github.com/as/ui/win"
)

// Table is a dimensionless box for a tagged list of Planes.
type Table struct {
	Tag  *tag.Tag
	List []Plane
}

func (t *Table) Close() error {
	t.Tag.Close()
	for _, t := range t.List {
		if t == nil {
			continue
		}
		if t, ok := t.(io.Closer); ok {
			t.Close()
		}
	}
	t.List = nil
	return nil
}

func (t *Table) Dirty() bool {
	for _, v := range t.List {
		if v.Dirty() {
			return true
		}
	}
	return false
}

func (t *Table) IDPoint(pt image.Point) (id int) {
	for id = 0; id < len(t.List); id++ {
		if pt.In(t.List[id].Bounds()) {
			break
		}
	}
	return id
}

func (t *Table) ID(w Plane) (id int) {
	for id = 0; id < len(t.List); id++ {
		if w != nil && t.List[id] != nil && w == t.List[id] {
			break
		}
	}
	return id
}

func (t *Table) Len() int {
	return len(t.List)
}
func (t *Table) Label() *win.Win {
	return t.Tag.Label
}

func (t *Table) Kid(n int) Plane {
	return t.List[n]
}
func (t *Table) Kids() []Plane {
	return t.List
}
func (t *Table) Refresh() {
	t.Tag.Refresh()
	for _, v := range t.List {
		v.Refresh()
	}
}

func (t *Table) Upload() {
	type Uploader interface {
		Upload()
		Dirty() bool
	}
	t.Tag.Upload()
	for _, t := range t.List {
		if t, ok := t.(Uploader); ok {
			t.Upload()
		}
	}
}

func (t *Table) FindName(name string) *tag.Tag {
	for _, v := range t.List {
		switch v := v.(type) {
		case *Col:
			t := v.FindName(name)
			if t != nil {
				return t
			}
		case *tag.Tag:
			if v.FileName() == name {
				return v
			}

		}
	}
	return nil
}

func (t *Table) Lookup(pid interface{}) Plane {
	type Named interface {
		Plane
		FileName() string
	}

	kids := t.Kids()
	if len(kids) == 0 {
		return nil
	}
	switch pid := pid.(type) {
	case int:
		if pid >= len(kids) {
			pid = len(kids) - 1
		}
		return t.Kids()[pid]
	case string:
		for i, v := range t.Kids() {
			if v, ok := v.(Named); ok {
				if v.FileName() == pid {
					return t.Kids()[i]
				}
			}
		}
	case image.Point:
		return ptInAny(pid, t.Kids()...)
	case interface{}:
		panic("")
	}
	return nil
}

// attach inserts w in position id, shifting the original forwards
func (t *Table) attach(w Plane, id int) {
	if id >= len(t.List) {
		t.List = append(t.List, w)
		return
	}
	t.List = append(t.List[:id], append([]Plane{w}, t.List[id:]...)...)
}

// detach (logical)
func (t *Table) detach(id int) Plane {
	if id < 0 || id >= len(t.List) {
		return nil
	}
	w := t.List[id]
	copy(t.List[id:], t.List[id+1:])
	t.List = t.List[:len(t.List)-1]
	return w
}

func (t *Table) PrintList() {
	for i, v := range t.List {
		fmt.Printf("%d: %#v\n", i, v)
	}
}

func ptInAny(pt image.Point, list ...Plane) (x Plane) {
	for i, w := range list {
		if ptInPlane(pt, w) {
			return list[i]
		}
	}
	return nil
}

func ptInPlane(pt image.Point, p Plane) bool {
	if p == nil {
		return false
	}
	return pt.In(p.Bounds())
}
