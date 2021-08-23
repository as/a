
// +build windows

package win32

import (
	"syscall"
	"unsafe"

	"github.com/as/shiny/event/size"
	"github.com/as/shiny/geom"
	"github.com/as/shiny/screen"
)

type Size = size.Event

var (
	SizeEvent func(hwnd syscall.Handle, e size.Event)
)

func sendSizeEvent(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	wp := (*WindowPos)(unsafe.Pointer(lParam))
	if wp.Flags&SwpNosize != 0 {
		return 0
	}
	sendSize(hwnd)
	return 0
}

func sendSize(hwnd syscall.Handle) {
	r := &Rectangle{}
	if err := GetClientRect(hwnd, r); err != nil {
		panic(err) // TODO(andlabs)
	}

	dx, dy := int(r.Dx()), int(r.Dy())
	screen.SendSize(size.Event{
		WidthPx:     dx,
		HeightPx:    dy,
		WidthPt:     geom.Pt(dx),
		HeightPt:    geom.Pt(dy),
		PixelsPerPt: 1,
	})
}
