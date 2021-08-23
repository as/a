
// +build windows


package win32

import (
	"sync"
	"syscall"
)

const windowClass = "shiny_Window"
const (
	CS_OWNDC = 32
)

const (
	msgCreateWindow = WmUser + iota
	msgMainCallback
	msgShow
	msgQuit
	msgLast
)

var (
	hDefaultIcon   syscall.Handle
	hDefaultCursor syscall.Handle
	hThisInstance  syscall.Handle
)

var currentUserWM userWM

func initCommon() (err error) {
	hDefaultIcon, err = LoadIcon(0, IdiApplication)
	if err != nil {
		return err
	}
	hDefaultCursor, err = LoadCursor(0, IdcArrow)
	if err != nil {
		return err
	}
	// TODO(andlabs) hThisInstance
	return nil
}

func initWindowClass() (err error) {
	wcname, err := syscall.UTF16PtrFromString(windowClass)
	if err != nil {
		return err
	}
	_, err = RegisterClass(&WindowClass{
		Style:         CS_OWNDC,
		LpszClassName: wcname,
		LpfnWndProc:   syscall.NewCallback(windowWndProc),
		HIcon:         hDefaultIcon,
		HCursor:       hDefaultCursor,
		HInstance:     hThisInstance,
		HbrBackground: syscall.Handle(ColorBtnface + 1),
	})
	return err
}

// userWM is used to generate private (WM_USER and above) window message IDs
// for use by screenWindowWndProc and windowWndProc.
type userWM struct {
	sync.Mutex
	id uint32
}

func (m *userWM) next() uint32 {
	m.Lock()
	if m.id == 0 {
		m.id = msgLast
	}
	r := m.id
	m.id++
	m.Unlock()
	return r
}
