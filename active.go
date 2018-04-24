package main

import (
	"image"

	"github.com/as/ui/tag"
)

var (
	actCol *Col
	actTag *tag.Tag
	act    tag.Window
)

func actinit(g *Grid) {
	// This in particular needs to go
	actCol = g.List[0].(*Col)
	actTag = actCol.List[0].(*tag.Tag)
	act = actTag.Body
}

func active2(pt image.Point, list ...Plane) (x Plane) {
	for i, w := range list {
		if w == nil {
			continue
		}
		r := w.Loc()
		if pt.In(r) {
			return list[i]
		}
	}
	return nil
}

func activelabel(pt image.Point, t *tag.Tag) bool {
	if t.Win == nil {
		panic("FLAG: label is nil")
	}
	if pt.In(t.Win.Loc()) {
		actTag = t
		act = t.Win
		return true
	}
	return false
}

func activate(pt image.Point, w Plane) {
	switch w := w.(type) {
	case *Grid:
		if activelabel(pt, w.Tag) {
			return
		}
		x := active2(pt, w.List...)
		switch x := x.(type) {
		case *tag.Tag:
			panic("tag not allowed in column anymore")
			actCol = w.Col
			actTag = x
			act = x.Win
		case *Col:
			activate(pt, x)
		default:
			//panic(fmt.Sprintf("activate: unknown plane: %T", x))
		}
	case *Col:
		actCol = w
		if activelabel(pt, w.Tag) {
			return
		}
		x := active2(pt, w.List...)
		activate(pt, x)
	case *tag.Tag:
		actTag = w
		if w.Body != nil {
			activate(pt, active2(pt, w.Body, w.Win))
		}
	case tag.Window:
		act = w
	}
}
