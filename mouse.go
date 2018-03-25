package main

import (
	"image"

	"github.com/as/text"
	"github.com/as/ui/tag"
)

// mouseMove(pt image.Point) // defined in mouse_other.go and mouse_linux.go

func cursorNop(p image.Point) {}

func shouldCursor(p Plane) (fn func(image.Point)) {
	switch p.(type) {
	case Named:
		return cursorNop
	default:
		return moveMouse
	}
}
func ajump2(ed text.Editor, cursor bool) {
	fn := moveMouse
	if !cursor {
		fn = nil
	}
	if ed, ok := ed.(text.Jumper); ok {
		ed.Jump(fn)
	}
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
