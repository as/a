package main

import (
	"fmt"
	"image"

	"github.com/as/frame/tag"
)

var (
	actCol *Col
	actTag *tag.Tag
	act    *tag.Invertable
)

type Plane interface {
	Loc() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
}

func eq(a, b Plane) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Loc() == b.Loc()
}

func sizeof(r image.Rectangle) image.Point {
	return r.Max.Sub(r.Min)
}

func active2(pt image.Point, list ...Plane) (x Plane) {
	for i, w := range list {
		r := w.Loc()
		if pt.In(r) {
			return list[i]
		}
	}
	return nil
}

// Put
func active(pt image.Point, act Plane, list ...Plane) (x Plane) {
	if tag.Buttonsdown != 0 {
		return act
	}
	if act != nil {
		list = append([]Plane{act}, list...)
	}
	for i, w := range list {
		r := w.Loc()
		if pt.In(r) {
			return list[i]
		}
	}
	return act
}

func activate(pt image.Point, w Plane) {
	if tag.Buttonsdown != 0 {
		return
	}
	switch w := w.(type) {
	case *Grid:
		x := active2(pt, w.List...)
		switch x := x.(type) {
		case *tag.Tag:
			actCol = w.Col
			actTag = x
			act = x.Wtag
		case *Col:
			activate(pt, x)
		default:
			panic(fmt.Sprintf("activate: unknown plane: %T", x))
		}
	case *Col:
		actCol = w
		x := active2(pt, w.List...)
		if eq(x, w.List[0]) {
			actTag = x.(*tag.Tag)
			act = x.(*tag.Tag).Wtag
		} else {
			activate(pt, x)
		}
	case *tag.Tag:
		actTag = w
		activate(pt, active2(pt, w.W, w.Wtag))
	case *tag.Invertable:
		act = w
	}
}
