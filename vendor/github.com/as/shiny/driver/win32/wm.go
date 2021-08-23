
// +build windows


package win32

import "syscall"

// Edit |tr ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz
// Edit ,x,Wm.,|tr abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ
// Edit ,x,Wm.,x,M_,c,m,
const (
	WmSetfocus         = 7
	WmKillfocus        = 8
	WmPaint            = 15
	WmClose            = 16
	WmWindowposchanged = 71
	WmKeydown          = 256
	WmKeyup            = 257
	WmSyskeydown       = 260
	WmSyskeyup         = 261
	WmMousemove        = 512
	WmMousewheel       = 522
	WmLbuttondown      = 513
	WmLbuttonup        = 514
	WmRbuttondown      = 516
	WmRbuttonup        = 517
	WmMbuttondown      = 519
	WmMbuttonup        = 520
	WmUser             = 0x0400
)

const (
	WsOverlapped   = 0x00000000
	WsMaximizebox  = 0x00010000
	WsMinimizebox  = 0x00020000
	WsThickframe   = 0x00040000
	WsSysmenu      = 0x00080000
	WsDlgframe     = 0x00400000
	WsBorder       = 0x00800000
	WsCaption      = 0x00c00000
	WsClipchildren = 0x02000000
	WsClipsiblings = 0x04000000
	WsDisabled     = 0x08000000
	WsVisible      = 0x10000000
	WsChild        = 0x40000000

	WsOverlappedWindow = WsOverlapped | WsCaption | WsSysmenu | WsThickframe | WsMinimizebox | WsMaximizebox
)

const (
	VkShift   = 16
	VkControl = 17
	VkMenu    = 18
	VkLwin    = 0x5B
	VkRwin    = 0x5C
)

const (
	MkLbutton = 0x0001
	MkMbutton = 0x0010
	MkRbutton = 0x0002
)

const (
	ColorBtnface = 15
)

const (
	IdiApplication = 32512
	IdiError       = 32513
	IdiQuestion    = 32514
	IdiWarning     = 32515
	IdiAsterisk    = 32516
	IdiWinlogo     = 32517
	IdiShield      = 32518
)

const (
	IdcAppstarting = (32650)
	IdcArrow       = (32512)
	IdcIbeam       = (32513)
	IdcWait        = (32514)
	IdcCross       = (32515)
	IdcUparrow     = (32516)
	IdcSize        = (32640)
	IdcIcon        = (32641)
	IdcSizenwse    = (32642)
	IdcSizenesw    = (32643)
	IdcSizewe      = (32644)
	IdcSizens      = (32645)
	IdcSizeall     = (32646)
	IdcNo          = (32648)
	IdcHand        = (32649)
	IdcHelp        = (32651)
)

const (
	CwUseDefault  = 0x80000000 - 0x100000000
	SwShowdefault = 10
	HwndMessage   = syscall.Handle(^uintptr(2)) // -3
	SwpNosize     = 0x0001
)

type Msg struct {
	HWND    syscall.Handle
	Message uint32
	Wp      uintptr
	Lp      uintptr
	Time    uint32
	Pt      Point
}

type WindowClass struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
}

type WindowPos struct {
	HWND            syscall.Handle
	HWNDInsertAfter syscall.Handle
	X               int32
	Y               int32
	Cx              int32
	Cy              int32
	Flags           uint32
}
