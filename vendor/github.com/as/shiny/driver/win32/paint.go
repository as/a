// +build windows

package win32

import (
	"syscall"

	"github.com/as/shiny/event/paint"
	"github.com/as/shiny/screen"
)

type Paint = paint.Event

var PaintEvent func(hwnd syscall.Handle, e paint.Event)

func sendPaint(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	screen.SendPaint(Paint{})
	return DefWindowProc(hwnd, uMsg, wParam, lParam)
}
