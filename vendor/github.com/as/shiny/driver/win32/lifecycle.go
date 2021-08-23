// +build windows

package win32

import (
	"fmt"
	"syscall"

	"github.com/as/shiny/event/lifecycle"
)

type Lifecycle = lifecycle.Event

var LifecycleEvent func(hwnd syscall.Handle, e lifecycle.Stage)

func sendFocus(h syscall.Handle, msg uint32, wp, lp uintptr) (res uintptr) {
	switch msg {
	case WmSetfocus:
		LifecycleEvent(h, lifecycle.StageFocused)
	case WmKillfocus:
		LifecycleEvent(h, lifecycle.StageVisible)
	default:
		panic(fmt.Sprintf("unexpected focus message: %d", msg))
	}
	return DefWindowProc(h, msg, wp, lp)
}

func sendClose(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	LifecycleEvent(hwnd, lifecycle.StageDead)
	return 0
}

func sendShow(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	LifecycleEvent(hwnd, lifecycle.StageVisible)
	ShowWindow(hwnd, SwShowdefault)
	sendSize(hwnd)
	return 0
}
