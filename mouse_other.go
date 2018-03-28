// +build windows !linux

package main

import (
	"image"

	"github.com/as/cursor"
	"github.com/as/ms/win"
)

func moveMouse(pt image.Point) {
	logf("pt=%s", pt)
	logf("client=%s", win.ClientAbs())
	logf("pt+client=%s", win.ClientAbs().Min.Add(pt))
	cursor.MoveTo(win.ClientAbs().Min.Add(pt))
}
