package win

import (
	"errors"
	"fmt"
	"image"
	"strings"
	"syscall"
	"unsafe"
)

var (
	ErrNoWindows  = errors.New("pid has no windows")
	ErrZeroWindow = errors.New("zero window")
	ErrBadRect    = errors.New("bad rectangle")
)

// Window is a Windows GUI window
type Window uintptr

// Open opens a process and returns its first window, if the process
// has no windows, ErrNoWindow is returned
func Open(pid int) (Window, error) {
	w, err := FromPID(pid)
	if err != nil {
		return 0, err
	}
	if len(w) == 0 {
		return 0, ErrNoWindows
	}
	w0 := Window(w[0])
	if w0 == 0 {
		return 0, ErrZeroWindow
	}
	return w0, nil
}

// Client returns the Window's client area. The client area is the
// rectangle that can be painted by an application and excludes
// the window border. The value uses absolute coordinates and
// returns a canonicalized rectangle.
func (w Window) Client() (image.Rectangle, error) {
	r32 := rect32{}
	ret, _, err := getClientRect.Call(uintptr(w), uintptr(unsafe.Pointer(&r32)))
	if ret == 0 {
		return image.ZR, err
	}
	wr, err := w.Bounds()
	if err != nil {
		return image.ZR, err
	}
	cr := r32.Rect()
	return wr.Add(Border(wr, cr)), nil
}

// Client returns window bounds. This includes the border. The value
// uses absolute coordinates and returns a canonicalized rectangle.
func (w Window) Bounds() (image.Rectangle, error) {
	r32 := rect32{}
	ret, _, err := getWindowRect.Call(uintptr(w), uintptr(unsafe.Pointer(&r32)))
	if ret == 0 {
		return image.ZR, err
	}
	return r32.Rect(), nil
}

func (w Window) Reshape(r image.Rectangle) error {
	return move(int(w), r, true)
}

var (
	u32 = syscall.MustLoadDLL("user32.dll")
	k32 = syscall.MustLoadDLL("kernel32.dll")

	getClientRect   = u32.MustFindProc("GetClientRect")
	getWindowRect   = u32.MustFindProc("GetWindowRect")
	getActiveWindow = u32.MustFindProc("GetActiveWindow")
	//windowFromPoint = u32.MustFindProc("WindowFromPoint")
	moveWindow = u32.MustFindProc("MoveWindow")
	//setWindowPos = u32.MustFindProc("SetWindowPos")

	pEnumWindows              = u32.MustFindProc("EnumWindows")
	pEnumChildWindows         = u32.MustFindProc("EnumChildWindows")
	pGetWindowTextW           = u32.MustFindProc("GetWindowTextW")
	pSendMessage              = u32.MustFindProc("SendMessageW")
	pGetWindowThreadProcessId = u32.MustFindProc("GetWindowThreadProcessId")
	pFindWindowEx             = u32.MustFindProc("FindWindowExW")
	pGetClassName             = u32.MustFindProc("GetClassNameW")
	pGetWindowRect            = u32.MustFindProc("GetWindowRect")
	pGetWindowInfo            = u32.MustFindProc("GetWindowInfo")
	pGetForegroundWindow      = u32.MustFindProc("GetForegroundWindow")
	pIsIconic                 = u32.MustFindProc("IsIconic")
	pIsZoomed                 = u32.MustFindProc("IsZoomed")
	pIsWindowEnabled          = u32.MustFindProc("IsWindowEnabled")
	pIsWindowVisible          = u32.MustFindProc("IsWindowVisible")
	pGetProcessImageFileName  = k32.MustFindProc("K32GetProcessImageFileNameW")
	pOpenProcess              = k32.MustFindProc("OpenProcess")
	pGetProcessId             = k32.MustFindProc("GetProcessId")
	pMessageBeep              = u32.MustFindProc("MessageBeep")
)

