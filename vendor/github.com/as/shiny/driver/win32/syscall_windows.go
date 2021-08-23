package win32

//sys	AlphaBlend(dcdest syscall.Handle, xoriginDest int32, yoriginDest int32, wDest int32, hDest int32, dcsrc syscall.Handle, xoriginSrc int32, yoriginSrc int32, wsrc int32, hsrc int32, ftn uintptr) (err error) = msimg32.AlphaBlend
//sys	BitBlt(dcdest syscall.Handle, xdest int32, ydest int32, width int32, height int32, dcsrc syscall.Handle, xsrc int32, ysrc int32, rop uint32) (err error) = gdi32.BitBlt
//sys	CreateCompatibleBitmap(dc syscall.Handle, width int32, height int32) (bitmap syscall.Handle, err error) = gdi32.CreateCompatibleBitmap
//sys	CreateCompatibleDC(dc syscall.Handle) (newdc syscall.Handle, err error) = gdi32.CreateCompatibleDC
//sys	CreateDIBSection(dc syscall.Handle, bmi *_BITMAPINFO, usage uint32, bits **byte, section syscall.Handle, offset uint32) (bitmap syscall.Handle, err error) = gdi32.CreateDIBSection
//sys	CreateSolidBrush(color _COLORREF) (brush syscall.Handle, err error) = gdi32.CreateSolidBrush
//sys	DeleteDC(dc syscall.Handle) (err error) = gdi32.DeleteDC
//sys	DeleteObject(object syscall.Handle) (err error) = gdi32.DeleteObject
//sys	FillRect(dc syscall.Handle, rc *_RECT, brush syscall.Handle) (err error) = user32.FillRect
//sys	ModifyWorldTransform(dc syscall.Handle, x *_XFORM, mode uint32) (err error) = gdi32.ModifyWorldTransform
//sys	SelectObject(dc syscall.Handle, gdiobj syscall.Handle) (newobj syscall.Handle, err error) = gdi32.SelectObject
//sys	SetGraphicsMode(dc syscall.Handle, mode int32) (oldmode int32, err error) = gdi32.SetGraphicsMode
//sys	SetWorldTransform(dc syscall.Handle, x *_XFORM) (err error) = gdi32.SetWorldTransform
//sys	StretchBlt(dcdest syscall.Handle, xdest int32, ydest int32, wdest int32, hdest int32, dcsrc syscall.Handle, xsrc int32, ysrc int32, wsrc int32, hsrc int32, rop uint32) (err error) = gdi32.StretchBlt
//sys	GetDeviceCaps(dc syscall.Handle, index int32) (ret int32) = gdi32.GetDeviceCaps

//sys	GetConsoleWindow() (h syscall.Handle) = kernel32.GetConsoleWindow
//sys	GetDC(hwnd syscall.Handle) (dc syscall.Handle, err error) = user32.GetDC
//sys	ReleaseDC(hwnd syscall.Handle, dc syscall.Handle) (err error) = user32.ReleaseDC
//sys	sendMessage(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) = user32.SendMessageW
//sys	CreateWindowEx(exstyle uint32, className *uint16, windowText *uint16, style uint32, x int32, y int32, width int32, height int32, parent syscall.Handle, menu syscall.Handle, hInstance syscall.Handle, lpParam uintptr) (hwnd syscall.Handle, err error) = user32.CreateWindowExW
//sys	DefWindowProc(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) = user32.DefWindowProcW
//sys	DestroyWindow(hwnd syscall.Handle) (err error) = user32.DestroyWindow
//sys	DispatchMessage(msg *_MSG) (ret int32) = user32.DispatchMessageW
//sys	GetClientRect(hwnd syscall.Handle, rect *_RECT) (err error) = user32.GetClientRect
//sys	GetWindowRect(hwnd syscall.Handle, rect *_RECT) (err error) = user32.GetWindowRect
//sys   GetKeyboardLayout(threadID uint32) (locale syscall.Handle) = user32.GetKeyboardLayout
//sys   GetKeyboardState(lpKeyState *byte) (err error) = user32.GetKeyboardState
//sys	GetKeyState(virtkey int32) (keystatus int16) = user32.GetKeyState
//sys	GetMessage(msg *_MSG, hwnd syscall.Handle, msgfiltermin uint32, msgfiltermax uint32) (ret int32, err error) [failretval==-1] = user32.GetMessageW
//sys	LoadCursor(hInstance syscall.Handle, cursorName uintptr) (cursor syscall.Handle, err error) = user32.LoadCursorW
//sys	LoadIcon(hInstance syscall.Handle, iconName uintptr) (icon syscall.Handle, err error) = user32.LoadIconW
//sys	MoveWindow(hwnd syscall.Handle, x int32, y int32, w int32, h int32, repaint bool) (err error) = user32.MoveWindow
//sys	PostMessage(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult bool) = user32.PostMessageW
//sys   PostQuitMessage(exitCode int32) = user32.PostQuitMessage
//sys	RegisterClass(wc *_WNDCLASS) (atom uint16, err error) = user32.RegisterClassW
//sys	ShowWindow(hwnd syscall.Handle, cmdshow int32) (wasvisible bool) = user32.ShowWindow
//sys	ScreenToClient(hwnd syscall.Handle, lpPoint *_POINT) (ok bool) = user32.ScreenToClient
//sys   ToUnicodeEx(wVirtKey uint32, wScanCode uint32, lpKeyState *byte, pwszBuff *uint16, cchBuff int32, wFlags uint32, dwhkl syscall.Handle) (ret int32) = user32.ToUnicodeEx
//sys	TranslateMessage(msg *_MSG) (done bool) = user32.TranslateMessage
