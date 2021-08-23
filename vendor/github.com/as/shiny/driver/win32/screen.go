// +build windows

package win32

import (
	"syscall"
	"unsafe"
)

// screenHWND is the handle to the "Screen window".
// The Screen window encapsulates all screen.Screen operations
// in an actual Windows window so they all run on the main thread.
// Since any messages sent to a window will be executed on the
// main thread, we can safely use the messages below.
var (
	screenHWND syscall.Handle
	screenMsgs = map[uint32]func(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) (lResult uintptr){}
)

func AddScreenMsg(fn func(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr)) uint32 {
	uMsg := currentUserWM.next()
	screenMsgs[uMsg] = func(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) uintptr {
		fn(hwnd, uMsg, wParam, lParam)
		return 0
	}
	return uMsg
}

func SendScreenMessage(uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	return SendMessage(screenHWND, uMsg, wParam, lParam)
}

func initScreenWindow() (err error) {
	const screenWindowClass = "shiny_ScreenWindow"
	swc, err := syscall.UTF16PtrFromString(screenWindowClass)
	if err != nil {
		return err
	}
	empty, err := syscall.UTF16PtrFromString("")
	if err != nil {
		return err
	}

	wc := WindowClass{
		LpszClassName: swc,
		LpfnWndProc:   syscall.NewCallback(screenWindowWndProc),
		HIcon:         hDefaultIcon,
		HCursor:       hDefaultCursor,
		HInstance:     hThisInstance,
		HbrBackground: syscall.Handle(ColorBtnface + 1),
	}
	_, err = RegisterClass(&wc)
	if err != nil {
		return err
	}

	const (
		//style = WsOverlappedWindow | WsVisible | WsChild
		style = WsOverlappedWindow
		def   = int32(CwUseDefault)
	)
	//screenHWND, err = CreateWindowEx(0, swc, empty, style, def, def, def, def, GetConsoleWindow(), 0, hThisInstance, 0)
	screenHWND, err = CreateWindowEx(0, swc, empty, style, def, def, def, def, HwndMessage, 0, hThisInstance, 0)
	return err
}

func screenWindowWndProc(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	switch uMsg {
	case msgCreateWindow:
		p := (*newWindowParams)(unsafe.Pointer(lParam))
		p.w, p.err = newWindow(p.opts)
	case msgMainCallback:
		go func() {
			mainCallback()
			SendScreenMessage(msgQuit, 0, 0)
		}()
	case msgQuit:
		PostQuitMessage(0)
	}
	fn := screenMsgs[uMsg]
	if fn != nil {
		return fn(hwnd, uMsg, wParam, lParam)
	}
	return DefWindowProc(hwnd, uMsg, wParam, lParam)
}