const (
	WM_GETTEXT                        = 0x000D
	WM_GETTEXTLENGTH                  = 0x000E
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

type rect32 struct {
	Min, Max point32
}
type point32 struct {
	X, Y int32
}

func (p point32) Point() image.Point {
	return image.Point{X: int(p.X), Y: int(p.Y)}
}
func (r rect32) Rect() image.Rectangle {
	r2 := image.Rectangle{Min: r.Min.Point(), Max: r.Max.Point()}.Canon()
	return r2
}

/*
func rect() (image.Rectangle, bool) {
	r32 := rect32{}
	ret, _, _ := getWindowRect.Call(active(), uintptr(unsafe.Pointer(&r32)))
	ok := ret != 0
	return r32.Rect(), ok
}
func clientRect() (image.Rectangle, bool) {
	r32 := rect32{}
	ret, _, _ := getClientRect.Call(active(), uintptr(unsafe.Pointer(&r32)))
	ok := ret != 0
	return r32.Rect(), ok
}
func active() uintptr {
	w, err := FromPID(os.Getpid())
	if len(w) == 0 {
		return 0
	}
	return w[0]
}
func rect(wid int) image.Rectangle {
	var r Rect
	rp := uintptr(unsafe.Pointer(&r))
	e, _, _ := pGetWindowRect.Call(uintptr(uint32(wid)), rp)
	if (e == 0) {
		return nil
	}
	return &r
}
*/
func activewin() int {
	r, _, _ := pGetForegroundWindow.Call()
	return int(r)
}

func minimized(wid int) bool {
	r, _, _ := pIsIconic.Call(uintptr(uint32(wid)), 0, 0)
	return r != 0
}

func maximized(wid int) bool {
	r, _, _ := pIsZoomed.Call(uintptr(uint32(wid)), 0, 0)
	return r != 0

}
func visible(wid int) bool {
	r, _, _ := pIsWindowEnabled.Call(uintptr(uint32(wid)), 0, 0)
	return r != 0
}

func enabled(wid int) bool {
	r, _, _ := pIsWindowVisible.Call(uintptr(uint32(wid)), 0, 0)
	return r != 0
}

func findwin(pwid, cwid int) int {
	r, _, _ := pFindWindowEx.Call(uintptr(uint32(pwid)),
		uintptr(uint32(cwid)),
		0, 0)
	return int(r)
}

func classof(wid int) string {
	b := make([]uint16, 1024)
	p := uintptr(unsafe.Pointer(&b[0]))
	pGetClassName.Call(uintptr(uint32(wid)), p, 1024)
	return string(syscall.UTF16ToString(b))
}

type Noise uint32

const (
	BeepBeep     Noise = 0xFFFFFFFF
	BeepStop     Noise = 0x00000010
	BeepQuestion Noise = 0x00000020
	BeepWarning  Noise = 0x00000030
	BeepAsterisk Noise = 0x00000040
)

func Beep(id Noise) {
	pMessageBeep.Call(uintptr(id))
}

// pidof gets a pid of a handle created by openprocess
func pidof(h int) int {
	r, _, _ := pGetProcessId.Call(uintptr(uint32(h)))
	return int(r)
}

func openproc(afl int, dup bool, pid int) (int, error) {
	d := 0
	if dup {
		d = 1
	}
	r, _, _ := pOpenProcess.Call(uintptr(uint32(afl)),
		uintptr(int(d)),
		uintptr(uint32(pid)))

	if pid != pidof(int(r)) {
		return 0, fmt.Errorf("openproc: %d != %d", r, pidof(int(r)))
	}
	return int(r), nil
}

func pidname(pid int) (name string, err error) {
	b := make([]uint16, 1024)

	crapfd, err := openproc(0x1000|0x0400, false, pid)
	if err != nil {
		return "", err
	}
	h := uintptr(uint32(crapfd))
	p := uintptr(unsafe.Pointer(&b[0]))
	r1, _, e1 := syscall.Syscall(pGetProcessImageFileName.Addr(), 3, h, p, 1024*2)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	if r1 > 0 {
		path := string(syscall.UTF16ToString(b))
		sp := strings.LastIndex(path, "\\")
		return path[sp+1:], nil
	}
	return "none", err
}

func ridpid(wid int) (r int, err error) {
	var pid uint32
	h := uintptr(syscall.Handle(wid))
	p := uintptr(unsafe.Pointer(&pid))

	r1, _, e1 := syscall.Syscall(pGetWindowThreadProcessId.Addr(), 2, h, uintptr(p), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return int(pid), err
}

func wintext(wid int) (text string, err error) {
	tl := SendMessage(uintptr(uint32(wid)), WM_GETTEXTLENGTH, 0, 0) * 2
	b := make([]uint16, tl+20)
	r := SendMessage(uintptr(wid), WM_GETTEXT, tl, uintptr(unsafe.Pointer(&b[0])))
	if r <= 0 {
		return "", fmt.Errorf("wintext: %d", r)
	}
	text = syscall.UTF16ToString(b)
	return text, nil
}

type box struct {
	p    int // Process id
	wids []uintptr
}

var callback = syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
	args := (*box)(unsafe.Pointer(p))
	pid := args.p
	pid2, _ := ridpid(int(h))
	if pid2 != pid {
		return 1
	}
	args.wids = append(args.wids, uintptr(h))
	return 1
})

func fromPID(pid int) (wids []uintptr, err error) {
	b := &box{pid, make([]uintptr, 0, 1024)}
	r0, _, e := pEnumWindows.Call(
		callback,
		uintptr(unsafe.Pointer(b)),
	)
	if r0 == 0 {
		return nil, e
	}
	wids = make([]uintptr, len(b.wids))
	copy(wids, b.wids)
	return wids, nil
}

func SendMessage(h uintptr, msg uint, wp, lp uintptr) uintptr {
	r, _, _ := pSendMessage.Call(uintptr(h), uintptr(msg), wp, lp)
	return r
}

func move(wid int, to image.Rectangle, paint bool) (err error) {
	p := 0
	if paint {
		p = 1
	}
	r1, _, e1 := moveWindow.Call(uintptr(wid),
		uintptr(to.Min.X),
		uintptr(to.Min.Y),
		uintptr(to.Max.X),
		uintptr(to.Max.Y),
		uintptr(p),
	)
	if r1 == 0 {
		err = e1
	}
	return
}

/*
func EnumWindows() (wids []uintptr,  err error) {
	wids = make([]int, 0, 100)
	n := 0
	syscall.Syscall(pEnumWindows.Addr(), 2, cb, 0, 0)
	wids = wids[:n]
	wids2 := make([]int, len(wids))
	copy(wids2, wids)
	return wids2, nil
}

func EnumChildWindows(rwid int) ([]int,  error) {
	wids = make([]int, 0, 100)
	n := 0
	syscall.Syscall(pEnumChildWindows.Addr(),
		3, uintptr(syscall.Handle(rwid)), cb, 0)
	wids = wids[:n]
	wids2 := make([]int, len(wids))
	copy(wids2, wids)
	return wids2, nil
}

func walk(wid int) {
	cwids, _ := EnumChildWindows(wid)
	for _, w := range cwids {
		// printrid(w)
		walk(w)
	}
}
*/
