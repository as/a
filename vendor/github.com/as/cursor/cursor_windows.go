package cursor

import (
	"image"
	"syscall"
)

var (
	user32       = syscall.MustLoadDLL("user32.dll")
	setCursorPos = user32.MustFindProc("SetCursorPos")
)

func MoveTo(p image.Point) bool {
	return moveTo(p)
}

func moveTo(p image.Point) bool {
	r, _, _ := setCursorPos.Call(uintptr(p.X), uintptr(p.Y))
	return r == 1
}
