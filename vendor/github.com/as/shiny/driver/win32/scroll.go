// +build windows

package win32

import (
	"syscall"

	"github.com/as/shiny/event/mouse"
	"github.com/as/shiny/screen"
)

type Scroll = mouse.Event

func sendScrollEvent(hwnd syscall.Handle, _ uint32, wp, lp uintptr) (lResult uintptr) {

	// Convert from screen to window coordinates.
	p := Point{int32(uint16(lp)), int32(uint16(lp >> 16))}
	ScreenToClient(hwnd, &p)

	e := mouse.Event{
		X:         float32(p.X),
		Y:         float32(p.Y),
		Modifiers: keyModifiers(),
		Direction: mouse.DirStep,
		Button:    mouse.ButtonWheelDown,
	}
	if int16(wp>>16) > 0 {
		e.Button = mouse.ButtonWheelUp
	}

	screen.SendScroll(e)
	return
}
