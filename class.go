package main

import (
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

type Kind int

func (k Kind) List() bool {
	return k&KindCol != 0 || k&KindGrid != 0
}

const (
	KindUnknown Kind = 1 << iota
	KindPlane
	KindBody
	KindTag
	KindCol
	KindGrid
)

func KindOf(e interface{}) (Kind, Plane) {
	if e == g.Label() || e == g {
		return KindGrid, g
	}
	for _, c := range g.Kids() {
		c, _ := c.(*col.Col)
		if c == nil {
			continue
		}
		if e == c.Tag || e == c.Label() {
			return KindCol, c
		}
	}
	switch t := e.(type) {
	case *tag.Tag:
		return KindTag, t
	case col.Plane:
		return KindPlane, t
	}
	return KindUnknown, nil
}
