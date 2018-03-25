// +build linux

package main

import (
	"image"

	"github.com/as/cursor"
)

func moveMouse(pt image.Point) {
	cursor.MoveRelativeTo(pt)
}
