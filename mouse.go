// +build !linux
// +build !darwin

package main

import (
	"image"

	"github.com/as/cursor"
	"github.com/as/ms/win"
)

func moveMouse(pt image.Point) {
	cursor.MoveTo(win.ClientAbs().Min.Add(pt))
}
