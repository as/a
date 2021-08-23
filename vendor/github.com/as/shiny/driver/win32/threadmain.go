// +build windows

package win32

import (
	"fmt"
	"runtime"
)

var mainCallback func()

func Main(f func()) (retErr error) {
	runtime.LockOSThread()

	if err := initCommon(); err != nil {
		return err
	}

	if err := initScreenWindow(); err != nil {
		return err
	}
	defer DestroyWindow(screenHWND)

	if err := initWindowClass(); err != nil {
		return err
	}

	// Prime the pump.
	mainCallback = f
	PostMessage(screenHWND, msgMainCallback, 0, 0)

	// Main message pump.
	m := &Msg{}
	for {
		done, err := GetMessage(m, 0, 0, 0)
		if err != nil {
			return fmt.Errorf("win32 GetMessage failed: %v", err)
		}
		if done == 0 { // WM_QUIT
			break
		}
		TranslateMessage(m)
		DispatchMessage(m)
	}

	return nil
}
