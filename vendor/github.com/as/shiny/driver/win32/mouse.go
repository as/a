
// +build windows

package win32

import (
	"syscall"

	"github.com/as/shiny/event/mouse"
	"github.com/as/shiny/screen"
)

type Mouse = mouse.Event

var MouseEvent func(hwnd syscall.Handle, e mouse.Event)

type mouseevent struct {
	dir mouse.Direction
	but mouse.Button
}

func (m *mouseevent) send(hwnd syscall.Handle, msg uint32, wp, lp uintptr) (lResult uintptr) {
	screen.SendMouse(mouse.Event{
		Direction: m.dir,
		Button:    m.but,
		X:         float32(uint16(lp)),
		Y:         float32(uint16(lp >> 16)),
		Modifiers: keyModifiers(),
	})
	return 0
}
func sendMouseEvent(hwnd syscall.Handle, msg uint32, wp, lp uintptr) (lResult uintptr) {
	return mousetab[msg].send(hwnd, msg, wp, lp)
}

var mousetab = [...]mouseevent{
	WmLbuttondown: {mouse.DirPress, mouse.ButtonLeft},
	WmMbuttondown: {mouse.DirPress, mouse.ButtonMiddle},
	WmRbuttondown: {mouse.DirPress, mouse.ButtonRight},
	WmLbuttonup:   {mouse.DirRelease, mouse.ButtonLeft},
	WmMbuttonup:   {mouse.DirRelease, mouse.ButtonMiddle},
	WmRbuttonup:   {mouse.DirRelease, mouse.ButtonRight},
	WmMousemove:   {},
}
