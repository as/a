
// +build windows

package win32

import (
	"syscall"
	"unsafe"

	"github.com/as/shiny/screen"
)

type newWindowParams struct {
	opts *screen.NewWindowOptions
	w    syscall.Handle
	err  error
}

var windowMsgs = map[uint32]func(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) (lResult uintptr){
	WmSetfocus:         sendFocus,
	WmKillfocus:        sendFocus,
	WmPaint:            sendPaint,
	msgShow:            sendShow,
	WmWindowposchanged: sendSizeEvent,
	WmClose:            sendClose,

	WmLbuttondown: mousetab[WmLbuttondown].send,
	WmLbuttonup:   mousetab[WmLbuttonup].send,
	WmMbuttondown: mousetab[WmMbuttondown].send,
	WmMbuttonup:   mousetab[WmMbuttonup].send,
	WmRbuttondown: mousetab[WmRbuttondown].send,
	WmRbuttonup:   mousetab[WmRbuttonup].send,
	WmMousemove:   mousetab[WmMousemove].send,
	WmMousewheel:  sendScrollEvent,

	WmKeydown: keytab.sendDown,
	WmKeyup:   keytab.sendUp,
	// TODO case WmSyskeydown, WmSyskeyup:

	// TODO(as): This will probably break something, let's not
	//WmInputlangchange: changeLanguage,
}

func NewWindow(opts *screen.NewWindowOptions) (syscall.Handle, error) {
	var p newWindowParams
	p.opts = opts
	SendScreenMessage(msgCreateWindow, 0, uintptr(unsafe.Pointer(&p)))
	return p.w, p.err
}

func AddWindowMsg(fn func(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr)) uint32 {
	uMsg := currentUserWM.next()
	windowMsgs[uMsg] = func(hwnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) uintptr {
		fn(hwnd, uMsg, wParam, lParam)
		return 0
	}
	return uMsg
}

func SendMessage(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	return sendMessage(hwnd, uMsg, wParam, lParam)
}

// Resize makes hwnd client rectangle opts.Width by opts.Height in size.
func Resize(h syscall.Handle, p Point) error {
	if p.X == 0 || p.Y == 0 {
		return nil
	}

	var cr Rectangle
	if err := GetClientRect(h, &cr); err != nil {
		return err
	}

	var wr Rectangle
	if err := GetWindowRect(h, &wr); err != nil {
		return err
	}

	wr.Max.X = wr.Dx() - (cr.Max.X - int32(p.X))
	wr.Max.Y = wr.Dy() - (cr.Max.Y - int32(p.Y))

	return Reshape(h, wr)
}

// Reshape makes hwnd client rectangle opts.Width by opts.Height in size.
func Reshape(h syscall.Handle, r Rectangle) error {
	return MoveWindow(h, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), false)
}

// Show shows a newly created window.
// It sends the appropriate lifecycle events, makes the window appear
// on the screen, and sends an initial size event.
//
// This is a separate step from NewWindow to give the driver a chance
// to setup its internal state for a window before events start being
// delivered.
func Show(hwnd syscall.Handle) {
	SendMessage(hwnd, msgShow, 0, 0)
}

func Release(hwnd syscall.Handle) {
	DestroyWindow(hwnd)
}

func newWindow(opts *screen.NewWindowOptions) (syscall.Handle, error) {
	// TODO(brainman): convert windowClass to *uint16 once (in initWindowClass)
	wcname, err := syscall.UTF16PtrFromString(windowClass)
	if err != nil {
		return 0, err
	}
	title, err := syscall.UTF16PtrFromString(opts.GetTitle())
	if err != nil {
		return 0, err
	}

	// h := syscall.Handle(0)
	// if opts.Overlay{
	//		h = GetConsoleWindow()
	//	}
	hwnd, err := CreateWindowEx(0,
		wcname, title,
		WsOverlappedWindow,
		CwUseDefault, CwUseDefault,
		CwUseDefault, CwUseDefault,
		0, // was console handle in experiment
		0, hThisInstance, 0)
	if err != nil {
		return 0, err
	}
	// TODO(andlabs): use proper nCmdShow
	// TODO(andlabs): call UpdateWindow()

	return hwnd, nil
}

func windowWndProc(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	fn := windowMsgs[uMsg]
	if fn != nil {
		return fn(hwnd, uMsg, wParam, lParam)
	}
	return DefWindowProc(hwnd, uMsg, wParam, lParam)
}
