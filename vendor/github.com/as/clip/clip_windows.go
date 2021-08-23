// +build windows
package clip

import (
	"bytes"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	CFText int = iota + 1
	CFBitmap
	CFMetaFile
	CFSymlink
	CFDif
	CFTiff
	CFOemText
	CFDib
	CFPalette
	CFPenData
	CFRIFF
	CFWave
	CFUnicode
	CFMetaFileEx
)

//go:generate go run C:\go\src\syscall\mksyscall_windows.go  -output zclip_windows.go clip_windows.go

const (
	gmMovable      = 2
	defaultBufSize = 1024 * 128
)

type Clip struct {
	gh                    syscall.Handle
	wp                    []byte
	readch, writech, done chan request
}

func New() (*Clip, error) {
	c := &Clip{
		readch:  make(chan request),
		writech: make(chan request),
		done:    make(chan request),
	}
	err := c.loop() // runs in its own goroutine
	if err != nil {
		return nil, err
	}
	return c, nil
}

func free(c *Clip) {
	GlobalFree(c.gh)
}

func (c *Clip) loop() error {
	initc := make(chan error)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		initc <- nil
		for {
			select {
			case r := <-c.readch:
				n, err := c.read(r.p)
				r.replyc <- reply{n, err}
			case r := <-c.writech:
				n, err := c.write(r.p)
				r.replyc <- reply{n, err}
			case r := <-c.done:
				err := c.close()
				r.replyc <- reply{err: err}
				return
			}
		}
		panic("never happens")
	}()
	return <-initc
}

type request struct {
	p      []byte
	replyc chan reply
}

type reply struct {
	n   int
	err error
}

func (c *Clip) Read(p []byte) (n int, err error) {
	rep := make(chan reply)
	c.readch <- request{p, rep}
	r := <-rep
	return r.n, r.err
}

func (c *Clip) Write(p []byte) (n int, err error) {
	rep := make(chan reply)
	c.writech <- request{p, rep}
	r := <-rep
	return r.n, r.err
}

func (c *Clip) Close() (err error) {
	rep := make(chan reply)
	c.done <- request{nil, rep}
	r := <-rep
	return r.err
}

func (c *Clip) close() (err error) {
	return CloseClipboard()
}

func (c *Clip) write(p []byte) (n int, err error) {
	defer func() {
		if err != nil {
			fmt.Println("write err", err)
		}
	}()
	if err = OpenClipboard(0); err != nil {
		return 0, fmt.Errorf("OpenClipboard: %s", err)
	}
	defer CloseClipboard()
	gh, err := GlobalAlloc(gmMovable, len(p)+2)
	if err != nil {
		return 0, err
	}
	h, err := GlobalLock(gh)
	if err != nil {
		return 0, fmt.Errorf("GlobalLock: %s", err)
	}
	buf := (*(*[1<<31 - 1]byte)(unsafe.Pointer(h)))[:len(p)]
	n = copy(buf[:], p)
	if err = GlobalUnlock(c.gh); err != nil {
		return 0, fmt.Errorf("GlobalUnlock: %s", err)
	}
	if err = SetClipboardData(CFUnicode, gh); err != nil {
		return 0, fmt.Errorf("SetClipboardData: %s", err)
	}
	return n, err
	// return n, CloseClipboard()
}

func (c *Clip) read(p []byte) (n int, err error) {
	OpenClipboard(0)
	defer CloseClipboard()
	ha, err := GetClipboardData(CFUnicode)
	if err != nil {
		return n, err
	}
	glen, err := GlobalSize(ha)
	if err != nil {
		return
	}
	h, err := GlobalLock(ha)
	if err != nil {
		return
	}
	buf := (*(*[1<<31 - 1]byte)(unsafe.Pointer(h)))[:glen]
	defer GlobalUnlock(ha)
	if n = bytes.Index(buf[:], []byte("\x00\x00")); n < 0 {
		n = len(buf)
	} else {
		n += 2
	}
	return copy(p, buf[:n]), err
}

// user32.dll
//sys	OpenClipboard(ha syscall.Handle) (err error) = user32.OpenClipboard
//sys	CloseClipboard() (err error) = user32.CloseClipboard
//sys	EmptyClipboard() (err error) = user32.EmptyClipboard
//sys	GetClipboardData(fmt int) (ha syscall.Handle, err error) = user32.GetClipboardData
//sys	SetClipboardData(fmt int, ha syscall.Handle) (err error) = user32.SetClipboardData

// kernel32.dll
//sys	GlobalAlloc(flag int, size int) (gh syscall.Handle, err error) = kernel32.GlobalAlloc
//sys	GlobalLock(gh syscall.Handle) (h syscall.Handle, err error) = kernel32.GlobalLock
//sys	GlobalUnlock(gh syscall.Handle) (err error) [failretval==syscall.InvalidHandle] = kernel32.GlobalUnlock
//sys	GlobalFree(gh syscall.Handle) (err error) = kernel32.GlobalFree
//sys	GlobalSize(gh syscall.Handle) (size int, err error) = kernel32.GlobalSize
